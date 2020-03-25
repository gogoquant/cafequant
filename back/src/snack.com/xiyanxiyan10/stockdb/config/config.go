package config

import (
	"encoding/base64"
	"fmt"
	"github.com/go-ini/ini"
	"snack.com/xiyanxiyan10/stockdb/types"
	"strings"
	"time"
)

type logConfig struct {
	Enable   bool           `ini:"-"`
	Timezone string         `ini:"timezone"`
	Console  bool           `ini:"console"`
	File     bool           `ini:"file"`
	Location *time.Location `ini:"-"`
}

var (
	config        = make(map[string]string)
	openMethods   = make(map[string]bool)
	defaultOption = types.Option{}
	logConf       = logConfig{}
)

func GetConfig() map[string]string {
	return config
}

func GetDefaultOption() types.Option {
	return defaultOption
}

func GetOpenMethods() map[string]bool {
	return openMethods
}

func GetLogConf() logConfig {
	return logConf
}

func loadConfig(path string) {
	if path == "" {
		path = "stockdb.ini"
	}
	conf, err := ini.Load(path)
	if err != nil {
		fmt.Printf("Load config file error: %s", err.Error())
		return
	}
	_ = conf.Section("log").MapTo(&logConf)
	logConf.Enable = logConf.Console || logConf.File
	if loc, err := time.LoadLocation(logConf.Timezone); err != nil || loc == nil {
		logConf.Location = time.Local
	} else {
		logConf.Location = loc
	}
	_ = conf.Section("default").MapTo(&defaultOption)
	if defaultOption.Period < minPeriod {
		defaultOption.Period = minPeriod
	}
	for _, s := range conf.Sections() {
		name := s.Name()
		for _, k := range s.Keys() {
			config[strings.ToLower(name+"."+k.Name())] = k.Value()
		}
	}
	for _, m := range strings.Split(config["http.openmethods"], ",") {
		openMethods[strings.TrimSpace(m)] = true
	}
	config["http.auth"] = "Basic " + base64.StdEncoding.EncodeToString([]byte(config["http.auth"]))
}
