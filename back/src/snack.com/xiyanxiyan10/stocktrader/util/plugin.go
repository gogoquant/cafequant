package util

import (
	"fmt"
	"plugin"
)

// HotPlugin ...
func HotPlugin(path, funcname string) (interface{}, error) {
	handler, err := plugin.Open(path)
	if err != nil {
		return nil, err
	}
	s, err := handler.Lookup(funcname)
	if err != nil {
		return nil, err
	}
	if s == nil {
		return nil, fmt.Errorf("HotPlugin fail interface is nil")
	}
	return s, nil
}
