package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/utking/spaces/internal/adapters/db"
	"github.com/utking/spaces/internal/adapters/web/go_echo/helpers"
	"github.com/utking/spaces/internal/application/domain"
	"xorm.io/builder"
)

const naString = "N/A"

// GetUsers retrieves a list of users from the database based on the provided request parameters.
func (a *Adapter) GetUsers(ctx context.Context, req *domain.UserRequest) ([]domain.User, error) {
	var dbItems []db.User

	sqlBuilder := builder.Dialect(sqlDialect).
		Select(
			"user.id as id",
			"user.username",
			"a.item_name as role_name",
			"user.email as email",
			"user.status as status",
		).
		From(db.User{}.TableName()).
		LeftJoin("auth_assignment a", "a.user_id = user.id").
		OrderBy("username")

	if req != nil {
		if req.Username != "" {
			sqlBuilder = sqlBuilder.Where(builder.Like{"username", req.Username})
		}

		if req.Email != "" {
			sqlBuilder = sqlBuilder.Where(builder.Like{"user.email", req.Email})
		}

		if req.Status != nil {
			sqlBuilder = sqlBuilder.Where(builder.Eq{"user.status": *req.Status})
		}
	}

	sql, err := sqlBuilder.ToBoundSQL()
	if err != nil {
		return nil, err
	}

	err = a.db.SelectContext(ctx, &dbItems, sql)
	if err != nil {
		return nil, err
	}

	items := make([]domain.User, len(dbItems))

	for i, item := range dbItems {
		items[i] = domain.User{
			ID:       item.ID,
			Username: item.Username,
			Email:    item.Email,
			Status:   item.Status,
		}

		if item.RoleName != nil {
			items[i].RoleName = *item.RoleName
		} else {
			items[i].RoleName = naString // Default role if not set
		}
	}

	return items, nil
}

// GetUsersCount retrieves the count of users from the database based on the provided request parameters.
func (a *Adapter) GetUsersCount(ctx context.Context, req *domain.UserRequest) (int64, error) {
	var count int64

	sqlBuilder := builder.Dialect(sqlDialect).
		Select("COUNT(1) as count").
		From(db.User{}.TableName())

	if req != nil {
		if req.Username != "" {
			sqlBuilder = sqlBuilder.Where(builder.Like{"username", req.Username})
		}

		if req.Email != "" {
			sqlBuilder = sqlBuilder.Where(builder.Like{"user.email", req.Email})
		}

		if req.Status != nil {
			sqlBuilder = sqlBuilder.Where(builder.Eq{"user.status": *req.Status})
		}
	}

	sql, err := sqlBuilder.ToBoundSQL()
	if err != nil {
		return 0, err
	}

	err = a.db.GetContext(ctx, &count, sql)

	return count, err
}

// GetUser retrieves a user by ID from the database.
func (a *Adapter) GetUser(ctx context.Context, id string) (*domain.User, error) {
	var dbItem db.User

	sqlBuilder := builder.Dialect(sqlDialect).
		Select(
			"user.id as id", "user.username", "user.email as email",
			"user.status as status",
			"a.item_name as role_name",
			"user.created_at",
			"user.updated_at",
		).
		From(dbItem.TableName()).
		LeftJoin("auth_assignment a", "a.user_id = user.id").
		Where(builder.Eq{"user.id": id})

	sqlStr, err := sqlBuilder.ToBoundSQL()
	if err != nil {
		return nil, err
	}

	err = a.db.GetContext(ctx, &dbItem, sqlStr)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("user with id %s not found", id)
		}

		return nil, err
	}

	item := &domain.User{
		ID:        dbItem.ID,
		Username:  dbItem.Username,
		Email:     dbItem.Email,
		Status:    dbItem.Status,
		CreatedAt: dbItem.CreatedAt,
		UpdatedAt: dbItem.UpdatedAt,
	}

	if dbItem.RoleName != nil {
		item.RoleName = *dbItem.RoleName
	} else {
		item.RoleName = naString // Default role if not set
	}

	return item, nil
}

// CreateUser creates a new user in the database.
// It returns the ID of the newly created user and a token for the user verification.
// The password is hashed before storing it in the database.
// The auth_key is generated randomly and stored in the database.
// The created_at and updated_at fields are set to the current time.
func (a *Adapter) CreateUser(ctx context.Context, req *domain.User) (id, token string, err error) {
	var dbItem db.User

	id = helpers.GenerateUUID()
	dbItem.Username = req.Username
	dbItem.Email = req.Email
	dbItem.PasswordHash, _ = domain.GetPasswordHash(req.Password)
	dbItem.ActivationToken = domain.GenerateRandomString(32)

	// validate the request
	if vErr := dbItem.Validate(); vErr != nil {
		return "", "", vErr
	}

	insertMap := builder.Eq{
		"id":                       id,
		"username":                 dbItem.Username,
		"email":                    dbItem.Email,
		"password_hash":            dbItem.PasswordHash,
		"status":                   domain.UserInactive,
		"auth_key":                 domain.GenerateRandomString(32),
		"account_activation_token": dbItem.ActivationToken,
	}

	sqlBuilder := builder.Dialect(sqlDialect).Into(dbItem.TableName()).Insert(insertMap)

	sqlStr, args, err := sqlBuilder.ToSQL()
	if err != nil {
		return "", "", err
	}

	_, err = a.db.ExecContext(ctx, sqlStr, args...)
	if err != nil {
		// Check for unique constraint violation
		if sqliteUniqViolation(err) {
			return "", "", errors.New("user with such name or email already exists")
		}

		return "", "", err
	}

	// Set the user role to "user" by default
	if uErr := a.setUserRole(ctx, id, domain.UserRoleUser); uErr != nil {
		return "", "", errors.New("failed to set user role")
	}

	return id, dbItem.ActivationToken, err
}

// setUserRole sets the role for a user in the database.
func (a *Adapter) setUserRole(ctx context.Context, userID, roleName string) error {
	sqlBuilder := builder.Dialect(sqlDialect).
		Insert(builder.Eq{
			"user_id":    userID,
			"item_name":  roleName,
			"created_at": time.Now().Format(time.DateTime),
		}).
		Into("auth_assignment")

	sql, err := sqlBuilder.ToBoundSQL()
	if err != nil {
		return errors.New("SQL error setting user role")
	}

	_, err = a.db.ExecContext(ctx, sql)
	if err != nil {
		return err
	}

	return nil
}

// updateUserRole updates the role for a user in the database.
func (a *Adapter) updateUserRole(ctx context.Context, userID, roleName string) error {
	// check if the auth_assignment record exists
	var count int64

	sqlBuilder := builder.Dialect(sqlDialect).
		Select("COUNT(1) as count").
		From("auth_assignment").
		Where(builder.Eq{"user_id": userID})

	sql, err := sqlBuilder.ToBoundSQL()
	if err != nil {
		return fmt.Errorf("SQL error checking user role: %w", err)
	}

	err = a.db.GetContext(ctx, &count, sql)
	if err != nil {
		return fmt.Errorf("failed to check user role: %w", err)
	}

	if count == 0 {
		// if the record does not exist, insert a new one
		return a.setUserRole(ctx, userID, roleName)
	}

	// if the record exists, update it
	sqlBuilder = builder.Dialect(sqlDialect).
		Update(builder.Eq{"item_name": roleName}).
		From("auth_assignment").
		Where(builder.Eq{"user_id": userID})

	sql, err = sqlBuilder.ToBoundSQL()
	if err != nil {
		return fmt.Errorf("SQL for updating user role: %w", err)
	}

	if _, err = a.db.ExecContext(ctx, sql); err != nil {
		return fmt.Errorf("failed to update user role: %w", err)
	}

	return nil
}

// UpdateUser updates an existing user in the database.
// It returns the number of rows affected by the update.
// The password, if non-empty, is hashed before storing
// it in the database.
func (a *Adapter) UpdateUser(ctx context.Context, id string, req *domain.UserUpdate) (int64, error) {
	var dbItem db.User

	dbItem.Email = req.Email
	dbItem.Status = req.Status

	if req.Password != "" {
		dbItem.PasswordHash, _ = domain.GetPasswordHash(req.Password)
	}

	// validate the request
	if err := dbItem.ValidateUpdate(); err != nil {
		return 0, err
	}

	updateMap := builder.Eq{
		// "username": dbItem.Username, - not allowed to update
		"email":      dbItem.Email,
		"status":     dbItem.Status,
		"updated_at": time.Now().Format(time.DateTime),
	}

	if dbItem.PasswordHash != "" {
		updateMap["password_hash"] = dbItem.PasswordHash
	}

	sqlBuilder := builder.Dialect(sqlDialect).
		Update(updateMap).
		From(dbItem.TableName()).
		Where(builder.Eq{"id": id})

	sql, args, err := sqlBuilder.ToSQL()
	if err != nil {
		return 0, err
	}

	result, err := a.db.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, err
	}

	// Update user role in `auth_assignment` table if role name is provided
	if req.RoleName != "" {
		if err = a.updateUserRole(ctx, id, req.RoleName); err != nil {
			return 0, errors.New("failed to update user role")
		}
	}

	return result.RowsAffected()
}

// ChangePassword updates the password for a user in the database.
func (a *Adapter) ChangePassword(ctx context.Context, id, newPassword string) error {
	if newPassword == "" {
		return errors.New("password cannot be empty")
	}

	var dbItem db.User

	dbItem.PasswordHash, _ = domain.GetPasswordHash(newPassword)

	// validate the password length
	if len(newPassword) < 8 || len(newPassword) > 32 {
		return errors.New("new password length must be between 8 and 32 characters")
	}

	sqlBuilder := builder.Dialect(sqlDialect).
		Update(builder.Eq{"password_hash": dbItem.PasswordHash}).
		From(dbItem.TableName()).
		Where(builder.Eq{"id": id})

	sql, args, err := sqlBuilder.ToSQL()
	if err != nil {
		return errors.New("SQL error updating user password")
	}

	result, err := a.db.ExecContext(ctx, sql, args...)
	if err != nil {
		return errors.New("failed to update user password")
	}

	_, err = result.RowsAffected()

	return err
}

// DeleteUser marks the user as deleted by setting the status to UserDeleted.
// This is a soft delete, the user is not actually removed from the database.
func (a *Adapter) DeleteUser(ctx context.Context, id string) error {
	sqlBuilder := builder.Dialect(sqlDialect).
		Update(builder.Eq{"status": domain.UserInactive}).
		From(db.User{}.TableName()).
		Where(builder.Eq{"id": id})

	sql, err := sqlBuilder.ToBoundSQL()
	if err != nil {
		return errors.New("SQL error deleting user")
	}

	_, err = a.db.ExecContext(ctx, sql)

	return err
}

// GetUserByUsername retrieves a user by username from the database.
// The respose contains only the username, email, and password hash.
// The user MUST be active (status = UserActive) and an admin.
func (a *Adapter) GetUserByUsername(ctx context.Context, username string) (*domain.User, error) {
	var dbItem db.User

	sqlBuilder := builder.Dialect(sqlDialect).
		Select("user.id", "user.username", "password_hash", "a.item_name as role_name").
		From(dbItem.TableName()).
		LeftJoin("auth_assignment a", "a.user_id = user.id").
		Where(builder.And(
			builder.Eq{"user.username": username},
			builder.Eq{"user.status": domain.UserActive},
		))

	sql, err := sqlBuilder.ToBoundSQL()
	if err != nil {
		return nil, errors.New("error validating user")
	}

	err = a.db.GetContext(ctx, &dbItem, sql)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by username [%s]: %w", username, err)
	}

	item := &domain.User{
		ID:           dbItem.ID,
		Username:     dbItem.Username,
		PasswordHash: dbItem.PasswordHash,
	}

	if dbItem.RoleName != nil {
		item.RoleName = *dbItem.RoleName
	} else {
		item.RoleName = naString // Default role if not set
	}

	return item, nil
}

// SetUserVerified sets the user as verified by updating the account_activation_token to NULL.
func (a *Adapter) SetUserVerified(ctx context.Context, token string) (*domain.User, error) {
	if token == "" {
		return nil, errors.New("verification token is required")
	}

	// Select the user with the given token and status UserInactive
	var dbItem db.User

	sqlBuilder := builder.Dialect(sqlDialect).
		Select("id").
		From(dbItem.TableName()).
		Where(
			builder.And(
				builder.Eq{"account_activation_token": token},
				builder.Eq{"status": domain.UserInactive},
			),
		)

	sqlStr, err := sqlBuilder.ToBoundSQL()
	if err != nil {
		return nil, fmt.Errorf("SQL error getting user by token: %w", err)
	}

	err = a.db.GetContext(ctx, &dbItem, sqlStr)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("user with the given token was not found")
		}

		return nil, fmt.Errorf("failed to get user by token [%s]: %w", token, err)
	}

	// Update the user by setting account_activation_token to NULL
	sqlBuilder = builder.Dialect(sqlDialect).
		Update(builder.Eq{
			"account_activation_token": nil,
			"status":                   domain.UserActive,
		}).
		From(db.User{}.TableName()).
		Where(builder.Eq{"id": dbItem.ID})

	sqlStr, err = sqlBuilder.ToBoundSQL()
	if err != nil {
		return nil, err
	}

	result, err := a.db.ExecContext(ctx, sqlStr)
	if err != nil {
		return nil, err
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return nil, errors.New("user with the given token was not found")
	}

	return a.GetUser(ctx, dbItem.ID)
}

// GetUserAuthKey retrieves the auth_key for a user by their ID.
func (a *Adapter) GetUserAuthKey(ctx context.Context, userID string) ([]byte, error) {
	var dbItem db.User

	sqlBuilder := builder.Dialect(sqlDialect).
		Select("auth_key").
		From(dbItem.TableName()).
		Where(builder.Eq{"id": userID})

	sqlStr, err := sqlBuilder.ToBoundSQL()
	if err != nil {
		return []byte{}, fmt.Errorf("SQL error getting user auth_key: %w", err)
	}

	err = a.db.GetContext(ctx, &dbItem, sqlStr)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return []byte{}, fmt.Errorf("user with id %s not found", userID)
		}

		return []byte{}, fmt.Errorf("failed to get user auth_key: %w", err)
	}

	if dbItem.AuthKey == "" {
		return []byte{}, errors.New("auth_key is not set")
	}

	return []byte(dbItem.AuthKey), nil
}

// UpdateUserAuthKey updates the auth_key for a user by their ID.
func (a *Adapter) UpdateUserAuthKey(ctx context.Context, userID string, newEncKey []byte) error {
	if len(newEncKey) == 0 {
		return errors.New("new auth_key cannot be empty")
	}

	sqlBuilder := builder.Dialect(sqlDialect).
		Update(builder.Eq{"auth_key": string(newEncKey)}).
		From(db.User{}.TableName()).
		Where(builder.Eq{"id": userID})

	sqlStr, args, err := sqlBuilder.ToSQL()
	if err != nil {
		return fmt.Errorf("SQL error updating user auth_key: %w", err)
	}

	_, err = a.db.ExecContext(ctx, sqlStr, args...)
	if err != nil {
		return fmt.Errorf("failed to update user auth_key: %w", err)
	}

	return nil
}
