package main

import (
	"MongoDB-Proj/Player"
	"MongoDB-Proj/game/data"
	"fmt"
	"strconv"
)

func TestSave(i int) {
	playerIndex := strconv.Itoa(i)
	playerKey := "player@qq@" + playerIndex
	playerName := "name_" + playerIndex
	player := Player.NewPlayer("htzx", "player", playerKey)
	player.Save(data.TablePlayerInfoKey, data.PlayerInfo{
		PlayerId:   int64(i),
		PlayerName: playerName,
	})
	player.Save(data.TableReputationKey, data.ReputationData{map[int32]int32{1: 10, 2: 2, 3: 3}})
}

func TestGetRaw(i int) {
	playerIndex := strconv.Itoa(i)
	playerKey := "player@qq@" + playerIndex
	player := Player.NewPlayer("htzx", "player", playerKey)
	if r, err := player.Data.GetRaw(); err == nil {
		repuInfo := &data.ReputationData{}
		if e := r.Lookup(data.TableReputationKey).Unmarshal(repuInfo); e == nil {
			fmt.Printf("TestGetRaw: %+v\n", repuInfo)
		} else {
			fmt.Printf("TestGetRaw: err:%s\n", e.Error())
		}
		playerInfo := &data.PlayerInfo{}
		if e := r.Lookup(data.TablePlayerInfoKey).Unmarshal(playerInfo); e == nil {
			fmt.Printf("TestGetRaw: %+v\n", playerInfo)
		} else {
			fmt.Printf("TestGetRaw: err:%s\n", e.Error())
		}
	} else {
		fmt.Printf("TestGetRaw: err:%s\n", err.Error())
	}
}

func TestGetFieldRaw(i int) {
	playerIndex := strconv.Itoa(i)
	playerKey := "player@qq@" + playerIndex
	player := Player.NewPlayer("htzx", "player", playerKey)
	info := &data.PlayerInfo{}
	if err := player.Data.GetField(data.TablePlayerInfoKey, info); err == nil {
		fmt.Printf("TestGetFieldRaw: %+v\n", info)
	} else {
		fmt.Printf("TestGetFieldRaw: err:%s", err.Error())
	}
}

func TestDel(i int) {
	playerIndex := strconv.Itoa(i)
	playerKey := "player@qq@" + playerIndex
	player := Player.NewPlayer("htzx", "player", playerKey)
	if c, err := player.Data.Delete(); err == nil {
		fmt.Printf("TestDel: %d\n", c)
	} else {
		fmt.Printf("TestDel: err:%s", err.Error())
	}
}

func main() {
	for i := 0; i < 1; i++ {
		TestSave(i)
		TestGetRaw(i)
		TestDel(i)
		TestGetRaw(i)
	}
	//go TestSave()
	//go TestGetRaw()
	//go TestGetFieldRaw()

	ret := make(chan bool)
	<-ret
	fmt.Println("donw")
}
