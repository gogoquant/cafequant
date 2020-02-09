package config

import (
	"github.com/go-ini/ini"
	"log"
	"os"
	"strings"
)

var confs = make(map[string]string)

func init() {
	configFile := os.Getenv("QUANT_CONFIG")
	conf, err := ini.InsensitiveLoad(configFile)
	if err != nil {
		conf, err = ini.InsensitiveLoad("/tmp/config.ini")
		if err != nil {
			log.Fatalln("Load config.ini error:", err)
		}
	}
	keys := conf.Section("").KeyStrings()
	for _, k := range keys {
		confs[k] = conf.Section("").Key(k).String()
	}
	if confs["logstimezone"] == "" {
		confs["logstimezone"] = "Local"
	}
}

// String ...
func String(key string) string {
	return confs[strings.ToLower(key)]
}
