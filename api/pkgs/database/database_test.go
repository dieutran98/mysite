package database

import (
	"fmt"
	"log/slog"
	"mysite/pkgs/logger"
	"mysite/testing/dbtest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	pool, resource, err := dbtest.SetupDatabaseForTesting()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	defer func() {
		if err := Close(); err != nil {
			slog.Error("failed to close database", logger.AttrError(err))
		}
		if err := dbtest.PurgeResource(pool, resource); err != nil {
			slog.Error("failed to purge resource", logger.AttrError(err))
		}
	}()
	m.Run()
}

func TestSetupDatabase(t *testing.T) {
	assert.NoError(t, SetupDatabase())
}
