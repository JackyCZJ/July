package store

import (
	"context"
	"fmt"
	"time"

	"github.com/jackyczj/July/log"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/x/bsonx"

	"github.com/rs/xid"

	cacheClient "github.com/jackyczj/July/cache"

	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/jackyczj/July/utils"

	"github.com/jackyczj/July/pkg/auth"
	"go.mongodb.org/mongo-driver/bson"
)

type UserInformation struct {
	Id        int32     `json:"id,omitempty" bson:"id,omitempty"`
	Username  string    `json:"username" validate:"min=1,max=32" bson:"username,omitempty"`
	Password  string    `json:"-" validate:"min=1,max=32" bson:"password,omitempty"`
	Email     string    `json:"email,omitempty"  bson:"email,omitempty"`
	Avatar    string    `json:"avatar" bson:"avatar,omitempty"`
	Phone     string    `json:"phone,omitempty" bson:"phone,omitempty"`
	Role      int       `json:"role,omitempty" bson:"role,omitempty"`
	Gander    int       `json:"gander,omitempty" bson:"gander,omitempty"`
	Addresses []Address `json:"addresses,omitempty" bson:"addresses,omitempty"`
}

//创建用户
func (u *UserInformation) Create() error {
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
	u.Id = int32(xid.New().Pid())
	_, err = Client.db.Collection("user").InsertOne(ctx, u)
	if err != nil {
		return err
	}
	var cart Cart
	cart.Owner = u.Id
	cart.CreateAt = time.Now()
	cart.UpdateAt = cart.CreateAt
	_, err = Client.db.Collection("cart").InsertOne(ctx, cart)
	if err != nil {
		return err
	}
	cacheClient.SetCc("user."+fmt.Sprint(u.Id), u, time.Hour*24)
	return nil
}

//获取用户信息
func (u *UserInformation) GetUser() (*UserInformation, error) {
	m, err := utils.StructToBson(u)
	if err != nil {
		return nil, err
	}
	err = Client.db.Collection("user").FindOne(context.TODO(), m).Decode(&u)
	if err != nil {
		return nil, err
	}
	return u, nil
}

//删除一位用户
func (u *UserInformation) Delete() error {
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
	cacheClient.DelCc("user." + fmt.Sprint(u.Id))
	return nil
}

//修改用户信息
func (u *UserInformation) Set(filed string, value interface{}) error {
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
	cacheClient.SetCc("user."+fmt.Sprint(u.Id), u, time.Hour*24)
	return nil
}

func (u *UserInformation) ResetPassword(new string) error {
	filter := bson.D{{
		Key: "id", Value: u.Id,
	}}
	en, err := auth.Encrypt(new)
	if err != nil {
		log.Logworker.Fatal(err)
		return err
	}
	update := bson.D{
		{
			Key: "$set", Value: bson.D{{Key: "password", Value: en}},
		},
	}
	result := Client.db.Collection("user").FindOneAndUpdate(context.TODO(), filter, update)
	return result.Err()
}

func UserList(pageNumber int, PerPage int) ([]UserInformation, int, error) {
	var page int
	if pageNumber > 0 {
		page = (pageNumber - 1) * PerPage
	} else {
		page = 0
	}
	opt := options.FindOptions{}

	opt.SetSkip(int64(page))
	opt.SetLimit(int64(PerPage))
	var userList []UserInformation

	count, _ := Client.db.Collection("user").CountDocuments(context.TODO(), bson.M{})

	result, err := Client.db.Collection("user").Find(context.TODO(), bson.M{}, &opt)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return userList, 0, nil
		}
		return nil, 0, err
	}
	for result.Next(context.TODO()) {
		var user UserInformation
		err := result.Decode(&user)
		if err != nil {
			return nil, 0, err
		}
		userList = append(userList, user)
	}
	return userList, int(count), nil
}

func (u *UserInformation) ChangeRole(role int) error {
	filter := u
	update := bson.D{
		{
			"$set",
			bson.D{{"role", role}},
		},
	}
	_, err := Client.db.Collection("user").UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return err
	}

	return nil
}

func UserExist(username string) bool {
	filter := bson.D{
		{
			Key:   "username",
			Value: username,
		}}
	r := Client.db.Collection("user").FindOne(context.TODO(), filter)
	if r.Err() == mongo.ErrNoDocuments {
		return false
	}
	if r.Err() != nil {
		return false
	}
	return true
}

func EmailExist(email string) bool {
	filter := bson.D{
		{
			Key:   "email",
			Value: email,
		}}
	r := Client.db.Collection("user").FindOne(context.TODO(), filter)
	if r.Err() == mongo.ErrNoDocuments {
		return false
	}
	if r.Err() != nil {
		return false
	}
	return true
}

func (u *UserInformation) GetUserAddressList() []Address {
	return u.Addresses
}
