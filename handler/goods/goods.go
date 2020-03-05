package goods

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/jackyczj/July/handler"
	"github.com/jackyczj/July/store"
	"github.com/labstack/echo/v4"
	"github.com/rs/xid"
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
		Keyword string `query:"keyword"`
		Page    int    `query:"page"`
		PerPage int    `query:"PerPage"`
	}{}
	_ = ctx.Bind(search)
	key := search.Keyword
	page := search.Page
	perPage := search.PerPage
	data, err := store.Search(key, page, perPage)
	if err != nil {
		return handler.ErrorResp(ctx, err, 404)
	}

	return ctx.JSON(200, data)
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
	id, err := strconv.Atoi(pid)
	if err != nil {
		return handler.ErrorResp(ctx, err, 500)
	}
	p := store.Product{ProductId: int32(id)}
	err = p.Get()
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
	good.Owner = u.Username
	good.ProductId = int32(xid.New().Pid())

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
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 10)
	if err != nil {
		return handler.ErrorResp(ctx, err, 500)
	}
	p.ProductId = int32(id)
	err = p.Get()
	if err != nil {
		return handler.ErrorResp(ctx, err, 500)
	}
	if p.Owner != u.Username {
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
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 10)
	if err != nil {
		return handler.ErrorResp(ctx, err, 500)
	}
	p.ProductId = int32(id)
	err = p.Get()
	if err != nil {
		return handler.ErrorResp(ctx, err, 404)
	}
	if p.Owner != u.Username {
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

func List(ctx echo.Context) error {
	return nil
}
