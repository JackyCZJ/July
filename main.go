package main

import (
	"github.com/jackyczj/NoGhost/store"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func main() {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())
	//e.Use(middleware.KeyAuthWithConfig(middleware.KeyAuthConfig{
	//	Skipper:   skipper,   // 跳过验证条件 在 auth.go 定义
	//	Validator: validator, // 处理验证结果 在 auth.go 定义
	//}))
	store.Client.Init()
	defer store.Client.Close()
	e.Use(middleware.CORS())
	e.Logger.Fatal("Service start at port:", e.Start(":2333"))
}
