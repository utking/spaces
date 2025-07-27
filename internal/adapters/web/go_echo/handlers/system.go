package handlers

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"gogs.utking.net/utking/spaces/internal/adapters/web/go_echo/helpers"
	"gogs.utking.net/utking/spaces/internal/ports"
)

func getSystemStatsWrapper(
	api ports.SystemStatsService,
	userAPI ports.UsersService,
) echo.HandlerFunc {
	return func(c echo.Context) error {
		const userID = ""
		var (
			code    = http.StatusOK
			errList []error
		)

		stats, err := api.GetStats(c.Request().Context(), userID)
		if err != nil {
			errList = append(errList, fmt.Errorf("failed to get system stats: %w", err))
		}

		ctx := c.Request().Context()
		diskUse, _ := userAPI.GetDiskUsage(ctx, userID)

		for _, e := range errList {
			err = errors.Join(err, e)
		}

		return c.Render(
			code,
			"system/stats.html",
			map[string]interface{}{
				"Title":   "System Stats",
				"Error":   helpers.ErrorMessage(err),
				"Stats":   stats,
				"DiskUse": diskUse,
			},
		)
	}
}
