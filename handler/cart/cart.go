package cart

import (
	"github.com/jackyczj/July/handler"
	"github.com/labstack/echo/v4"
)

var res handler.ResponseStruct

//🛒 Add , add something into cart
func Add(ctx echo.Context) error {

	return nil
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
