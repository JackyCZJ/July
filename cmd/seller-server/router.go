package main

import (
	Auth "github.com/jackyczj/July/auth"
	"github.com/jackyczj/July/handler/goods"
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
	Manager := e.Group("/manage")
	//店铺数据
	//Manager.GET("/statistic")

	//商户登陆
	Account := Manager.Group("/seller")
	{
		Account.POST("/login", user.Login)
		Account.POST("/register", user.Register)
	}

	//Goods
	//		 upload
	//		 List
	//		 manage
	//产品管理
	Goods := Manager.Group("/product")
	{
		//新增商品
		Goods.POST("/add", goods.Add)
		Goods.POST("/delete/:id", goods.Delete)
		Goods.POST("/update", goods.Edit)
		Goods.GET("/list", goods.List)
		Goods.GET("/search/:keyword", goods.Search)

	}
	//类别管理
	//Category := Manager.Group("/category")
	//{
	//	Category.GET("/:id", cate.Get)
	//	Category.POST("/add",cate.Add)
	//}

	//订单管理
	// Order manage
	//Order := Manager.Group("/order")
	{
		//Order.GET("/detail/orderNo?id", order.detail)
		//Order.GET("/search")
		//Order.GET("/list?pagenum=")
		//Order.POST("/send_goods")
	}
}
