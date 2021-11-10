package data

import (
	"go.mongodb.org/mongo-driver/bson"
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

func (d *DataMgr) Unmarshal(r bson.Raw, key string, v interface{}) error {
	return r.Lookup(key).Unmarshal(v)
}

type PlayerInfo struct {
	PlayerId   int64
	PlayerName string
}

type ReputationData struct {
	RepuMap map[int32]int32
}
