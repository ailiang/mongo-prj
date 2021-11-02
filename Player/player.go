package Player

import "MongoDB-Proj/db"

type Player struct {
	Data *db.SaveData
}

func NewPlayer(key string) *Player {
	return &Player{Data: db.NewSaveData(key)}
}

func (p *Player) Save(field string, value string) {
	p.Data.Add(field, value)
	p.Data.Save()
}

func (p *Player) Load() {

}
