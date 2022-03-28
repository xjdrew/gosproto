package sproto

import (
	"reflect"
	"sort"
	"testing"
)

/*
.NestData {
	A 1 : string
	B 3 : boolean
	C 5 : integer
	D 6 : double
}

.SimpleMapItem {
	K 1 : integer
	V 2 : string
}

.StructMapItem {
	Key 3 : integer
	Value 5 : NestData
}

.MapMsg {
	SimpleMap 0 : *SimpleMapItem()
	StructMap 1 : *StructMapItem()
	MainIndexMap 2 : *NestData(A)
}
*/
type NestData struct {
	A string  `sproto:"string,1"`
	B bool    `sproto:"boolean,3"`
	C int     `sproto:"integer,5"`
	D float64 `sproto:"double,6"`
}

type SimpleMapItem struct {
	K int    `sproto:"integer,1"`
	V string `sproto:"string,2"`
}

type SimpleMapPtrItem struct {
	K *int    `sproto:"integer,1"`
	V *string `sproto:"string,2"`
}

type StructMapItem struct {
	Key   int       `sproto:"integer,3"`
	Value *NestData `sproto:"struct,5"`
}

type MapMsg struct {
	SimpleMap        map[int]string       `sproto:"struct,0,array,key=1,value=2,subtype=simpleMapItem"`
	SimpleMapPtr     map[int]string       `sproto:"struct,1,array,key=1,value=2,subtype=simpleMapPtrItem"` // 支持item中定义为指针类型，但map定义为值类型
	StructMap        map[int]*NestData    `sproto:"struct,2,array,key=3,value=5,subtype=structMapItem"`
	MainIndexMap     map[string]*NestData `sproto:"struct,3,array,key=1"` // value=all
	simpleMapItem    *SimpleMapItem
	simpleMapPtrItem *SimpleMapPtrItem
	structMapItem    *StructMapItem
}

type ArrayMsg struct {
	SimpleMap    []*SimpleMapItem `sproto:"struct,0,array"`
	StructMap    []*StructMapItem `sproto:"struct,1,array"`
	MainIndexMap []*NestData      `sproto:"struct,2,array"`
}

var (
	mapMsg = MapMsg{
		SimpleMap: map[int]string{
			1: "v1",
			2: "v2",
		},
		SimpleMapPtr: map[int]string{
			1: "v1",
			2: "v2",
		},
		StructMap: map[int]*NestData{
			11: {
				A: "11va",
				B: true,
				C: 123,
				D: 0.123,
			},
			12: {
				A: "12va",
				B: false,
				C: 456,
				D: 0.456,
			},
		},
		MainIndexMap: map[string]*NestData{
			"11va": {
				A: "11va",
				B: true,
				C: 123,
				D: 0.123,
			},
			"12va": {
				A: "12va",
				B: false,
				C: 456,
				D: 0.456,
			},
		},
	}
	arrayMsg = ArrayMsg{
		SimpleMap: []*SimpleMapItem{
			{K: 1, V: "v1"},
			{K: 2, V: "v2"},
		},
		StructMap: []*StructMapItem{
			{
				Key: 11,
				Value: &NestData{
					A: "11va",
					B: true,
					C: 123,
					D: 0.123,
				},
			},
			{
				Key: 12,
				Value: &NestData{
					A: "12va",
					B: false,
					C: 456,
					D: 0.456,
				},
			},
		},
		MainIndexMap: []*NestData{
			{
				A: "11va",
				B: true,
				C: 123,
				D: 0.123,
			},
			{
				A: "12va",
				B: false,
				C: 456,
				D: 0.456,
			},
		},
	}
)

func TestMapEncode(t *testing.T) {
	mapMsgData, err := Encode(&mapMsg)
	if err != nil {
		t.Error(err)
		return
	}

	mapMsg2 := MapMsg{}
	used, err := Decode(mapMsgData, &mapMsg2)
	if err != nil {
		t.Error(err)
		return
	}
	if used != len(mapMsgData) {
		t.Error("decode failed: unexpected used")
		return
	}

	if !reflect.DeepEqual(mapMsg, mapMsg2) {
		t.Error("mapMsg is not equal to mapMsg2")
		return
	}
}

func TestArrayToMap(t *testing.T) {
	arrayMsgData, err := Encode(&arrayMsg)
	if err != nil {
		t.Error(err)
		return
	}

	mapMsg2 := MapMsg{}
	used, err := Decode(arrayMsgData, &mapMsg2)
	if err != nil {
		t.Error(err)
		return
	}
	if used != len(arrayMsgData) {
		t.Error("decode failed: unexpected used")
		return
	}

	if !reflect.DeepEqual(mapMsg, mapMsg2) {
		t.Error("mapMsg is not equal to mapMsg2")
		return
	}
}

func TestMapToArray(t *testing.T) {
	mapMsgData, err := Encode(&mapMsg)
	if err != nil {
		t.Error(err)
		return
	}

	arrayMsg2 := ArrayMsg{}
	used, err := Decode(mapMsgData, &arrayMsg2)
	if err != nil {
		t.Error(err)
		return
	}
	if used != len(mapMsgData) {
		t.Error("decode failed: unexpected used")
		return
	}

	// 数组排序
	sort.Slice(arrayMsg2.SimpleMap, func(i, j int) bool {
		return arrayMsg2.SimpleMap[i].K < arrayMsg2.SimpleMap[j].K
	})

	sort.Slice(arrayMsg2.StructMap, func(i, j int) bool {
		return arrayMsg2.StructMap[i].Key < arrayMsg2.StructMap[j].Key
	})

	sort.Slice(arrayMsg2.MainIndexMap, func(i, j int) bool {
		return arrayMsg2.MainIndexMap[i].A < arrayMsg2.MainIndexMap[j].A
	})

	if !reflect.DeepEqual(arrayMsg, arrayMsg2) {
		t.Error("arrayMsg is not equal to arrayMsg2")
		return
	}
}
