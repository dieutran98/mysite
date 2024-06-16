package database

import (
	"context"
	"mysite/models/pgmodel"
	"testing"

	"github.com/friendsofgo/errors"
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
	if err := user.Insert(ctx, tx, boil.Infer()); err != nil {
		return errors.Wrap(err, "failed insert user data")
	}

	userTest, err := pgmodel.UserAccounts(pgmodel.UserAccountWhere.ID.EQ(user.ID)).One(ctx, tx)
	if err != nil {
		return errors.Wrap(err, "failed get user data inserted")
	}

	if userTest.UserName != "test" || userTest.Password != "test" {
		return errors.New("invalid data inserted")
	}
	return nil
}
