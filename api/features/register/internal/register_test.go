package internal

import (
	"context"
	"fmt"
	"testing"

	"mysite/models/model"
	"mysite/models/pgmodel"
	"mysite/pkgs/database"
	"mysite/repositories/useraccount"
	dbtest "mysite/testing/database"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

type mockService struct {
	useraccount.UserAccountRepo
	GetUserAccountByUserNameFunc func() (*pgmodel.UserAccount, error)
	ActiveUserFunc               func() error
	InsertFunc                   func() error
}

func (m mockService) GetUserAccountByUserName(ctx context.Context, tx boil.ContextTransactor, userName string) (*pgmodel.UserAccount, error) {
	return m.GetUserAccountByUserNameFunc()
}

func (m mockService) Insert(ctx context.Context, tx boil.ContextTransactor, user pgmodel.UserAccount) error {
	return m.InsertFunc()
}

func (m mockService) ActiveUser(ctx context.Context, tx boil.ContextTransactor, pgUser pgmodel.UserAccount) error {
	return m.ActiveUserFunc()
}

func TestMain(m *testing.M) {
	pool, resource, err := dbtest.SetupDatabaseForTesting()
	if err != nil {
		return
	}

	defer func() {
		database.Close()
		if err := dbtest.PurgeResource(pool, resource); err != nil {
			fmt.Println("failed to purge resource")
		}
	}()
	m.Run()
}

func TestRegister(t *testing.T) {
	require.NoError(t, database.SetupDatabase())
	ctx := dbtest.SetTestTransactionCtx(context.Background())

	{ // register success, active user
		svc := service{
			repo: mockService{
				ActiveUserFunc:               func() error { return nil },
				InsertFunc:                   func() error { return nil },
				GetUserAccountByUserNameFunc: func() (*pgmodel.UserAccount, error) { return &pgmodel.UserAccount{IsActive: false}, nil },
			},
			req: RegisterRequest{
				Password: "secret",
				UserName: "test@gamil.com",
			},
		}
		require.NoError(t, svc.Register(ctx))
	}
	{ // register success, create user
		svc := service{
			repo: mockService{
				ActiveUserFunc:               func() error { return nil },
				InsertFunc:                   func() error { return nil },
				GetUserAccountByUserNameFunc: func() (*pgmodel.UserAccount, error) { return nil, nil },
			},
			req: RegisterRequest{
				Password: "secret",
				UserName: "test@gamil.com",
			},
		}
		require.NoError(t, svc.Register(ctx))
	}

	{ // register failed, active user failed
		svc := service{
			repo: mockService{
				ActiveUserFunc:               func() error { return errors.New("active user failed") },
				InsertFunc:                   func() error { return nil },
				GetUserAccountByUserNameFunc: func() (*pgmodel.UserAccount, error) { return &pgmodel.UserAccount{IsActive: false}, nil },
			},
			req: RegisterRequest{
				Password: "secret",
				UserName: "test@gamil.com",
			},
		}
		require.Error(t, svc.Register(ctx))
	}

	{ // register failed, create user failed
		svc := service{
			repo: mockService{
				ActiveUserFunc:               func() error { return nil },
				InsertFunc:                   func() error { return errors.New("insert error failed") },
				GetUserAccountByUserNameFunc: func() (*pgmodel.UserAccount, error) { return nil, nil },
			},
			req: RegisterRequest{
				Password: "secret",
				UserName: "test@gamil.com",
			},
		}
		require.Error(t, svc.Register(ctx))
	}

	{ // register failed, user exist
		svc := service{
			repo: mockService{
				ActiveUserFunc: func() error { return nil },
				InsertFunc:     func() error { return nil },
				GetUserAccountByUserNameFunc: func() (*pgmodel.UserAccount, error) {
					return &pgmodel.UserAccount{IsActive: true, IsDeleted: false}, nil
				},
			},
			req: RegisterRequest{
				Password: "secret",
				UserName: "test@gamil.com",
			},
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
