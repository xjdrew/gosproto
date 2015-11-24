package sproto

import (
	"errors"
)

var (
	ErrNonPtr    = errors.New("sproto: Encode called with Non-Ptr")
	ErrNonStruct = errors.New("sproto: Encode called with Non-Ptr")
	ErrNil       = errors.New("sproto: Encode called with nil")
)

func String(v string) *string {
	return &v
}

func Bool(v bool) *bool {
	return &v
}

func Int8(v int8) *int8 {
	return &v
}

func Uint8(v uint8) *uint8 {
	return &v
}

func Int16(v int16) *int16 {
	return &v
}

func Uint16(v uint16) *uint16 {
	return &v
}

func Int32(v int32) *int32 {
	return &v
}

func Uint32(v uint32) *uint32 {
	return &v
}

func Int64(v int64) *int64 {
	return &v
}

func Uint64(v uint64) *uint64 {
	return &v
}

func Int(v int) *int {
	return &v
}

func Uint(v uint) *uint {
	return &v
}
