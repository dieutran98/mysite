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

func generateTestData(ctx context.Context, tx boil.ContextTransactor) error {
	userAccount := pgmodel.UserAccount{
		UserName:  "userName",
		Password:  "password",
		IsActive:  true,
		IsDeleted: false,
	}

	return userAccount.Insert(ctx, tx, boil.Infer())
}

func TestGetUserAccountByUserName(t *testing.T) {
	t.Parallel()
	require.NoError(t, database.SetupDatabase())
	repo := NewRepo()
	ctx := dbtest.SetTestTransactionCtx(context.Background())

	{ // found user
		var userAccount *pgmodel.UserAccount
		err := database.NewBoilerTransaction(ctx, func(ctx context.Context, tx boil.ContextTransactor) error {
			if err := generateTestData(ctx, tx); err != nil {
				return errors.Wrap(err, "failed generate data")
			}

			var err error
			userAccount, err = repo.GetUserAccountByUserName(ctx, tx, "userName")
			if err != nil {
				return errors.Wrap(err, "failed get userAccount")
			}

			return nil
		})

		require.NoError(t, err)
		require.Equal(t, "userName", userAccount.UserName)
		require.Equal(t, "password", userAccount.Password)
	}
	{ // not found user
		var userAccount *pgmodel.UserAccount
		err := database.NewBoilerTransaction(context.Background(), func(ctx context.Context, tx boil.ContextTransactor) error {

			var err error
			userAccount, err = repo.GetUserAccountByUserName(ctx, tx, "userName1")
			if err != nil {
				return errors.Wrap(err, "failed get userAccount")
			}

			return nil
		})
		require.NoError(t, err)
		require.Nil(t, userAccount)
	}

}

func TestGetActiveUserAccountByName(t *testing.T) {
	t.Parallel()
	require.NoError(t, database.SetupDatabase())
	repo := NewRepo()
	ctx := dbtest.SetTestTransactionCtx(context.Background())

	{ // found user
		var userAccount *pgmodel.UserAccount
		err := database.NewBoilerTransaction(ctx, func(ctx context.Context, tx boil.ContextTransactor) error {
			if err := generateTestData(ctx, tx); err != nil {
				return errors.Wrap(err, "failed generate data")
			}

			var err error
			userAccount, err = repo.GetActiveUserAccountByName(ctx, tx, "userName")
			if err != nil {
				return errors.Wrap(err, "failed get userAccount")
			}

			return nil
		})

		require.NoError(t, err)
		require.Equal(t, "userName", userAccount.UserName)
		require.Equal(t, "password", userAccount.Password)
	}
	{ // not found user
		var userAccount *pgmodel.UserAccount
		err := database.NewBoilerTransaction(context.Background(), func(ctx context.Context, tx boil.ContextTransactor) error {

			var err error
			userAccount, err = repo.GetActiveUserAccountByName(ctx, tx, "userName1")
			if err != nil {
				return errors.Wrap(err, "failed get userAccount")
			}

			return nil
		})
		require.Error(t, err)
		require.Nil(t, userAccount)
	}

}

func TestGetActiveUserAccountById(t *testing.T) {
	t.Parallel()
	require.NoError(t, database.SetupDatabase())
	repo := NewRepo()
	ctx := dbtest.SetTestTransactionCtx(context.Background())

	{ // found user
		var userAccount *pgmodel.UserAccount
		err := database.NewBoilerTransaction(ctx, func(ctx context.Context, tx boil.ContextTransactor) error {
			testUser := pgmodel.UserAccount{
				UserName: "userName",
				Password: "password",
				IsActive: true,
			}
			testUser.Insert(ctx, tx, boil.Infer())

			var err error
			userAccount, err = repo.GetActiveUserAccountById(ctx, tx, testUser.ID)
			if err != nil {
				return errors.Wrap(err, "failed get userAccount")
			}
			return nil
		})

		require.NoError(t, err)
		require.Equal(t, "userName", userAccount.UserName)
		require.Equal(t, "password", userAccount.Password)
	}
	{ // not found user
		var userAccount *pgmodel.UserAccount
		err := database.NewBoilerTransaction(context.Background(), func(ctx context.Context, tx boil.ContextTransactor) error {

			var err error
			userAccount, err = repo.GetActiveUserAccountById(ctx, tx, "userName1")
			if err != nil {
				return errors.Wrap(err, "failed get userAccount")
			}

			return nil
		})
		require.Error(t, err)
		require.Nil(t, userAccount)
	}

}
