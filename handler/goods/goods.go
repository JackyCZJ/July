package goods

import (
	"fmt"
	"io"
	"net/http"
	"os"
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
	return nil
}

func Index(ctx echo.Context) error {
	return nil
}

func Get(ctx echo.Context) error {
	return nil
}

/*
	Name        string `json:"name"`
	Image	    file
	Information Type   `json:"info"`
	Price       int    `json:"price"`
	Off         int    `json:"off"`
*/

func GetFileContentType(out *os.File) (string, error) {

	// Only the first 512 bytes are used to sniff the content type.
	buffer := make([]byte, 512)

	_, err := out.Read(buffer)
	if err != nil {
		return "", err
	}

	// Use the net/http package's handy DectectContentType function. Always returns a valid
	// content-type by returning "application/octet-stream" if no others seemed to match.
	contentType := http.DetectContentType(buffer)

	return contentType, nil
}

var u store.UserInformation

func Add(ctx echo.Context) error {
	id := ctx.Get("user_id").(uint16)
	u.Id = id
	u, err := u.GetUser()
	if err != nil {
		return handler.Response(ctx, handler.ResponseStruct{
			Code:    0,
			Message: err.Error(),
			Data:    nil,
		})
	}
	if u.Role != 2 {
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
		return handler.ErrorResp(ctx, err)
	}
	good.Owner = u.Username
	good.ProductId = xid.New().Pid()

	file, err := ctx.FormFile("image")
	if err != nil {
		return err
	}
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()
	// Destination
	dst, err := os.Create(file.Filename)
	if err != nil {
		return err
	}
	defer dst.Close()
	buffer := make([]byte, 512)

	_, err = dst.Read(buffer)
	if err != nil {
		return err
	}
	contentType := http.DetectContentType(buffer)
	switch contentType {
	case "image/png", "image/gif", "image/jpeg", "image/jpg":
	default:
		return handler.ErrorResp(ctx, fmt.Errorf("InVail file type. "))
	}
	// Copy
	if _, err = io.Copy(dst, src); err != nil {
		return err
	}
	good.ImageUri = file.Filename
	if err = good.Add(); err != nil {
		return handler.ErrorResp(ctx, err)
	}
	return nil
}

/*
	ProductId	uint16
*/
func Delete(ctx echo.Context) error {
	var err error
	uid := ctx.Get("user_id").(uint16)
	u.Id = uid
	u, err := u.GetUser()
	if err != nil {
		return handler.ErrorResp(ctx, err)
	}
	p := new(store.Product)
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 10)
	if err != nil {
		return handler.ErrorResp(ctx, err)
	}
	p.ProductId = uint16(id)
	err = p.Get()
	if err != nil {
		return handler.ErrorResp(ctx, err)
	}
	if p.Owner != u.Username {
		return handler.Response(ctx, handler.ResponseStruct{
			Code:    0,
			Message: "Error delete, not your good ",
			Data:    nil,
		})
	}
	if err = p.Delete(); err != nil {
		return handler.ErrorResp(ctx, err)
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
	uid := ctx.Get("user_id").(uint16)
	u.Id = uid
	u, err := u.GetUser()
	if err != nil {
		return handler.ErrorResp(ctx, err)
	}
	p := new(store.Product)
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 10)
	if err != nil {
		return handler.ErrorResp(ctx, err)
	}
	p.ProductId = uint16(id)
	err = p.Get()
	if err != nil {
		return handler.ErrorResp(ctx, err)
	}
	if p.Owner != u.Username {
		return handler.Response(ctx, handler.ResponseStruct{
			Code:    0,
			Message: "Error Update, not your good ",
			Data:    nil,
		})
	}
	if err = p.Update(); err != nil {
		return handler.ErrorResp(ctx, err)

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
