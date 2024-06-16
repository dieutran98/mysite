package database

import (
	"context"
	"mysite/models/pgmodel"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

func TestNewBoilerTransaction(t *testing.T) {
	assert.NoError(t, SetupDatabase())

	assert.NoError(t, NewBoilerTransaction(context.Background(), testExecuteQueries))

}

func testExecuteQueries(ctx context.Context, tx boil.ContextTransactor) error {
	user := pgmodel.UserAccount{
		UserName: "test",
		Password: "test",
	}
	return user.Insert(ctx, tx, boil.Infer())
}
