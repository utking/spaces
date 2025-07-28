package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/utking/spaces/internal/adapters/web/go_echo/helpers"
	"github.com/utking/spaces/internal/application/domain"
	"github.com/utking/spaces/internal/ports"
)

func getSecretsWrapper(
	api ports.SecretService,
	userAPI ports.UsersService,
) echo.HandlerFunc {
	return func(c echo.Context) error {
		var (
			code   = http.StatusOK
			item   = new(domain.Secret)
			items  = make([]domain.Secret, 0)
			query  = new(domain.SecretSearchRequest)
			userID = GetUserID(c, userAPI)
			err    error
		)

		_ = c.Bind(query)

		itemReq := &domain.SecretSearchRequest{Tag: query.Tag, SecretID: query.SecretID}

		tags, _ := api.GetTags(c.Request().Context(), userID)
		if query.Tag != "" {
			items, err = api.GetItems(c.Request().Context(), userID, itemReq)
		}

		if query.SecretID != "" {
			encKey, _ := userAPI.GetAuthKey(c.Request().Context(), userID)
			item, err = api.GetItem(c.Request().Context(), userID, query.SecretID)
			// decode the secret
			if err == nil {
				if item.Password, err = decryptString(
					c.Request().Context(),
					api,
					encKey,
					item.EncodedSecret,
				); err != nil {
					err = errors.New("error while decoding the secret. If the encryption key has changed, you need to re-encrypt your secrets")
				}
			}
			// decode the username
			if err == nil {
				if item.Username, err = decryptString(
					c.Request().Context(),
					api,
					encKey,
					item.EncodedUsername,
				); err != nil {
					err = errors.New("error while decoding the secret. If the encryption key has changed, you need to re-encrypt your secrets")
				}
			}
		}

		if err != nil {
			code = http.StatusInternalServerError
		}

		return c.Render(
			code,
			"secrets/index.html",
			map[string]interface{}{
				"Title":      "Secrets",
				"Items":      items,
				"ItemsCount": len(items),
				"Item":       item,
				"Tags":       tags,
				"Error":      helpers.ErrorMessage(err),
				"Query":      query,
				"TagsCount":  len(tags),
			},
		)
	}
}

// putSecretUpdateWrapper is a wrapper for the secrets update handler.
// It handles the request and response for updating a secret.
func putSecretUpdateWrapper(
	api ports.SecretService,
	userAPI ports.UsersService,
) echo.HandlerFunc {
	return func(c echo.Context) error {
		var (
			encErr error
			secret = new(domain.SecretRequest)
			userID = GetUserID(c, userAPI)
		)

		_ = c.Bind(secret)
		encKey, keyErr := userAPI.GetAuthKey(c.Request().Context(), userID)
		if keyErr != nil {
			return c.JSON(
				http.StatusInternalServerError,
				map[string]interface{}{
					"Error": "Could not retrieve the encryption key",
				},
			)
		}

		updateReq := &domain.Secret{
			Name:        secret.Name,
			URL:         secret.URL,
			Description: secret.Description,
			Tags:        secret.Tags,
		}

		if updateReq.EncodedSecret, encErr = encryptString(
			c.Request().Context(), api, encKey, secret.PasswordSecretValue,
		); encErr != nil {
			return c.JSON(
				http.StatusInternalServerError,
				map[string]interface{}{"Error": "Could not encrypt the secret"})
		}

		if updateReq.EncodedUsername, encErr = encryptString(
			c.Request().Context(), api, encKey, secret.UsernameSecretValue,
		); encErr != nil {
			return c.JSON(
				http.StatusInternalServerError,
				map[string]interface{}{"Error": "Could not encrypt the username"})
		}

		if valErr := updateReq.Validate(); valErr != nil {
			return c.JSON(
				http.StatusBadRequest,
				map[string]interface{}{
					"Error": helpers.ErrorMessage(valErr),
				},
			)
		}

		_, err := api.Update(c.Request().Context(), userID, secret.SecretID, updateReq)

		if err != nil {
			return c.JSON(
				http.StatusInternalServerError,
				map[string]interface{}{
					"Error": helpers.ErrorMessage(err),
				},
			)
		}

		tag := ""
		if len(secret.Tags) > 0 {
			tag = secret.Tags[0]
		}

		return c.JSON(
			http.StatusOK,
			map[string]interface{}{
				"ID":  secret.SecretID,
				"Tag": tag,
			},
		)
	}
}

// getSecretCreateWrapper is a wrapper for the secret create handler.
// It renders a form for creating a new secret.
func getSecretCreateWrapper(
	api ports.SecretService,
	userAPI ports.UsersService,
) echo.HandlerFunc {
	return func(c echo.Context) error {
		var (
			code  = http.StatusOK
			query = new(domain.SecretSearchRequest)
		)

		_ = c.Bind(query)

		userID := GetUserID(c, userAPI)
		items, _ := api.GetItems(c.Request().Context(), userID, query)
		tags, err := api.GetTags(c.Request().Context(), userID)

		if err != nil {
			code = http.StatusInternalServerError
		}

		return c.Render(
			code,
			"secrets/create.html",
			map[string]interface{}{
				"Title":      "Create Secret",
				"Tags":       tags,
				"Items":      items,
				"ItemsCount": len(items),
				"Query":      query,
				"Error":      helpers.ErrorMessage(err),
			},
		)
	}
}

// postSecretCreateWrapper is a wrapper for the secret create handler.
// It handles the request and response for creating a secret.
// JSON response contains the error message if any and the created secret ID.
func postSecretCreateWrapper(
	api ports.SecretService,
	userAPI ports.UsersService,
) echo.HandlerFunc {
	return func(c echo.Context) error {
		var (
			encErr error
			code   = http.StatusOK
			secret = new(domain.SecretRequest)
			userID = GetUserID(c, userAPI)
		)

		encKey, keyErr := userAPI.GetAuthKey(c.Request().Context(), userID)
		if keyErr != nil {
			return c.JSON(
				http.StatusInternalServerError,
				map[string]interface{}{
					"Error": "Could not retrieve the encryption key",
				},
			)
		}

		err := c.Bind(secret)
		if err != nil {
			return c.JSON(
				http.StatusBadRequest,
				map[string]interface{}{
					"Error": helpers.ErrorMessage(err),
				},
			)
		}

		createReq := &domain.Secret{
			Name:        secret.Name,
			URL:         secret.URL,
			Description: secret.Description,
			Tags:        secret.Tags,
			// Password:    secret.PasswordSecretValue,
			// Username:    secret.UsernameSecretValue,
		}

		if createReq.EncodedSecret, encErr = encryptString(
			c.Request().Context(), api, encKey, secret.PasswordSecretValue,
		); encErr != nil {
			return c.JSON(
				http.StatusInternalServerError,
				map[string]interface{}{"Error": "Could not encrypt the secret"})
		}

		if createReq.EncodedUsername, encErr = encryptString(
			c.Request().Context(), api, encKey, secret.UsernameSecretValue,
		); encErr != nil {
			return c.JSON(
				http.StatusInternalServerError,
				map[string]interface{}{"Error": "Could not encrypt the secret"})
		}

		if valErr := createReq.Validate(); valErr != nil {
			return c.JSON(
				http.StatusBadRequest,
				map[string]interface{}{
					"Error": helpers.ErrorMessage(valErr),
				},
			)
		}

		secretID, err := api.Create(c.Request().Context(), userID, createReq)
		if err != nil {
			code = http.StatusInternalServerError
		}

		tag := ""
		if len(secret.Tags) > 0 {
			tag = secret.Tags[0]
		}

		return c.JSON(
			code,
			map[string]interface{}{
				"Error": helpers.ErrorMessage(err),
				"ID":    secretID,
				"Tag":   tag,
			},
		)
	}
}

// deleteSecretWrapper is a wrapper for the secret delete handler.
// It handles the request and response for deleting a secret.
// JSON response contains the error message if any.
func deleteSecretWrapper(
	api ports.SecretService,
	userAPI ports.UsersService,
) echo.HandlerFunc {
	return func(c echo.Context) error {
		var (
			secretID = helpers.GetIDParam(c)
			code     = http.StatusOK
		)

		userID := GetUserID(c, userAPI)
		err := api.Delete(c.Request().Context(), userID, secretID)

		if err != nil {
			code = http.StatusInternalServerError
		}

		return c.JSON(
			code,
			map[string]interface{}{
				"Error": helpers.ErrorMessage(err),
			},
		)
	}
}

// getExportSecretsWrapper is a wrapper-handler for showing the the secrets export page.
func getExportSecretsWrapper() echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.Render(
			http.StatusOK,
			"secrets/export.html",
			map[string]interface{}{
				"Title":    "Export Secrets",
				"Password": domain.GenerateRandomString(32),
			},
		)
	}
}

// postExportSecretsWrapper is a wrapper-handler for exporting secrets.
func postExportSecretsWrapper(
	api ports.SecretService,
	userAPI ports.UsersService,
	secretsAPI ports.SecretService,
) echo.HandlerFunc {
	return func(c echo.Context) error {
		const fileName = "secrets_export.json"
		var (
			eFile  *os.File
			err    error
			req    = new(domain.SecretExportRequest)
			userID = GetUserID(c, userAPI)
		)

		authKey, akErr := userAPI.GetAuthKey(c.Request().Context(), userID)
		if akErr != nil {
			return c.Render(
				http.StatusInternalServerError,
				"secrets/export.html",
				map[string]interface{}{
					"Title": "Export Secrets",
					"Error": "Could not retrieve the encryption key to proceed with export",
					"Query": req,
				},
			)
		}

		_ = c.Bind(req)
		if err = req.Validate(); err == nil {
			secrets, sErr := api.GetItemsMap(c.Request().Context(), userID, nil)
			if sErr != nil {
				return c.Render(
					http.StatusInternalServerError,
					"secrets/export.html",
					map[string]interface{}{
						"Title": "Export Secrets",
						"Error": helpers.ErrorMessage(sErr),
						"Query": req,
					},
				)
			}

			var (
				dErr      error
				decErrors []string
			)

			// get deciphered secrets for every secret that is not empty
			for idx := range secrets {
				// use reference to change the item in the slice
				item := &secrets[idx]

				if len(item.EncodedPassword) == 0 {
					continue // Skip if nonce or secret is not set
				}

				if item.Password, dErr = decryptString(
					c.Request().Context(),
					secretsAPI,
					authKey,
					item.EncodedPassword,
				); dErr != nil {
					decErrors = append(
						decErrors,
						fmt.Sprintf("Error decrypting secret %s: %v", item.Name, dErr),
					)
				}

				if item.Username, dErr = decryptString(
					c.Request().Context(),
					secretsAPI,
					authKey,
					item.EncodedUsername,
				); dErr != nil {
					decErrors = append(
						decErrors,
						fmt.Sprintf("Error decrypting secret %s: %v", item.Name, dErr),
					)
				}
			}

			// stop if there are any decryption errors on export
			if len(decErrors) > 0 {
				return c.Render(
					http.StatusInternalServerError,
					"secrets/export.html",
					map[string]interface{}{
						"Title": "Export Secrets",
						"Error": strings.Join(decErrors, "; "),
						"Query": req,
					},
				)
			}

			eFile, err = saveSecretsToFile(
				c.Request().Context(),
				secrets,
				[]byte(req.Password),
				secretsAPI,
			)

			if err != nil {
				return c.Render(
					http.StatusInternalServerError,
					"secrets/export.html",
					map[string]interface{}{
						"Title": "Export Secrets",
						"Error": helpers.ErrorMessage(err),
						"Query": req,
					},
				)
			}

			defer func() {
				_ = os.Remove(eFile.Name()) // Clean up the temporary file after sending
			}()

			return c.Attachment(eFile.Name(), fileName)
		}

		req.Password = domain.GenerateRandomString(32)

		return c.Render(
			http.StatusBadRequest,
			"secrets/export.html",
			map[string]interface{}{
				"Title": "Export Secrets",
				"Error": helpers.ErrorMessage(err),
				"Query": req,
			},
		)
	}
}

// saveSecretsToFile is a helper function that saves the secrets data to a file.
// The file content is encrypted using the user's auth_key.
func saveSecretsToFile(
	ctx context.Context,
	items []domain.SecretExportItem,
	authKey []byte,
	secretsAPI ports.SecretService,
) (*os.File, error) {
	data, err := json.Marshal(items)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal secrets data: %w", err)
	}

	eFile, err := os.CreateTemp("", "secrets_export*.json")
	if err != nil {
		return nil, fmt.Errorf("failed to create temp file: %w", err)
	}

	// Encrypt the data using the user's auth key
	nonce, encData, encErr := secretsAPI.Encrypt(
		ctx,
		&domain.SecretEncodeRequest{
			PlainText: data,
		},
		authKey,
	)

	if encErr != nil {
		return nil, errors.New("failed to encrypt secrets data")
	}

	_, err = eFile.Write(nonce)
	if err != nil {
		return nil, fmt.Errorf("failed to write nonce to file: %w", err)
	}

	_, err = eFile.Write(encData)

	if err != nil {
		return nil, fmt.Errorf("failed to write data to file: %w", err)
	}

	if err = eFile.Close(); err != nil {
		return nil, fmt.Errorf("failed to close file: %w", err)
	}

	return eFile, nil
}

// getSearchSecretsWrapper is a wrapper for the search secrets handler.
func getSearchSecretsWrapper(
	api ports.SecretService,
	userAPI ports.UsersService,
) echo.HandlerFunc {
	type SecretItem struct {
		Tag  string `json:"tag"`
		ID   string `json:"id"`
		Text string `json:"text"`
	}

	return func(c echo.Context) error {
		term := c.QueryParam("term")
		if term == "" {
			return c.JSON(
				http.StatusOK,
				map[string]interface{}{
					"items": []SecretItem{},
				},
			)
		}

		userID := GetUserID(c, userAPI)
		items, _ := api.SearchItemsByTerm(c.Request().Context(), userID, &domain.SecretRequest{
			Name:        term,
			Username:    term,
			URL:         term,
			Description: term,
			RequestPageMeta: domain.RequestPageMeta{
				Limit: 10,
			},
		})

		filteredItems := make([]SecretItem, 0, len(items))
		for _, item := range items {
			filteredItems = append(filteredItems, SecretItem{
				ID:   item.ID,
				Tag:  item.Tags[0], // Assuming at least one tag exists
				Text: item.Name,
			})
		}

		return c.JSON(
			http.StatusOK,
			map[string]interface{}{
				"items": filteredItems,
			},
		)
	}
}

// getSecretsRotateKeyWrapper is a wrapper for the secrets rotate key handler.
func getSecretsRotateKeyWrapper() echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.Render(
			http.StatusOK,
			"secrets/rotate-key.html",
			map[string]interface{}{
				"Title": "Rotate Encryption Key",
			},
		)
	}
}

// postSecretsRotateKeyWrapper is a wrapper for the secrets rotate key handler.
func postSecretsRotateKeyWrapper(
	api ports.SecretService,
	userAPI ports.UsersService,
	logger ports.LoggingService,
) echo.HandlerFunc {
	return func(c echo.Context) error {
		userID := GetUserID(c, userAPI)
		if userID == "" {
			return c.Render(
				http.StatusUnauthorized,
				"secrets/rotate-key.html",
				map[string]interface{}{
					"Title": "Rotate Encryption Key",
					"Error": "Unauthorized",
				},
			)
		}

		// Rotate the user's secrets
		err := rotateUserSecrets(c, api, userAPI, userID, logger)
		if err != nil {
			return c.Render(
				http.StatusInternalServerError,
				"secrets/rotate-key.html",
				map[string]interface{}{
					"Title": "Rotate Encryption Key",
					"Error": helpers.ErrorMessage(err),
				},
			)
		}

		return c.Render(
			http.StatusOK,
			"secrets/rotate-key.html",
			map[string]interface{}{
				"Title":   "Rotate Encryption Key",
				"Success": "Encryption key rotated successfully",
			},
		)
	}
}

// rotateUserSecrets rotates the encryption key for a user's secrets.
func rotateUserSecrets(
	ctx echo.Context,
	api ports.SecretService,
	userAPI ports.UsersService,
	userID string,
	logger ports.LoggingService,
) error {
	// Rotation process
	// 0. get the current encryption key for the user
	// 1. generate a new encryption key
	// 2. get all secrets for the user
	// 3. re-encrypt each secret with the new key and store in a map
	// 4. in a transaction, update the secrets in the database with the new encrypted values
	// 5. update the encryption key in the user profile if the update is successful

	encKey, keyErr := userAPI.GetAuthKey(ctx.Request().Context(), userID)
	if keyErr != nil {
		return fmt.Errorf("failed to retrieve encryption key for user %s: %w", userID, keyErr)
	}

	newEncKey := []byte(domain.GenerateRandomString(32))
	if len(newEncKey) == 0 {
		return errors.New("failed to generate a new encryption key")
	}

	secrets, iErr := api.GetItems(ctx.Request().Context(), userID, nil)
	if iErr != nil {
		return fmt.Errorf("failed to retrieve secrets for user %s: %w", userID, iErr)
	}

	reencryptedSecrets := make(map[string]domain.EncryptSecret, len(secrets))

	for _, secret := range secrets {
		// get the secret item (with the already decoded password)
		item, err := api.GetItem(ctx.Request().Context(), userID, secret.ID)
		if err != nil {
			return fmt.Errorf("failed to get secret item %s: %w", secret.ID, err)
		}

		var (
			decoded         string
			encodedUsername []byte
			encodedPassword []byte
			encErr          error
			decErr          error
		)

		// decode the secret
		if decoded, decErr = decryptString(ctx.Request().Context(), api, encKey, item.EncodedSecret); decErr != nil {
			return errors.New("error decoding the secret")
		}

		// encode the secret with the new encryption key and store it in the map with the nonce and secret ID
		encodedPassword, encErr = encryptString(ctx.Request().Context(), api, newEncKey, decoded)
		if encErr != nil {
			return fmt.Errorf("failed to encrypt secret %s with new key: %w", secret.ID, encErr)
		}

		// decode the username
		if decoded, decErr = decryptString(ctx.Request().Context(), api, encKey, item.EncodedUsername); decErr != nil {
			return errors.New("error decoding the username")
		}

		// encode the username with the new encryption key and store it in the map with the nonce and secret ID
		encodedUsername, encErr = encryptString(ctx.Request().Context(), api, newEncKey, decoded)
		if encErr != nil {
			return fmt.Errorf("failed to encrypt username %s with new key: %w", secret.ID, encErr)
		}

		reencryptedSecrets[secret.ID] = domain.EncryptSecret{
			ID:       secret.ID,
			Password: encodedPassword,
			Username: encodedUsername, // Assuming the password is the same as the secret
		}
	}

	err := api.UpdateEncryptedSecrets(ctx.Request().Context(), userID, reencryptedSecrets)
	if err != nil {
		return fmt.Errorf("failed to update secrets in transaction: %w", err)
	}

	// Update the user's encryption key.
	// INFO: if the update fails, the secrets will be unavailable
	// and will have to be restored from backup.
	err = userAPI.UpdateAuthKey(ctx.Request().Context(), userID, newEncKey)
	if err != nil {
		logger.Error(
			ctx.Request().Context(),
			"Failed to update user encryption key",
			ports.NewLoggerBag("error", err.Error()),
			ports.NewLoggerBag("user_id", userID),
		)

		return fmt.Errorf("failed to update user encryption key: %w", err)
	}

	logger.Info(
		ctx.Request().Context(),
		"User encryption key rotated successfully",
		ports.NewLoggerBag("user_id", userID),
	)

	return nil
}

func encryptString(
	ctx context.Context,
	api ports.SecretService,
	encKey []byte,
	plainText string,
) ([]byte, error) {
	// Empty string encryption is a shortcut to avoid unnecessary processing.
	if len(plainText) == 0 {
		return nil, nil
	}

	nonce, encoded, err := api.Encrypt(
		ctx,
		&domain.SecretEncodeRequest{
			PlainText: []byte(plainText),
		},
		encKey,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to encrypt string: %w", err)
	}

	result := make([]byte, 0, len(nonce)+len(encoded))
	result = append(result, nonce...)
	result = append(result, encoded...)

	return result, nil
}

func decryptString(
	ctx context.Context,
	api ports.SecretService,
	encKey []byte,
	encoded []byte,
) (string, error) {
	if len(encoded) == 0 {
		return "", nil // No data to decrypt
	}

	if len(encoded) <= 12 {
		return "", errors.New("invalid encoded data length")
	}

	decoded, err := api.Decrypt(
		ctx,
		encoded[:12], // Nonce is the first 12 bytes
		encoded[12:], // Encoded data starts after the nonce
		encKey,
	)

	if err != nil {
		return "", fmt.Errorf("failed to decrypt string: %w", err)
	}

	return string(decoded), nil
}
