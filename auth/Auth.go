package auth

import (
	"strings"

	"github.com/go-redis/cache"

	cacheClient "github.com/jackyczj/July/cache"

	"github.com/jackyczj/July/handler/user"
	"github.com/jackyczj/July/log"
	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"
)

// skipper 这些不需要token
func Skipper(c echo.Context) bool {
	method := c.Request().Method
	path := c.Path()
	// 先处理非GET方法，除了登录，现实中还可能有一些 webhooks
	switch path {
	case
		"/user/login",
		"/user/register",
		"/captcha":
		return true
	}
	// 从这里开始必须是GET方法
	if method != "GET" {
		return false
	}
	if path == "" {
		return true
	}
	resource := strings.Split(path, "/")[1]
	switch resource {
	case
		// 公开信息，把需要公开的资源每个一行写这里
		"goods",
		"swagger",
		"public":
		return true
	}
	return false
}

// Validator 校验token是否合法，顺便根据token在 context中赋值 user id
func Validator(token string, c echo.Context) (bool, error) {
	// 调试后门
	log.Logworker.SugaredLogger.Debug("token:", token)
	if viper.GetString("runmode") == "debug" {
		c.Set("user_id", 1)
		return true, nil
	}
	// 寻找token
	var t = new(user.Token)
	err := cacheClient.GetCc("token:"+token, t)
	if err == cache.ErrCacheMiss {
		return false, nil
	} else if err != nil {
		return false, err
	}
	// 设置用户
	c.Set("user_id", t.UserID)

	return true, nil
}
