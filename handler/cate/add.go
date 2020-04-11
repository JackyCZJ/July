package cate

import (
	"fmt"

	"github.com/jackyczj/July/handler"
	"github.com/jackyczj/July/log"
	"github.com/jackyczj/July/store"
	"github.com/labstack/echo/v4"
)

func Add(ctx echo.Context) error {
	a := struct {
		New    string `json:"new"`
		Parent string `json:"parent"`
	}{}
	err := ctx.Bind(&a)
	if err != nil {
		log.Logworker.Error(err)
		return handler.ErrorResp(ctx, err, 500)
	}
	if a.Parent == "root" {
		c := store.Cate{}
		c.Name = a.New
		c.Parent = a.Parent
		err := c.InsertCate()
		if err != nil {
			return handler.ErrorResp(ctx, err, 500)
		}
	} else {
		c := store.Cate{}
		c.Parent = a.Parent
		c.Name = a.New
		err = c.InsertCate()
		if err != nil {
			fmt.Println(err)
			return handler.ErrorResp(ctx, err, 500)
		}
	}
	return handler.Response(ctx, handler.ResponseStruct{
		Code:    0,
		Message: "",
		Data:    nil,
	})
}
