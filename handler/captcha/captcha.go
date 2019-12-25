package captcha

import (
	"math/rand"
	"time"

	"github.com/mojocn/base64Captcha"

	"github.com/jackyczj/July/cache"

	"github.com/jackyczj/July/handler"

	"github.com/labstack/echo/v4"
)

func init() {
	//init rand seed
	rand.Seed(time.Now().UnixNano())
}

const (
	// Default number of digits in captcha solution.
	DefaultLen      = 6
	DefaultDotCount = 30
	// The number of captchas created that triggers garbage collection used
	// by default store.
	CollectNum     = 100
	DefaultMaxSkew = 3
	// Expiration time of captchas used by default store.
	Expiration = 10 * time.Minute
	// Standard width and height of a captcha image.
	StdWidth  = 240
	StdHeight = 80
)

type ConfigJsonBody struct {
	Id          string `json:"id"`
	VerifyValue string `json:"verify_value"`
	DriverDigit *base64Captcha.DriverDigit
}

var store Store

//implementing Store interface
type Store struct {
}

func (Store) Set(id string, value string) {
	cache.SetCc(id, value, 10*time.Minute)
}

func (Store) Get(id string, clear bool) string {
	var s string
	err := cache.GetCc(id, s)
	if err != nil {
		return ""
	}
	if clear {
		cache.DelCc(id)
	}
	return s
}
func (s Store) Verify(id, answer string, clear bool) bool {
	return answer == s.Get(id, clear)
}

func Generate(e echo.Context) error {
	return generateCaptchaHandler(e)
}

func Verify(e echo.Context) error {
	return captchaVerifyHandle(e)
}

// base64Captcha create http handler
func generateCaptchaHandler(e echo.Context) error {
	//parse request parameters

	postParameters := new(ConfigJsonBody)
	err := e.Bind(&postParameters)
	if err != nil {
		return err
	}
	postParameters.DriverDigit = &base64Captcha.DriverDigit{
		Height:   StdHeight,
		Width:    StdWidth,
		Length:   DefaultLen,
		MaxSkew:  DefaultMaxSkew,
		DotCount: DefaultDotCount,
	}

	//create base64 encoding captcha
	c := base64Captcha.NewCaptcha(postParameters.DriverDigit, store)
	//or you can just write the captcha content to the httpResponseWriter.
	//before you put the captchaId into the response COOKIE.
	//captcaInterfaceInstance.WriteTo(w)

	//set json response
	id, b64s, err := c.Generate()

	res := handler.ResponseStruct{
		Code: 1,
		Data: struct {
			Id   string `json:"id"`
			Data string `json:"data"`
		}{
			id,
			b64s,
		},
		Message: "success",
	}
	if err != nil {
		res.Code = 0
		res.Data = nil
		res.Message = err.Error()
	}
	return handler.Response(e, res)
}

// base64Captcha verify http handler
func captchaVerifyHandle(ctx echo.Context) error {

	//parse request parameters
	//接收客户端发送来的请求参数
	var postParameters ConfigJsonBody
	err := ctx.Bind(&postParameters)
	if err != nil {
		return err
	}
	//verify the captcha
	//比较图像验证码

	verifyResult := store.Verify(postParameters.Id, postParameters.VerifyValue, false)
	//set json response
	//设置json响应
	res := handler.ResponseStruct{
		Code:    0,
		Message: "captcha failed",
		Data:    "验证失败",
	}
	if verifyResult {
		res.Code = 1
		res.Data = "验证通过"
		res.Message = "captcha verified"
	}
	return handler.Response(ctx, res)
}
