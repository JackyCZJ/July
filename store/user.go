package store

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/rs/xid"

	"github.com/go-redis/cache"

	cacheClient "github.com/jackyczj/July/cache"

	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/jackyczj/July/utils"

	"github.com/jackyczj/July/pkg/auth"
	"go.mongodb.org/mongo-driver/bson"
	bson2 "gopkg.in/mgo.v2/bson"
)

type UserInformation struct {
	Id       string `json:"id,omitempty"`
	Username string `json:"username" validate:"min=1,max=32"`
	Password string `json:"password,omitempty" validate:"min=1,max=32"`
	Email    string `json:"email,omitempty"`
	Role     int    `json:"role,omitempty"`
	Gander   int    `json:"gander,omitempty"`
	Phone    string `json:"phone,omitempty"`
	sync.RWMutex
}

func (u *UserInformation) Create() error {
	u.Lock()
	defer u.Unlock()
	var err error
	u.Password, err = auth.Encrypt(u.Password)
	if err != nil {
		return err
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

	m := make(bson.M)
	d, err := json.Marshal(u)
	if err != nil {
		return nil, err
	}
	err = bson2.UnmarshalJSON(d, &m)
	if err != nil {
		return nil, err
	}
	ctx := context.TODO()
	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()
	err = Client.db.Collection("user").FindOne(ctx, m).Decode(&u)
	if err != nil {
		return nil, err
	}
	cacheClient.SetCc("user."+u.Username, u, time.Hour*24)
	return u, nil
}

func (u *UserInformation) GetId() (string, error) {
	u.RLock()
	defer u.RUnlock()
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()
	m, err := utils.JsonToBson(u)
	if err != nil {
		return "", err
	}
	result := Client.db.Collection("user").FindOne(ctx, m)
	if result.Err() != nil {
		return "", err
	}
	s, err := result.DecodeBytes()
	if err != nil {
		return "", err
	}
	id := s.Lookup("Id").String()

	return id, nil
}

func (u *UserInformation) Delete() error {
	u.Lock()
	defer u.Unlock()
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()
	m, err := utils.JsonToBson(u)
	if err != nil {
		return err
	}
	result := Client.db.Collection("user").FindOneAndDelete(ctx, m)
	if result.Err() != nil {
		return err
	}
	cacheClient.DelCc("user." + u.Username)
	return nil
}

func (u *UserInformation) Set() error {
	u.Lock()
	defer u.Unlock()
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()
	m, err := utils.JsonToBson(u)
	if err != nil {
		return err
	}
	op := options.FindOneAndUpdate()
	op.SetProjection(bson.D{{Key: "id", Value: u.Id}})

	result := Client.db.Collection("user").FindOneAndUpdate(ctx, u, m)
	if result.Err() != nil {
		return result.Err()
	}
	err = Client.db.Collection("user").FindOne(ctx, m).Decode(&u)
	if err != nil {
		return err
	}
	cacheClient.SetCc("user."+u.Username, u, time.Hour*24)
	return nil
}
