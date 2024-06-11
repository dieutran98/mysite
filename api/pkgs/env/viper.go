package env

import (
	"log/slog"

	"github.com/fsnotify/fsnotify"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
)

type viperConfigure interface {
	SetConfigName(name string)
	SetConfigType(typ string)
	AddConfigPath(path string)
	OnConfigChange(run func(e fsnotify.Event))
	WatchConfig()
	ReadInConfig() error
	Unmarshal(rawVal interface{}, opts ...viper.DecoderConfigOption) error
}

type viperConfig struct {
	viperCfg viperConfigure
}

func newViperConfig() viperConfig {
	return viperConfig{
		viperCfg: viper.New(),
	}
}

func (v viperConfig) setConfigFile() error {
	v.viperCfg.SetConfigName("env.develop") // name of config file (without extension)
	v.viperCfg.SetConfigType("json")
	v.viperCfg.AddConfigPath("./config")

	v.viperCfg.OnConfigChange(v.onConfigChangeFunc)
	v.viperCfg.WatchConfig()

	return v.viperCfg.ReadInConfig()
}

func (v viperConfig) mappingStruct() error {
	return v.viperCfg.Unmarshal(&internalEnv, viperUnmarshalOption)
}

func (v viperConfig) onConfigChangeFunc(e fsnotify.Event) {
	slog.Info("env changed.", "fileName", e.Name)
	v.viperCfg.ReadInConfig()
	if err := v.mappingStruct(); err != nil {
		slog.Error("failed to unMarshal env: ", "error", err.Error())
	}
}

func viperUnmarshalOption(c *mapstructure.DecoderConfig) {
	if c == nil {
		return
	}
	c.TagName = "json"
}
