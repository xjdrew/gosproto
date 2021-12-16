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
	WireVarintName  = "integer" // int, uint, int8, uint8, int16, uint16, int32, uint32, int64, uint64
	WireBooleanName = "boolean" // bool
	WireStringName  = "string"  // string
	WireBytesName   = "binary"  // []byte
	WireDoubleName  = "double"  // double
	WireStructName  = "struct"  // struct
)

const (
	TagMin = 0
	TagMax = 32766
)

var (
	mutex sync.Mutex
	stMap = make(map[reflect.Type]*SprotoType)
)

type headerEncoder func(st *SprotoField, v reflect.Value) (header uint16, isNil bool)

// TODO: 避免 encoder 分配内存
type encoder func(st *SprotoField, v reflect.Value) []byte
type decoder func(val *uint16, data []byte, st *SprotoField, v reflect.Value) error

type SprotoField struct {
	field *reflect.StructField // go StructField

	Wire     string
	Tag      int
	Array    bool
	KeyTag   int    // -1 表示无效值
	ValueTag int    // -1 表示无效值
	SubType  string // 仅当 ValueTag != 1 时有效

	st *SprotoType // for struct types only

	headerEnc headerEncoder
	enc       encoder
	dec       decoder
}

func parseTag(s string) (tag int, err error) {
	tag, err = strconv.Atoi(s)
	if err != nil {
		return
	}
	if tag < TagMin || tag > TagMax {
		err = fmt.Errorf("tag(%d) overflow", tag)
		return
	}
	return
}

// parse filed meta information
func (sf *SprotoField) parse(s string) error {
	sf.KeyTag = -1
	sf.ValueTag = -1

	// wire,tag,options...
	fields := strings.Split(s, ",")
	if len(fields) < 2 {
		return fmt.Errorf("sproto: parse(%s) tag must have 2 or more fields", s)
	}

	sf.Wire = fields[0]
	switch sf.Wire {
	case WireVarintName, WireBooleanName, WireStringName, WireBytesName, WireDoubleName, WireStructName:
	default:
		return fmt.Errorf("sproto: parse(%s) unknown wire type: %s", s, sf.Wire)
	}

	var tag int
	var err error
	tag, err = parseTag(fields[1])
	if err != nil {
		return fmt.Errorf("sproto: parse(%s) parse tag option faield: %s", s, err)
	}
	sf.Tag = tag

	// optional options
	for i := 2; i < len(fields); i++ {
		f := fields[i]
		switch {
		case f == "array":
			sf.Array = true
		case strings.HasPrefix(f, "key="):
			tag, err = parseTag(f[len("key="):])
			if err != nil {
				return fmt.Errorf("parse(%s) parse key option failed:%s", s, err)
			}
			sf.KeyTag = tag
		case strings.HasPrefix(f, "value="):
			tag, err = parseTag(f[len("value="):])
			if err != nil {
				return fmt.Errorf("parse(%s) parse value option failed:%s", s, err)
			}
			sf.ValueTag = tag
		case strings.HasPrefix(f, "subtype="):
			sf.SubType = f[len("subtype="):]
		default:
			return fmt.Errorf("sproto: parse(%s) unknown option: %s", s, f)
		}
	}

	if sf.KeyTag != -1 && !sf.Array {
		return fmt.Errorf("sproto: parse(%s) failed: KeyTag depends on Array", s)
	}

	if sf.ValueTag != -1 {
		if sf.KeyTag == -1 {
			return fmt.Errorf("sproto: parse(%s) failed: ValueTag depends on KeyTag", s)
		}
		if sf.SubType == "" {
			return fmt.Errorf("sproto: parse(%s) failed: ValueTag depends on SubType", s)
		}
	}

	return nil
}

func (sf *SprotoField) assertWire(expectedWire string, expectedArray bool) error {
	if expectedWire != "" && sf.Wire != expectedWire {
		return fmt.Errorf("sproto: field(%s) expect %s but get %s", sf.field.Name, expectedWire, sf.Wire)
	}
	if sf.Array != expectedArray {
		n := "not"
		if expectedArray {
			n = ""
		}
		return fmt.Errorf("sproto: field(%s) should %s be array", sf.field.Name, n)
	}
	return nil
}

// 校验 meta 元信息是否与 map 类型匹配
func (sf *SprotoField) initMapElemType(mapType reflect.Type, valueType reflect.Type) (err error) {
	if valueType.Kind() != reflect.Ptr {
		err = fmt.Errorf("sproto: field(%s) illegal type(%s), expect reflect.Ptr", sf.field.Name, valueType.Kind().String())
		return
	}

	elemType := valueType.Elem()
	if elemType.Kind() != reflect.Struct {
		err = fmt.Errorf("sproto: field(%s) illegal type(%s), expect reflect.Struct", sf.field.Name, elemType.Kind().String())
		return
	}

	// check elemType
	var stype *SprotoType
	if stype, err = getSprotoTypeLocked(elemType); err != nil {
		return
	}

	keyField := stype.FieldByTag(sf.KeyTag)
	if keyField == nil {
		err = fmt.Errorf("sproto: field(%s) key type(%s) no tag(%d) ", sf.field.Name, stype.Type.Name(), sf.KeyTag)
		return
	}

	if keyField.Wire != WireVarintName && keyField.Wire != WireStringName {
		err = fmt.Errorf("sproto: field(%s) illegal key type(%s), map key must be integer or string", sf.field.Name, keyField.Wire)
		return
	}

	if !isSameBaseType(mapType.Key(), keyField.field.Type) {
		err = fmt.Errorf("sproto: field(%s) key type unmatched (%s != %s)", sf.field.Name, mapType.Key(), keyField.field.Type)
		return
	}

	if sf.ValueTag != -1 {
		valueField := stype.FieldByTag(sf.ValueTag)
		if valueField == nil {
			err = fmt.Errorf("sproto: field(%s) value type(%s) no tag(%d) ", sf.field.Name, stype.Type.Name(), sf.ValueTag)
			return
		}

		if mapType.Elem() != valueField.field.Type {
			err = fmt.Errorf("sproto: field(%s) value type unmatched (%s != %s)", sf.field.Name, mapType.Elem().Name(), valueField.field.Type)
			return
		}
	}
	sf.st = stype
	return
}

func (sf *SprotoField) initEncAndDec(structType reflect.Type, f *reflect.StructField) error {
	var stype reflect.Type
	var err error
	t1 := f.Type
	if t1.Kind() == reflect.Ptr {
		t1 = t1.Elem()
	}

	switch t1.Kind() {
	case reflect.Bool:
		sf.headerEnc = headerEncodeBool
		sf.dec = decodeBool
		err = sf.assertWire(WireBooleanName, false)
	case reflect.Int8, reflect.Uint8, reflect.Int16, reflect.Uint16,
		reflect.Int32, reflect.Uint32, reflect.Int64, reflect.Uint64,
		reflect.Int, reflect.Uint:
		sf.headerEnc = headerEncodeInt
		sf.enc = encodeInt
		sf.dec = decodeInt
		err = sf.assertWire(WireVarintName, false)
	case reflect.Float64:
		sf.headerEnc = headerEncodeDefault
		sf.enc = encodeDouble
		sf.dec = decodeDouble
	case reflect.String:
		sf.headerEnc = headerEncodeDefault
		sf.enc = encodeString
		sf.dec = decodeString
		err = sf.assertWire(WireStringName, false)
	case reflect.Struct:
		stype = t1
		sf.headerEnc = headerEncodeDefault
		sf.enc = encodeStruct
		sf.dec = decodeStruct
		err = sf.assertWire(WireStructName, false)
	case reflect.Slice:
		switch t2 := t1.Elem(); t2.Kind() {
		case reflect.Bool:
			sf.headerEnc = headerEncodeDefault
			sf.enc = encodeBoolSlice
			sf.dec = decodeBoolSlice
			err = sf.assertWire(WireBooleanName, true)
		case reflect.Uint8:
			sf.headerEnc = headerEncodeDefault
			// allowed to be "string" as well as "binary", for compatibility
			if sf.Wire == WireBytesName || sf.Wire == WireStringName {
				sf.enc = encodeBytes
				sf.dec = decodeBytes
				err = sf.assertWire("", false)
			} else {
				sf.enc = encodeIntSlice
				sf.dec = decodeIntSlice
				err = sf.assertWire(WireVarintName, true)
			}
		case reflect.Int8, reflect.Int16, reflect.Uint16,
			reflect.Int32, reflect.Uint32, reflect.Int64, reflect.Uint64,
			reflect.Int, reflect.Uint:
			sf.headerEnc = headerEncodeDefault
			sf.enc = encodeIntSlice
			sf.dec = decodeIntSlice
			err = sf.assertWire(WireVarintName, true)
		case reflect.Float64:
			sf.headerEnc = headerEncodeDefault
			sf.enc = encodeDoubleSlice
			sf.dec = decodeDoubleSlice
			err = sf.assertWire(WireDoubleName, true)
		case reflect.String:
			sf.headerEnc = headerEncodeDefault
			sf.enc = encodeStringSlice
			sf.dec = decodeStringSlice
			err = sf.assertWire(WireStringName, true)
		case reflect.Ptr:
			switch t3 := t2.Elem(); t3.Kind() {
			case reflect.Struct:
				stype = t3
				sf.headerEnc = headerEncodeDefault
				sf.enc = encodeStructSlice
				sf.dec = decodeStructSlice
				err = sf.assertWire(WireStructName, true)
			default:
				err = fmt.Errorf("sproto: field(%s) no coders for %s -> %s -> %s", sf.field.Name, t1.Kind().String(), t2.Kind().String(), t3.Kind().String())
			}
		default:
			err = fmt.Errorf("sproto: field(%s) no coders for %s -> %s", sf.field.Name, t1.Kind().String(), t2.Kind().String())
		}
	case reflect.Map:
		err = sf.assertWire(WireStructName, true)
		if err != nil {
			break
		}

		var valueType reflect.Type
		if sf.ValueTag == -1 {
			valueType = t1.Elem()
		} else {
			valueField, ok := structType.FieldByName(sf.SubType)
			if !ok {
				err = fmt.Errorf("sproto: field(%s) no subtype filed(%s)", sf.field.Name, sf.SubType)
				break
			}
			valueType = valueField.Type
		}
		err = sf.initMapElemType(t1, valueType)
		if err != nil {
			break
		}
		sf.headerEnc = headerEncodeDefault
		sf.enc = encodeMap
		sf.dec = decodeMap
	default:
		err = fmt.Errorf("sproto: field(%s) no coders for %s", sf.field.Name, t1.Kind().String())
	}

	if err != nil {
		return err
	}

	if stype != nil {
		if sf.st, err = getSprotoTypeLocked(stype); err != nil {
			return err
		}
	}
	return nil
}

func (sf *SprotoField) init(structType reflect.Type, f *reflect.StructField) error {
	sf.field = f

	tagString := f.Tag.Get("sproto")
	if tagString == "" {
		sf.Tag = -1
		return nil
	}

	if err := sf.parse(tagString); err != nil {
		return err
	}
	if err := sf.initEncAndDec(structType, f); err != nil {
		return err
	}
	return nil
}

type SprotoType struct {
	Type reflect.Type // go internal type

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

func GetSprotoType(t reflect.Type) (*SprotoType, error) {
	if t.Kind() != reflect.Struct {
		return nil, fmt.Errorf("sproto: type must have kind struct")
	}
	mutex.Lock()
	sp, err := getSprotoTypeLocked(t)
	mutex.Unlock()
	return sp, err
}

func getSprotoTypeLocked(t reflect.Type) (*SprotoType, error) {
	if st, ok := stMap[t]; ok {
		return st, nil
	}

	st := new(SprotoType)
	stMap[t] = st

	st.Type = t
	numField := t.NumField()
	st.Fields = make([]*SprotoField, numField)
	st.order = make([]int, numField)
	st.tagMap = make(map[int]int)

	for i := 0; i < numField; i++ {
		sf := new(SprotoField)
		f := t.Field(i)
		if err := sf.init(t, &f); err != nil {
			delete(stMap, t)
			return nil, err
		}

		st.Fields[i] = sf
		st.order[i] = i
		if sf.Tag >= 0 {
			// check repeated tag
			if _, ok := st.tagMap[sf.Tag]; ok {
				return nil, fmt.Errorf("sproto: field(%s.%s) tag repeated", st.Type.Name(), sf.field.Name)
			}
			st.tagMap[sf.Tag] = i
		}
	}

	// Re-order prop.order
	sort.Sort(st)
	return st, nil
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
