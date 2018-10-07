package marry

import (
	goback "github.com/xiyanxiyan10/gobacktest"
)

var (
	marryStore = map[string]goback.MarryHandler{}
)


// MarryStore
func MarryStore()map[string]goback.MarryHandler{
	return marryStore
}
