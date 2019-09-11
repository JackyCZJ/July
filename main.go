package main

import "github.com/labstack/echo"

func main() {
	e := echo.New()
	e.Logger.Fatal("Service start at port:",e.Start(":2333"))
}
