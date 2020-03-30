package auth

import (
	"net/http"
	"strconv"

	"github.com/casbin/casbin/v2"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type (
	// Config defines the config for CasbinAuth middleware.
	Config struct {
		// Skipper defines a function to skip middleware.
		Skipper middleware.Skipper

		// Enforcer CasbinAuth main rule.
		// Required.
		Enforcer *casbin.Enforcer
	}
)

var (
	// DefaultConfig is the default CasbinAuth middleware config.
	DefaultConfig = Config{
		Skipper: middleware.DefaultSkipper,
	}
)

func MiddlewareWithConfig(config Config) echo.MiddlewareFunc {
	// Defaults
	if config.Skipper == nil {
		config.Skipper = DefaultConfig.Skipper
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if config.Skipper(c) {
				return next(c)
			}

			if pass, err := config.CheckPermission(c); err == nil && pass {
				return next(c)
			} else if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
			}

			return echo.ErrForbidden
		}
	}
}

func Middleware(ce *casbin.Enforcer) echo.MiddlewareFunc {
	c := DefaultConfig
	c.Enforcer = ce
	return MiddlewareWithConfig(c)
}

func (a *Config) CheckPermission(e echo.Context) (bool, error) {
	user := e.Get("user_id")
	role := e.Get("role")
	if user == nil || role == nil {
		user = "guest"
		_, err := a.Enforcer.AddRoleForUser(user.(string), "0")
		if err != nil {
			return false, err
		}
	}
	role = strconv.Itoa(role.(int))

	_, err := a.Enforcer.AddRoleForUser(user.(string), role.(string))
	if err != nil {
		return false, err
	}
	method := e.Request().Method
	path := e.Request().URL.Path
	return a.Enforcer.Enforce(role, path, method)
}
