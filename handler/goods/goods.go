package goods

import (
	"fmt"
	"strconv"

	"github.com/jackyczj/July/log"

	"github.com/jackyczj/July/handler"
	"github.com/jackyczj/July/store"
	"github.com/labstack/echo/v4"
)

/*
	商品名
	商品类别
	商品库存
	商品价格
	商品折扣
	商品图片
	商品介绍
*/

func Search(ctx echo.Context) error {
	search := struct {
		Keyword string          `json:"keyword" query:"keyword"`
		Page    int             `json:"page" query:"page"`
		PerPage int             `json:"per_page" query:"PerPage"`
		Total   int             `json:"total"`
		Data    []store.Product `json:"data"`
	}{}
	search.Keyword = ctx.Param("keyword")
	search.Page, _ = strconv.Atoi(ctx.Param("page"))
	search.PerPage, _ = strconv.Atoi(ctx.Param("pageSize"))
	data, total, err := store.Search(search.Keyword, search.Page, search.PerPage)
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

func Index(ctx echo.Context) error {
	//var goodsList []store.Product
	data, err := store.GetRandom()
	if err != nil {
		fmt.Println(err.Error())
		return ctx.JSON(200, nil)
	}
	return handler.Response(ctx, handler.ResponseStruct{
		Code:    0,
		Message: "",
		Data:    data,
	})
}

func Get(ctx echo.Context) error {
	pid := ctx.Param("id")

	p := store.Product{ProductId: pid}
	err := p.Get()
	if err != nil {
		return handler.ErrorResp(ctx, err, 404)
	}

	return handler.Response(ctx, handler.ResponseStruct{
		Code:    0,
		Message: "",
		Data:    p,
	})
}

var u store.UserInformation

func Add(ctx echo.Context) error {
	id := ctx.Get("user_id").(int32)
	u.Id = id
	u, err := u.GetUser()
	if err != nil {
		return handler.ErrorResp(ctx, err, 500)
	}
	switch u.Role {
	case 2, 3:
	default:
		return handler.Response(
			ctx, handler.ResponseStruct{
				Code:    0,
				Message: "you are not seller",
				Data:    nil,
			})
	}
	good := new(store.Product)
	err = ctx.Bind(&good)
	if err != nil {
		log.Logworker.Error(err)
		return handler.ErrorResp(ctx, err, 500)
	}
	var s store.Shop
	s.Owner = u.Id
	err = s.GetByOwner()
	if err != nil {
		return handler.ErrorResp(ctx, err, 500)
	}
	good.Owner = s.Id
	if err = good.Add(); err != nil {
		return handler.ErrorResp(ctx, err, 500)
	}
	return handler.Response(ctx, handler.ResponseStruct{
		Code:    0,
		Message: "",
		Data:    nil,
	})
}

/*
	ProductId	int32
*/
func Delete(ctx echo.Context) error {
	var err error
	uid := ctx.Get("user_id").(int32)
	u.Id = uid
	u, err := u.GetUser()
	if err != nil {
		return handler.ErrorResp(ctx, err, 403)
	}
	p := new(store.Product)
	id := ctx.Param("id")
	p.ProductId = id
	err = p.Get()
	if err != nil {
		return handler.ErrorResp(ctx, err, 500)
	}
	if !CheckOwner(p.Owner, *u) {
		return handler.Response(ctx, handler.ResponseStruct{
			Code:    0,
			Message: "Error delete, not your good ",
			Data:    nil,
		})
	}
	if err = p.Delete(); err != nil {
		return handler.ErrorResp(ctx, err, 500)
	}
	return nil
}

/*
	Name        string `json:"name"`
	Image	    file
	Information Type   `json:"info"`
	Price       int    `json:"price"`
	Off         int    `json:"off"`
*/
func Edit(ctx echo.Context) error {
	var err error
	uid := ctx.Get("user_id").(int32)
	u.Id = uid
	u, err := u.GetUser()
	if err != nil {
		return handler.ErrorResp(ctx, err, 403)
	}
	p := new(store.Product)
	id := ctx.Param("id")
	p.ProductId = id
	err = p.Get()
	if err != nil {
		return handler.ErrorResp(ctx, err, 404)
	}

	var s store.Shop
	s.Id = p.Owner
	err = s.Get()
	if err != nil {
		return handler.ErrorResp(ctx, err, 404)
	}
	if s.Owner != u.Id {
		return handler.Response(ctx, handler.ResponseStruct{
			Code:    0,
			Message: "Error Update, not your good ",
			Data:    nil,
		})
	}
	if err = p.Update(); err != nil {
		return handler.ErrorResp(ctx, err, 500)

	}
	return handler.Response(ctx, handler.ResponseStruct{
		Code:    1,
		Message: "ok ",
		Data:    nil,
	})
}

func Suggestion(ctx echo.Context) error {
	keyword := ctx.Param("keyword")
	answer := store.Suggestion(keyword)
	return handler.Response(ctx, handler.ResponseStruct{
		Code:    0,
		Message: "",
		Data:    answer,
	})
}

func Comment(ctx echo.Context) error {
	commentId := ctx.Param("id")
	c := store.Comment{}
	err := ctx.Bind(&c)
	if err != nil {
		return ctx.JSON(500, err)
	}
	err = store.AddComment(commentId, c)
	if err != nil {
		return ctx.JSON(500, err)
	}
	return ctx.JSON(200, nil)
}

func DelComment(ctx echo.Context) error {
	commentId := ctx.Param("id")
	user := store.UserInformation{}
	user.Id = ctx.Get("user_id").(int32)
	u, err := user.GetUser()
	if err != nil {
		return ctx.JSON(500, err)
	}
	err = store.DeleteComment(commentId, u.Username)
	if err != nil {
		return ctx.JSON(500, err)
	}
	return ctx.JSON(200, nil)
}

func CheckOwner(owner string, u store.UserInformation) bool {
	var s store.Shop
	s.Id = owner
	err := s.Get()
	if err != nil {
		return false
	}
	if s.Owner != u.Id {
		return false
	}
	return true
}

func ProductListByShopId(ctx echo.Context) error {
	shop := ctx.Param("id")
	page := ctx.QueryParam("page")
	p, _ := strconv.Atoi(page)
	if p == 0 {
		p = 1
	}
	data, total := store.GetListByShop(shop, false, p)
	resp := struct {
		Total int64           `json:"total"`
		Data  []store.Product `json:"data"`
	}{
		Total: total,
		Data:  data,
	}
	return handler.Response(ctx, handler.ResponseStruct{
		Code:    0,
		Message: "",
		Data:    resp,
	})
}

func ProductListByShop(ctx echo.Context) error {
	page := ctx.Param("id")
	id := ctx.Get("user_id").(int32)
	p, _ := strconv.Atoi(page)
	var s store.Shop
	s.Owner = id
	err := s.GetByOwner()
	if err != nil {
		return err
	}
	u := store.UserInformation{}
	u.Id = id

	data, total := store.GetListByShop(s.Id, CheckOwner(s.Id, u), p)
	resp := struct {
		Total int64           `json:"total"`
		Data  []store.Product `json:"data"`
	}{
		Total: total,
		Data:  data,
	}
	return handler.Response(ctx, handler.ResponseStruct{
		Code:    0,
		Message: "",
		Data:    resp,
	})
}

func SearchInShop(ctx echo.Context) error {
	keyword := ctx.QueryParam("keyword")
	page := ctx.QueryParam("page")
	Type := ctx.QueryParam("type")
	fmt.Println(keyword, page, Type)
	p, _ := strconv.Atoi(page)
	if p < 1 {
		p = 1
	}
	id := ctx.Get("user_id").(int32)
	var s store.Shop
	s.Owner = id
	err := s.GetByOwner()
	if err != nil {
		return err
	}
	switch Type {
	default:
		data, total := store.SearchInShop(s.Id, keyword, CheckOwner(s.Id, u), p)
		resp := struct {
			Total int64           `json:"total"`
			Data  []store.Product `json:"data"`
		}{
			Total: total,
			Data:  data,
		}
		return handler.Response(ctx, handler.ResponseStruct{
			Code:    0,
			Message: "",
			Data:    resp,
		})
	case "productId":
		data, total := store.SearchInShopById(s.Id, keyword, CheckOwner(s.Id, u), p)
		resp := struct {
			Total int64           `json:"total"`
			Data  []store.Product `json:"data"`
		}{
			Total: total,
			Data:  data,
		}
		return handler.Response(ctx, handler.ResponseStruct{
			Code:    0,
			Message: "",
			Data:    resp,
		})
	}

}
