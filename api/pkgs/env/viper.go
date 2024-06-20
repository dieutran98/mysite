package env

import (
	"log/slog"
	"mysite/pkgs/logger"

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
	SetDefault(key string, value interface{})
}

type viperConfig struct {
	viperCfg viperConfigure
}

func newViperConfig() viperConfig {
	return viperConfig{
		viperCfg: viper.New(),
	}
}

func (v viperConfig) setDefault() error {
	v.viperCfg.SetDefault("database.connmaxlifeidle", 30)
	v.viperCfg.SetDefault("database.transactiontimeout", 30)
	v.viperCfg.SetDefault("database.connmaxlifetime", 30)
	v.viperCfg.SetDefault("database.connmaxopen", 100)
	v.viperCfg.SetDefault("database.sslmode", "disable")
	v.viperCfg.SetDefault("database.port", "5432")
	return nil
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
	if err := v.viperCfg.ReadInConfig(); err != nil {
		slog.Error("failed to ReadInConfig: ", logger.AttrError(err))
		return
	}
	if err := v.mappingStruct(); err != nil {
		slog.Error("failed to unMarshal env: ", logger.AttrError(err))
		return
	}
}

func viperUnmarshalOption(c *mapstructure.DecoderConfig) {
	if c == nil {
		return
	}
	c.TagName = "json"
}
