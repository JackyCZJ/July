package store

import (
	"context"
	"log"

	"github.com/spf13/viper"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Mgo struct {
	db *mongo.Database
}

var Client Mgo

func openDB(url string, database string) *mongo.Database {
	client := openClient(url)
	db := client.Database(database)
	return db
}

func InitDB() *mongo.Database {
	viper.SetDefault("mgo.database", "user")
	viper.SetDefault("mgo.url", "mongodb://mongo1:27017,mongo2:27018,mongo3:27019/?replicaSet=rs0")
	database := viper.GetString("mgo.database")
	url := viper.GetString("mgo.url")
	return openDB(url, database)
}

func openClient(url string) *mongo.Client {
	client, err := mongo.NewClient(options.Client().ApplyURI(url))
	if err != nil {
		log.Fatal(err)
	}
	ctx := context.Background()
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	return client
}

func (m *Mgo) Init() {
	Client = Mgo{
		db: InitDB(),
	}
}

func (m *Mgo) Close() {
	err := Client.db.Client().Disconnect(context.TODO())
	if err != nil {
		panic(err)
	}
}
