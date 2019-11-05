package goods

import "github.com/labstack/echo/v4"

type Product struct {
	Name  string `json:"name"`
	Type  Type   `json:"type"`
	Price int    `json:"price"`
	Off   int    `json:"off"`
	Owner string `json:"owner"`
}

type Type struct {
	Category string `json:"category"` //产品分类
	Brand    string `json:"brand"`    //产品品牌
}

func Search(ctx echo.Context) error {
	return nil
}

func Index(ctx echo.Context) error {
	return nil
}

func Get(ctx echo.Context) error {
	return nil
}

func Add(ctx echo.Context) error {
	return nil
}

func Delete(ctx echo.Context) error {
	return nil

}
