package store

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/jackyczj/July/log"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/x/bsonx"

	"github.com/go-redis/cache/v7"
	"github.com/rs/xid"

	cacheClient "github.com/jackyczj/July/cache"

	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/jackyczj/July/utils"

	"github.com/jackyczj/July/pkg/auth"
	"go.mongodb.org/mongo-driver/bson"
)

type UserInformation struct {
	Id           string    `json:"id,omitempty" bson:"id,omitempty"`
	Username     string    `json:"username" validate:"min=1,max=32" bson:"username,omitempty"`
	Password     string    `json:"password,omitempty" validate:"min=1,max=32" bson:"password,omitempty"`
	Email        string    `json:"email,omitempty"  bson:"email,omitempty"`
	Role         int       `json:"role,omitempty" bson:"role,omitempty"`
	Gander       int       `json:"gander,omitempty" bson:"gander,omitempty"`
	Addresses    []Address `json:"addresses,omitempty" bson:"addresses,omitempty"`
	sync.RWMutex `bson:"-"`
}

type Address struct {
	Name     string `json:"name"`
	Addr     string `json:"addr"`
	Phone    string `json:"phone"`
	Postal   string `json:"postal"`
	Province string `json:"province"`
}

func (u *UserInformation) Create() error {
	u.Lock()
	defer u.Unlock()
	var err error
	u.Password, err = auth.Encrypt(u.Password)
	if err != nil {
		return err
	}
	_, err = Client.db.Collection("user").Indexes().CreateOne(context.TODO(),
		mongo.IndexModel{
			Keys:    bsonx.Doc{{Key: "username", Value: bsonx.String("")}},
			Options: options.Index().SetUnique(true),
		})
	if err != nil {
		log.Logworker.Error(err.Error())
	}
	ctx, cancel := context.WithTimeout(context.TODO(), 3*time.Second)
	defer cancel()
	u.Id = xid.New().String()
	_, err = Client.db.Collection("user").InsertOne(ctx, u)
	if err != nil {
		return err
	}
	fmt.Println("New member register , id:", u.Id)
	cacheClient.SetCc("user."+u.Username, u, time.Hour*24)
	return nil
}

func (u *UserInformation) GetUser() (*UserInformation, error) {
	u.RLock()
	defer u.RUnlock()
	err := cacheClient.GetCc("user."+u.Username, &u)
	switch err {
	case nil:
		return u, nil
	default:
		return nil, err
	case cache.ErrCacheMiss:
	}

	m, err := utils.StructToBson(u)
	if err != nil {
		return nil, err
	}
	err = Client.db.Collection("user").FindOne(context.TODO(), m).Decode(&u)
	if err != nil {
		return nil, err
	}
	cacheClient.SetCc("user."+u.Username, u, time.Hour*24)
	return u, nil
}

func (u *UserInformation) Delete() error {
	u.Lock()
	defer u.Unlock()
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()
	m, err := utils.StructToBson(u)
	if err != nil {
		return err
	}
	result, err := Client.db.Collection("user").DeleteOne(ctx, m)
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return fmt.Errorf("Delete Faild ")
	}
	cacheClient.DelCc("user." + u.Username)
	return nil
}

func (u *UserInformation) Set(filed string, value interface{}) error {
	u.Lock()
	defer u.Unlock()
	if err != nil {
		return err
	}
	op := options.FindOneAndUpdate()
	op.SetProjection(bson.D{{Key: "id", Value: u.Id}})
	update := bson.D{{Key: "$set", Value: bson.D{{Key: filed, Value: value}}}}
	result := Client.db.Collection("user").FindOneAndUpdate(context.TODO(), u, update)
	if result.Err() != nil {
		return result.Err()
	}
	err = Client.db.Collection("user").FindOne(context.TODO(), u).Decode(&u)
	if err != nil {
		return err
	}
	cacheClient.SetCc("user."+u.Username, u, time.Hour*24)
	return nil
}
