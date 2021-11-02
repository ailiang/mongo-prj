package main

import (
	"MongoDB-Proj/Player"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
)

type Data struct {
	s string
	i int
}

func main() {
	player := Player.NewPlayer("p1")
	player.Save("field1", "value1")
	player.Save("field2", "value2")

	doc, err := bson.Marshal(&bson.D{{"value", "d"}, {"key", 1}})
	if err == nil {
		var ret Data
		err = bson.Unmarshal(doc, &ret)
		if err == nil {
			fmt.Printf("%+v", ret)
		} else {
			fmt.Println(err)
		}
	}

}
