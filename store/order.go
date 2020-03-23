package store

import (
	"context"
	"time"

	cache2 "github.com/go-redis/cache"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/jackyczj/July/utils"
	"go.mongodb.org/mongo-driver/bson"

	"github.com/jackyczj/July/cache"
)

/*
	OrderNo 	订单号
	Seller  	卖家
	Buyer		卖家
	Payment 	价格
	PaymentType 支付方式
	ShippingTo	邮寄地址
	IsClose		订单是否已关闭
	CreateAt 	订单创建时间

*/
type Order struct {
	OrderNo     string    `json:"OrderNo" bson:"_id"`
	Seller      string    `json:"seller" bson:"seller,omitempty"` //买家商户Id
	Buyer       int32     `json:"buyer" bson:"buyer"`             //买家用户id
	Payment     float64   `json:"payment" bson:"payment" `
	PaymentType int       `json:"payment_type" bson:"payment_type,omitempty" `
	ShippingTo  int       `json:"shipping_to" bson:"shipping_to,omitempty"`
	Item        []Item    `json:"item" bson:"item,omitempty"`
	CreateTime  time.Time `json:"create_time" bson:"create_time,omitempty"`
	Status      string    `json:"status" bson:"status,omitempty"`
	TrackingNum string    `json:"tracking_num" bson:"tracking_num,omitempty"`
	IsClose     bool      `json:"is_close"`
	EndTime     time.Time `json:"end_time"`
	SendTime    time.Time `json:"send_time"`
}

type Item struct {
	ProductId int32 `json:"product_id" bson:"_id"` //商品id
	Count     int   `json:"count"`
}

var err error

//订单创建
func (o *Order) Create() error {
	ctx, cancel := context.WithTimeout(context.TODO(), 3*time.Second)
	defer cancel()
	_, err = Client.db.Collection("order").InsertOne(ctx, o)
	if err != nil {
		return err
	}
	cache.SetCc("Order."+o.OrderNo, o, 1*time.Hour)
	return nil
}

//订单删除
func (o *Order) Delete() error {
	ctx, cancel := context.WithTimeout(context.TODO(), 3*time.Second)
	defer cancel()
	m, err := utils.StructToBson(o)
	if err != nil {
		return err
	}
	result := Client.db.Collection("order").FindOneAndDelete(ctx, m)
	if result.Err() != nil {
		return result.Err()
	}
	cache.DelCc("Order." + o.OrderNo)
	return nil
}

//订单修改
func (o *Order) Update(filed string, value interface{}) error {
	op := options.FindOneAndUpdate()
	op.SetProjection(bson.D{{Key: "OrderNo", Value: o.OrderNo}})
	update := bson.D{{Key: "$set", Value: bson.D{{Key: filed, Value: value}}}}
	result := Client.db.Collection("order").FindOneAndUpdate(context.TODO(), o, update)
	if result.Err() != nil {
		return result.Err()
	}
	err = result.Decode(&o)
	if err != nil {
		return err
	}
	return nil
}

func (o *Order) Get() error {
	err := cache.GetCc("Order."+o.OrderNo, &o)
	if err != cache2.ErrCacheMiss && err != nil {
		return nil
	}
	ob, err := primitive.ObjectIDFromHex(o.OrderNo)
	if err != nil {
		return err
	}
	filter := bson.D{
		{
			"_id", ob,
		},
	}
	result := Client.db.Collection("order").FindOne(context.TODO(), filter)
	if result.Err() != nil {
		return result.Err()
	}
	cache.SetCc("Order."+o.OrderNo, o, 1*time.Hour)
	return nil
}

func (o *Order) UpdateAll() error {
	op := options.FindOneAndUpdate()
	op.SetProjection(bson.D{{Key: "OrderNo", Value: o.OrderNo}})
	update := bson.D{{Key: "$set", Value: o}}
	result := Client.db.Collection("order").FindOneAndUpdate(context.TODO(), o, update)
	if result.Err() != nil {
		return result.Err()
	}
	err = result.Decode(&o)
	if err != nil {
		return err
	}
	return nil
}

func OrderList(userId int32, role int) []Order {
	filter := bson.D{}
	u := UserInformation{}
	u.Id = userId

	switch role {
	case 1:
		filter = bson.D{
			{"buyer", userId},
		}
	case 2:
		s := Shop{}
		s.Owner = u.Id
		_ = s.GetByOwner()
		filter = bson.D{
			{
				Key: "$or", Value: bson.D{
					{"buyer", userId},
					{"seller", s.Id},
				},
			},
		}
	}
	var OrderList []Order
	result, err := Client.db.Collection("order").Find(context.TODO(), filter)
	if err != nil {
		return OrderList
	}

	if result.Next(context.TODO()) {
		var o Order
		if err = result.Decode(&o); err != nil {
			return OrderList
		}
		OrderList = append(OrderList, o)
	}
	return OrderList
}
