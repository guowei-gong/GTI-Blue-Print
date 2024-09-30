package etc

import (
	"github.com/gti-blue-print/config"
	"github.com/gti-blue-print/config/core/value"
	"github.com/gti-blue-print/config/file/core"
)

var globalConfigurator config.Configurator

func init() {
	// TODO: 从环境变量读取 path
	globalConfigurator = config.NewConfigurator(config.WithSources(core.NewSource("./config", config.ReadOnly)))
}

// Get 获取配置值
func Get(pattern string, def ...interface{}) value.Value {
	return globalConfigurator.Get(pattern, def...)
}
