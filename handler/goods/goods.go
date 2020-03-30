package goods

import (
	"fmt"
	"net/http"

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
	err := ctx.Bind(&search)
	if err != nil {
		return handler.ErrorResp(ctx, err, 404)
	}
	key := search.Keyword
	page := search.Page
	perPage := search.PerPage
	data, total, err := store.Search(key, page, perPage)
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
		return handler.Response(ctx, handler.ResponseStruct{
			Code:    0,
			Message: err.Error(),
			Data:    nil,
		})
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
		return handler.ErrorResp(ctx, err, 404)
	}
	var s store.Shop
	s.Owner = u.Id
	err = s.GetByOwner()
	if err != nil {
		return handler.ErrorResp(ctx, err, 404)
	}
	good.Owner = s.Id
	form, err := ctx.MultipartForm()
	if err != nil {
		return err
	}
	files := form.File["images"]
	var pathArray []string
	if len(files) > 4 {
		return handler.ErrorResp(ctx, fmt.Errorf("too many images , less than or qual 4 plz. "), 403)
	}
	for _, file := range files {
		src, err := file.Open()
		if file.Size > 150 {
			return handler.ErrorResp(ctx, fmt.Errorf("image too big "), 403)
		}
		if err != nil {
			return err
		}
		buffer := make([]byte, 512)
		_, err = src.Read(buffer)
		if err != nil {
			return err
		}
		contentType := http.DetectContentType(buffer)
		switch contentType {
		case "image/png", "image/gif", "image/jpeg", "image/jpg":
		default:
			return handler.ErrorResp(ctx, fmt.Errorf("InVail file type. "), 403)
		}
		// Copy
		path, err := store.Upload(src, file.Filename)
		if err != nil {
			return err
		}
		_ = src.Close()
		pathArray = append(pathArray, path)
	}
	good.ImageUri = pathArray
	if err = good.Add(); err != nil {
		return handler.ErrorResp(ctx, err, 500)
	}
	return nil
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
