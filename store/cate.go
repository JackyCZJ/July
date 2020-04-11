package store

import (
	"context"
	"fmt"

	"github.com/jackyczj/July/log"

	"go.mongodb.org/mongo-driver/bson"
)

type Cate struct {
	Name   string `bson:"_id" json:"name"`
	Parent string `json:"parent" bson:"parent"`
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

func GetCateByParent(name string) []Cate {
	fmt.Println(name)
	filter := bson.D{
		{
			"parent", name,
		},
	}
	c := make([]Cate, 0)
	r, err := Client.db.Collection("cate").Find(context.TODO(), filter)
	if err != nil {
		return nil
	}
	for r.Next(context.TODO()) {
		var a Cate
		_ = r.Decode(&a)
		c = append(c, a)
	}
	return c
}

func GetCateTree() []Cate {
	var cateTree []Cate
	r, err := Client.db.Collection("cate").Find(context.TODO(), bson.D{{Key: "parent", Value: "root"}})
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

func (c *Cate) DeleteFromCate(del string) error {
	update := bson.D{
		{
			Key:   "name",
			Value: del,
		},
	}
	_, err := Client.db.Collection("cate").DeleteOne(context.TODO(), update)
	if err != nil {
		return err
	}
	return nil
}
