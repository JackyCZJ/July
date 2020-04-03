package store

import (
	"context"
	"time"

	"github.com/jackyczj/July/log"

	"go.uber.org/zap"

	"github.com/spf13/viper"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Mgo struct {
	db     *mongo.Database
	client *mongo.Client
}

var Client Mgo

func openDB(url string, database string) *mongo.Database {
	client := openClient(url)
	db := client.Database(database)
	return db
}

func InitDB() *mongo.Database {
	viper.SetDefault("mgo.database", "July")
	viper.SetDefault("mgo.url", "mongodb://mongo1:27017,mongo2:27018,mongo3:27019/?replicaSet=rs0")
	database := viper.GetString("mgo.database")
	url := viper.GetString("mgo.url")
	return openDB(url, database)
}

func openClient(url string) *mongo.Client {
	client, err := mongo.NewClient(options.Client().ApplyURI(url))
	if err != nil {
		log.Logworker.Fatal(err)
	}
	ctx := context.TODO()
	err = client.Connect(ctx)
	if err != nil {
		log.Logworker.Fatal(err)
	}
	return client
}

func (m *Mgo) Init() {
	viper.SetDefault("mgo.database", "July")
	viper.SetDefault("mgo.url", "mongodb://mongo1:27017,mongo2:27018,mongo3:27019/?replicaSet=rs0")
	database := viper.GetString("mgo.database")
	m.client = openClient(viper.GetString("mgo.url"))
	m.db = m.client.Database(database)
	createIndex("good", "name")
	Client = *m
}

func (m *Mgo) Close() {
	err := Client.db.Client().Disconnect(context.TODO())
	if err != nil {
		panic(err)
	}
}

func createIndex(collection string, keys string) {
	db := Client.db.Collection(collection)
	opts := options.CreateIndexes().SetMaxTime(10 * time.Second)

	indexView := db.Indexes()
	model := bson.M{
		keys: 1,
	}
	// 创建索引
	result, err := indexView.CreateOne(
		context.Background(),
		mongo.IndexModel{
			Keys:    model,
			Options: options.Index().SetUnique(false),
		},
		opts,
	)
	if result == "" || err != nil {
		log.Logworker.Error("EnsureIndex error", zap.String("error", err.Error()))
	}
}
