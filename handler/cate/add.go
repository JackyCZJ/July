package cate

import (
	"github.com/jackyczj/July/handler"
	"github.com/jackyczj/July/store"
	"github.com/labstack/echo/v4"
)

var addreq struct {
	New    string `json:"new"`
	Parent string `json:"parent"`
}

func Add(ctx echo.Context) error {
	err := ctx.Bind(&addreq)
	if err != nil {
		return handler.ErrorResp(ctx, err, 500)
	}
	if addreq.Parent == "root" {
		c := store.Cate{}
		c.Name = addreq.New
		err := c.InsertCate()
		if err != nil {
			return handler.ErrorResp(ctx, err, 500)
		}
	} else {
		c := store.Cate{}
		c.Name = addreq.Parent
		err = c.AddToCate(addreq.New)
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
