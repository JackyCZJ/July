package store

import (
	"context"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/jackyczj/July/utils"
	"go.mongodb.org/mongo-driver/bson"

	"github.com/jackyczj/July/cache"
)

type Order struct {
	OrderNo      string    `json:"OrderNo"`
	Seller       string    `json:"seller"`
	Buyer        string    `json:"buyer"`
	Payment      int32     `json:"payment"`
	PaymentType  int       `json:"payment_type"`
	ShippingTo   Address   `json:"shipping_to" bson:"shipping_to"`
	Item         []string  `json:"item"`
	CreateTime   time.Time `json:"create_time"`
	IsClose      bool      `json:"is_close"`
	EndTime      time.Time `json:"end_time"`
	SendTime     time.Time `json:"send_time"`
	sync.RWMutex `bson:"-"`
}

var err error

func (o *Order) Create() error {
	o.Lock()
	defer o.Unlock()
	ctx, cancel := context.WithTimeout(context.TODO(), 3*time.Second)
	defer cancel()
	_, err = Client.db.Collection("Order").InsertOne(ctx, o)
	if err != nil {
		return err
	}
	cache.SetCc("Order."+o.OrderNo, o, 1*time.Hour)
	return nil
}

func (o *Order) Delete() error {
	o.Lock()
	defer o.Unlock()
	ctx, cancel := context.WithTimeout(context.TODO(), 3*time.Second)
	defer cancel()
	m, err := utils.StructToBson(o)
	if err != nil {
		return err
	}
	result := Client.db.Collection("Order").FindOneAndDelete(ctx, m)
	if result.Err() != nil {
		return result.Err()
	}
	cache.DelCc("Order." + o.OrderNo)
	return nil
}

func (o *Order) Update(filed string, value interface{}) error {
	o.Lock()
	defer o.Unlock()
	op := options.FindOneAndUpdate()
	op.SetProjection(bson.D{{Key: "OrderNo", Value: o.OrderNo}})
	update := bson.D{{Key: "$set", Value: bson.D{{Key: filed, Value: value}}}}
	result := Client.db.Collection("Order").FindOneAndUpdate(context.TODO(), o, update)
	if result.Err() != nil {
		return result.Err()
	}
	return nil
}
