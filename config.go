package config

var globalConfigurator Configurator

// SetConfigurator 设置配置器
func SetConfigurator(configurator Configurator) {
	if globalConfigurator != nil {
		globalConfigurator.Close()
	}
	globalConfigurator = configurator
}
