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

var (
	internalEnv appEnv
)

func ReadEnv(config Configure) error {
	if err := setConfigFile(config); err != nil {
		return errors.Wrap(err, "failed set config file")
	}

	if err := mappingStruct(config); err != nil {
		return errors.Wrap(err, "failed to mapping env to struct")
	}

	return nil
}

func GetEnv() appEnv {
	cloneEnv := internalEnv
	return cloneEnv
}
