package db

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"sync"
	"time"
)

const dbUri = "mongodb://localhost:27017"

var dbInstance *DbManager
var dbInstanceOnce sync.Once

type DbManager struct {
	Client *mongo.Client
}

func NewClient() *mongo.Client {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(dbUri))
	if err != nil {
		return nil
	}
	return client
}

func GetDbManager() *DbManager {
	dbInstanceOnce.Do(func() {
		dbInstance = &DbManager{
			Client: NewClient(),
		}
	})
	return dbInstance
}
