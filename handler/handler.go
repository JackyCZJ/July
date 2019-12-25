package handler

import "github.com/labstack/echo/v4"

type ResponseStruct struct {
	Code    int         `json:"code"`
	Message string      `json:"msg"`
	Data    interface{} `json:"data"`
}

func Response(ctx echo.Context, responseStruct ResponseStruct) error {
	return ctx.JSON(200, responseStruct)
}
