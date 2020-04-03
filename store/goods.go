package store

import (
	"context"
	"fmt"
	"time"

	cache2 "github.com/go-redis/cache"

	"github.com/jackyczj/July/cache"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/jackyczj/July/utils"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Product struct {
	ProductId   string    `json:"product_id" bson:"_id"`                    //商品id
	Name        string    `json:"name" bson:"name,omitempty"`               //商品名
	ImageUri    []string  `json:"image_uri" bson:"image_uri,omitempty"`     //商品图片url
	Description string    `json:"description" bson:"description,omitempty"` //商品介绍
	Information Type      `json:"info" bson:"info,omitempty"`               //品牌
	Price       float64   `json:"price" bson:"price,omitempty"`             //价格
	Store       int       `json:"store" bson:"store,omitempty"`             //库存
	Off         int       `json:"off" bson:"off,omitempty"`                 //折扣
	Owner       string    `json:"owner" bson:"owner,omitempty"`             //拥有者,商店ID
	CreateAt    time.Time `json:"create_at" bson:"create_at,omitempty"`     //创建时间
	Shelves     bool      `json:"shelves" bson:"shelves,omitempty"`         //是否上架
}

type Type struct {
	Category string `json:"category"` //产品分类
	Brand    string `json:"brand"`    //产品品牌
}

func (p *Product) Add() error {
	id := primitive.NewObjectID()
	p.ProductId = id.Hex()
	return Client.client.UseSession(context.TODO(), func(sessionContext mongo.SessionContext) error {
		if err = sessionContext.StartTransaction(); err != nil {
			return err
		}
		_, err = Client.db.Collection("good").InsertOne(sessionContext, p)
		if err != nil {
			fmt.Println(err)
			return sessionContext.AbortTransaction(sessionContext)
		} else {
			return sessionContext.CommitTransaction(sessionContext)
		}
	})
}

func (p *Product) Delete() error {
	return Client.client.UseSession(context.TODO(), func(sessionContext mongo.SessionContext) error {
		if err = sessionContext.StartTransaction(); err != nil {
			return err
		}
		m, err := utils.StructToBson(p)
		if err != nil {
			return err
		}
		result, err := Client.db.Collection("good").DeleteOne(sessionContext, m)
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

func (p *Product) Set(filed string, value interface{}) error {
	return Client.client.UseSession(context.TODO(), func(sessionContext mongo.SessionContext) error {
		if err = sessionContext.StartTransaction(); err != nil {
			return err
		}
		op := options.FindOneAndUpdate()
		op.SetProjection(bson.D{{Key: "product_id", Value: p.ProductId}})
		update := bson.D{{Key: "$set", Value: bson.D{{Key: filed, Value: value}}}}
		result := Client.db.Collection("good").FindOneAndUpdate(context.TODO(), p, update)
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

func (p *Product) Get() error {
	filter := bson.D{{Key: "_id", Value: p.ProductId}}
	result := Client.db.Collection("good").FindOne(context.TODO(), filter)
	if result.Err() != nil {
		return result.Err()
	}
	if err = result.Decode(&p); err != nil {
		return err
	}
	return nil
}

func (p *Product) Update() error {
	return Client.client.UseSession(context.TODO(), func(sessionContext mongo.SessionContext) error {
		if err = sessionContext.StartTransaction(); err != nil {
			return err
		}
		var pf Product
		pf.ProductId = p.ProductId
		result := Client.db.Collection("good").FindOneAndUpdate(sessionContext, pf, p)
		if result.Err() != nil {
			return sessionContext.AbortTransaction(sessionContext)
		}
		if err = result.Decode(&p); err != nil {
			return sessionContext.AbortTransaction(sessionContext)
		}
		return sessionContext.CommitTransaction(sessionContext)
	})
}

func GetRandom() ([]bson.M, error) {
	var pList []bson.M
	result, err := Client.db.Collection("good").Aggregate(context.Background(),
		mongo.Pipeline{
			bson.D{
				{Key: "$match",
					Value: bson.D{
						{Key: "shelves",
							Value: bson.D{
								{Key: "$eq", Value: true},
							},
						},
					},
				},
			},
			bson.D{
				{
					Key: "$sample",
					Value: bson.D{
						{Key: "size", Value: 12},
					},
				},
			},
		})
	if err != nil {
		return nil, err
	}

	for result.Next(context.TODO()) {
		var res bson.M
		_ = result.Decode(&res)
		pList = append(pList, res)
	}

	return pList, nil
}

func Search(key string, pageNumber int, PerPage int) ([]Product, int, error) {
	var pList []Product
	filter := bson.D{
		{Key: "$and",
			Value: bson.A{
				bson.D{{Key: "name", Value: primitive.Regex{Pattern: key, Options: "i"}}},
				bson.D{{Key: "shelves", Value: bson.D{{Key: "$ne", Value: false}}}},
			},
		},
	}
	opt := options.FindOptions{}
	var page int
	if pageNumber > 0 {
		page = (pageNumber - 1) * PerPage
	} else {
		page = 0
	}
	opt.SetSkip(int64(page))
	opt.SetLimit(int64(PerPage))
	Total, _ := Client.db.Collection("good").CountDocuments(context.TODO(), filter)
	if pageNumber > 0 {
		cachekey := fmt.Sprintf("%v.%v.%v", key, pageNumber-1, PerPage)
		var lastid string
		err := cache.GetCc(cachekey, &lastid)

		switch err {
		case nil:
			filter = bson.D{
				{Key: "$and",
					Value: bson.A{
						bson.D{{Key: "name", Value: primitive.Regex{Pattern: key, Options: "i"}}},
						bson.D{{Key: "_id", Value: bson.D{{Key: "$gt", Value: lastid}}}},
						bson.D{{Key: "shelves", Value: bson.D{{Key: "$ne", Value: false}}}},
					},
				},
			}
		case cache2.ErrCacheMiss:
		}
	}
	result, err := Client.db.Collection("good").Find(context.TODO(), filter, &opt)
	if err != nil {
		return nil, 0, err
	}
	var lastId string
	for result.Next(context.TODO()) {
		var p Product
		err := result.Decode(&p)
		if err != nil {
			return nil, 0, err
		}
		pList = append(pList, p)
		lastId = p.ProductId
	}
	if pageNumber > 0 {
		cachekey := fmt.Sprintf("%v.%v.%v", key, pageNumber, PerPage)
		cache.SetCc(cachekey, lastId, 5*time.Minute)
	}

	return pList, int(Total), nil
}

func GetListByShop(shop string, role bool, page int) ([]Product, int64) {
	var pList []Product
	filter := bson.D{{Key: "owner", Value: shop}}
	if role {
		filter = bson.D{{
			Key: "$and",
			Value: bson.D{
				{Key: "owner", Value: shop},
				{Key: "shelves", Value: false}},
		}}
	}
	opt := options.Find()
	if page > 0 {
		page = (page - 1) * 10
	} else {
		page = 0
	}
	fmt.Println(shop, page)
	opt.SetSkip(int64(page))
	opt.SetLimit(10)
	total, _ := Client.db.Collection("good").CountDocuments(context.TODO(), filter)
	pagenumber := total / 10
	result, err := Client.db.Collection("good").Find(context.TODO(), filter, opt)
	if err != nil {
		return nil, 0
	}
	for result.Next(context.TODO()) {
		var p Product
		err := result.Decode(&p)
		if err != nil {
			return nil, 0
		}
		pList = append(pList, p)
	}
	return pList, pagenumber
}

func Suggestion(keyword string) []string {
	//filter := bson.D{{
	//	Key: "name", Value: primitive.Regex{Pattern: keyword, Options: ""}},
	//}
	var resultGroup []string

	if keyword == "" {
		return resultGroup
	}

	result, err := Client.db.Collection("good").Aggregate(context.TODO(), mongo.Pipeline{
		bson.D{
			{Key: "$match",
				Value: bson.D{{
					Key: "name", Value: primitive.Regex{Pattern: keyword, Options: "i"}},
				},
			},
		},
		bson.D{
			{
				Key: "$sample",
				Value: bson.D{
					{Key: "size", Value: 6},
				},
			},
		},
	})
	if err != nil {
		return []string{}
	}
	for result.Next(context.TODO()) {
		var p Product
		err := result.Decode(&p)
		if err != nil {
			return nil
		}
		resultGroup = append(resultGroup, p.Name)
	}
	return resultGroup
}

//Comment module
type Comment struct {
	Username string `json:"username" bson:"username"`
	Rank     int    `json:"rank" bson:"rank"`
	Content  string `json:"content"`
}

func GetComment(id string) []Comment {
	filter := bson.D{
		{
			Key: "product_id", Value: id,
		},
	}
	comment := struct {
		ProductId   string    `bson:"product_id"`
		CommentList []Comment `bson:"comment_list"`
	}{}

	result := Client.db.Collection("comment").FindOne(context.TODO(), filter)

	if err := result.Decode(&comment); err != nil {
		return nil
	}
	return comment.CommentList
}

func AddComment(id string, c Comment) error {
	filter := bson.D{
		{
			Key: "product_id", Value: id,
		},
	}
	update := bson.D{
		{
			Key: "$push",
			Value: bson.D{
				{Key: "comment_list", Value: c},
			},
		},
	}
	return Client.db.Collection("comment").FindOneAndUpdate(context.TODO(), filter, update).Err()
}

func DeleteComment(id string, username string) error {
	filter := bson.D{
		{
			Key: "product_id", Value: id,
		},
	}
	update := bson.D{{
		Key: "$pull",
		Value: bson.D{{
			Key: "username", Value: username,
		}},
	},
	}
	return Client.db.Collection("comment").FindOneAndUpdate(context.TODO(), filter, update).Err()
}
