package store

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/jackyczj/July/utils"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/bsonx"

	"go.mongodb.org/mongo-driver/mongo"
)

type Shop struct {
	Id          string    `bson:"_id,omitempty"`
	Name        string    `json:"name" bson:"name"`
	Owner       string    `json:"owner" bson:"owner,omitempty"`
	Description string    `json:"description" bson:"description,omitempty"`
	IsClose     bool      `json:"is_close" bson:"is_close,omitempty"`
	CreateAt    time.Time `json:"create_at" bson:"create_at,omitempty"`
	CloseAt     time.Time `json:"close_at" bson:"close_at,omitempty"`
	IsDelete    bool      `json:"is_delete" bson:"is_delete,omitempty"`
}

//创建一家店铺的存储方法
func (s *Shop) Create() error {
	_, err = Client.db.Collection("shop").Indexes().CreateOne(context.TODO(), mongo.IndexModel{
		Keys:    bsonx.Doc{{Key: "name", Value: bsonx.String("")}},
		Options: options.Index().SetUnique(true),
	})

	return Client.client.UseSession(context.TODO(), func(sessionContext mongo.SessionContext) error {
		if err = sessionContext.StartTransaction(); err != nil {
			return err
		}
		_, err = Client.db.Collection("shop").InsertOne(sessionContext, s)
		if err != nil {
			return sessionContext.AbortTransaction(sessionContext)
		} else {
			return sessionContext.CommitTransaction(sessionContext)
		}
	})
}

//删除一家店铺
func (s *Shop) Delete() error {
	return Client.client.UseSession(context.TODO(), func(sessionContext mongo.SessionContext) error {
		if err = sessionContext.StartTransaction(); err != nil {
			return err
		}
		m, err := utils.StructToBson(s)
		if err != nil {
			return err
		}
		result, err := Client.db.Collection("shop").DeleteOne(sessionContext, m)
		if err != nil {
			err = sessionContext.AbortTransaction(sessionContext)
			return err
		} else if result.DeletedCount == 0 {
			return fmt.Errorf("Delete Faild ")
		} else {
			return sessionContext.CommitTransaction(sessionContext)
		}
	})
}

//修改店铺信息
func (s *Shop) Set(filed string, value interface{}) error {
	return Client.client.UseSession(context.TODO(), func(sessionContext mongo.SessionContext) error {
		if err = sessionContext.StartTransaction(); err != nil {
			return err
		}
		op := options.FindOneAndUpdate()
		op.SetProjection(bson.D{{Key: "name", Value: s.Name}})
		update := bson.D{{Key: "$set", Value: bson.D{{Key: filed, Value: value}}}}
		result := Client.db.Collection("shop").FindOneAndUpdate(context.TODO(), s, update)
		if result.Err() != nil {
			return result.Err()
		}
		if err != nil {
			return sessionContext.AbortTransaction(sessionContext)
		} else {
			return sessionContext.CommitTransaction(sessionContext)
		}
	})
}

func ShopList(pageNumber int, PerPage int) ([]Shop, int, error) {
	var page int
	if pageNumber > 0 {
		page = (pageNumber - 1) * PerPage
	} else {
		page = 0
	}
	opt := options.FindOptions{}

	opt.SetSkip(int64(page))
	opt.SetLimit(int64(PerPage))
	var shopList []Shop

	filter := bson.D{
		{Key: "is_delete", Value: false},
	}
	count, _ := Client.db.Collection("shop").CountDocuments(context.TODO(), filter)

	result, err := Client.db.Collection("shop").Find(context.TODO(), filter, &opt)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return shopList, 0, nil
		}
		return nil, 0, err
	}
	for result.Next(context.TODO()) {
		var shop Shop
		err := result.Decode(&shop)
		if err != nil {
			return nil, 0, err
		}
		shopList = append(shopList, shop)
	}
	return shopList, int(count), nil
}

func SearchShop(keyword string, pageNumber int, PerPage int) ([]Shop, int, error) {
	filter := bson.D{{
		Key: "$or", Value: []bson.D{
			{{Key: "name", Value: primitive.Regex{Pattern: keyword, Options: ""}}},
			{{Key: "owner", Value: primitive.Regex{Pattern: keyword, Options: ""}}},
		}},
	}
	var shopList []Shop
	opt := options.FindOptions{}
	var page int
	if pageNumber > 0 {
		page = (pageNumber - 1) * PerPage
	} else {
		page = 0
	}
	opt.SetSkip(int64(page))
	opt.SetLimit(int64(PerPage))
	Total, _ := Client.db.Collection("shop").CountDocuments(context.TODO(), filter)
	result, err := Client.db.Collection("shop").Find(context.TODO(), filter, &opt)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return shopList, 0, nil
		}
		return nil, 0, err
	}
	for result.Next(context.TODO()) {
		var shop Shop
		err := result.Decode(&shop)
		if err != nil {
			return nil, 0, err
		}
		shopList = append(shopList, shop)
	}
	return shopList, int(Total), nil
}
