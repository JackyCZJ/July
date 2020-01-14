package store

import "time"

type Product struct {
	ProductId   string `json:"product_id" bson:"_id"`
	Name        string `json:"name"`
	ImageUri    string `json:"image_uri"`
	Information Type   `json:"info"`
	Price       int    `json:"price"`
	Off         int    `json:"off"`
	Owner       string `json:"owner"`
	CreateAt    time.Time
	IsDelete    bool
}

type Type struct {
	Category string `json:"category"` //产品分类
	Brand    string `json:"brand"`    //产品品牌
}
