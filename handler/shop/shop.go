package shop

import (
	"github.com/jackyczj/July/handler"
	"github.com/jackyczj/July/store"
	"github.com/labstack/echo/v4"
)

func Add(ctx echo.Context) error {
	s := store.Shop{}
	err := ctx.Bind(&s)
	if err != nil {
		return err
	}
	return s.Create()

}

func Delete(ctx echo.Context) error {
	return nil
}

func List(ctx echo.Context) error {
	shopList := struct {
		Page    int          `query:"page"`
		PerPage int          `query:"PerPage"`
		Total   int          `json:"total"`
		Data    []store.Shop `json:"data"`
	}{}
	_ = ctx.Bind(shopList)
	data, total, err := store.ShopList(shopList.Page, shopList.PerPage)
	if err != nil {
		return handler.ErrorResp(ctx, err, 404)
	}
	shopList.Data = data
	shopList.Total = total

	return handler.Response(ctx, handler.ResponseStruct{
		Code:    0,
		Message: "",
		Data:    shopList,
	})
}

func Search(ctx echo.Context) error {
	search := struct {
		Keyword string       `query:"keyword"`
		Page    int          `query:"page"`
		PerPage int          `query:"PerPage"`
		Total   int          `json:"total"`
		Data    []store.Shop `json:"data"`
	}{}
	_ = ctx.Bind(search)
	key := search.Keyword
	page := search.Page
	perPage := search.PerPage
	data, total, err := store.SearchShop(key, page, perPage)
	if err != nil {
		return handler.ErrorResp(ctx, err, 404)
	}
	search.Total = total
	search.Data = data
	resp := handler.ResponseStruct{
		Code:    0,
		Message: "",
		Data:    search,
	}
	return handler.Response(ctx, resp)
}
