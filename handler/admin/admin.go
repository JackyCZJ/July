package admin

import (
	"github.com/jackyczj/July/handler"
	"github.com/jackyczj/July/store"
	"github.com/labstack/echo/v4"
)

func ShopList(ctx echo.Context) error {
	shopList := struct {
		Page    int          `query:"page" json:"page"`
		PerPage int          `query:"per_page" json:"per_page"`
		Total   int          `json:"total"`
		Data    []store.Shop `json:"data"`
	}{}
	err := ctx.Bind(&shopList)
	if err != nil {
		return err
	}
	s, count, err := store.ShopList(shopList.Page, shopList.PerPage)
	if err != nil {
		return err
	}
	shopList.Total = count
	shopList.Data = s
	return handler.Response(ctx, handler.ResponseStruct{
		Code:    0,
		Message: "",
		Data:    shopList,
	})
}

func UserList() {
	userList := struct {
		Page    int                     `query:"page" json:"page"`
		PerPage int                     `query:"per_page" json:"per_page"`
		Total   int                     `json:"total"`
		Data    []store.UserInformation `json:"data"`
	}{}

}

func UserBan() {

}

func ShopBan() {

}

func Comment() {

}
