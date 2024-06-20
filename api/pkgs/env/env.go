package env

import (
	"github.com/pkg/errors"
)

type AppEnv struct {
	Database database `json:"database"`
	Jwt      jwt      `json:"jwt"`
}

type database struct {
	Database           string `json:"database" validate:"required"`
	HostName           string `json:"hostName" validate:"required"`
	User               string `json:"user" validate:"required"`
	Password           string `json:"password" validate:"required"`
	SslMode            string `json:"sslMode"`
	Port               string `json:"port"`
	ConnMaxLifeTime    int    `json:"connMaxLifeTime"`
	ConnMaxLifeIdle    int    `json:"connMaxLifeIdle"`
	ConnMaxOpen        int    `json:"connMaxOpen"`
	TransactionTimeout int    `json:"transactionTimeout"`
}

type jwt struct {
	AccessKey  string `json:"accessKey"`
	RefreshKey string `json:"refreshKey"`
}

type configure interface {
	setConfigFile() error
	mappingStruct() error
	setDefault() error
}

type EnvOption func(appEnv *AppEnv)

var (
	internalEnv AppEnv
)

var newConfigure = func() configure {
	return newViperConfig()
}

func ReadEnv(opts ...EnvOption) error {
	if len(opts) > 0 {
		for _, opt := range opts {
			opt(&internalEnv)
		}
		return nil
	}

	configure := newConfigure()

	if err := configure.setDefault(); err != nil {
		return errors.Wrap(err, "failed to set default config")
	}

	if err := configure.setConfigFile(); err != nil {
		return errors.Wrap(err, "failed set config file")
	}

	if err := configure.mappingStruct(); err != nil {
		return errors.Wrap(err, "failed to mapping env to struct")
	}

	return nil
}

func GetEnv() AppEnv {
	cloneEnv := internalEnv
	return cloneEnv
}
