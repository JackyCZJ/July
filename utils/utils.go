package utils

import (
	"github.com/jackyczj/fdfs_client"
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

func FastDfs(addr []string, maxCount int) (client *fdfs_client.Client, err error) {
	client, err = fdfs_client.NewClient(addr, maxCount)
	return
}
