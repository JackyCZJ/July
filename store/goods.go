package store

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/jackyczj/July/utils"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/bsonx"
)

type Product struct {
	ProductId    uint16    `json:"product_id" bson:"_id"`
	Name         string    `json:"name" bson:"name,omitempty"`
	ImageUri     string    `json:"image_uri" bson:"image_uri,omitempty"`
	Information  Type      `json:"info" bson:"info,omitempty"`
	Price        int       `json:"price" bson:"price,omitempty"`
	Off          int       `json:"off" bson:"off,omitempty"`
	Owner        string    `json:"owner" bson:"owner,omitempty"`
	CreateAt     time.Time `json:"create_at" bson:"create_at,omitempty"`
	IsDelete     bool      `json:"is_delete" bson:"is_delete,omitempty"`
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
		Keys:    bsonx.Doc{{Key: "product_id", Value: bsonx.Int32(0)}},
		Options: options.Index().SetUnique(true),
	})

	return Client.client.UseSession(context.TODO(), func(sessionContext mongo.SessionContext) error {
		if err = sessionContext.StartTransaction(); err != nil {
			return err
		}
		_, err = Client.db.Collection("good").InsertOne(sessionContext, p)
		if err != nil {
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
