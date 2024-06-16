package migration

import (
	"embed"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/pkg/errors"
)

//go:embed ddl/*
var f embed.FS

func Migrate(dbUrl string) error {
	srcDriver, err := iofs.New(f, "ddl")
	if err != nil {
		return errors.Wrap(err, "failed to create src driver")
	}

	m, err := migrate.NewWithSourceInstance("iofs", srcDriver, dbUrl)
	if err != nil {
		return errors.Wrap(err, "failed create migrate instance")
	}

	if err := m.Up(); err != nil {
		return errors.Wrap(err, "failed migrate")
	}
	return nil
}
