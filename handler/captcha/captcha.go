package captcha

import (
	"encoding/json"
	"log"
	"time"

	"github.com/labstack/echo/v4"

	"github.com/mojocn/base64Captcha"
)

const (
	// Default number of digits in captcha solution.
	DefaultLen = 6
	// The number of captchas created that triggers garbage collection used
	// by default store.
	CollectNum = 100
	// Expiration time of captchas used by default store.
	Expiration = 10 * time.Minute
	// Standard width and height of a captcha image.
	StdWidth  = 240
	StdHeight = 80
)

type ConfigJsonBody struct {
	Id              string
	CaptchaType     string
	VerifyValue     string
	ConfigAudio     base64Captcha.ConfigAudio
	ConfigCharacter base64Captcha.ConfigCharacter
	ConfigDigit     base64Captcha.ConfigDigit
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
	//接收客户端发送来的请求参数

	decoder := json.NewDecoder(e.Request().Body)
	var postParameters ConfigJsonBody
	err := decoder.Decode(&postParameters)
	if err != nil {
		log.Println(err)
	}
	defer e.Request().Body.Close()

	//create base64 encoding captcha

	var config interface{}
	switch postParameters.CaptchaType {
	case "audio":
		config = postParameters.ConfigAudio
	case "character":
		config = postParameters.ConfigCharacter
	default:
		config = postParameters.ConfigDigit
	}
	captchaId, captcaInterfaceInstance := base64Captcha.GenerateCaptcha(postParameters.Id, config)
	base64blob := base64Captcha.CaptchaWriteToBase64Encoding(captcaInterfaceInstance)

	//or you can just write the captcha content to the httpResponseWriter.
	//before you put the captchaId into the response COOKIE.
	//captcaInterfaceInstance.WriteTo(w)

	//set json response
	body := map[string]interface{}{"code": 1, "data": base64blob, "captchaId": captchaId, "msg": "success"}
	return e.JSON(200, body)
}

// base64Captcha verify http handler
func captchaVerifyHandle(ctx echo.Context) error {

	//parse request parameters
	//接收客户端发送来的请求参数
	decoder := json.NewDecoder(ctx.Request().Body)
	var postParameters ConfigJsonBody
	err := decoder.Decode(&postParameters)
	if err != nil {
		log.Println(err)
	}
	defer ctx.Request().Body.Close()
	//verify the captcha
	//比较图像验证码
	verifyResult := base64Captcha.VerifyCaptcha(postParameters.Id, postParameters.VerifyValue)

	//set json response
	//设置json响应
	body := map[string]interface{}{"code": "error", "data": "验证失败", "msg": "captcha failed"}
	if verifyResult {
		body = map[string]interface{}{"code": "success", "data": "验证通过", "msg": "captcha verified"}
	}
	return ctx.JSON(200, body)
}
