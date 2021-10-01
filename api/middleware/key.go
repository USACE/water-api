package middleware

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func KeyAuth(validKey string) echo.MiddlewareFunc {
	return middleware.KeyAuthWithConfig(
		middleware.KeyAuthConfig{
			KeyLookup: "query:key",
			Validator: func(key string, c echo.Context) (bool, error) {
				return key == validKey, nil
			},
			// Custom error handler; Do not expose any information
			// other than that the request was unauthorized
			ErrorHandler: func(e error, c echo.Context) error {
				return echo.NewHTTPError(http.StatusUnauthorized)
			},
		},
	)
}
