package internal

import (
	"context"
	"fmt"
	"testing"

	"mysite/models/model"
	"mysite/models/pgmodel"
	"mysite/pkgs/database"
	"mysite/repositories/useraccountrepo"
	"mysite/testing/dbtest"
	"mysite/testing/mocking/pkgs/authmock"
	"mysite/testing/mocking/repomock"
	"mysite/utils/ptrconv"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {
	pool, resource, err := dbtest.SetupDatabaseForTesting()
	if err != nil {
		return
	}

	defer func() {
		if r := recover(); r != nil {
			database.Close()
			if err := dbtest.PurgeResource(pool, resource); err != nil {
				fmt.Println("failed to purge resource")
			}
		}
		database.Close()
		if err := dbtest.PurgeResource(pool, resource); err != nil {
			fmt.Println("failed to purge resource")
		}
	}()
	m.Run()
}

func TestRegister(t *testing.T) {
	t.Parallel()
	require.NoError(t, database.SetupDatabase())
	ctx := dbtest.SetTestTransactionCtx(context.Background())

	{ // register success, active user
		userAccountMock := repomock.NewUserAccountMock()
		userAccountMock.ActiveUserFunc = func() error { return nil }
		userAccountMock.InsertFunc = func() error { return nil }
		userAccountMock.GetUserAccountByUserNameFunc = func() (*pgmodel.UserAccount, error) { return &pgmodel.UserAccount{IsActive: false}, nil }

		authMock := authmock.NewMockService()
		authMock.HashPasswordFunc = func() (string, error) { return "token", nil }
		svc := service{
			repo: userAccountMock,
			req: RegisterRequest{
				Password: "secret",
				UserName: "test@gamil.com",
			},
			authSvc: authMock,
		}
		require.NoError(t, svc.Register(ctx))
	}
	{ // register success, create user

		authMock := authmock.NewMockService()
		authMock.HashPasswordFunc = func() (string, error) { return "token", nil }
		svc := service{
			repo: useraccountrepo.NewRepo(),
			req: RegisterRequest{
				Password: "secret",
				UserName: "test@gamil.com",
				Name:     ptrconv.String("testing"),
			},
			authSvc: authMock,
		}
		require.NoError(t, svc.Register(ctx))
	}

	{ // register failed, active user failed
		userAccountMock := repomock.NewUserAccountMock()
		userAccountMock.ActiveUserFunc = func() error { return errors.New("active user failed") }
		userAccountMock.InsertFunc = func() error { return nil }
		userAccountMock.GetUserAccountByUserNameFunc = func() (*pgmodel.UserAccount, error) { return &pgmodel.UserAccount{IsActive: false}, nil }

		authMock := authmock.NewMockService()
		authMock.HashPasswordFunc = func() (string, error) { return "token", nil }
		svc := service{
			repo: userAccountMock,
			req: RegisterRequest{
				Password: "secret",
				UserName: "test@gamil.com",
			},
			authSvc: authMock,
		}
		require.Error(t, svc.Register(ctx))
	}

	{ // register failed, create user failed
		userAccountMock := repomock.NewUserAccountMock()
		userAccountMock.ActiveUserFunc = func() error { return nil }
		userAccountMock.InsertFunc = func() error { return errors.New("insert error failed") }
		userAccountMock.GetUserAccountByUserNameFunc = func() (*pgmodel.UserAccount, error) { return nil, nil }

		authMock := authmock.NewMockService()
		authMock.HashPasswordFunc = func() (string, error) { return "token", nil }
		svc := service{
			repo: userAccountMock,
			req: RegisterRequest{
				Password: "secret",
				UserName: "test@gamil.com",
			},
			authSvc: authMock,
		}
		require.Error(t, svc.Register(ctx))
	}

	{ // register failed, user exist
		userAccountMock := repomock.NewUserAccountMock()
		userAccountMock.ActiveUserFunc = func() error { return nil }
		userAccountMock.InsertFunc = func() error { return nil }
		userAccountMock.GetUserAccountByUserNameFunc = func() (*pgmodel.UserAccount, error) {
			return &pgmodel.UserAccount{IsActive: true, IsDeleted: false}, nil
		}

		authMock := authmock.NewMockService()
		authMock.HashPasswordFunc = func() (string, error) { return "token", nil }
		svc := service{
			repo: userAccountMock,
			req: RegisterRequest{
				Password: "secret",
				UserName: "test@gamil.com",
			},
			authSvc: authMock,
		}
		require.Error(t, svc.Register(ctx))
	}
	{ // register failed, get user failed
		userAccountMock := repomock.NewUserAccountMock()
		userAccountMock.ActiveUserFunc = func() error { return nil }
		userAccountMock.InsertFunc = func() error { return nil }
		userAccountMock.GetUserAccountByUserNameFunc = func() (*pgmodel.UserAccount, error) {
			return nil, errors.New("failed get user")
		}

		authMock := authmock.NewMockService()
		authMock.HashPasswordFunc = func() (string, error) { return "token", nil }
		svc := service{
			repo: userAccountMock,
			req: RegisterRequest{
				Password: "secret",
				UserName: "test@gamil.com",
			},
			authSvc: authMock,
		}
		require.Error(t, svc.Register(ctx))
	}

	{ // register failed, validate failed
		userAccountMock := repomock.NewUserAccountMock()
		userAccountMock.ActiveUserFunc = func() error { return nil }
		userAccountMock.InsertFunc = func() error { return nil }
		userAccountMock.GetUserAccountByUserNameFunc = func() (*pgmodel.UserAccount, error) {
			return nil, errors.New("failed get user")
		}

		authMock := authmock.NewMockService()
		authMock.HashPasswordFunc = func() (string, error) { return "token", nil }
		svc := service{
			repo: userAccountMock,
			req: RegisterRequest{
				Password: "",
				UserName: "test",
			},
			authSvc: authMock,
		}
		require.Error(t, svc.Register(ctx))
	}
}

func TestNewParams(t *testing.T) {
	{ // success
		result, err := NewParams(model.RegisterRequest{
			Password: "secret",
			UserName: "test@gmail.com",
		})
		require.NoError(t, err)
		require.NotNil(t, result)
		require.Equal(t, "secret", result.Password)
		require.Equal(t, "test@gmail.com", result.UserName)
	}
	{ // empty data
		result, err := NewParams(model.RegisterRequest{})
		require.NoError(t, err)
		require.NotNil(t, result)
	}
	{ // success, has only passowrd
		result, err := NewParams(model.RegisterRequest{
			Password: "secret",
		})
		require.NoError(t, err)
		require.NotNil(t, result)
		require.Equal(t, "secret", result.Password)
		require.Empty(t, result.UserName)
	}
}
