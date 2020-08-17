package plugin

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/robertkrimen/otto/registry"
	"snack.com/xiyanxiyan10/stocktrader/config"
)

var (
	scriptsMap = make(map[string]string)
	scripts    = []string{}
	entry      = registry.Register(func() string {
		return strings.Join(scripts, "")
	})
)

// Get 获取插件脚本
func Get(name string) string {
	return scriptsMap[name]
}

// Load 加载插件脚本
func Load() {
	filepath.Walk(config.String("plugin"), func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() || !strings.HasSuffix(path, ".js") {
			return err
		}
		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()
		data, _ := ioutil.ReadAll(file)
		scripts = append(scripts, string(data))
		base := filepath.Base(path)
		bases := strings.Split(base, ".")
		scriptsMap[bases[0]] = string(data)
		return nil
	})
}
