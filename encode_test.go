package sproto

import (
	"testing"

	"bytes"
	"encoding/json"
	"reflect"

	"math"

	"github.com/xjdrew/gosproto"
)

// 不影响已有功能测试：Ptr结构可以被编码
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

// 不影响已有功能测试：Ptr结构可以在包含nil的情况下被编码
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

// 新特性测试: 等价的Val结构和Ptr结构编码结果一致（不考虑ptr为nil）
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

// 新特性测试:预期valMsgData可以被相同的val结构体接收
func TestValueDecode(t *testing.T) {
	Reset()
	valMsgData, err := Encode(&valMSG)
	if err != nil {
		t.Error(err, valMsgData)
		return
	}

	// 预期valMsgData可以被相同的val结构体接收
	valMsg2 := ValMSG{}
	_, err = Decode(valMsgData, &valMsg2)
	if err != nil {
		t.Error(err)
		return
	}

	// 预期结果应该保持一致
	if !reflect.DeepEqual(valMSG, valMsg2) {
		t.Error("ValMsg expect equal to ValMsg2")
		return
	}
}

// 新特性测试：预期ptrMsgData可以被等价的Val结构体接收（补充测试）
func TestValueEncodeToPtr(t *testing.T) {
	Reset()
	ptrMsgData, err := Encode(&ptrMsg)
	if err != nil {
		t.Error(err)
		return
	}

	// 预期可以被等价的val结构体接收
	valMsg2 := ValMSG{}
	_, err = Decode(ptrMsgData, &valMsg2)
	if err != nil {
		t.Error(err)
		return
	}

	// 通过比对Json结果间接比对接收结果
	ptrJson, _ := json.Marshal(ptrMsg)
	valJson, _ := json.Marshal(valMsg2)
	if !bytes.Equal(ptrJson, valJson) {
		t.Error("Expect ptrJson Euqal To valJson")
		return
	}
}

// 新特性测试：包含部分等价nil的消息可以被正确接收并设置为默认值
func TestValueDecodeNil(t *testing.T) {
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

	// 测试Val对nil的接收
	valMsg2 := ValMSG{Bool: true, Int: 1024}
	_, err = Decode(ptrMsgData, &valMsg2)
	if err != nil {
		t.Error(err)
		return
	}
	// 测试默认值
	if valMsg2.Bool != false {
		t.Error("Expect Bool equal to false")
		return
	}
	if valMsg2.Int != 0 {
		t.Error("Expect Int equal to 0")
		return
	}
}

// 与原装的sproto进行比对，看编码结果是否一致
func TestOldAndNewEncode(t *testing.T) {
	Reset()
	// 云风的版本
	oldMsgData, err := sproto.Encode(&ptrMsg)
	if err != nil {
		t.Error(err)
	}
	// 新版本
	newMsgData, err := Encode(&ptrMsg)
	if err != nil {
		t.Error(err)
	}
	if !reflect.DeepEqual(oldMsgData, newMsgData) {
		t.Error("新方案与旧方案的编码结果不一致。")
	}

	Reset()
	// 云风的版本
	ptrMsg.Int = nil
	ptrMsg.String = nil
	oldMsgData, err = sproto.Encode(&ptrMsg)
	if err != nil {
		t.Error(err)
	}
	// 新版本
	newMsgData, err = Encode(&ptrMsg)
	if err != nil {
		t.Error(err)
	}
	if !reflect.DeepEqual(oldMsgData, newMsgData) {
		t.Error("包含nil的情况下新方案与旧方案的编码结果不一致。")
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
	IntNeg      *int          `sproto:"integer,9,name=IntNeg"`
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
		Int:         Int(math.MaxInt64),
		IntNeg:      Int(math.MinInt32 + 1),
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
		Int:         math.MaxInt64,
		IntNeg:      math.MinInt32 + 1,
		String:      "Hello",
		Bool:        true,
		ByteSlice:   []byte("World"),
		BoolSlice:   []bool{true, true, false, true, false},
		IntSlice:    []int{123, 321, 1234567},
		StringSlice: []string{"FOO", "BAR"},
		Struct: &HoldValMSG{
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
	Struct      *HoldValMSG   `sproto:"struct,3,name=Struct"`
	ByteSlice   []byte        `sproto:"string,4,name=ByteSlice"`
	BoolSlice   []bool        `sproto:"boolean,5,array,name=BoolSlice"`
	IntSlice    []int         `sproto:"integer,6,array,name=IntSlice"`
	StringSlice []string      `sproto:"string,7,array,name=StringSlice"`
	StructSlice []*HoldValMSG `sproto:"struct,8,array,name=StructSlice"`
	IntNeg      int           `sproto:"integer,9,name=IntNeg"`
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
