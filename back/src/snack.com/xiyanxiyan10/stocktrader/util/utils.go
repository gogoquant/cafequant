package util

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

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
