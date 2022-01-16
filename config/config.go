package config

import (
	"os"
	"strings"

	"github.com/go-ini/ini"
	log "github.com/sirupsen/logrus"
)

var confMap = make(map[string]string)

// Init ...
func Init(path string) error {
	// Log as JSON instead of the default ASCII formatter.
	log.SetFormatter(&log.JSONFormatter{})

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	log.SetOutput(os.Stdout)

	// Only log the warning severity or above.
	log.SetLevel(log.InfoLevel)
	log.SetReportCaller(true)
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
