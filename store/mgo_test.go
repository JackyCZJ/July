package store

import (
	"context"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

func TestOpenDB(t *testing.T) {
	db := openDB("mongodb://mongo1:27017,mongo2:27018,mongo3:27019/?replicaSet=rs0", "testing")

	collection := db.Collection("user")
	res, err := collection.InsertOne(context.Background(), UserInformation{Username: "world"})
	if err != nil {
		t.Fatal(err)
	}
	id := res.InsertedID

	t.Log(id)

	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
		defer cancel()
		m := bson.M{
			"username": "world",
		}
		d, err := collection.DeleteMany(ctx, m)
		if err != nil {
			t.Fatal()
		}
		t.Log(d)

	}()
}
