package cate

import (
	"github.com/jackyczj/July/handler"
	"github.com/jackyczj/July/store"
	"github.com/labstack/echo/v4"
)

func List(ctx echo.Context) error {
	id := ctx.Param("id")
	if id == "0" {
		return handler.Response(ctx, handler.ResponseStruct{
			Code:    0,
			Message: "",
			Data:    store.GetCateTree(),
		})
	}

	return handler.Response(ctx, handler.ResponseStruct{
		Code:    0,
		Message: "",
		Data:    store.GetCateByParent(id),
	})
}
