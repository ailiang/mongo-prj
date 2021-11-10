package Player

import (
	"MongoDB-Proj/db"
	"go.mongodb.org/mongo-driver/bson"
)

type Player struct {
	Data *db.SaveData
}

func NewPlayer(dbName string, collection string, key string) *Player {
	return &Player{Data: db.NewSaveData(dbName, collection, key)}
}

func (p *Player) Save(field string, value interface{}) {
	p.Data.Add(field, value)
	p.Data.Save()
}

func (p *Player) Get() bson.D {
	return p.Data.Get()
}

func (p *Player) GetRaw() bson.Raw {
	return p.Data.GetRaw()
}
