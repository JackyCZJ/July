package main

import (
	"github.com/jackyczj/July/handler/user"
	"github.com/labstack/echo"
)

var Router echo.Router

func init() {
	Router.Add(echo.POST, "login", user.Login)

}
