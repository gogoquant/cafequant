package util

import (
	"encoding/json"
	"fmt"
)

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
