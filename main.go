package main

import (
	"MongoDB-Proj/Player"
	"MongoDB-Proj/game/data"
	"MongoDB-Proj/protoc"
	"fmt"
	"github.com/golang/protobuf/proto"
	"go.mongodb.org/mongo-driver/bson"
	"reflect"
	"strconv"
)

func TestSave(i int, val interface{}) {
	playerIndex := strconv.Itoa(i)
	playerKey := "player@qq@" + playerIndex
	playerName := "name_" + playerIndex
	player := Player.NewPlayer("htzx", "player", playerKey)
	player.Save(data.TablePlayerInfoKey, data.PlayerInfo{
		PlayerId:   int64(i),
		PlayerName: playerName,
	})
	player.Save(data.TableReputationKey, data.ReputationData{map[int32]int32{1: 10, 2: 2, 3: 3}})
	player.Save("binData", val)
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

var (
	tests = []struct {
		name          string
		pb            proto.Message
		equivalentPbs []proto.Message
	}{
		{
			name: "simple message",
			pb: &protoc.SimpleMessage{
				StringField: "foo",
				Int32Field:  32525,
				Int64Field:  1531541553141312315,
				FloatField:  21541.3242,
				DoubleField: 21535215136361617136.543858,
				BoolField:   true,
				EnumField:   protoc.Enum_VAL_2,
			},
			equivalentPbs: []proto.Message{
				&protoc.RepeatedFieldMessage{
					StringField: []string{"foo"},
					Int32Field:  []int32{32525},
					Int64Field:  []int64{1531541553141312315},
					FloatField:  []float32{21541.3242},
					DoubleField: []float64{21535215136361617136.543858},
					BoolField:   []bool{true},
					EnumField:   []protoc.Enum{protoc.Enum_VAL_2},
				},
			},
		},
	}
)

func testMarshalUnmarshal() {
	rb := bson.NewRegistryBuilder()
	rb.RegisterCodec(reflect.TypeOf((*proto.Message)(nil)).Elem(), protoc.NewProtobufCodec())
	reg := rb.Build()

	for _, testCase := range tests {
		b, err := bson.MarshalWithRegistry(reg, testCase.pb)
		if err != nil {
			fmt.Printf("bson.MarshalWithRegistry error = %v \n", err)
		}
		fmt.Printf("%+v", b)

		//TestSave(1, testCase.pb)
		//for _, equivalentPb := range append(testCase.equivalentPbs, testCase.pb) {
		//	out := reflect.New(reflect.TypeOf(equivalentPb).Elem()).Interface().(proto.Message)
		//	if err = bson.UnmarshalWithRegistry(reg, b, &out); err != nil {
		//		fmt.Printf("bson.UnmarshalWithRegistry error = %v\n", err)
		//	}
		//	if !proto.Equal(equivalentPb, out) {
		//		fmt.Printf("failed: in=%#q, out=%#q\n", equivalentPb, out)
		//	}
		//}
	}
}

func main() {

	testMarshalUnmarshal()
	<-make(chan interface{}, 0)
	fmt.Println("down")
}
