package useraccount

import (
	"fmt"
	"mysite/pkgs/database"
	databasetesting "mysite/testing/database"
	"testing"
)

func TestMain(m *testing.M) {
	pool, resource, err := databasetesting.SetupDatabaseForTesting()
	if err != nil {
		return
	}

	defer func() {
		database.Close()
		if err := databasetesting.PurgeResource(pool, resource); err != nil {
			fmt.Println("failed to purge resource")
		}
	}()
	m.Run()
}
