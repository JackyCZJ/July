package store

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/bsonx"

	"github.com/jackyczj/July/utils"

	"go.mongodb.org/mongo-driver/mongo"
)

type item struct {
	ProductId string `bson:"product_id"`
	Count     int    `bson:"count,omitempty"`
}

type Cart struct {
	Owner    int32     `bson:"owner"`
	Item     []item    `bson:"item,omitempty"`
	CreateAt time.Time `bson:"create_at,omitempty"`
	UpdateAt time.Time `bson:"update_at,omitempty"`
}

func CartAdd(id int32, product Product, count int) error {
	var stash Cart
	stash.Owner = id
	filter, err := utils.StructToBson(stash)
	if err != nil {
		return err
	}
	_, err = Client.db.Collection("cart").Indexes().CreateOne(context.TODO(), mongo.IndexModel{
		Keys:    bsonx.Doc{{Key: "owner", Value: bsonx.Int32(1)}},
		Options: options.Index().SetUnique(true),
	})
	i := item{}
	i.ProductId = product.ProductId
	i.Count = count
	if err != nil {
		return err
	}
	update := bson.D{
		{
			Key: "$push",
			Value: bson.D{
				{Key: "item", Value: i},
			},
		},
	}
	result, err := Client.db.Collection("cart").UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return err
	}
	if result.MatchedCount != 1 {
		return fmt.Errorf("Add Cart Error ")
	}
	return nil
}

func CartDel(id int32, product ...Product) error {
	var stash Cart
	stash.Owner = id
	filter, err := utils.StructToBson(stash)
	if err != nil {
		return err
	}
	var it item
	for p := range product {
		it.ProductId = product[p].ProductId
		update := bson.D{
			{
				Key: "$pull",
				Value: bson.D{
					{Key: "item", Value: it},
				},
			},
		}
		result, err := Client.db.Collection("cart").UpdateOne(context.TODO(), filter, update)
		if err != nil {
			return err
		}
		if result.MatchedCount == 0 {
			return fmt.Errorf("Add Cart Error ")
		}
	}
	return nil
}

func CartClear(id int32) error {
	var cart Cart
	cart.Owner = id
	filter, err := utils.StructToBson(cart)
	if err != nil {
		return err
	}
	update := bson.D{
		{
			Key: "$set",
			Value: bson.D{
				{
					"item", []item{},
				},
			},
		},
	}
	return Client.db.Collection("cart").FindOneAndUpdate(context.TODO(), filter, update).Err()
}

func CartList(id int32) (*Cart, error) {
	var cart Cart
	cart.Owner = id
	filter, err := utils.StructToBson(cart)
	if err != nil {
		return nil, err
	}
	err = Client.db.Collection("cart").FindOne(context.TODO(), filter).Decode(&cart)
	if err != nil {
		return nil, err
	}
	return &cart, nil
}
