package utils

import (
	"encoding/json"

	uuid "github.com/satori/go.uuid"
	"go.mongodb.org/mongo-driver/bson"
	bson2 "gopkg.in/mgo.v2/bson"
)

func NewUUID() string {
	id := uuid.NewV1()
	return id.String()
}

func JsonToBson(j interface{}) (bson.M, error) {
	m := make(bson.M)
	d, err := json.Marshal(j)
	if err != nil {
		return nil, err
	}
	err = bson2.UnmarshalJSON(d, &m)
	if err != nil {
		return nil, err
	}
	return m, nil
}
