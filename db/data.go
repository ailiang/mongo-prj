package db

import (
	"container/list"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"hash/crc32"
	"sync"
)

const DataSaveKey = "DataKey"

type SaveData struct {
	Key    string
	Fields []string
	Values []interface{}

	DelFields []string

	DB         string
	COLLECTION string
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
	for {
		data := g.popFront()
		if data != nil {
			g.save(data)
		}
	}
}
func (g *SaveGoroutine) save(data *SaveData) error {
	cli := GetDbManager().GetClient()
	if cli == nil {
		panic("cli nil")
	}
	col := cli.Database(data.DB).Collection(data.COLLECTION)
	if col == nil {
		panic("col nil")
	}
	return MongoUpdateOneWithColl(col, data.Key, data.Fields, data.Values)
}

func (g *SaveGoroutine) popFront() *SaveData {
	g.cond.L.Lock()
	if g.dataList.Len() == 0 {
		g.cond.Wait()
	}
	front := g.dataList.Front()
	g.dataList.Remove(front)
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
	if g.dataList.Len() == 1 {
		g.cond.Signal()
	}
	g.cond.L.Unlock()
}

func NewSaveData(db string, collection string, key string) *SaveData {
	return &SaveData{
		Key:        key,
		Fields:     nil,
		Values:     nil,
		DelFields:  nil,
		DB:         db,
		COLLECTION: collection,
	}
}

func (d *SaveData) Add(filed string, value interface{}) error {
	d.Fields = append(d.Fields, filed)
	d.Values = append(d.Values, value)
	return nil
}

func (d *SaveData) Save() {
	m := GetSaveManager()
	id := int(crc32.ChecksumIEEE([]byte(d.Key))) % m.gNum
	m.goroutines[id].pushFront(d)
}

func (d *SaveData) Get() bson.D {
	cli := GetDbManager().GetClient()
	if cli == nil {
		panic("cli nil")
	}
	col := cli.Database(d.DB).Collection(d.COLLECTION)
	if col == nil {
		panic("col nil")
	}
	r, err := MongoGetOneWithColl(col, d.Key)
	if err != nil {
		fmt.Println("get err:", err)
		return nil
	} else {
		fmt.Printf("get result: %+v", r)
		return r
	}
}

func (d *SaveData) GetRaw() (bson.Raw, error) {
	cli := GetDbManager().GetClient()
	if cli == nil {
		panic("cli nil")
	}
	col := cli.Database(d.DB).Collection(d.COLLECTION)
	if col == nil {
		panic("col nil")
	}
	return MongoGetOneRawWithColl(col, d.Key)
}

func (d *SaveData) GetField(fieldName string, ret interface{}) error {
	cli := GetDbManager().GetClient()
	if cli == nil {
		panic("cli nil")
	}
	col := cli.Database(d.DB).Collection(d.COLLECTION)
	if col == nil {
		panic("col nil")
	}
	return MongoGetOneFiledRawWithColl(col, d.Key, fieldName, ret)
}

func (d *SaveData) Delete() (int64, error) {
	cli := GetDbManager().GetClient()
	if cli == nil {
		panic("cli nil")
	}
	col := cli.Database(d.DB).Collection(d.COLLECTION)
	if col == nil {
		panic("col nil")
	}
	return MongoDelOneWithColl(col, d.Key)
}
