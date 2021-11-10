package main

import (
	"MongoDB-Proj/Player"
	"MongoDB-Proj/game/data"
	"fmt"
)

func TestSave() {
	player := Player.NewPlayer("htzx", "player", "player@qq@001")
	player.Save(data.TablePlayerInfoKey, data.PlayerInfo{
		PlayerId:   101,
		PlayerName: "player101",
	})
	player.Save(data.TableReputationKey, data.ReputationData{map[int32]int32{1: 10, 2: 2, 3: 3}})
}

func TestGet() {
	player := Player.NewPlayer("htzx", "player", "player@qq@001")
	r := player.GetRaw()
	info := &data.ReputationData{}
	if e := data.GetDataMgr().Unmarshal(r, data.TableReputationKey, info); e == nil {
		fmt.Printf("%+v", info)
	} else {
		fmt.Printf("err:%s", e.Error())
	}
}

func main() {
	TestSave()
	TestGet()

	ret := make(chan bool)
	<-ret
	fmt.Println("donw")
}
