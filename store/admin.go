package store

import "context"

//返回总人数，总订单数，总商品数
type Status struct {
	User  int `json:"user"`
	Order int `json:"order"`
	Goods int `json:"goods"`
	Daily int `json:"daily"`
}

func StatusGet() (*Status, error) {
	Order, err := Client.db.Collection("Order").CountDocuments(context.TODO(), nil)
	if err != nil {
		return nil, err
	}
	Goods, err := Client.db.Collection("goods").CountDocuments(context.TODO(), nil)
	if err != nil {
		return nil, err
	}
	User, err := Client.db.Collection("user").CountDocuments(context.TODO(), nil)
	if err != nil {
		return nil, err
	}
	var s Status
	s.Order = int(Order)
	s.Goods = int(Goods)
	s.User = int(User)
	return &s, nil
}
