package marry

import (
	goback "github.com/dirkolbrich/gobacktest"
)

var (
	marryStore = map[string]goback.MarryHandler{}
)


// MarryStore
func MarryStore()map[string]goback.MarryHandler{
	return marryStore
}
