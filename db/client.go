package db

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func MongoUpdateOneWithColl(col *mongo.Collection, key string, fields []string, values []interface{}) error {
	if len(fields) != len(values) {
		panic("db save k!=v")
	}
	filterE := bson.E{DataSaveKey, key}
	filterD := bson.D{filterE}
	updateValueD := make([]bson.E, 0, len(fields)+1)
	updateValueD = append(updateValueD, filterE)
	for i := 0; i < len(fields); i++ {
		updateValueD = append(updateValueD, bson.E{fields[i], values[i]})
	}
	updateE := bson.E{"$set", updateValueD}
	ops := options.Update().SetUpsert(true)
	result, err := col.UpdateOne(context.TODO(), filterD, bson.D{updateE}, ops)
	if err == nil {
		fmt.Printf("%+v", result)
	} else {
		println(err.Error())
	}
	return err
}

func MongoGetOneWithColl(col *mongo.Collection, key string) (reply bson.D, err error) {
	filterE := bson.E{DataSaveKey, key}
	filterD := bson.D{filterE}
	var r bson.D
	err = col.FindOne(context.TODO(), filterD).Decode(&r)
	if err == nil {
		reply = r
	}
	return
}

func MongoGetOneRawWithColl(col *mongo.Collection, key string) (r bson.Raw, err error) {
	filterE := bson.E{DataSaveKey, key}
	filterD := bson.D{filterE}
	return col.FindOne(context.TODO(), filterD).DecodeBytes()
}

func MongoGetOneFiledRawWithColl(col *mongo.Collection, key string, filedName string, ret interface{}) error {
	filterE := bson.E{DataSaveKey, key}
	filterD := bson.D{filterE}
	projection := bson.D{{filedName, 1}, {"_id", 0}}
	opts := options.FindOne().SetProjection(projection)
	if r, e := col.FindOne(context.TODO(), filterD, opts).DecodeBytes(); e == nil {
		return r.Lookup(filedName).Unmarshal(ret)
	} else {
		return e
	}
}
