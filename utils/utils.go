package utils

import (
	uuid "github.com/satori/go.uuid"
	"go.mongodb.org/mongo-driver/bson"
)

func NewUUID() string {
	id := uuid.NewV1()
	return id.String()
}

func StructToBson(j interface{}) (bson.M, error) {
	m := make(bson.M)
	d, err := bson.Marshal(j)
	if err != nil {
		return nil, err
	}
	err = bson.Unmarshal(d, &m)
	if err != nil {
		return nil, err
	}
	return m, nil
}
