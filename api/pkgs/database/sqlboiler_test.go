package database

import (
	"context"
	"log/slog"
	"mysite/entities"
	"mysite/pkgs/logger"
	"os"
	"testing"

	"github.com/friendsofgo/errors"
	"github.com/stretchr/testify/assert"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

func TestNewBoilerTransaction(t *testing.T) {
	logger.SetLogger(os.Stdout)
	logger.SetLogLevel(slog.LevelDebug)
	assert.NoError(t, SetupDatabase())

	{
		// success
		assert.NoError(t, NewBoilerTransaction(context.Background(), testSuccessQueries))
	}

	{
		// failed, expect roll back
		assert.Error(t, NewBoilerTransaction(context.Background(), testFailedQueries))

	}

}

func testSuccessQueries(ctx context.Context, tx boil.ContextTransactor) error {
	user := entities.UserAccount{
		UserName: "test",
		Password: "test",
	}
	if err := user.Insert(ctx, tx, boil.Infer()); err != nil {
		return errors.Wrap(err, "failed insert user data")
	}

	userTest, err := entities.UserAccounts(entities.UserAccountWhere.ID.EQ(user.ID)).One(ctx, tx)
	if err != nil {
		return errors.Wrap(err, "failed get user data inserted")
	}

	if userTest.UserName != "test" || userTest.Password != "test" {
		return errors.New("invalid data inserted")
	}
	return nil
}

func testFailedQueries(ctx context.Context, tx boil.ContextTransactor) error {
	return errors.New("some error")
}
