package config

import (
	"github.com/go-ini/ini"
	"log"
	"strings"
)

var confMap = make(map[string]string)

// Init ...
func Init(path string) error {
	conf, err := ini.InsensitiveLoad(path)
	if err != nil {
		log.Printf("Load %s error: %s\n", path, err.Error())
		return err
	}
	keys := conf.Section("").KeyStrings()
	for _, k := range keys {
		confMap[k] = conf.Section("").Key(k).String()
	}
	if confMap["logstimezone"] == "" {
		confMap["logstimezone"] = "Local"
	}
	return nil
}

// String ...
func String(key string) string {
	return confMap[strings.ToLower(key)]
}
