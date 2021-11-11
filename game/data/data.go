package data

import (
	"sync"
)

const (
	TablePlayerInfoKey = "PlayerInfoKey"
	TableReputationKey = "PlayerReputationKey"
)

var once sync.Once
var inst *DataMgr

type DataMgr struct {
	DataMap map[string]interface{}
}

func GetDataMgr() *DataMgr {
	once.Do(func() {
		inst = &DataMgr{
			DataMap: map[string]interface{}{
				TablePlayerInfoKey: &PlayerInfo{},
				TableReputationKey: &ReputationData{},
			},
		}
	})
	return inst
}

type PlayerInfo struct {
	PlayerId   int64
	PlayerName string
}

type ReputationData struct {
	RepuMap map[int32]int32
}
