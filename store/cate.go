package store

import (
	"context"

	"github.com/jackyczj/July/log"

	"go.mongodb.org/mongo-driver/bson"
)

type Cate struct {
	Name     string   `json:"name" bson:"name"`
	Children []string `json:"children"`
}

func (c *Cate) InsertCate() error {
	_, err := Client.db.Collection("cate").InsertOne(context.TODO(), c)
	if err != nil {
		return err
	}
	return nil
}

func (c *Cate) DeleteCate() error {
	_, err := Client.db.Collection("cate").DeleteOne(context.TODO(), c)
	if err != nil {
		return err
	}
	return nil
}

func (c *Cate) Get() error {
	r := Client.db.Collection("cate").FindOne(context.TODO(), c)
	if r.Err() != nil {
		return r.Err()
	}
	if err := r.Decode(&c); err != nil {
		return err
	}
	return nil
}

func GetCateTree() []Cate {
	var cateTree []Cate
	r, err := Client.db.Collection("cate").Find(context.TODO(), bson.M{})
	if err != nil {
		log.Logworker.Fatal(err)
	}
	for r.Next(context.TODO()) {
		var c Cate
		_ = r.Decode(&c)
		cateTree = append(cateTree, c)
	}
	return cateTree
}

func (c *Cate) AddToCate(newC string) error {
	update := bson.D{
		{
			Key: "$push",
			Value: bson.D{
				{Key: "children", Value: newC},
			},
		},
	}
	_, err := Client.db.Collection("cate").UpdateOne(context.TODO(), c, update)
	if err != nil {
		return err
	}
	return nil
}

func (c *Cate) DeleteFromCate(del string) error {
	update := bson.D{
		{
			Key: "$pull",
			Value: bson.D{
				{Key: "children", Value: del},
			},
		},
	}
	_, err := Client.db.Collection("cate").UpdateOne(context.TODO(), c, update)
	if err != nil {
		return err
	}
	return nil
}
