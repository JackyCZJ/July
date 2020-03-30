package auth

import (
	"strings"

	"github.com/go-redis/cache"

	cacheClient "github.com/jackyczj/July/cache"

	"github.com/jackyczj/July/handler/user"
	"github.com/labstack/echo/v4"
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
	if strings.HasPrefix(path, "/user/checkUsername") {
		return true
	}
	if strings.HasPrefix(path, "/api/v1/Goods/") {
		return true
	}
	switch path {
	case "",
		"/api/v1/Shop/List",
		"/api/v1/Goods/index":
		return true
	}
	resource := strings.Split(path, "/")[1]
	switch resource {
	case
		// 公开信息，把需要公开的资源每个一行写这里
		"swagger",
		"public",
		"image":
		return true
	}
	return false
}

// Validator 校验token是否合法，顺便根据token在 context中赋值 user id
func Validator(token string, c echo.Context) (bool, error) { // 寻找token
	var t = new(user.Token)
	err := cacheClient.GetCc("token:"+token, t)
	if err == cache.ErrCacheMiss {
		return false, nil
	} else if err != nil {
		return false, err
	}
	// 设置用户
	c.Set("user_id", t.UserID)
	c.Set("role", t.Role)
	return true, nil
}
