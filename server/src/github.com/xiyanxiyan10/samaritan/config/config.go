package config

import (
	"github.com/go-ini/ini"
	"log"
	"os"
	"strings"
)

var confs = make(map[string]string)

func init() {
	config_path := os.Getenv("QUANT_CONFIG")
	conf, err := ini.InsensitiveLoad(config_path)
	if err != nil {
		log.Fatalf("Load config.ini from (%s) error (%s):", config_path, err.Error())
		return
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

// GetConf ...
func GetConf(key string) string {
	return confs[strings.ToLower(key)]
}

// GetConfs ...
func GetConfs() map[string]string {
	return confs
}
