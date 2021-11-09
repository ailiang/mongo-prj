package db

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"sync"
)

const dbUri = "mongodb://10.236.254.121:27017"

var dbInstance *DbManager
var dbInstanceOnce sync.Once

type DbManager struct {
	Client *mongo.Client
}

func NewClient() *mongo.Client {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(dbUri))
	if err != nil {
		panic(err)
	}

	if e := client.Ping(context.TODO(), readpref.Primary()); e != nil {
		panic(e)
	}

	return client
}

func (db *DbManager) GetClient() *mongo.Client {
	return db.Client
}

func GetDbManager() *DbManager {
	dbInstanceOnce.Do(func() {
		dbInstance = &DbManager{
			Client: NewClient(),
		}
	})
	return dbInstance
}
