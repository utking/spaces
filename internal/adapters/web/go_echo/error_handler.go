package goecho

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"gogs.utking.net/utking/spaces/internal/adapters/web/go_echo/helpers"
)

// HTTPErrorHandler is a custom error handler for the Echo framework.
func HTTPErrorHandler(err error, c echo.Context) {
	templateFile := "errors/error.html"

	code := http.StatusInternalServerError

	if c.Response().Committed {
		return
	}

	// Check if there is an existing code to use
	var httpError *echo.HTTPError
	if errors.As(err, &httpError) {
		code = httpError.Code
	}

	if httpError != nil {
		switch httpError.Code {
		case http.StatusForbidden, http.StatusNotFound, http.StatusInternalServerError:
			templateFile = fmt.Sprintf("errors/%d.html", code)
		case http.StatusMethodNotAllowed:
			// if content type is JSON, return JSON error response
			if c.Request().Header.Get(echo.HeaderContentType) == echo.MIMEApplicationJSON {
				// Return JSON error response
				_ = c.JSON(
					http.StatusMethodNotAllowed,
					map[string]interface{}{
						"Error": helpers.ErrorMessage(err),
					})

				return
			}
			// Otherwise, return HTML error response
			templateFile = fmt.Sprintf("errors/%d.html", code)
		default:
			templateFile = "errors/error.html"
		}
	}

	_ = c.Render(
		code,
		templateFile,
		map[string]interface{}{
			"Title": "Something went wrong",
			"Error": err,
		})
}
