package main

import (
	Auth "github.com/jackyczj/July/auth"
	"github.com/jackyczj/July/handler/user"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func Load(e *echo.Echo) {
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.KeyAuthWithConfig(middleware.KeyAuthConfig{
		Skipper:    Auth.Skipper,
		Validator:  Auth.Validator,
		KeyLookup:  "header:" + echo.HeaderAuthorization,
		AuthScheme: "Bearer",
	}))

	// init config
	Account := e.Group("/seller")
	{
		Account.POST("/login", user.Login)
		Account.POST("/register", user.Register)
	}

	//Goods
	//		 upload
	//		 List
	//		 manage

	// Order manage
}
