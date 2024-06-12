package env

import (
	"github.com/pkg/errors"
)

type appEnv struct {
	Database database `json:"database"`
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

type configure interface {
	setConfigFile() error
	mappingStruct() error
	setDefault() error
}

var (
	internalEnv appEnv
)

var newConfigure = func() configure {
	return newViperConfig()
}

func ReadEnv() error {
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

func GetEnv() appEnv {
	cloneEnv := internalEnv
	return cloneEnv
}
