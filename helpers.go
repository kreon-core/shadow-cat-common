package tul

import (
	"reflect"
	"slices"
)

var numericZeros = []any{ //nolint:gochecknoglobals // common utility
	int(0),
	int8(0),
	int16(0),
	int32(0),
	int64(0),
	uint(0),
	uint8(0),
	uint16(0),
	uint32(0),
	uint64(0),
	float32(0),
	float64(0),
}

func IsEmpty(obj any) bool {
	if obj == nil || obj == "" || obj == false {
		return true
	}

	if slices.Contains(numericZeros, obj) {
		return true
	}

	objValue := reflect.ValueOf(obj)
	//nolint:exhaustive // we only care about these kinds
	switch objValue.Kind() {
	case reflect.Map:
		fallthrough
	case reflect.Slice, reflect.Array:
		return objValue.Len() == 0
	case reflect.Struct:
		return reflect.DeepEqual(obj, ZeroOf(obj))
	case reflect.Pointer:
		return objValue.IsNil()
	}

	return false
}

func IsZero(obj any) bool {
	if obj == nil || obj == "" || obj == false {
		return true
	}

	if slices.Contains(numericZeros, obj) {
		return true
	}

	return reflect.DeepEqual(obj, ZeroOf(obj))
}

func ZeroOf(obj any) any {
	if obj == nil {
		return nil
	}

	return reflect.Zero(reflect.TypeOf(obj)).Interface()
}
