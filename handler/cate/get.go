package cate

import (
	"github.com/jackyczj/July/handler"
	"github.com/jackyczj/July/store"
	"github.com/labstack/echo/v4"
)

func Get(ctx echo.Context) error {
	key := ctx.Param("id")
	if key != "" {
		c := store.Cate{Name: key}
		err := c.Get()
		if err != nil {
			return handler.ErrorResp(ctx, err, 404)
		}
		return handler.Response(ctx, handler.ResponseStruct{
			Code:    0,
			Message: "",
			Data:    c,
		})
	}
	c := store.GetCateTree()
	return handler.Response(ctx, handler.ResponseStruct{
		Code:    0,
		Message: "",
		Data:    c,
	})
}
