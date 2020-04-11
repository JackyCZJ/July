package cate

import (
	"github.com/jackyczj/July/handler"
	"github.com/jackyczj/July/store"
	"github.com/labstack/echo/v4"
)

type delreq struct {
	Name   string `json:"name"`
	Parent string `json:"parent"`
}

func Delete(ctx echo.Context) error {
	var d delreq
	err := ctx.Bind(&d)
	if err != nil {
		return handler.ErrorResp(ctx, err, 500)
	}
	if d.Parent == "root" {
		c := store.Cate{}
		c.Name = d.Name
		err := c.DeleteCate()
		if err != nil {
			return handler.ErrorResp(ctx, err, 500)
		}
	} else {
		c := store.Cate{}
		c.Name = d.Parent
		err = c.DeleteFromCate(d.Name)
		if err != nil {
			return handler.ErrorResp(ctx, err, 500)
		}
	}
	return handler.Response(ctx, handler.ResponseStruct{
		Code:    0,
		Message: "",
		Data:    nil,
	})
}
