package cart

import (
	"github.com/jackyczj/July/handler"
	"github.com/jackyczj/July/log"
	"github.com/jackyczj/July/store"
	"github.com/labstack/echo/v4"
)

type Request struct {
	Product string `json:"product"`
	Count   int    `json:"count"`
}

//🛒 Add , add something into cart
func Add(ctx echo.Context) error {
	r := new(Request)
	err := ctx.Bind(&r)
	if err != nil {
		log.Logworker.Error(err)
		return handler.ErrorResp(ctx, err, 500)
	}
	err = store.CartAdd(ctx.Get("user_id").(int32), r.Product, r.Count)
	if err != nil {
		log.Logworker.Error(err)
		return handler.ErrorResp(ctx, err, 500)
	}
	return handler.Response(ctx, handler.ResponseStruct{
		Code:    0,
		Message: "",
		Data:    nil,
	})
}

//🛒 Delete , with id it will delete what the goods stand for , without id it will clear up 🛒
func Delete(ctx echo.Context) error {
	order := ctx.Param("id")
	id := ctx.Get("user_id")
	err := store.CartDel(id.(int32), order)
	if err != nil {
		return handler.ErrorResp(ctx, err, 500)
	}
	return handler.Response(ctx, handler.ResponseStruct{
		Code:    0,
		Message: "",
		Data:    nil,
	})

}

func Clear(ctx echo.Context) error {
	id := ctx.Get("user_id").(int32)
	err := store.CartClear(id)
	if err != nil {
		return handler.ErrorResp(ctx, err, 500)
	}
	return handler.Response(ctx, handler.ResponseStruct{
		Code:    0,
		Message: "",
		Data:    nil,
	})
}

//🛒 List

func List(ctx echo.Context) error {
	id := ctx.Get("user_id").(int32)
	cart, err := store.CartList(id)
	if err != nil {
		return handler.ErrorResp(ctx, err, 500)
	}
	return handler.Response(ctx, handler.ResponseStruct{
		Code:    0,
		Message: "",
		Data:    cart,
	})
}
