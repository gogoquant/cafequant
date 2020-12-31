package config

import (
	"github.com/go-ini/ini"
	"log"
	"os"
	"strings"
)

var confMap = make(map[string]string)

// Init ...
func init() {
	configFile := os.Getenv("QUANT_CONFIG")
	conf, err := ini.InsensitiveLoad(configFile)
	if err != nil {
		conf, err = ini.InsensitiveLoad("/tmp/config.ini")
		if err != nil {
			log.Panicln("Load config.ini error:", err)
		}
	}
	keys := conf.Section("").KeyStrings()
	for _, k := range keys {
		confMap[k] = conf.Section("").Key(k).String()
	}
	if confMap["logstimezone"] == "" {
		confMap["logstimezone"] = "Local"
	}
}

// String ...
func String(key string) string {
	return confMap[strings.ToLower(key)]
}
