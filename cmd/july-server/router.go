package main

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo-contrib/prometheus"

	"github.com/jackyczj/July/handler/cate"

	"github.com/jackyczj/July/handler/order"

	"github.com/jackyczj/July/handler/file"

	"github.com/jackyczj/July/log"

	"github.com/casbin/casbin/v2"
	"github.com/jackyczj/July/handler/goods"

	Auth "github.com/jackyczj/July/auth"
	adminHandler "github.com/jackyczj/July/handler/admin"
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
	p := prometheus.NewPrometheus("echo", nil)
	p.Use(e)
	//c := middleware.CSRFConfig{
	//	Skipper:      middleware.DefaultSkipper,
	//	TokenLength:  32,
	//	TokenLookup:  "header:" + echo.HeaderXCSRFToken,
	//	ContextKey:   "csrf",
	//	CookieName:   "_csrf",
	//	CookieMaxAge: 86400,
	//}
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"https://localhost:2333", "http://localhost:2333", "http://localhost:3000", "http://localhost:3001"},
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
		Account.GET("/checkUsername/:key", user.CheckUsername)
		Account.GET("/checkEmail/:key", user.CheckEmail)
		Account.POST("/register", user.Register)
	}
	e.GET("/image/:filename", file.Image)
	e.POST("/upload", file.Upload)

	//Todo: ⬇️
	api := e.Group("/api/v1")
	User := api.Group("/Users")
	{
		User.GET("/information", user.Get)
		User.POST("/information", user.Set)
		User.GET("/address", user.Address)
	}
	Goods := api.Group("/Goods")
	{
		fmt.Println("➡️ inject goods api")

		Goods.GET("/index", goods.Index) //index goods list
		Goods.GET("/:id", goods.Get)
		Goods.GET("/search/:keyword/:page/:pageSize", goods.Search)
		//search hint
		Goods.GET("/suggestion/:keyword", goods.Suggestion)
		Goods.GET("/suggestion/", goods.Suggestion)
		//Comment api ,but without auth
		Goods.POST("/comment/:id", goods.Comment)
		Goods.DELETE("/comment/:id", goods.DelComment)
	}
	Cart := api.Group("/Cart")
	{
		fmt.Println("➡️ inject cart api")
		Cart.GET("/List", cartHandler.List) //获取所有购物车
		Cart.POST("/Add", cartHandler.Add)  //添加入购物车
		Cart.PUT("/Edit", cartHandler.Add)  //修改某样商品的数量
		Cart.DELETE("/:id", cartHandler.Delete)
		Cart.POST("/clear", cartHandler.Clear) //清空购物车
	}

	Shop := api.Group("/Shop")
	{
		//Shop
		fmt.Println("➡️ inject shop api")
		Shop.GET("/search/:keyword/:page/:pageSize", shopHandler.Search)
		Shop.GET("/order/list", order.List)
		Shop.GET("/list", shopHandler.List)
		Shop.POST("/:id", shopHandler.Add)
		Shop.GET("/status", shopHandler.Status)
		Shop.GET("/:id/List/*", goods.ProductListByShopId)

		Shop.POST("/product/add", goods.Add)
		Shop.DELETE("/product/:id", goods.Delete)
		Shop.PUT("/product/:id", goods.Edit)
		Shop.GET("/product/:id", goods.ProductListByShop)
		Shop.GET("/product/search/*", goods.SearchInShop)

		//取现
		//shopping.POST("/takeMoney",shopHandler.TakeMoney)
	}

	Cate := api.Group("/category")
	{
		Cate.GET("/:id", cate.Get)
		Cate.POST("/add", cate.Add)
		Cate.GET("/List/:id", cate.List)
		Cate.DELETE("/delete", cate.Delete)
	}

	Order := api.Group("/Order")
	{
		fmt.Println("➡️ inject order api")
		Order.POST("/Create", order.Create)
		Order.GET("/List", order.List)
		Order.GET("/Get/:id", order.Get)
		Order.DELETE("/Delete/:id", order.Delete)
		Order.POST("/Pay/:id", order.Pay)
		Order.PUT("/Edit/:id", order.Edit)
		Order.PUT("/Send/*", order.Transmit)
		Order.POST("/Confirm/:id", order.Confirm)
		Order.DELETE("/Cancel/:id", order.Cancel)
	}

	admin := e.Group("/admin")
	{
		fmt.Println("➡️ inject admin api")
		admin.POST("/auth/login", user.Login)
		admin.GET("/auth/self", user.Get)
		admin.GET("/status", adminHandler.Status)
		shop := admin.Group("/shop")
		{
			shop.GET("/List", adminHandler.ShopList)
			shop.DELETE("/:id", shopHandler.Delete)
			shop.PUT("/ban/:id", adminHandler.ShopBan)
		}
		adminUser := admin.Group("/users")
		{
			adminUser.GET("/list", adminHandler.UserList)
		}
	}

}
