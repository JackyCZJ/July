package store

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

type Wallet struct {
	Uid           int32       `json:"uid"`
	Balance       float64     `json:"balance"`      //余额
	PayPassword   string      `json:"pay_password"` //支付密码
	Name          string      `json:"name"`
	AlipayAccount string      `json:"alipay_account"`
	Log           []WalletLog `json:"log"`
}

type WalletLog struct { //钱包记录
	LogId           string    `json:"log_id"` //log 的id
	Usage           string    `json:"usage"`  //用途
	Money           float64   `json:"money"`  //金额
	TransactionType string    `json:"transaction_type"`
	Change          bool      `json:"change"`      //更改，true为增，false为减
	CreateTime      time.Time `json:"create_time"` //创建日期
}

func (*Wallet) AddWalletLog(log WalletLog, uid int32) error {
	filter := bson.D{
		{
			Key: "uid", Value: uid,
		},
	}
	update := bson.D{
		{Key: "$push",
			Value: bson.D{
				{Key: "log", Value: log},
			},
		},
	}
	return Client.db.Collection("wallet").FindOneAndUpdate(context.TODO(), filter, update).Err()
}

//初始支付密码为登录密码
func (w *Wallet) SetPayPassword(OldPassword string, NewPassword string) {

	//filter := bson.D{
	//	{
	//		Key: "uid", Value: w.Uid,
	//	},
	//}
	//update:= bson.D{
	//	{
	//		Key: "pay_password",
	//	}
	//}

}
