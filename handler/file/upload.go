package file

import (
	"github.com/jackyczj/July/log"
	"github.com/jackyczj/July/store"
	"github.com/labstack/echo/v4"
)

//name: "download (3).png"
//lastModified: 1571127952000
//lastModifiedDate: Tue Oct 15 2019 16:25:52 GMT+0800 (China Standard Time) {}
//webkitRelativePath: ""
//size: 3172
//type: "image/png"
//uid: "rc-upload-1586246126130-5"

type ImageReq struct {
	Name string `json:"name"`
	Type string `json:"type"`
	Uid  string `json:"uid"`
}

type UploadResp struct {
	Name   string `json:"name"`
	Status string `json:"status"`
	Url    string `json:"url"`
}

func Upload(ctx echo.Context) error {
	file, err := ctx.FormFile("file")
	if err != nil {
		log.Logworker.Error(err)
		return ctx.JSON(500, err)
	}
	var i ImageReq
	err = ctx.Bind(&i)
	if err != nil {
		log.Logworker.Error(err)
		return ctx.JSON(500, err)
	}
	var u UploadResp
	u.Name = i.Name

	f, err := file.Open()
	if err != nil {
		log.Logworker.Error(err)
		u.Status = "error"
		return ctx.JSON(500, u)
	}
	url, err := store.Upload(f, i.Uid)
	defer f.Close()
	u.Url = "http://localhost:2333/image/" + "" + url
	u.Status = "done"
	return ctx.JSON(200, u)
}
