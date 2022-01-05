package util

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"
)

const (
	TimeLayout = "2006-01-02 15:04:05"
)

func CatchException(handle func(e interface{})) {
	if err := recover(); err != nil {
		e := printStackTrace(err)
		handle(e)
	}
}

func printStackTrace(err interface{}) string {
	buf := new(bytes.Buffer)
	fmt.Fprintf(buf, "%v\n", err)
	for i := 1; ; i++ {
		pc, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}
		fmt.Fprintf(buf, "%s:%d (0x%x)\n", file, line, pc)
	}
	return buf.String()
}

// TimeUnix2Str ...
func TimeUnix2Str(t int64) string {
	return time.Unix(t, 0).Format(TimeLayout)
}

// TimeUnix2Str
func TimeStr2Unix(in string) (int64, error) {
	times, err := time.Parse(TimeLayout, in)
	if err != nil {
		return 0, err
	}
	return times.Unix(), nil
}

// DeepCopy ...
func DeepCopy(dst, src interface{}) error {
	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(src); err != nil {
		return err
	}
	return gob.NewDecoder(bytes.NewBuffer(buf.Bytes())).Decode(dst)
}

// File2Map convert file json to map[string]string
func File2Map(path string) (map[string]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	content, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}
	m := make(map[string]string)
	err = json.Unmarshal([]byte(content), &m)
	if err != nil {
		return nil, err
	}
	return m, nil
}

// SafefloatDivide ...
func SafefloatDivide(lft, rht float64) float64 {
	if rht == 0.0 {
		return lft
	}
	return lft / rht
}

func converError(val interface{}, t string) error {
	return fmt.Errorf("conver error, the %T{%v} can not conver to a %v", val, val, t)
}

// DeepCopyStruct ...
func DeepCopyStruct(source, target interface{}) {
	data, _ := json.Marshal(source)
	json.Unmarshal(data, target)
}

func stringToBool(val string) (bool, error) {
	switch val {
	case "1", "t", "T", "true", "TRUE", "True", "ok", "OK", "yes", "YES":
		return true, nil
	case "0", "f", "F", "false", "FALSE", "False", "":
		return false, nil
	}
	return false, converError(val, "bool")
}

// Struct2Json ...
func Struct2Json(m interface{}) string {
	dataType, _ := json.Marshal(m)
	dataString := string(dataType)
	return dataString
}

// String : Conver "val" to a String
func String(val interface{}) (string, error) {
	switch ret := val.(type) {
	case string:
		return ret, nil
	case []byte:
		return string(ret), nil
	default:
		str := fmt.Sprintf("%+v", val)
		if val == nil || len(str) == 0 {
			return "", fmt.Errorf("conver.String(), the %+v is empty", val)
		}
		return str, nil
	}
}

// StringMust : Must Conver "val" to a String
func StringMust(val interface{}, def ...string) string {
	ret, err := String(val)
	if err != nil && len(def) > 0 {
		return def[0]
	}
	return ret
}

// Bool : Conver "val" to a Bool
func Bool(val interface{}) (bool, error) {
	if val == nil {
		return false, nil
	}
	switch ret := val.(type) {
	case bool:
		return ret, nil
	case int, int8, int16, int32, int64, float32, float64, uint, uint8, uint16, uint32, uint64:
		return ret != 0, nil
	case []byte:
		return stringToBool(string(ret))
	case string:
		return stringToBool(ret)
	default:
		return false, converError(val, "bool")
	}
}

// BoolMust : Must Conver "val" to a Bool
func BoolMust(val interface{}, def ...bool) bool {
	ret, err := Bool(val)
	if err != nil && len(def) > 0 {
		return def[0]
	}
	return ret
}

// Bytes : Conver "val" to []byte
func Bytes(val interface{}) ([]byte, error) {
	switch ret := val.(type) {
	case []byte:
		return ret, nil
	default:
		str, err := String(val)
		return []byte(str), err
	}
}

// BytesMust : Must Conver "val" to []byte
func BytesMust(val interface{}, def ...[]byte) []byte {
	ret, err := Bytes(val)
	if err != nil && len(def) > 0 {
		return def[0]
	}
	return ret
}

// Float32 : Conver "val" to a Float32
func Float32(val interface{}) (float32, error) {
	switch ret := val.(type) {
	case float32:
		return ret, nil
	case int:
		return float32(ret), nil
	case int8:
		return float32(ret), nil
	case int16:
		return float32(ret), nil
	case int32:
		return float32(ret), nil
	case int64:
		return float32(ret), nil
	case uint:
		return float32(ret), nil
	case uint8:
		return float32(ret), nil
	case uint16:
		return float32(ret), nil
	case uint32:
		return float32(ret), nil
	case uint64:
		return float32(ret), nil
	case float64:
		return float32(ret), nil
	case bool:
		if ret {
			return 1.0, nil
		}
		return 0.0, nil
	default:
		str := strings.Replace(strings.TrimSpace(StringMust(val)), " ", "", -1)
		f, err := strconv.ParseFloat(str, 32)
		return float32(f), err
	}
}

// Float32Must : Must Conver "val" to Float32
func Float32Must(val interface{}, def ...float32) float32 {
	ret, err := Float32(val)
	if err != nil && len(def) > 0 {
		return def[0]
	}
	return ret
}

// Float64 : Conver "val" to a Float64
func Float64(val interface{}) (float64, error) {
	switch ret := val.(type) {
	case float64:
		return ret, nil
	case int:
		return float64(ret), nil
	case int8:
		return float64(ret), nil
	case int16:
		return float64(ret), nil
	case int32:
		return float64(ret), nil
	case int64:
		return float64(ret), nil
	case uint:
		return float64(ret), nil
	case uint8:
		return float64(ret), nil
	case uint16:
		return float64(ret), nil
	case uint32:
		return float64(ret), nil
	case uint64:
		return float64(ret), nil
	case float32:
		return float64(ret), nil
	case bool:
		if ret {
			return 1.0, nil
		}
		return 0.0, nil
	default:
		str := strings.Replace(strings.TrimSpace(StringMust(val)), " ", "", -1)
		return strconv.ParseFloat(str, 64)
	}
}

// Float64Must : Must Conver "val" to Float64
func Float64Must(val interface{}, def ...float64) float64 {
	ret, err := Float64(val)
	if err != nil && len(def) > 0 {
		return def[0]
	}
	return ret
}

// Int : Conver "val" to a rounded Int
func Int(val interface{}) (int, error) {
	i, err := Int64(val)
	if err != nil {
		return 0, err
	}

	return int(i), nil
}

// IntMust : Must Conver "val" to a rounded Int
func IntMust(val interface{}, def ...int) int {
	ret, err := Int(val)
	if err != nil && len(def) > 0 {
		return def[0]
	}
	return ret
}

// Int32 : Conver "val" to a rounded Int32
func Int32(val interface{}) (int32, error) {
	i, err := Int64(val)
	if err != nil {
		return 0, err
	}

	return int32(i), nil
}

// Int32Must : Must Conver "val" to a rounded Int32
func Int32Must(val interface{}, def ...int32) int32 {
	ret, err := Int32(val)
	if err != nil && len(def) > 0 {
		return def[0]
	}
	return ret
}

// Int64 : Conver "val" to a rounded Int64
func Int64(val interface{}) (int64, error) {
	str, err := String(val)
	if err != nil {
		return 0, err
	}

	return strconv.ParseInt(str, 10, 64)
}

// Int64Must : Must Conver "val" to a rounded Int64
func Int64Must(val interface{}, def ...int64) int64 {
	ret, err := Int64(val)
	if err != nil && len(def) > 0 {
		return def[0]
	}
	return ret
}
