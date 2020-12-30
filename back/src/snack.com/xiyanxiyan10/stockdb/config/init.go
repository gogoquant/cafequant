package config

import (
	"flag"
	"fmt"
	"github.com/hprose/hprose-golang/io"
	"os"
	"snack.com/xiyanxiyan10/stockdb/log"
	"snack.com/xiyanxiyan10/stockdb/types"
)

const (
	version         = "0.2.2"
	minPeriod int64 = 3
)

func init() {
	io.Register(types.Response{}, "Response", "json")
	confPath := os.Getenv("STOCKDB_CONFIG")
	if confPath == "" {
		confPath = "/tmp/stockdb.ini"
	}
	flag.Parse()
	loadConfig(confPath)
	log.Log(log.InfoLog, fmt.Sprintf("StockDB Version %s running at %s", version, config["http.bind"]))
}
