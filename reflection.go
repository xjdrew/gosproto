package sproto

import "reflect"

// 对齐value和target的值/指针类型。
// 如果类型已经一致则返回value本身；
// 如果target是指针类型但value是值类型，则返回value对应的指针类型值；
// 如果是反过来的情况则返回value对应的值类型值。
func adjustTypePtr(value reflect.Value, target reflect.Type) reflect.Value {
	if target.Kind() == reflect.Ptr && value.Kind() != reflect.Ptr {
		if value.CanAddr() {
			return value.Addr()
		} else {
			clone := reflect.New(value.Type())
			clone.Elem().Set(value)
			return clone
		}
	}
	if target.Kind() != reflect.Ptr && value.Kind() == reflect.Ptr {
		return value.Elem()
	}
	return value
}

// 等价于 field.Set(value) ，自动处理指针类型与值类型相互赋值的情况；
// 当field是值类型，value是指针类型时： field = *value
// 当field是指针类型，value是值类型时： field = &value
// 注意field和value的基础类型必须相同，否则panic
func setValue(field reflect.Value, value reflect.Value) {
	field.Set(adjustTypePtr(value, field.Type()))
}

func isSameBaseType(base reflect.Type, typ reflect.Type) bool {
	return base == typ || (typ.Kind() == reflect.Ptr && base == typ.Elem())
}
