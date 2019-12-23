package handler

import "github.com/labstack/echo"

type ResponseStruct struct {
	Code    int         `json:"code"`
	Message string      `json:"msg"`
	Data    interface{} `json:"data"`
}

func Response(ctx echo.Context, responseStruct ResponseStruct) error {
	return ctx.JSON(200, responseStruct)
}

func ErrorResponse(ctx echo.Context, err error) error {
	rs := ResponseStruct{
		Code:    0,
		Message: err.Error(),
		Data:    nil,
	}
	return ctx.JSON(200, rs)
}
