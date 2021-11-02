package db

import (
	"container/list"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"hash/crc32"
	"log"
	"sync"
)

const (
	DB_NAME    = "GAME"
	TABLE_NAME = "PTB"
)

type SaveData struct {
	Key    string
	Fields []string
	Values []string

	DelFields []string
}

type SaveGoroutine struct {
	dataList   *list.List
	dataLocker sync.Mutex
	cond       *sync.Cond
}

type SaveManager struct {
	goroutines []*SaveGoroutine
	gNum       int
}

var saveInstance *SaveManager
var saveInstanceOnce sync.Once

func GetSaveManager() *SaveManager {
	saveInstanceOnce.Do(func() {
		saveInstance = &SaveManager{}
		saveInstance.gNum = 10
		saveInstance.goroutines = make([]*SaveGoroutine, 0, saveInstance.gNum)
		for i := 0; i < saveInstance.gNum; i++ {
			g := &SaveGoroutine{}
			g.dataList = list.New()
			g.cond = sync.NewCond(&g.dataLocker)
			saveInstance.goroutines = append(saveInstance.goroutines, g)
		}
	})
	return saveInstance
}
func (sm *SaveManager) Init() {
	for _, goroutine := range sm.goroutines {
		go goroutine.run()
	}
}

func init() {
	GetSaveManager().Init()
}
func (g *SaveGoroutine) run() {
	log.Printf("goroutine run")
	for {
		data := g.popFront()
		if data != nil {
			g.save(data)
		}
	}
}
func (g *SaveGoroutine) save(data *SaveData) error {
	cli := GetDbManager().Client
	col := cli.Database(DB_NAME).Collection(TABLE_NAME)
	filter := bson.D{{"key", data.Key}}

	update := bson.D{{"key", data.Key}}
	col.UpdateOne(context.TODO(), filter, update)
	return nil
}

func (g *SaveGoroutine) popFront() *SaveData {
	g.cond.L.Lock()
	if g.dataList.Len() == 0 {
		g.cond.Wait()
	}
	front := g.dataList.Front()
	g.cond.L.Unlock()
	if front == nil {
		return nil
	}
	data, ok := front.Value.(*SaveData)
	if ok {
		return data
	}
	return nil
}

func (g *SaveGoroutine) pushFront(data *SaveData) {
	g.cond.L.Lock()
	g.dataList.PushFront(data)
	if g.dataList.Len() == 0 {
		g.cond.Signal()
	}
	g.cond.L.Unlock()
}

func NewSaveData(key string) *SaveData {
	return &SaveData{
		Key:       key,
		Fields:    nil,
		Values:    nil,
		DelFields: nil,
	}
}

func (d *SaveData) Add(filed string, value string) error {
	d.Fields = append(d.Fields, filed)
	d.Values = append(d.Values, value)
	return nil
}

func (d *SaveData) Save() {
	m := GetSaveManager()
	id := int(crc32.ChecksumIEEE([]byte(d.Key))) % m.gNum
	m.goroutines[id].pushFront(d)
}
