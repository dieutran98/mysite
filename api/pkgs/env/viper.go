package env

import (
	"log/slog"

	"github.com/fsnotify/fsnotify"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
)

type Configure interface {
	SetConfigName(name string)
	SetConfigType(typ string)
	AddConfigPath(path string)
	OnConfigChange(run func(e fsnotify.Event))
	WatchConfig()
	ReadInConfig() error
	Unmarshal(rawVal interface{}, opts ...viper.DecoderConfigOption) error
}

func setConfigFile(configure Configure) error {
	configure.SetConfigName("env.develop") // name of config file (without extension)
	configure.SetConfigType("json")
	configure.AddConfigPath("./config")
	configure.OnConfigChange(onConfigChange(configure))
	configure.WatchConfig()
	return configure.ReadInConfig()
}

func mappingStruct(configure Configure) error {
	return configure.Unmarshal(&internalEnv, viperUnmarshalOption)
}

func viperUnmarshalOption(c *mapstructure.DecoderConfig) {
	if c == nil {
		return
	}
	c.TagName = "json"
}

func onConfigChange(config Configure) func(e fsnotify.Event) {
	return func(e fsnotify.Event) {
		slog.Info("env changed.", "fileName", e.Name)
		config.ReadInConfig()
		if err := mappingStruct(config); err != nil {
			slog.Error("failed to unMarshal env: ", "error", err.Error())
		}
	}
}
