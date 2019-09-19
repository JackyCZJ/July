package store

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/jackyczj/NoGhost/utils"

	"github.com/jackyczj/NoGhost/pkg/auth"
	"go.mongodb.org/mongo-driver/bson"
	bson2 "gopkg.in/mgo.v2/bson"
)

type UserInformation struct {
	Id       string `json:"_id,omitempty"`
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

	id, err := Client.db.Collection("user").InsertOne(context.TODO(), u)
	if err != nil {
		return err
	}
	fmt.Println("New member register , id:", id.InsertedID)
	return nil
}

func (u *UserInformation) GetUser() (*UserInformation, error) {
	u.RLock()
	defer u.RUnlock()
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()
	m := make(bson.M)
	d, err := json.Marshal(u)
	if err != nil {
		return nil, err
	}
	err = bson2.UnmarshalJSON(d, &m)
	if err != nil {
		return nil, err
	}

	err = Client.db.Collection("user").FindOne(ctx, m).Decode(&u)
	if err != nil {
		return nil, err
	}

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
	a := s.Lookup("_id").String()
	var id objectId
	if err := json.Unmarshal([]byte(a), &id); err != nil {
		return "", err
	}
	return id.Id, nil
}

type objectId struct {
	Id string `json:"$oid"`
}

func (u *UserInformation) del() error {
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
	return nil
}
