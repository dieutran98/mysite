package useraccountrepo

import (
	"context"
	"mysite/models/pgmodel"
	"mysite/pkgs/database"
	dbtest "mysite/testing/dbtest"
	"testing"

	"github.com/friendsofgo/errors"
	"github.com/stretchr/testify/require"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

func TestActiveUser(t *testing.T) {
	t.Parallel()
	require.NoError(t, database.SetupDatabase())
	repo := NewRepo()
	ctx := dbtest.SetTestTransactionCtx(context.Background())

	{ // success active user
		userAccount := pgmodel.UserAccount{
			UserName:  "userName",
			Password:  "password",
			IsActive:  false,
			IsDeleted: false,
		}

		var result *pgmodel.UserAccount
		err := database.NewBoilerTransaction(ctx, func(ctx context.Context, tx boil.ContextTransactor) error {
			if err := repo.Insert(ctx, tx, &userAccount); err != nil {
				return errors.Wrap(err, "failed insert userAccount")
			}

			var err error
			result, err = repo.GetUserAccountByUserName(ctx, tx, "userName")
			if err != nil {
				return errors.Wrap(err, "failed GetUserAccountByUserName")
			}

			if err := repo.ActiveUser(ctx, tx, *result); err != nil {
				return errors.Wrap(err, "failed to active user")
			}

			result, err = repo.GetUserAccountByUserName(ctx, tx, "userName")
			if err != nil {
				return errors.Wrap(err, "failed GetUserAccountByUserName")
			}

			return nil
		})

		require.NoError(t, err)
		require.Equal(t, "userName", result.UserName)
		require.Equal(t, "password", result.Password)
		require.True(t, result.IsActive)
		require.False(t, result.IsDeleted)
		require.Nil(t, result.DeletedAt.Ptr())

	}
}
