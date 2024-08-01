package internal

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"mysite/dtos"
	"mysite/entities"
	"mysite/pkgs/database"
	"mysite/repositories/useraccountrepo"
	"mysite/testing/dbtest"
	"mysite/testing/mocking/pkgmock"
	"mysite/testing/mocking/repomock"
	"mysite/utils/ptrconv"

	"github.com/stretchr/testify/require"
	"github.com/volatiletech/sqlboiler/v4/boil"
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
		userAccountMock := &repomock.UserAccountRepoMock{}
		userAccountMock.GetUserAccountByUserNameFunc = func(ctx context.Context, tx boil.ContextTransactor, userName string) (*entities.UserAccount, error) {
			return &entities.UserAccount{IsActive: false}, nil
		}
		userAccountMock.ActiveUserFunc = func(ctx context.Context, tx boil.ContextTransactor, pgUser entities.UserAccount) error { return nil }
		userAccountMock.InsertFunc = func(ctx context.Context, tx boil.ContextTransactor, user *entities.UserAccount) error { return nil }

		authMock := &pkgmock.AuthServiceMock{}
		authMock.HashPasswordFunc = func(password string) (string, error) { return "token", nil }
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

		authMock := &pkgmock.AuthServiceMock{}
		authMock.HashPasswordFunc = func(password string) (string, error) { return "token", nil }
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
		userAccountMock := &repomock.UserAccountRepoMock{}
		userAccountMock.ActiveUserFunc = func(ctx context.Context, tx boil.ContextTransactor, pgUser entities.UserAccount) error {
			return errors.New("failed active user")
		}
		userAccountMock.GetUserAccountByUserNameFunc = func(ctx context.Context, tx boil.ContextTransactor, userName string) (*entities.UserAccount, error) {
			return &entities.UserAccount{IsActive: false}, nil
		}

		authMock := &pkgmock.AuthServiceMock{}
		authMock.HashPasswordFunc = func(password string) (string, error) { return "token", nil }
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
		userAccountMock := &repomock.UserAccountRepoMock{}
		userAccountMock.ActiveUserFunc = func(ctx context.Context, tx boil.ContextTransactor, pgUser entities.UserAccount) error { return nil }
		userAccountMock.InsertFunc = func(ctx context.Context, tx boil.ContextTransactor, user *entities.UserAccount) error {
			return errors.New("Insert error failed")
		}
		userAccountMock.GetUserAccountByUserNameFunc = func(ctx context.Context, tx boil.ContextTransactor, userName string) (*entities.UserAccount, error) {
			return nil, nil
		}

		authMock := &pkgmock.AuthServiceMock{}
		authMock.HashPasswordFunc = func(password string) (string, error) { return "token", nil }
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
		userAccountMock := &repomock.UserAccountRepoMock{}
		userAccountMock.GetUserAccountByUserNameFunc = func(ctx context.Context, tx boil.ContextTransactor, userName string) (*entities.UserAccount, error) {
			return &entities.UserAccount{IsActive: true, IsDeleted: false}, nil
		}

		authMock := &pkgmock.AuthServiceMock{}
		authMock.HashPasswordFunc = func(password string) (string, error) { return "token", nil }
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
		userAccountMock := &repomock.UserAccountRepoMock{}
		userAccountMock.GetUserAccountByUserNameFunc = func(ctx context.Context, tx boil.ContextTransactor, userName string) (*entities.UserAccount, error) {
			return nil, errors.New("get user account failed")
		}

		authMock := &pkgmock.AuthServiceMock{}
		authMock.HashPasswordFunc = func(password string) (string, error) { return "token", nil }
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
		userAccountMock := &repomock.UserAccountRepoMock{}
		userAccountMock.ActiveUserFunc = func(ctx context.Context, tx boil.ContextTransactor, pgUser entities.UserAccount) error { return nil }
		userAccountMock.InsertFunc = func(ctx context.Context, tx boil.ContextTransactor, user *entities.UserAccount) error { return nil }
		userAccountMock.GetUserAccountByUserNameFunc = func(ctx context.Context, tx boil.ContextTransactor, userName string) (*entities.UserAccount, error) {
			return &entities.UserAccount{IsActive: true, IsDeleted: false}, nil
		}

		authMock := &pkgmock.AuthServiceMock{}
		authMock.HashPasswordFunc = func(password string) (string, error) { return "token", nil }
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
		result, err := NewParams(dtos.RegisterRequest{
			Password: "secret",
			UserName: "test@gmail.com",
		})
		require.NoError(t, err)
		require.NotNil(t, result)
		require.Equal(t, "secret", result.Password)
		require.Equal(t, "test@gmail.com", result.UserName)
	}
	{ // empty data
		result, err := NewParams(dtos.RegisterRequest{})
		require.NoError(t, err)
		require.NotNil(t, result)
	}
	{ // success, has only passowrd
		result, err := NewParams(dtos.RegisterRequest{
			Password: "secret",
		})
		require.NoError(t, err)
		require.NotNil(t, result)
		require.Equal(t, "secret", result.Password)
		require.Empty(t, result.UserName)
	}
}
