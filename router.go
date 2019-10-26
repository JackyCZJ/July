package main

import (
	Auth "github.com/jackyczj/July/auth"
	"github.com/jackyczj/July/handler/captcha"
	"github.com/jackyczj/July/handler/user"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func Load(e *echo.Echo) {
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())
	//e.Use(middleware.CSRF())
	e.Use(middleware.KeyAuthWithConfig(middleware.KeyAuthConfig{
		Skipper:   Auth.Skipper,
		Validator: Auth.Validator,
	}))

	// init config
	e.POST("/login", user.Login)
	e.POST("/register", user.Register)
	e.GET("/captcha", captcha.Generate)
	e.POST("/captcha", captcha.Verify)

}
