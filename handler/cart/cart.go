package cart

import (
	"github.com/jackyczj/July/handler"
	"github.com/jackyczj/July/store"
	"github.com/labstack/echo/v4"
)

type Request struct {
	Product store.Product `json:"product"`
	Count   int           `json:"count"`
}

//🛒 Add , add something into cart
func Add(ctx echo.Context) error {
	r := new(Request)
	err := ctx.Bind(&r)
	if err != nil {
		return handler.ErrorResp(ctx, err, 500)
	}
	err = store.CartAdd(ctx.Get("user_id").(int32), r.Product, r.Count)
	if err != nil {
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
	return nil

}

func Clear(ctx echo.Context) error {
	return nil
}

//🛒 List
func List(ctx echo.Context) error {
	return nil
}

func Get(ctx echo.Context) error {
	return nil
}
