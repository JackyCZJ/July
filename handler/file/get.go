package file

import (
	"github.com/jackyczj/July/store"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func Image(ctx echo.Context) error {
	file := ctx.Param("filename")
	id, err := primitive.ObjectIDFromHex(file)
	if err != nil {
		return echo.NewHTTPError(404, err)
	}
	err = store.Download(id, ctx.Response().Writer)
	if err != nil {
		return echo.NewHTTPError(404, err)
	}
	return ctx.NoContent(200)
}
