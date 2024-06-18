package useraccount

import (
	"context"
	"mysite/models/pgmodel"
	"mysite/pkgs/database"
	dbtest "mysite/testing/database"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

func TestInsert(t *testing.T) {
	t.Parallel()
	require.NoError(t, database.SetupDatabase())
	repo := New()
	ctx := dbtest.SetTestTransactionCtx(context.Background())

	{ // insert success
		userAccount := pgmodel.UserAccount{
			UserName:  "userName",
			Password:  "password",
			IsActive:  true,
			IsDeleted: false,
		}
		var result *pgmodel.UserAccount
		err := database.NewBoilerTransaction(ctx, func(ctx context.Context, tx boil.ContextTransactor) error {
			if err := repo.Insert(ctx, tx, userAccount); err != nil {
				return errors.Wrap(err, "failed insert userAccount")
			}

			var err error
			result, err = repo.GetUserAccountByUserName(ctx, tx, userAccount.UserName)
			if err != nil {
				return errors.Wrap(err, "failed GetUserAccountByUserName")
			}

			return nil
		})

		require.NoError(t, err)
		require.Equal(t, "userName", result.UserName)
		require.Equal(t, "password", result.Password)

	}
}
