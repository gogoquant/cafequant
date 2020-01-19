package config

import (
	"flag"
	"fmt"
	"os"

	"github.com/hprose/hprose-golang/io"
)

const (
	version         = "0.2.2"
	minPeriod int64 = 3
)

func init() {
	io.Register(response{}, "Response", "json")
	confPath := os.Getenv("STOCK_CONFIG")
	if confPath == ""{
		confPath = "/tmp/stockdb.ini"
	}
	flag.Parse()
	loadConfig(confPath)
	log(logInfo, fmt.Sprintf("StockDB Version %s running at %s", version, config["http.bind"]))
}
