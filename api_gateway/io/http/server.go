package http

import (
	"encoding/json"
	"github.com/golobby/container/v3"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"log"
	"my-stocks/api-gateway/app"
	"my-stocks/api-gateway/config"
)

func addRoutes(ctr container.Container, e *echo.Echo) {
	e.Use(middleware.Logger())
	g := e.Group("/auth")
	var authService app.AuthService
	_ = ctr.Resolve(&authService)
	var userService app.UserService
	_ = ctr.Resolve(&userService)
	authController := NewAuthController(&authService, &userService)
	g.POST("/register", authController.Register)
	g.POST("/login", authController.Login)

	userController := NewUserController(&userService)
	g = e.Group("/profile", AuthMiddleware(&authService, &userService))
	g.GET("/", userController.Profile)
	g.DELETE("/logout", authController.Logout)

	var coinService app.CoinService
	_ = ctr.Resolve(&coinService)
	coinController := NewCoinController(&coinService)
	g = e.Group("/coins", AuthMiddleware(&authService, &userService))
	g.GET("/", coinController.Paginate)

	data, err := json.MarshalIndent(e.Routes(), "", "  ")
	if err != nil {
		log.Println(err)
	}
	log.Println(string(data))
}

func Start(cfg *config.App, ctr container.Container) {
	e := echo.New()
	addRoutes(ctr, e)
	e.Logger.Fatal(e.Start(cfg.HttpServerAddress))
}
