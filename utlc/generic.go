package utlc

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

func IsZero(obj any) bool {
	if obj == nil {
		return true
	}

	objValue := reflect.ValueOf(obj)
	//nolint:exhaustive // we only care about these kinds
	switch objValue.Kind() {
	case reflect.Map:
		fallthrough
	case reflect.Slice, reflect.Array:
		return objValue.Len() == 0
	case reflect.Pointer:
		return objValue.IsNil()
	case reflect.String:
		return obj == ""
	case reflect.Bool:
		return obj == false
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		fallthrough
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		fallthrough
	case reflect.Float32, reflect.Float64:
		return slices.Contains(numericZeros, obj)
	}

	return false
}

func OrElse[T any](obj *T, defaultVal T) T {
	if obj == nil {
		return defaultVal
	}

	if IsZero(*obj) {
		return defaultVal
	}

	return *obj
}
