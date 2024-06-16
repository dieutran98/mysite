package database

import (
	"fmt"
	"mysite/testing/database"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	pool, resource, err := database.SetupDatabaseForTesting()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	defer func() {
		if err := database.PurgeResource(pool, resource); err != nil {
			fmt.Println("failed to purge resource")
		}
		Close()
	}()
	m.Run()
}

func TestSetupDatabase(t *testing.T) {
	assert.NoError(t, SetupDatabase())
}
