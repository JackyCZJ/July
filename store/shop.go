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
	Name        string    `json:"name" bson:"name,omitempty"`
	Owner       int32     `json:"owner" bson:"owner,omitempty"`
	Manager     []int32   `json:"manager" bson:"manager,omitempty"`
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

//商店列表
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

//搜索商店
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

//根据ID获取商店信息
func (s *Shop) Get() error {
	i, err := primitive.ObjectIDFromHex(s.Id)
	if err != nil {
		return err
	}
	filter := bson.D{
		{
			"_id", i,
		},
	}
	result := Client.db.Collection("shop").FindOne(context.TODO(), filter)
	if result.Err() != nil {
		return result.Err()
	}
	if err := result.Decode(&s); err != nil {
		return err
	}
	return nil
}

//根据拥有者Id获取商店
func (s *Shop) GetByOwner() error {
	filter := bson.D{
		{
			"owner", s.Owner,
		},
	}
	result := Client.db.Collection("shop").FindOne(context.TODO(), filter)
	if result.Err() != nil {
		return result.Err()
	}
	if err := result.Decode(&s); err != nil {
		return err
	}
	return nil
}

type ShopStatus struct {
	Shop
	UnPay     int     `json:"un_pay"`
	UnShip    int     `json:"un_ship"`
	Shipping  int     `json:"shipping"`
	Received  int     `json:"received"`
	Commented int     `json:"commented"`
	Money     float64 `json:"money"`
	OrderList []Order `json:"order_list"`
}

const UnPay = "0"
const UnShip = "1"
const SHIPPING = "2"
const RECEIVED = "3"
const COMMENTED = "4"

//获取商店状态
func (s *Shop) Status() (*ShopStatus, error) {
	filter := bson.D{
		{"owner", s.Owner},
	}
	var status ShopStatus
	var OrderList []Order
	if s.Get() != nil {
		return nil, err
	}
	status.Shop = *s
	result, err := Client.db.Collection("Order").Find(context.TODO(), filter)
	if err != nil {
		if err != mongo.ErrNoDocuments {

		}
		return nil, err
	}
	err = result.Decode(&OrderList)
	if err != nil {
		return nil, err
	}
	for i := range OrderList {
		switch OrderList[i].Status {
		case UnPay:
			status.UnPay += 1
		case UnShip:
			status.UnShip += 1
		case SHIPPING:
			status.Shipping += 1
		case RECEIVED:
			status.Money += OrderList[i].Payment
			status.Received += 1
		case COMMENTED:
			status.Money += OrderList[i].Payment
			status.Commented += 1
		default:
		}
	}
	status.OrderList = OrderList

	return &status, nil
}

//Shop modify ;current can only modify description ，name ,and Close status ,
func (s *Shop) ShopModify() error {
	i, err := primitive.ObjectIDFromHex(s.Id)
	if err != nil {
		return err
	}
	filter := bson.D{
		{
			"_id", i,
		},
	}
	update := bson.D{{
		Key: "$set", Value: bson.D{
			{"description", s.Description},
			{"is_close", s.IsClose},
			{"name", s.Name},
		},
	}}
	result := Client.db.Collection("shop").FindOneAndUpdate(context.TODO(), filter, update)
	if result.Err() != nil {
		return result.Err()
	}
	return nil
}
