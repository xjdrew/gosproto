package sproto_test

import (
	"testing"

	"reflect"

	"github.com/xjdrew/gosproto"
)

type Person struct {
	Name     *string   `sproto:"bytes,0,name=name"`
	Age      *int      `sproto:"varint,1,name=age"`
	Marital  *bool     `sproto:"boolean,2,name=marital"`
	Children []*Person `sproto:"struct,3,array,name=children"`
}

type Data struct {
	Numbers   []int64 `sproto:"varint,0,array,name=numbers"`
	Bools     []bool  `sproto:"boolean,1,array,name=bools"`
	Number    *int    `sproto:"varint,2,name=number"`
	BigNumber *int64  `sproto:"varint,3,name=bignumber"`
}

type TestCase struct {
	Name   string
	Struct interface{}
	Data   []byte
}

var testCases []*TestCase = []*TestCase{
	&TestCase{
		Name: "SimpleStruct",
		Struct: &Person{
			Name:    sproto.String("Alice"),
			Age:     sproto.Int(13),
			Marital: sproto.Bool(false),
		},
		Data: []byte{
			0x03, 0x00, // (fn = 3)
			0x00, 0x00, // (id = 0, value in data part)
			0x1C, 0x00, // (id = 1, value = 13)
			0x02, 0x00, // (id = 2, value = false)
			0x05, 0x00, 0x00, 0x00, // (sizeof "Alice")
			0x41, 0x6C, 0x69, 0x63, 0x65, // ("Alice")
		},
	},
	&TestCase{
		Name: "StructArray",
		Struct: &Person{
			Name: sproto.String("Bob"),
			Age:  sproto.Int(40),
			Children: []*Person{
				&Person{
					Name: sproto.String("Alice"),
					Age:  sproto.Int(13),
				},
				&Person{
					Name: sproto.String("Carol"),
					Age:  sproto.Int(5),
				},
			},
		},
		Data: []byte{
			0x04, 0x00, // (fn = 4)
			0x00, 0x00, // (id = 0, value in data part)
			0x52, 0x00, // (id = 1, value = 40)
			0x01, 0x00, // (skip id = 2)
			0x00, 0x00, // (id = 3, value in data part)
			0x03, 0x00, 0x00, 0x00, // (sizeof "Bob")
			0x42, 0x6F, 0x62, // ("Bob")
			0x26, 0x00, 0x00, 0x00, // (sizeof children)
			0x0F, 0x00, 0x00, 0x00, // (sizeof child 1)
			0x02, 0x00, //(fn = 2)
			0x00, 0x00, //(id = 0, value in data part)
			0x1C, 0x00, //(id = 1, value = 13)
			0x05, 0x00, 0x00, 0x00, // (sizeof "Alice")
			0x41, 0x6C, 0x69, 0x63, 0x65, //("Alice")
			0x0F, 0x00, 0x00, 0x00, // (sizeof child 2)
			0x02, 0x00, //(fn = 2)
			0x00, 0x00, //(id = 0, value in data part)
			0x0C, 0x00, //(id = 1, value = 5)
			0x05, 0x00, 0x00, 0x00, //(sizeof "Carol")
			0x43, 0x61, 0x72, 0x6F, 0x6C, //("Carol")
		},
	},
	&TestCase{
		Name: "NumberArray",
		Struct: &Data{
			Numbers: []int64{1, 2, 3, 4, 5},
		},
		Data: []byte{
			0x01, 0x00, // (fn = 1)
			0x00, 0x00, // (id = 0, value in data part)

			0x15, 0x00, 0x00, 0x00, // (sizeof numbers)
			0x04,                   //(sizeof int32)
			0x01, 0x00, 0x00, 0x00, //(1)
			0x02, 0x00, 0x00, 0x00, //(2)
			0x03, 0x00, 0x00, 0x00, //(3)
			0x04, 0x00, 0x00, 0x00, //(4)
			0x05, 0x00, 0x00, 0x00, //(5)
		},
	},
	&TestCase{
		Name: "BigNumberArray",
		Struct: &Data{
			Numbers: []int64{
				(1 << 32) + 1,
				(1 << 32) + 2,
				(1 << 32) + 3,
			},
		},
		Data: []byte{
			0x01, 0x00, // (fn = 1)
			0x00, 0x00, // (id = 0, value in data part)

			0x19, 0x00, 0x00, 0x00, // (sizeof numbers)
			0x08,                                           //(sizeof int32)
			0x01, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, //((1<<32) + 1)
			0x02, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, //((1<<32) + 2)
			0x03, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, //((1<<32) + 3)
		},
	},
	&TestCase{
		Name: "BoolArray",
		Struct: &Data{
			Bools: []bool{false, true, false},
		},
		Data: []byte{
			0x02, 0x00, // (fn = 2)
			0x01, 0x00, // (skip id)
			0x00, 0x00, // (id = 2, value in data part)

			0x03, 0x00, 0x00, 0x00, // (sizeof bools)
			0x00, //(false)
			0x01, //(true)
			0x00, //(false)
		},
	},
	&TestCase{
		Name: "Number",
		Struct: &Data{
			Number:    sproto.Int(100000),
			BigNumber: sproto.Int64(-10000000000),
		},
		Data: []byte{
			0x03, 0x00, // (fn = 3)
			0x03, 0x00, // (skip id 0/1)
			0x00, 0x00, // (id = 2, value in data part)
			0x00, 0x00, // (id = 3, value in data part)

			0x04, 0x00, 0x00, 0x00, //(sizeof number, data part)
			0xA0, 0x86, 0x01, 0x00, //(100000, 32bit integer)

			0x08, 0x00, 0x00, 0x00, //(sizeof bignumber, data part)
			0x00, 0x1C, 0xF4, 0xAB, 0xFD, 0xFF, 0xFF, 0xFF, //(-10000000000, 64bit integer)
		},
	},
}

func isEqualBytes(dst, src []byte) bool {
	sz := len(dst)
	if sz != len(src) {
		return false
	}
	for i := 0; i < sz; i++ {
		if dst[i] != src[i] {
			return false
		}
	}
	return true
}

func TestEncode(t *testing.T) {
	for _, tc := range testCases {
		output, err := sproto.Encode(tc.Struct)
		if err != nil {
			t.Fatalf("test case *%s* failed with error:%s", tc.Name, err)
		}
		if !isEqualBytes(output, tc.Data) {
			t.Log("encoded:", output)
			t.Log("expected:", tc.Data)
			t.Fatalf("test case %s failed", tc.Name)
		}
	}
}

func TestDecode(t *testing.T) {
	for _, tc := range testCases {
		sp := reflect.New(reflect.TypeOf(tc.Struct).Elem()).Interface()
		used, err := sproto.Decode(tc.Data, sp)
		if err != nil {
			t.Fatalf("test case *%s* failed with error:%s", tc.Name, err)
		}

		if used != len(tc.Data) {
			t.Fatalf("test case *%s* failed: data length mismatch", tc.Name)
		}

		output, err := sproto.Encode(sp)
		if err != nil {
			t.Fatalf("test case *%s* failed with error:%s", tc.Name, err)
		}
		if !isEqualBytes(output, tc.Data) {
			t.Log("encoded:", output)
			t.Log("expected:", tc.Data)
			t.Fatalf("test case %s failed", tc.Name)
		}
	}
}
