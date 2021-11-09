package main

import (
	"MongoDB-Proj/Player"
	"fmt"
	"reflect"
)

func main() {
	player := Player.NewPlayer("htzx", "player", "player@qq@001")
	//player.Save(data.TablePlayerInfoKey, data.PlayerInfo{
	//	PlayerId:   101,
	//	PlayerName: "player101",
	//})
	//player.Save(data.TableReputationKey, data.ReputationData{map[int32]int32{1: 10, 2: 2, 3: 3}})
	result := player.Get()

	fmt.Printf("result: %+v", result)
	fmt.Printf("%s", reflect.TypeOf(result))
	ret := make(chan bool)
	<-ret
	fmt.Println("donw")

}
