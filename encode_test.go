package sproto

import (
	"bytes"
	"encoding/json"
	"math"
	"reflect"
	"testing"
)

// 不影响已有功能测试：Ptr结构可以被编码
func TestPtrEncode(t *testing.T) {
	resetEncodeTestEnv()
	ptrMsgData, err := Encode(&ptrMsg)
	if err != nil {
		t.Error(err)
		return
	}
	// 测试解包结果
	resetEncodeTestEnv()
	ptrMsg2 := PtrMSG{}
	Decode(ptrMsgData, &ptrMsg2)
	if !reflect.DeepEqual(ptrMsg, ptrMsg2) {
		t.Error("ptrMsg is not equal to ptrMsg2")
	}
}

// 不影响已有功能测试：Ptr结构可以在包含nil的情况下被编码
func TestPtrNilEncode(t *testing.T) {
	// 测试对nil值的支持
	resetEncodeTestEnv()
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
	resetEncodeTestEnv()
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
		t.Error("ValMsgData expect to equal with PtrMsgData")
		return
	}

	resetEncodeTestEnv()
	// 预期val编码结果应该允许被等价结构的含ptr结构体接收
	ptrMsg2 := PtrMSG{}
	Decode(valMsgData, &ptrMsg2)
	if !reflect.DeepEqual(ptrMsg2, ptrMsg) {
		t.Error("预期val编码结果应该允许被等价结构的含ptr结构体接收")
	}
}

// 新特性测试:预期valMsgData可以被相同的val结构体接收
func TestValueDecode(t *testing.T) {
	resetEncodeTestEnv()
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
	resetEncodeTestEnv()
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
	resetEncodeTestEnv()
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
	resetEncodeTestEnv()
	// 云风的版本
	oldMsgData, err := Encode(&ptrMsg)
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

	resetEncodeTestEnv()
	// 云风的版本
	ptrMsg.Int = nil
	ptrMsg.String = nil
	oldMsgData, err = Encode(&ptrMsg)
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
	Int         *int      `sproto:"integer,0"`
	IntNeg      *int      `sproto:"integer,1"`
	String      *string   `sproto:"string,2"`
	Bool        *bool     `sproto:"boolean,3"`
	Double      *float64  `sproto:"double,4"`
	Binary      []byte    `sproto:"binary,5"`
	ByteSlice   []byte    `sproto:"string,6"`
	BoolSlice   []bool    `sproto:"boolean,7,array"`
	IntSlice    []int     `sproto:"integer,8,array"`
	DoubleSlice []float64 `sproto:"double,9,array"`
	StringSlice []string  `sproto:"string,10,array"`

	Struct      *HoldPtrMSG   `sproto:"struct,20"`
	StructSlice []*HoldPtrMSG `sproto:"struct,21,array"`
}

type HoldPtrMSG struct {
	Int         *int      `sproto:"integer,0"`
	IntNeg      *int      `sproto:"integer,1"`
	String      *string   `sproto:"string,2"`
	Bool        *bool     `sproto:"boolean,3"`
	Double      *float64  `sproto:"double,4"`
	Binary      []byte    `sproto:"binary,5"`
	ByteSlice   []byte    `sproto:"string,6"`
	BoolSlice   []bool    `sproto:"boolean,7,array"`
	IntSlice    []int     `sproto:"integer,8,array"`
	DoubleSlice []float64 `sproto:"double,9,array"`
	StringSlice []string  `sproto:"string,10,array"`
}

type ValMSG struct {
	Int         int       `sproto:"integer,0"`
	IntNeg      int       `sproto:"integer,1"`
	String      string    `sproto:"string,2"`
	Bool        bool      `sproto:"boolean,3"`
	Double      float64   `sproto:"double,4"`
	Binary      []byte    `sproto:"binary,5"`
	ByteSlice   []byte    `sproto:"string,6"`
	BoolSlice   []bool    `sproto:"boolean,7,array"`
	IntSlice    []int     `sproto:"integer,8,array"`
	DoubleSlice []float64 `sproto:"double,9,array"`
	StringSlice []string  `sproto:"string,10,array"`

	Struct      *HoldValMSG   `sproto:"struct,20"`
	StructSlice []*HoldValMSG `sproto:"struct,21,array"`
}

type HoldValMSG struct {
	Int         int       `sproto:"integer,0"`
	IntNeg      int       `sproto:"integer,1"`
	String      string    `sproto:"string,2"`
	Bool        bool      `sproto:"boolean,3"`
	Double      float64   `sproto:"double,4"`
	Binary      []byte    `sproto:"binary,5"`
	ByteSlice   []byte    `sproto:"string,6"`
	BoolSlice   []bool    `sproto:"boolean,7,array"`
	IntSlice    []int     `sproto:"integer,8,array"`
	DoubleSlice []float64 `sproto:"double,9,array"`
	StringSlice []string  `sproto:"string,10,array"`
}

var ptrMsg PtrMSG
var valMSG ValMSG

func resetEncodeTestEnv() {
	ptrMsg = PtrMSG{
		Int:         Int(math.MaxInt64),
		IntNeg:      Int(math.MinInt32 + 1),
		Double:      Double(math.MaxFloat64),
		String:      String("Hello"),
		Bool:        Bool(true),
		Binary:      []byte("binary"),
		ByteSlice:   []byte("World"),
		BoolSlice:   []bool{true, true, false, true, false},
		IntSlice:    []int{123, 321, 1234567},
		DoubleSlice: []float64{0.123, 0.321, 0.1234567},
		StringSlice: []string{"FOO", "BAR"},
		Struct: &HoldPtrMSG{
			Int:         Int(1),
			IntNeg:      Int(math.MinInt32 + 1),
			Double:      Double(0.1),
			String:      String("Hello"),
			Bool:        Bool(true),
			Binary:      []byte("binary"),
			ByteSlice:   []byte("World"),
			BoolSlice:   []bool{true, true, false, true, false},
			IntSlice:    []int{123, 321, 1234567},
			DoubleSlice: []float64{0.123, 0.321, 0.1234567},
			StringSlice: []string{"FOO", "BAR"},
		},
		StructSlice: []*HoldPtrMSG{
			{
				Int:         Int(2),
				IntNeg:      Int(math.MinInt32 + 1),
				Double:      Double(0.2),
				String:      String("Hello2"),
				Bool:        Bool(true),
				Binary:      []byte("binary2"),
				ByteSlice:   []byte("World2"),
				BoolSlice:   []bool{true, true, false, true, false},
				IntSlice:    []int{123, 321, 1234567},
				DoubleSlice: []float64{0.123, 0.321, 0.1234567},
				StringSlice: []string{"FOO2", "BAR2"},
			},
			{
				Int:         Int(3),
				IntNeg:      Int(math.MinInt32 + 2),
				Double:      Double(0.3),
				String:      String("Hello3"),
				Bool:        Bool(true),
				Binary:      []byte("binary3"),
				ByteSlice:   []byte("World3"),
				BoolSlice:   []bool{true, true, false, true, false},
				IntSlice:    []int{123, 321, 1234567},
				DoubleSlice: []float64{0.123, 0.321, 0.1234567},
				StringSlice: []string{"FOO3", "BAR3"},
			},
		},
	}
	valMSG = ValMSG{
		Int:         math.MaxInt64,
		IntNeg:      math.MinInt32 + 1,
		Double:      math.MaxFloat64,
		String:      "Hello",
		Bool:        true,
		Binary:      []byte("binary"),
		ByteSlice:   []byte("World"),
		BoolSlice:   []bool{true, true, false, true, false},
		IntSlice:    []int{123, 321, 1234567},
		DoubleSlice: []float64{0.123, 0.321, 0.1234567},
		StringSlice: []string{"FOO", "BAR"},
		Struct: &HoldValMSG{
			Int:         1,
			IntNeg:      math.MinInt32 + 1,
			Double:      0.1,
			String:      "Hello",
			Bool:        true,
			Binary:      []byte("binary"),
			ByteSlice:   []byte("World"),
			BoolSlice:   []bool{true, true, false, true, false},
			IntSlice:    []int{123, 321, 1234567},
			DoubleSlice: []float64{0.123, 0.321, 0.1234567},
			StringSlice: []string{"FOO", "BAR"},
		},
		StructSlice: []*HoldValMSG{
			{
				Int:         2,
				IntNeg:      math.MinInt32 + 1,
				Double:      0.2,
				String:      "Hello2",
				Bool:        true,
				Binary:      []byte("binary2"),
				ByteSlice:   []byte("World2"),
				BoolSlice:   []bool{true, true, false, true, false},
				IntSlice:    []int{123, 321, 1234567},
				DoubleSlice: []float64{0.123, 0.321, 0.1234567},
				StringSlice: []string{"FOO2", "BAR2"},
			},
			{
				Int:         3,
				IntNeg:      math.MinInt32 + 2,
				Double:      0.3,
				String:      "Hello3",
				Bool:        true,
				Binary:      []byte("binary3"),
				ByteSlice:   []byte("World3"),
				BoolSlice:   []bool{true, true, false, true, false},
				IntSlice:    []int{123, 321, 1234567},
				DoubleSlice: []float64{0.123, 0.321, 0.1234567},
				StringSlice: []string{"FOO3", "BAR3"},
			},
		},
	}
}
