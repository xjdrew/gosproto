package sproto_types

import (
	"reflect"
	"strings"
	"testing"

	sproto "github.com/xjdrew/gosproto"
)

func ptrString(s string) *string {
	return &s
}

func ptrInt(v int64) *int64 {
	return &v
}

func ptrFloat(v float64) *float64 {
	return &v
}

func TestTypes(t *testing.T) {
	cases := []struct {
		obj interface{}
	}{
		{&Person{
			Name:  ptrString("David"),
			Id:    ptrInt(123),
			Email: ptrString("aaa@example.com"),
			Phone: []*PersonPhoneNumber{
				{Number: ptrString("1234567"), Type: ptrInt(1)},
				{Number: ptrString("8765432"), Type: ptrInt(2)},
			},
			Height: ptrInt(178),
			Data:   []byte("extra data"),
			Weight: ptrFloat(64.3),
			Pics:   [][]byte{[]byte("image data1"), []byte("image data2")},
		}},

		{&CreditCard{
			CardNum: ptrString("987654321"),
			Owner:   &Person{Name: ptrString("Ken")},
		}},

		{&Bank{
			Cards: map[string]*Person{
				"11111": &Person{
					Name: ptrString("Keven"),
					Id:   ptrInt(100),
				},
				"22222": &Person{
					Name: ptrString("John"),
					Id:   ptrInt(101),
				},
			},
			Clients: map[int64]*Person{
				12345: &Person{Id: ptrInt(12345)},
				67890: &Person{Id: ptrInt(67890)},
			},
		}},

		{&SimpleItem{
			Key:   ptrInt(999),
			Value: ptrString("something"),
		}},

		{&NodeItem{
			Id: ptrInt(1),
			Node: &NodeItem{
				Id: ptrInt(2),
				Node: &NodeItem{
					Id:   ptrInt(3),
					Node: nil,
				},
			},
		}},

		{&ArraysStruct{
			IntArr:    []int64{1, 2, 3, 4, 5},
			BoolArr:   []bool{true, false, true},
			StrArr:    []string{"str1", "str2", "str3"},
			BinArr:    [][]byte{[]byte("bin1"), []byte("bin2"), []byte("bin3")},
			DoubleArr: []float64{1.23, 4.56, 7.89},
			StructArr: []*SimpleItem{
				{Key: ptrInt(1), Value: ptrString("v1")},
				{Key: ptrInt(2), Value: ptrString("v2")},
			},
		}},

		{&NestedMapItem{
			Id: ptrInt(33),
			Nested: map[int64]*string{
				123: ptrString("whatever"),
			},
		}},

		{&NestedArrayItem{
			Id: ptrInt(66),
			Nested: []*SimpleItem{
				{Key: ptrInt(77), Value: ptrString("77 value")},
				{Key: ptrInt(88), Value: ptrString("88 value")},
			},
		}},

		{&MapStruct{
			Map1: map[int64]*SimpleItem{
				123: &SimpleItem{Key: ptrInt(123), Value: ptrString("item1")},
				456: &SimpleItem{Key: ptrInt(456), Value: ptrString("item2")},
			},
			Map2: map[int64]*string{
				123: ptrString("item1"),
				456: ptrString("item2"),
			},
			Map3: map[int64]*NodeItem{
				999: &NodeItem{Id: ptrInt(1000), Node: &NodeItem{}},
				888: &NodeItem{Id: ptrInt(889), Node: &NodeItem{}},
			},
			Map4: map[int64]*NodeItem{
				999: &NodeItem{Id: ptrInt(999), Node: &NodeItem{}},
				888: &NodeItem{Id: ptrInt(888), Node: &NodeItem{}},
			},
		}},

		{&NestedMapStruct{
			NestedMap1: map[int64]map[int64]*string{
				444: map[int64]*string{
					555: ptrString("555"),
				},
				777: map[int64]*string{
					888: ptrString("888"),
				},
			},
			NestedMap2: map[int64]*NestedMapItem{
				444: &NestedMapItem{
					Id:     ptrInt(444),
					Nested: map[int64]*string{666: ptrString("666")},
				},
				777: &NestedMapItem{
					Id:     ptrInt(777),
					Nested: map[int64]*string{999: ptrString("999")},
				},
			},
			NestedArr: map[int64][]*SimpleItem{
				1: []*SimpleItem{{Key: ptrInt(1), Value: ptrString("1111")}},
				2: []*SimpleItem{{Key: ptrInt(2), Value: ptrString("2222")}},
			},
		}},
	}

	for _, c := range cases {
		tt := reflect.TypeOf(c.obj).Elem()
		newPtr := reflect.New(tt)
		newObj := newPtr.Interface()

		data, err := sproto.Encode(c.obj)
		if err != nil {
			t.Errorf("encode failed, obj: %+v, err: %s", c.obj, err)
			t.FailNow()
		}
		_, err = sproto.Decode(data, newObj)
		if err != nil {
			t.Errorf("decode failed, obj: %+v, err: %s", c.obj, err)
			t.FailNow()
		}
		if !reflect.DeepEqual(c.obj, newObj) {
			t.Errorf("decode failed, obj: %+v, err: %s", c.obj, err)
			t.FailNow()
		}
	}
}

func TestIndexKeyIsNil(t *testing.T) {
	obj := &Bank{
		Clients: map[int64]*Person{
			111: &Person{Id: nil},
		},
	}
	data, err := sproto.Encode(obj)
	if err != nil {
		t.Errorf("encode err: %v", err)
		return
	}

	dobj := &Bank{}
	_, err = sproto.Decode(data, dobj)
	if err == nil {
		t.Errorf("expected error")
		return
	}
	if !strings.Contains(err.Error(), "map key is nil") {
		t.Errorf("expected error contails \"map key is nil\", err: %s", err)
		return
	}
}
