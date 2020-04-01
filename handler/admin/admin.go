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

func Status(ctx echo.Context) error {
	s, err := store.StatusGet()
	if err != nil {
		return handler.ErrorResp(ctx, err, 404)
	}
	return handler.Response(ctx, handler.ResponseStruct{
		Code:    0,
		Message: "",
		Data:    s,
	})
}

func UserList(ctx echo.Context) error {
	var resq = struct {
		Page    int                     `query:"page" json:"page"`
		PerPage int                     `query:"per_Page" json:"per_Page"`
		Total   int                     `json:"total"`
		Data    []store.UserInformation `json:"data"`
	}{}
	err := ctx.Bind(&resq)
	if err != nil {
		resq.Page = 1
		resq.PerPage = 10
	}
	userList, total, err := store.UserList(resq.Page, resq.PerPage)
	if err != nil {
		return handler.ErrorResp(ctx, err, 500)
	}
	resq.Data = userList
	resq.Total = total
	return handler.Response(ctx, handler.ResponseStruct{
		Code:    0,
		Message: "",
		Data:    resq,
	})

}

func UserBan(ctx echo.Context) error {
	return nil
}

func ShopBan(ctx echo.Context) error {
	return nil
}
