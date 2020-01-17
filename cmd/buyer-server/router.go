package main

import (
	Auth "github.com/jackyczj/July/auth"
	"github.com/jackyczj/July/handler/captcha"
	cartHandler "github.com/jackyczj/July/handler/cart"
	shopHandler "github.com/jackyczj/July/handler/shop"
	"github.com/jackyczj/July/handler/user"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func Load(e *echo.Echo) {
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	//e.Use(middleware.CORS())
	//c := middleware.CSRFConfig{
	//	Skipper:      middleware.DefaultSkipper,
	//	TokenLength:  32,
	//	TokenLookup:  "header:" + echo.HeaderXCSRFToken,
	//	ContextKey:   "csrf",
	//	CookieName:   "_csrf",
	//	CookieMaxAge: 86400,
	//}
	//e.Use(middleware.CSRFWithConfig(c))
	e.Use(middleware.KeyAuthWithConfig(middleware.KeyAuthConfig{
		Skipper:    Auth.Skipper,
		Validator:  Auth.Validator,
		KeyLookup:  "header:" + echo.HeaderAuthorization,
		AuthScheme: "Bearer",
	}))

	// init config
	//Account := e.Group("/user/", middleware.CSRF())
	Account := e.Group("/user")
	{
		Account.POST("/login", user.Login)
		Account.POST("/register", user.Register)

	}

	Cap := e.Group("/captcha")
	{
		Cap.POST("/get", captcha.Generate)
		Cap.POST("/verify", captcha.Verify)
	}

	//Todo: ⬇️
	api := e.Group("/api/v1")

	//Goods := api.Group("/Goods")
	//{
	//Goods.GET("/Goods/:str",goodsHandler.Search) //search
	//Goods.GET("/Goods/index",goodsHandler.Index) //index goods list
	//}

	cart := api.Group("/cart")
	{
		cart.GET("/", cartHandler.List)
		cart.POST("/", cartHandler.Add)
		cart.DELETE("/:id", cartHandler.Delete)
		cart.DELETE("/", cartHandler.Delete)
	}

	shopping := api.Group("/shop")
	{
		shopping.GET("", shopHandler.List)
		shopping.POST("/:id", shopHandler.Add)
		shopping.DELETE("/:id", shopHandler.Delete)
	}

}