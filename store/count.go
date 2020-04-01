package store

import (
	"context"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
)

func COUNT() echo.MiddlewareFunc {
	return func(handlerFunc echo.HandlerFunc) echo.HandlerFunc {
		go Client.db.Collection("status").InsertOne(context.TODO(), bson.D{
			{
				"daily", 1,
			},
		})
		return handlerFunc
	}
}
