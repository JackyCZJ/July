package cate

import (
	"github.com/jackyczj/July/handler"
	"github.com/jackyczj/July/store"
	"github.com/labstack/echo/v4"
)

var delreq struct {
	New    string `json:"new"`
	Parent string `json:"parent"`
}

func Delete(ctx echo.Context) error {
	err := ctx.Bind(&delreq)
	if err != nil {
		return handler.ErrorResp(ctx, err, 500)
	}
	if addreq.Parent == "root" {
		c := store.Cate{}
		c.Name = delreq.New
		err := c.DeleteCate()
		if err != nil {
			return handler.ErrorResp(ctx, err, 500)
		}
	} else {
		c := store.Cate{}
		c.Name = delreq.Parent
		err = c.DeleteFromCate(delreq.New)
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
