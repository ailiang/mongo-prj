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

func TestGetRaw() {
	player := Player.NewPlayer("htzx", "player", "player@qq@001")
	if r, err := player.Data.GetRaw(); err == nil {
		info := &data.ReputationData{}
		if e := r.Lookup(data.TableReputationKey).Unmarshal(info); e == nil {
			fmt.Printf("TestGetRaw: %+v\n", info)
		} else {
			fmt.Printf("TestGetRaw: err:%s", e.Error())
		}
	}
}

func TestGetFieldRaw() {
	player := Player.NewPlayer("htzx", "player", "player@qq@001")
	info := &data.PlayerInfo{}
	if err := player.Data.GetField(data.TablePlayerInfoKey, info); err == nil {
		fmt.Printf("TestGetFieldRaw: %+v\n", info)
	} else {
		fmt.Printf("TestGetFieldRaw: err:%s", err.Error())
	}
}

func main() {
	TestSave()
	TestGetRaw()
	TestGetFieldRaw()

	ret := make(chan bool)
	<-ret
	fmt.Println("donw")
}
