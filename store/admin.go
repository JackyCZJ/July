package store

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"

	"go.mongodb.org/mongo-driver/mongo/options"

	"go.mongodb.org/mongo-driver/mongo"
)

//返回总人数，总订单数，总商品数
type Status struct {
	User  int `json:"user"`
	Order int `json:"order"`
	Goods int `json:"goods"`
	Daily int `json:"daily"`
}

func StatusGet() (*Status, error) {
	opts := options.Count().SetMaxTime(2 * time.Second)
	Order, err := Client.db.Collection("order").CountDocuments(context.TODO(), bson.M{}, opts)
	switch err {
	case mongo.ErrNilDocument:
		Order = 0
	case nil:
	default:
		return nil, err
	}
	Goods, err := Client.db.Collection("good").CountDocuments(context.TODO(), bson.M{}, opts)
	switch err {
	case mongo.ErrNilDocument:
		Goods = 0
	case nil:
	default:
		return nil, err
	}
	User, err := Client.db.Collection("user").CountDocuments(context.TODO(), bson.M{}, opts)
	switch err {
	case mongo.ErrNilDocument:
		User = 0
	case nil:
	default:
		return nil, err
	}
	var s Status
	s.Order = int(Order)
	s.Goods = int(Goods)
	s.User = int(User)
	return &s, nil
}
