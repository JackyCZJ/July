package store

import (
	"context"
	"time"

	"github.com/jackyczj/July/log"

	"go.mongodb.org/mongo-driver/mongo"

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
	Seller      string    `json:"seller" bson:"seller,omitempty"` //卖家商户Id
	Buyer       int32     `json:"buyer" bson:"buyer"`             //买家用户id
	Payment     float64   `json:"payment" bson:"payment" `
	PaymentType int       `json:"payment_type" bson:"payment_type,omitempty" `
	ShippingTo  int       `json:"shipping_to" bson:"shipping_to,omitempty"`
	Item        []Item    `json:"item" bson:"item,omitempty"`
	CreateAt    time.Time `json:"create_at" bson:"create_at,omitempty"`
	Status      string    `json:"status" bson:"status,omitempty"`
	TrackingNum string    `json:"tracking_num" bson:"tracking_num,omitempty"`
	IsClose     bool      `json:"is_close"`
	EndAt       time.Time `json:"end_at"`
	SendAt      time.Time `json:"send_at"`
}

type Item struct {
	ProductId string `json:"product_id" bson:"_id"` //商品id
	Count     int    `json:"count"`
}

var err error

//订单创建
func (o *Order) Create() error {
	id := primitive.NewObjectID()
	o.OrderNo = id.Hex()
	ctx, cancel := context.WithTimeout(context.TODO(), 3*time.Second)
	defer cancel()
	_, err = Client.db.Collection("order").InsertOne(ctx, o)
	if err != nil {
		return err
	}
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
	filter := bson.D{
		{
			"_id", o.OrderNo,
		},
	}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: filed, Value: value}}}}
	result := Client.db.Collection("order").FindOneAndUpdate(context.TODO(), filter, update)
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
	//err := cache.GetCc("Order."+o.OrderNo, &o)
	//if err != cache2.ErrCacheMiss && err != nil {
	//	return err
	//}
	//ob, err := primitive.ObjectIDFromHex(o.OrderNo)
	//if err != nil {
	//	return err
	//}
	filter := bson.D{
		{
			Key: "_id", Value: o.OrderNo,
		},
	}
	result := Client.db.Collection("order").FindOne(context.TODO(), filter)
	if result.Err() != nil {
		return result.Err()
	}
	err = result.Decode(&o)
	if err != nil {
		return err
	}
	//cache.SetCc("Order."+o.OrderNo, o, 1*time.Hour)
	return nil
}

func (o *Order) UpdateAll() error {
	op := options.FindOneAndUpdate()
	op.SetProjection(bson.D{{Key: "_id", Value: o.OrderNo}})
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

type humanOrder struct {
	Order
	Shop string `json:"shopName"`
}

func ShopOrderList(userId int32) []humanOrder {
	var result *mongo.Cursor
	var OrderList []humanOrder
	var o Shop
	o.Owner = userId
	err = o.GetByOwner()
	if err != nil {
		log.Logworker.Error(err)
	}
	result, err = Client.db.Collection("order").Find(context.TODO(),
		bson.D{
			{
				Key:   "seller",
				Value: o.Id,
			},
		},
	)
	if err != nil {
		log.Logworker.Error(err)
		return OrderList
	}

	for result.Next(context.TODO()) {
		var o humanOrder
		if err = result.Decode(&o.Order); err != nil {
			return OrderList
		}
		s := Shop{
			Id: o.Seller,
		}
		err = s.Get()
		if err != nil {
			log.Logworker.Error(err)
		}
		o.Shop = s.Name
		OrderList = append(OrderList, o)
	}
	return OrderList
}
func OrderList(userId int32) []humanOrder {
	var result *mongo.Cursor
	var OrderList []humanOrder
	result, err = Client.db.Collection("order").Find(context.TODO(),
		bson.D{
			{
				Key:   "buyer",
				Value: userId,
			},
		},
	)
	if err != nil {
		log.Logworker.Error(err)
		return OrderList
	}

	for result.Next(context.TODO()) {
		var o humanOrder
		if err = result.Decode(&o.Order); err != nil {
			return OrderList
		}
		s := Shop{
			Id: o.Seller,
		}
		err = s.Get()
		if err != nil {
			log.Logworker.Error(err)
		}
		o.Shop = s.Name
		OrderList = append(OrderList, o)
	}
	return OrderList
}
