package store

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/jackyczj/July/utils"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/bsonx"

	"go.mongodb.org/mongo-driver/mongo"
)

type Shop struct {
	Name         string    `json:"name" bson:"name"`
	Owner        string    `json:"owner" bson:"owner"`
	Description  string    `json:"description" bson:"description"`
	IsClose      bool      `json:"is_close" bson:"is_close"`
	CreateAt     time.Time `json:"create_at" bson:"create_at"`
	CloseAt      time.Time `json:"close_at" bson:"close_at"`
	sync.RWMutex `bson:"-"`
	IsDelete     bool `json:"is_delete" bson:"is_delete"`
}

func (s *Shop) Create() error {
	s.Lock()
	defer s.Unlock()
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

func (s *Shop) Delete() error {
	s.Lock()
	defer s.Unlock()
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

func (s *Shop) Set(filed string, value interface{}) error {
	s.Lock()
	defer s.Unlock()
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
