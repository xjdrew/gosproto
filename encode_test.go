package sproto

import (
	"testing"

	"go.szyhf.org/digo/log"

	"reflect"
)

func TestPtrEncode(t *testing.T) {
	Reset()
	ptrMsgData, err := Encode(&ptrMsg)
	if err != nil {
		t.Error(err)
		return
	}
	// 测试解包结果
	Reset()
	ptrMsg2 := PtrMSG{}
	Decode(ptrMsgData, &ptrMsg2)
	if !reflect.DeepEqual(ptrMsg, ptrMsg2) {
		t.Error("ptrMsg is not equal to ptrMsg2")
	}
}

func TestPtrNilEncode(t *testing.T) {
	// 测试对nil值的支持
	Reset()
	ptrMsg.Int = nil
	ptrMsg.Bool = nil
	ptrMsg.StructSlice = nil
	ptrMsg.Struct = nil
	ptrMsgData, err := Encode(&ptrMsg)
	if err != nil {
		t.Error(err)
		return
	}

	// 测试解包结果
	ptrMsg2 := PtrMSG{}
	Decode(ptrMsgData, &ptrMsg2)
	if !reflect.DeepEqual(ptrMsg, ptrMsg2) {
		t.Error("ptrMsg is not equal to ptrMsg2")
	}
}

func TestValueEncodeEqualToPtr(t *testing.T) {
	Reset()
	valMsgData, err := Encode(&valMSG)
	if err != nil {
		t.Error(err, valMsgData)
		return
	}

	ptrMsgData, err := Encode(&ptrMsg)
	if err != nil {
		t.Error(err)
		return
	}
	if !reflect.DeepEqual(valMsgData, ptrMsgData) {
		t.Error("ValMsgData exprect equal to PtrMsgData")
		log.Error(valMsgData)
		log.Error(ptrMsgData)
		return
	}

	Reset()
	// 预期val编码结果应该允许被等价结构的含ptr结构体接收
	ptrMsg2 := PtrMSG{}
	Decode(valMsgData, &ptrMsg2)
	if !reflect.DeepEqual(ptrMsg2, ptrMsg) {
		t.Error("预期val编码结果应该允许被等价结构的含ptr结构体接收")
	}
}

type PtrMSG struct {
	Int         *int          `sproto:"integer,0,name=Int"`
	String      *string       `sproto:"string,1,name=String"`
	Bool        *bool         `sproto:"boolean,2,name=Bool"`
	Struct      *HoldPtrMSG   `sproto:"struct,3,name=Struct"`
	ByteSlice   []byte        `sproto:"string,4,name=ByteSlice"`
	BoolSlice   []bool        `sproto:"boolean,5,array,name=BoolSlice"`
	IntSlice    []int         `sproto:"integer,6,array,name=IntSlice"`
	StringSlice []string      `sproto:"string,7,array,name=StringSlice"`
	StructSlice []*HoldPtrMSG `sproto:"struct,8,array,name=StructSlice"`
}

type HoldPtrMSG struct {
	Int         *int     `sproto:"integer,0,name=Int"`
	String      *string  `sproto:"string,1,name=String"`
	Bool        *bool    `sproto:"boolean,2,name=Bool"`
	ByteSlice   []byte   `sproto:"string,3,name=ByteSlice"`
	BoolSlice   []bool   `sproto:"boolean,4,array,name=BoolSlice"`
	IntSlice    []int    `sproto:"integer,5,array,name=IntSlice"`
	StringSlice []string `sproto:"string,6,array,name=StringSlice"`
}

var ptrMsg PtrMSG

func Reset() {
	ptrMsg = PtrMSG{
		Int:         Int(1),
		String:      String("Hello"),
		Bool:        Bool(true),
		ByteSlice:   []byte("World"),
		BoolSlice:   []bool{true, true, false, true, false},
		IntSlice:    []int{123, 321, 1234567},
		StringSlice: []string{"FOO", "BAR"},
		Struct: &HoldPtrMSG{
			Int:         Int(1),
			String:      String("Hello"),
			Bool:        Bool(true),
			ByteSlice:   []byte("World"),
			BoolSlice:   []bool{true, true, false, true, false},
			IntSlice:    []int{123, 321, 1234567},
			StringSlice: []string{"FOO", "BAR"},
		},
		StructSlice: []*HoldPtrMSG{
			&HoldPtrMSG{
				Int:         Int(2),
				String:      String("Hello2"),
				Bool:        Bool(true),
				ByteSlice:   []byte("World2"),
				BoolSlice:   []bool{true, true, false, true, false},
				IntSlice:    []int{123, 321, 1234567},
				StringSlice: []string{"FOO2", "BAR2"},
			},
			&HoldPtrMSG{
				Int:         Int(3),
				String:      String("Hello3"),
				Bool:        Bool(true),
				ByteSlice:   []byte("World3"),
				BoolSlice:   []bool{true, true, false, true, false},
				IntSlice:    []int{123, 321, 1234567},
				StringSlice: []string{"FOO3", "BAR3"},
			},
		},
	}
	valMSG = ValMSG{
		Int:         1,
		String:      "Hello",
		Bool:        true,
		ByteSlice:   []byte("World"),
		BoolSlice:   []bool{true, true, false, true, false},
		IntSlice:    []int{123, 321, 1234567},
		StringSlice: []string{"FOO", "BAR"},
		Struct: HoldValMSG{
			Int:         1,
			String:      "Hello",
			Bool:        true,
			ByteSlice:   []byte("World"),
			BoolSlice:   []bool{true, true, false, true, false},
			IntSlice:    []int{123, 321, 1234567},
			StringSlice: []string{"FOO", "BAR"},
		},
		StructSlice: []*HoldValMSG{
			&HoldValMSG{
				Int:         2,
				String:      "Hello2",
				Bool:        true,
				ByteSlice:   []byte("World2"),
				BoolSlice:   []bool{true, true, false, true, false},
				IntSlice:    []int{123, 321, 1234567},
				StringSlice: []string{"FOO2", "BAR2"},
			},
			&HoldValMSG{
				Int:         3,
				String:      "Hello3",
				Bool:        true,
				ByteSlice:   []byte("World3"),
				BoolSlice:   []bool{true, true, false, true, false},
				IntSlice:    []int{123, 321, 1234567},
				StringSlice: []string{"FOO3", "BAR3"},
			},
		},
	}
}

type ValMSG struct {
	Int         int           `sproto:"integer,0,name=Int"`
	String      string        `sproto:"string,1,name=String"`
	Bool        bool          `sproto:"boolean,2,name=Bool"`
	Struct      HoldValMSG    `sproto:"struct,3,name=Struct"`
	ByteSlice   []byte        `sproto:"string,4,name=ByteSlice"`
	BoolSlice   []bool        `sproto:"boolean,5,array,name=BoolSlice"`
	IntSlice    []int         `sproto:"integer,6,array,name=IntSlice"`
	StringSlice []string      `sproto:"string,7,array,name=StringSlice"`
	StructSlice []*HoldValMSG `sproto:"struct,8,array,name=StructSlice"`
}

type HoldValMSG struct {
	Int         int      `sproto:"integer,0,name=Int"`
	String      string   `sproto:"string,1,name=String"`
	Bool        bool     `sproto:"boolean,2,name=Bool"`
	ByteSlice   []byte   `sproto:"string,3,name=ByteSlice"`
	BoolSlice   []bool   `sproto:"boolean,4,array,name=BoolSlice"`
	IntSlice    []int    `sproto:"integer,5,array,name=IntSlice"`
	StringSlice []string `sproto:"string,6,array,name=StringSlice"`
}

var valMSG ValMSG
