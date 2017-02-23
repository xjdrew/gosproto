package sproto

import (
	"testing"

	"reflect"
)

func TestPtrEncode(t *testing.T) {
	ptrMsgData, err := Encode(&ptrMsg)
	if err != nil {
		t.Error(err)
		return
	}

	ptrMsg2 := PtrMSG{}
	Decode(ptrMsgData, &ptrMsg2)
	if !reflect.DeepEqual(ptrMsg, ptrMsg2) {
		t.Error("ptrMsg is not equal to ptrMsg2")
	}
}

func TestValueEncode(t *testing.T) {
	msgData, err := Encode(&valMSG)
	if err != nil {
		t.Error(err)
		return
	}

	valMsg2 := ValMSG{}
	Decode(msgData, &valMsg2)
	if !reflect.DeepEqual(valMSG, valMsg2) {
		t.Error("valMSG is not equal to valMsg2")
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

var ptrMsg = PtrMSG{
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

type ValMSG struct {
	Int         int           `sproto:"integer,0,name=Int"`
	String      string        `sproto:"string,1,name=String"`
	Bool        bool          `sproto:"boolean,2,name=Bool"`
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

var valMSG = ValMSG{
	Int:    1,
	String: "Hello",
	Bool:   true,
	StructSlice: []*HoldValMSG{
		&HoldValMSG{
			Int:         2,
			String:      "Foo",
			Bool:        true,
			ByteSlice:   []byte("World"),
			BoolSlice:   []bool{true, true, false, true, false},
			IntSlice:    []int{123, 321, 1234567},
			StringSlice: []string{"FOO", "BAR"},
		},
		&HoldValMSG{
			Int:         3,
			String:      "Foo2",
			Bool:        true,
			ByteSlice:   []byte("World"),
			BoolSlice:   []bool{true, true, false, true, false},
			IntSlice:    []int{123, 321, 1234567},
			StringSlice: []string{"FOO", "BAR"},
		},
	},
}
