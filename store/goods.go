package store

import (
	"context"
	"fmt"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"go.mongodb.org/mongo-driver/x/bsonx"

	"github.com/rs/xid"

	"github.com/jackyczj/July/utils"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Product struct {
	ProductId    uint16    `json:"product_id" bson:"_id"`                    //商品id
	Name         string    `json:"name" bson:"name,omitempty"`               //商品名
	ImageUri     string    `json:"image_uri" bson:"image_uri,omitempty"`     //商品图片url
	Description  string    `json:"description" bson:"description,omitempty"` //商品介绍
	Information  Type      `json:"info" bson:"info,omitempty"`               //品牌
	Price        int       `json:"price" bson:"price,omitempty"`             //价格
	Off          int       `json:"off" bson:"off,omitempty"`                 //折扣
	Owner        string    `json:"owner" bson:"owner,omitempty"`             //拥有者
	CreateAt     time.Time `json:"create_at" bson:"create_at,omitempty"`     //创建时间
	Shelves      bool      `json:"shelves" bson:"shelves,omitempty"`         //是否上架
	IsDelete     bool      `json:"is_delete" bson:"is_delete,omitempty"`     //是否已删除
	sync.RWMutex `bson:"_"`
}

type Type struct {
	Category string `json:"category"` //产品分类
	Brand    string `json:"brand"`    //产品品牌
}

func (p *Product) Add() error {
	p.Lock()
	defer p.Unlock()

	_, err = Client.db.Collection("good").Indexes().CreateOne(context.TODO(), mongo.IndexModel{
		Keys:    bsonx.Doc{{Key: "_id", Value: bsonx.Int32(0)}},
		Options: options.Index().SetUnique(true),
	})
	p.ProductId = xid.New().Pid()
	return Client.client.UseSession(context.TODO(), func(sessionContext mongo.SessionContext) error {
		if err = sessionContext.StartTransaction(); err != nil {
			return err
		}
		_, err = Client.db.Collection("good").InsertOne(sessionContext, p)
		if err != nil {
			fmt.Println(err)
			return sessionContext.AbortTransaction(sessionContext)
		} else {
			return sessionContext.CommitTransaction(sessionContext)
		}
	})
}

func (p *Product) Delete() error {
	p.Lock()
	defer p.Unlock()
	return Client.client.UseSession(context.TODO(), func(sessionContext mongo.SessionContext) error {
		if err = sessionContext.StartTransaction(); err != nil {
			return err
		}
		m, err := utils.StructToBson(p)
		if err != nil {
			return err
		}
		result, err := Client.db.Collection("good").DeleteOne(sessionContext, m)
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

func (p *Product) Set(filed string, value interface{}) error {
	p.Lock()
	defer p.Unlock()
	return Client.client.UseSession(context.TODO(), func(sessionContext mongo.SessionContext) error {
		if err = sessionContext.StartTransaction(); err != nil {
			return err
		}
		op := options.FindOneAndUpdate()
		op.SetProjection(bson.D{{Key: "product_id", Value: p.ProductId}})
		update := bson.D{{Key: "$set", Value: bson.D{{Key: filed, Value: value}}}}
		result := Client.db.Collection("good").FindOneAndUpdate(context.TODO(), p, update)
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

func (p *Product) Get() error {
	p.RLock()
	defer p.RUnlock()
	result := Client.db.Collection("good").FindOne(context.TODO(), p)
	if result.Err() != nil {
		return result.Err()
	}
	if err = result.Decode(&p); err != nil {
		return err
	}
	return nil
}

func (p *Product) Update() error {
	p.Lock()
	defer p.Unlock()
	return Client.client.UseSession(context.TODO(), func(sessionContext mongo.SessionContext) error {
		if err = sessionContext.StartTransaction(); err != nil {
			return err
		}
		var pf Product
		pf.ProductId = p.ProductId
		result := Client.db.Collection("good").FindOneAndUpdate(sessionContext, pf, p)
		if result.Err() != nil {
			return sessionContext.AbortTransaction(sessionContext)
		}
		if err = result.Decode(&p); err != nil {
			return sessionContext.AbortTransaction(sessionContext)
		}
		return sessionContext.CommitTransaction(sessionContext)
	})
}

func GetRandom() ([]Product, error) {
	var pList []Product
	result, err := Client.db.Collection("good").Aggregate(context.Background(),
		mongo.Pipeline{
			bson.D{{"$match", bson.D{
				{"shelves",
					bson.D{
						{"$eq", true},
					}},
			}}},
			bson.D{
				{
					"$sample",
					bson.D{
						{"size", 10},
					},
				},
			},
		})
	if err != nil {
		return nil, err
	}

	for result.Next(context.TODO()) {
		var res Product
		_ = result.Decode(&res)
		pList = append(pList, res)
	}

	return pList, nil
}

func Search(key string, pageNumber int, PerPage int) ([]Product, error) {
	var pList []Product
	filter := bson.D{
		{Key: "$and",
			Value: bson.A{
				bson.D{{Key: "name", Value: primitive.Regex{Pattern: key, Options: ""}}},
				bson.D{{Key: "shelves", Value: bson.D{{"$ne", false}}}},
			},
		},
	}
	opt := options.FindOptions{}
	var page int
	if pageNumber > 0 {
		page = (pageNumber - 1) * PerPage
	} else {
		page = 0
	}
	opt.SetSkip(int64(page))
	opt.SetLimit(int64(PerPage))
	result, err := Client.db.Collection("good").Find(context.TODO(), filter, &opt)
	if err != nil {
		return nil, err
	}
	for result.Next(context.TODO()) {
		var p Product
		err := result.Decode(&p)
		if err != nil {
			return nil, err
		}
		pList = append(pList, p)
	}
	return pList, nil

}

func SearchItself() {

}
