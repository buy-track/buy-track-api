package http

import (
	"github.com/labstack/echo/v4"
	"my-stocks/api-gateway/app"
	"net/http"
)

func AuthMiddleware(authService *app.AuthService, userService *app.UserService) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			token := c.Request().Header.Get("Authorization")
			if token == "" {
				return c.JSON(http.StatusUnauthorized, "unauthorized")
			}
			verified, err := authService.VerifyToken(token)
			if err != nil {
				return c.JSON(http.StatusUnauthorized, "unauthorized")
			}
			user, err := userService.GetById(verified)
			if err != nil {
				return c.JSON(http.StatusUnauthorized, "unauthorized")
			}
			c.Set("auth_token", token)
			c.Set("auth_user", user)
			return next(c)
		}
	}
}
