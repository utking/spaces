package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// getPasswordGeneratorWrapper is a wrapper for the secret generator handler.
func getPasswordGeneratorWrapper() echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.Render(
			http.StatusOK,
			"tools/secret-generator.html",
			map[string]interface{}{
				"Title": "Password Generator",
			},
		)
	}
}
