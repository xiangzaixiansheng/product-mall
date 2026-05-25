package tools

import (
	"encoding/json"
	"log/slog"
	"reflect"
	"strconv"
)

func StrToInt(val string) int {
	v1, err := strconv.ParseInt(val, 10, 64)
	if err != nil {
		slog.Error("strconv parse int failed", "error", err)
	}
	return int(v1)
}

func ToString(v any) string {
	t := reflect.TypeOf(v)
	var s string
	switch t.Kind() {
	case reflect.Int:
		s = strconv.FormatInt(int64(v.(int)), 10)
	case reflect.Int64:
		s = strconv.FormatInt(int64(v.(int64)), 10)
	case reflect.Int16:
		s = strconv.FormatInt(int64(v.(int16)), 10)
	case reflect.Int8:
		s = strconv.FormatInt(int64(v.(int8)), 10)
	case reflect.Uint:
		s = strconv.FormatUint(uint64(v.(uint)), 10)
	case reflect.Uint64:
		s = strconv.FormatUint(v.(uint64), 10)
	case reflect.Uint16:
		s = strconv.FormatUint(uint64(v.(uint16)), 10)
	case reflect.Uint8:
		s = strconv.FormatUint(uint64(v.(uint8)), 10)
	case reflect.Bool:
		s = strconv.FormatBool(v.(bool))
	case reflect.Float32:
		s = strconv.FormatFloat(float64(v.(float32)), 'f', -1, 64)
	case reflect.Float64:
		s = strconv.FormatFloat(v.(float64), 'f', -1, 64)
	case reflect.Map, reflect.Struct, reflect.Slice:
		s = ToJson(v)
	default:
		slog.Warn("unsupported type for ToString", "type", t.Kind().String())
	}
	return s
}

func ToJson(m any) string {
	data, err := json.Marshal(m)
	if err != nil {
		panic(err)
	}
	return string(data)
}

func TypeOf(data any) string {
	if data == nil {
		return "nil"
	}
	return reflect.TypeOf(data).Kind().String()
}
