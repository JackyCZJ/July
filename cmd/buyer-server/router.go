package main

import (
	"net/http"

	"github.com/jackyczj/July/log"

	"github.com/casbin/casbin/v2"
	"github.com/jackyczj/July/handler/goods"

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
	e.Use(middleware.CORS())
	//c := middleware.CSRFConfig{
	//	Skipper:      middleware.DefaultSkipper,
	//	TokenLength:  32,
	//	TokenLookup:  "header:" + echo.HeaderXCSRFToken,
	//	ContextKey:   "csrf",
	//	CookieName:   "_csrf",
	//	CookieMaxAge: 86400,
	//}
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"https://localhost:2333", "http://localhost:2333", "http://localhost:3000"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
		AllowMethods: []string{http.MethodGet, http.MethodHead, http.MethodPut, http.MethodPatch, http.MethodPost, http.MethodDelete},
	}))
	e.Use(middleware.KeyAuthWithConfig(middleware.KeyAuthConfig{
		Skipper:    Auth.Skipper,
		Validator:  Auth.Validator,
		KeyLookup:  "header:" + echo.HeaderAuthorization,
		AuthScheme: "Bearer",
	}))

	enforcer, err := casbin.NewEnforcer("conf/casbin_auth_model.conf", "conf/casbin_auth_policy.csv")
	if err != nil {
		log.Logworker.Fatal(err.Error())
	}

	e.Use(Auth.MiddlewareWithConfig(Auth.Config{
		Skipper:  Auth.Skipper,
		Enforcer: enforcer,
	}))

	// init config
	//Account := e.Group("/user/", middleware.CSRF())
	Account := e.Group("/user")
	{
		Account.POST("/login", user.Login)
		Account.POST("/logout", user.Logout)
		Account.POST("/register", user.Register)

	}

	Cap := e.Group("/captcha")
	{
		Cap.POST("/get", captcha.Generate)
		Cap.POST("/verify", captcha.Verify)
	}

	//Todo: ⬇️
	api := e.Group("/api/v1")

	Goods := api.Group("/Goods")
	{
		//Goods.GET("/Goods/:str",goodsHandler.Search) //search
		Goods.GET("/index", goods.Index) //index goods list
		Goods.GET("/:id", goods.Get)
	}

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
