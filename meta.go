package sproto

import (
	"fmt"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"sync"
)

const (
	WireVarintName  = "varint"  // int, uint, int8, uint8, int16, uint16, int32, uint32, int64, uint64
	WireBooleanName = "boolean" // bool
	WireBytesName   = "bytes"   // string, []byte
	WireStructName  = "struct"  // struct
)

var (
	mutex sync.Mutex
	spMap = make(map[reflect.Type]*SprotoType)
)

type headerEncoder func(st *SprotoField, v reflect.Value) (uint16, bool)
type encoder func(st *SprotoField, v reflect.Value) []byte
type decoder func(val *uint16, data []byte, st *SprotoField, v reflect.Value) error

type SprotoField struct {
	Name     string
	OrigName string
	Wire     string
	Tag      int
	Array    bool

	st *SprotoType // for struct types only

	index     []int // index sequence for Value.FieldByIndex
	headerEnc headerEncoder
	enc       encoder
	dec       decoder
}

// parse filed meta information
func (sf *SprotoField) parse(s string) {
	// children,object,3,array
	fields := strings.Split(s, ",")
	if len(fields) < 2 {
		panic("sproto: tag has 2 fields at least")
	}
	sf.Wire = fields[0]
	switch sf.Wire {
	case WireVarintName, WireBooleanName, WireBytesName, WireStructName:
	default:
		panic("sproto: unknown wire type: " + sf.Wire)
	}

	var err error
	sf.Tag, err = strconv.Atoi(fields[1])
	if err != nil {
		panic(err)
	}

	for i := 2; i < len(fields); i++ {
		f := fields[i]
		switch {
		case f == "array":
			sf.Array = true
		case strings.HasPrefix(f, "name="):
			sf.OrigName = f[len("name="):]
		}
	}
}

func (sf *SprotoField) assertWire(expectedWire string, expectedArray bool) {
	if sf.Wire != expectedWire {
		panic(fmt.Sprintf("sproto: field(%s) expect %s but get %s", sf.Name, expectedWire, sf.Wire))
	}
	if sf.Array != expectedArray {
		n := "not"
		if expectedArray {
			n = ""
		}
		panic(fmt.Sprintf("sproto: field(%s) should %s be array", sf.Name, n))
	}
}

func (sf *SprotoField) setEncAndDec(f *reflect.StructField) {
	var stype reflect.Type
	switch t1 := f.Type; t1.Kind() {
	case reflect.Ptr:
		switch t2 := t1.Elem(); t2.Kind() {
		case reflect.Bool:
			sf.headerEnc = headerEncodeBool
			sf.dec = decodeBool
			sf.assertWire(WireBooleanName, false)
		case reflect.Int8, reflect.Uint8, reflect.Int16, reflect.Uint16,
			reflect.Int32, reflect.Uint32, reflect.Int64, reflect.Uint64,
			reflect.Int, reflect.Uint:
			sf.headerEnc = headerEncodeInt
			sf.enc = encodeInt
			sf.dec = decodeInt
			sf.assertWire(WireVarintName, false)
		case reflect.String:
			sf.headerEnc = headerEncodeDefault
			sf.enc = encodeString
			sf.dec = decodeString
			sf.assertWire(WireBytesName, false)
		case reflect.Struct:
			stype = t1.Elem()
			sf.headerEnc = headerEncodeDefault
			sf.enc = encodeStruct
			sf.dec = decodeStruct
			sf.assertWire(WireStructName, false)
		default:
			panic("sproto: no coders for " + t1.Kind().String() + " -> " + t2.Kind().String())
		}
	case reflect.Slice:
		switch t2 := t1.Elem(); t2.Kind() {
		case reflect.Bool:
			sf.headerEnc = headerEncodeDefault
			sf.enc = encodeBoolSlice
			sf.dec = decodeBoolSlice
			sf.assertWire(WireBooleanName, true)
		case reflect.Uint8:
			sf.headerEnc = headerEncodeDefault
			if sf.Wire == WireBytesName {
				sf.enc = encodeBytes
				sf.dec = decodeBytes
			} else {
				sf.enc = encodeIntSlice
				sf.dec = decodeIntSlice
			}
			sf.assertWire(WireBytesName, true)
		case reflect.Int8, reflect.Int16, reflect.Uint16,
			reflect.Int32, reflect.Uint32, reflect.Int64, reflect.Uint64,
			reflect.Int, reflect.Uint:
			sf.headerEnc = headerEncodeDefault
			sf.enc = encodeIntSlice
			sf.dec = decodeIntSlice
			sf.assertWire(WireVarintName, true)
		case reflect.String:
			sf.headerEnc = headerEncodeDefault
			sf.enc = encodeStringSlice
			sf.dec = decodeStringSlice
			sf.assertWire(WireBytesName, true)
		case reflect.Ptr:
			switch t3 := t2.Elem(); t3.Kind() {
			case reflect.Struct:
				stype = t2.Elem()
				sf.headerEnc = headerEncodeDefault
				sf.enc = encodeStructSlice
				sf.dec = decodeStructSlice
				sf.assertWire(WireStructName, true)
			default:
				panic("sproto: no coders for " + t1.Kind().String() + " -> " + t2.Kind().String() + " -> " + t3.Kind().String())
			}
		default:
			panic("sproto: no coders for " + t1.Kind().String() + " -> " + t2.Kind().String())
		}
	default:
		panic("sproto: no coders for " + t1.Kind().String())
	}

	if stype != nil {
		sf.st = getSprotoTypeLocked(stype)
	}
}

func (sf *SprotoField) init(f *reflect.StructField) {
	sf.Name = f.Name
	sf.OrigName = f.Name

	tagString := f.Tag.Get("sproto")
	if tagString == "" {
		sf.Tag = -1
		return
	}

	sf.index = f.Index
	sf.parse(tagString)
	sf.setEncAndDec(f)
}

type SprotoType struct {
	Name   string // struct name
	Fields []*SprotoField
	tagMap map[int]int // tag -> fileds index
	order  []int       // list of struct field numbers in tag order
}

func (st *SprotoType) Len() int { return len(st.order) }
func (st *SprotoType) Less(i, j int) bool {
	return st.Fields[st.order[i]].Tag < st.Fields[st.order[j]].Tag
}
func (st *SprotoType) Swap(i, j int) {
	st.order[i], st.order[j] = st.order[j], st.order[i]
}

func (st *SprotoType) FieldByTag(tag int) *SprotoField {
	if index, ok := st.tagMap[tag]; ok {
		return st.Fields[index]
	}
	return nil
}

func GetSprotoType(t reflect.Type) *SprotoType {
	if t.Kind() != reflect.Struct {
		panic("sproto: type must have kind struct")
	}
	mutex.Lock()
	st := getSprotoTypeLocked(t)
	mutex.Unlock()
	return st
}

func getSprotoTypeLocked(t reflect.Type) *SprotoType {
	if st, ok := spMap[t]; ok {
		return st
	}

	st := new(SprotoType)
	spMap[t] = st

	st.Name = t.Name()

	numField := t.NumField()
	st.Fields = make([]*SprotoField, numField)
	st.order = make([]int, numField)
	st.tagMap = make(map[int]int)

	for i := 0; i < numField; i++ {
		sf := new(SprotoField)
		f := t.Field(i)
		sf.init(&f)

		st.Fields[i] = sf
		st.order[i] = i
		if sf.Tag >= 0 {
			st.tagMap[sf.Tag] = i
		}
	}

	// Re-order prop.order
	sort.Sort(st)
	return st
}

// Get the type and value of a pointer to a struct from interface{}
func getbase(sp interface{}) (t reflect.Type, v reflect.Value, err error) {
	if sp == nil {
		err = ErrNil
		return
	}

	t = reflect.TypeOf(sp)
	if t.Kind() != reflect.Ptr {
		err = ErrNonPtr
		return
	}

	if t.Elem().Kind() != reflect.Struct {
		err = ErrNonStruct
		return
	}

	v = reflect.ValueOf(sp)
	if v.IsNil() {
		err = ErrNil
		return
	}

	return
}
