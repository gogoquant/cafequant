package goplugin

import (
	"fmt"
	"snack.com/xiyanxiyan10/stocktrader/constant"
)

// Hello ...
func (g *GoPlugin) Hello(v ...interface{}) interface{} {
	e := g.exchanges[0]
	e.Log(constant.INFO, e.GetStockType(), 0.0, 0.0, fmt.Sprint("get val from js to go", v))
	return "success"
}

// goPluginGrid ...
type goPluginGrid struct {
	store []string
}
