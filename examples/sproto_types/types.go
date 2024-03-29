// Code generated by sprotodump
// source: types.sproto
// DO NOT EDIT!

/*
   Package sproto_types is a generated sproto package.
*/
package sproto_types

import (
	"reflect"

	"github.com/xjdrew/gosproto"
)

// avoids "imported but not used"
var _ reflect.Type

type Person struct {
	Name   *string              `sproto:"string,0"`
	Id     *int64               `sproto:"integer,1"`
	Email  *string              `sproto:"string,2"`
	Phone  []*PersonPhoneNumber `sproto:"struct,3,array"`
	Height *int64               `sproto:"integer,4"`
	Data   []byte               `sproto:"binary,5"`
	Weight *float64             `sproto:"double,6"`
	Pics   [][]byte             `sproto:"binary,7,array"`
}

type PersonPhoneNumber struct {
	Number *string `sproto:"string,0"`
	Type   *int64  `sproto:"integer,1"`
}

type CreditCard struct {
	CardNum *string `sproto:"string,0"`
	Owner   *Person `sproto:"struct,1"`
}

type Bank struct {
	Cards        map[string]*Person `sproto:"struct,0,array,key=0,value=1,subtype=mapItemCards"`
	Clients      map[int64]*Person  `sproto:"struct,1,array,key=1"`
	mapItemCards *CreditCard
}

type SimpleItem struct {
	Key   *int64  `sproto:"integer,3"`
	Value *string `sproto:"string,5"`
}

type NodeItem struct {
	Id   *int64    `sproto:"integer,0"`
	Node *NodeItem `sproto:"struct,1"`
}

type ArraysStruct struct {
	IntArr    []int64       `sproto:"integer,1,array"`
	BoolArr   []bool        `sproto:"boolean,2,array"`
	StrArr    []string      `sproto:"string,3,array"`
	BinArr    [][]byte      `sproto:"binary,4,array"`
	DoubleArr []float64     `sproto:"double,5,array"`
	StructArr []*SimpleItem `sproto:"struct,6,array"`
}

type NestedMapItem struct {
	Id            *int64            `sproto:"integer,0"`
	Nested        map[int64]*string `sproto:"struct,1,array,key=3,value=5,subtype=mapItemNested"`
	mapItemNested *SimpleItem
}

type NestedArrayItem struct {
	Id     *int64        `sproto:"integer,0"`
	Nested []*SimpleItem `sproto:"struct,1,array"`
}

type MapStruct struct {
	Map1        map[int64]*SimpleItem `sproto:"struct,9,array,key=3"`
	Map2        map[int64]*string     `sproto:"struct,10,array,key=3,value=5,subtype=mapItemMap2"`
	Map3        map[int64]*NodeItem   `sproto:"struct,20,array,key=0,value=1,subtype=mapItemMap3"`
	Map4        map[int64]*NodeItem   `sproto:"struct,21,array,key=0"`
	mapItemMap2 *SimpleItem
	mapItemMap3 *NodeItem
}

type NestedMapStruct struct {
	NestedMap1        map[int64]map[int64]*string `sproto:"struct,30,array,key=0,value=1,subtype=mapItemNestedMap1"`
	NestedMap2        map[int64]*NestedMapItem    `sproto:"struct,31,array,key=0"`
	NestedArr         map[int64][]*SimpleItem     `sproto:"struct,40,array,key=0,value=1,subtype=mapItemNestedArr"`
	mapItemNestedMap1 *NestedMapItem
	mapItemNestedArr  *NestedArrayItem
}

type ApiRequest struct {
	Ping *string `sproto:"string,0"`
}

type ApiResponse struct {
	Pong *string `sproto:"string,0"`
}

var Name string = "types"
var Protocols []*sproto.Protocol = []*sproto.Protocol{
	&sproto.Protocol{
		Type:       1,
		Name:       "types.api",
		MethodName: "Types.Api",
		Request:    reflect.TypeOf(&ApiRequest{}),
		Response:   reflect.TypeOf(&ApiResponse{}),
	},
}
