package internal

import (
	"context"
	"fmt"
	"testing"

	"mysite/models/model"
	"mysite/models/pgmodel"
	"mysite/pkgs/auth"
	"mysite/pkgs/database"
	"mysite/testing/dbtest"
	"mysite/testing/mocking/pkgs/authmock"
	"mysite/testing/mocking/repomock"
	"mysite/utils/httputil"

	"github.com/golang-jwt/jwt/v5"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
)

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
func TestLogin(t *testing.T) {
	t.Parallel()
	require.NoError(t, database.SetupDatabase())
	ctx := dbtest.SetTestTransactionCtx(context.Background())

	{ // login success
		repoMock := repomock.NewUserAccountMock()
		repoMock.GetUserAccountByUserNameFunc = func() (*pgmodel.UserAccount, error) {
			return &pgmodel.UserAccount{
				ID: "user-id",
			}, nil
		}

		authMock := authmock.NewMockService()
		authMock.ComparePasswordAndHashFunc = func() (match bool, err error) {
			return true, nil
		}
		authMock.CreateTokenFunc = func() (string, error) {
			return "token", nil
		}
		authMock.NewClaimsFunc = func() auth.Claims {
			return auth.Claims{
				RegisteredClaims: jwt.RegisteredClaims{
					Subject: "user-id",
				},
			}
		}

		svc := service{
			repo:    repoMock,
			authSvc: authMock,
			req: LoginRequest{
				Password: "password",
				UserName: "test@gmail.com",
			},
		}

		resp, err := svc.Login(ctx)
		require.NoError(t, err)
		require.NotNil(t, resp)
		require.Equal(t, "token", resp.AccessToken)
		require.Equal(t, "token", resp.RefreshToken)
	}
	{ // login failed, failed get user
		repoMock := repomock.NewUserAccountMock()
		repoMock.GetUserAccountByUserNameFunc = func() (*pgmodel.UserAccount, error) {
			return nil, errors.New("failed get user")
		}

		authMock := authmock.NewMockService()
		authMock.ComparePasswordAndHashFunc = func() (match bool, err error) {
			return true, nil
		}
		authMock.CreateTokenFunc = func() (string, error) {
			return "token", nil
		}
		authMock.NewClaimsFunc = func() auth.Claims {
			return auth.Claims{
				RegisteredClaims: jwt.RegisteredClaims{
					Subject: "user-id",
				},
			}
		}

		svc := service{
			repo:    repoMock,
			authSvc: authMock,
			req: LoginRequest{
				Password: "password",
				UserName: "test@gmail.com",
			},
		}

		resp, err := svc.Login(ctx)
		require.Error(t, err)
		require.ErrorIs(t, err, httputil.ErrUnauthorize)
		require.Nil(t, resp)
	}
	{ // login failed, failed get user
		repoMock := repomock.NewUserAccountMock()
		repoMock.GetUserAccountByUserNameFunc = func() (*pgmodel.UserAccount, error) {
			return nil, nil
		}

		authMock := authmock.NewMockService()
		authMock.ComparePasswordAndHashFunc = func() (match bool, err error) {
			return true, nil
		}
		authMock.CreateTokenFunc = func() (string, error) {
			return "token", nil
		}
		authMock.NewClaimsFunc = func() auth.Claims {
			return auth.Claims{
				RegisteredClaims: jwt.RegisteredClaims{
					Subject: "user-id",
				},
			}
		}

		svc := service{
			repo:    repoMock,
			authSvc: authMock,
			req: LoginRequest{
				Password: "password",
				UserName: "test@gmail.com",
			},
		}

		resp, err := svc.Login(ctx)
		require.Error(t, err)
		require.ErrorIs(t, err, httputil.ErrUnauthorize)
		require.Nil(t, resp)
	}
	{ // login failed, wrong password
		repoMock := repomock.NewUserAccountMock()
		repoMock.GetUserAccountByUserNameFunc = func() (*pgmodel.UserAccount, error) {
			return &pgmodel.UserAccount{
				ID: "user-id",
			}, nil
		}

		authMock := authmock.NewMockService()
		authMock.ComparePasswordAndHashFunc = func() (match bool, err error) {
			return false, nil
		}
		authMock.CreateTokenFunc = func() (string, error) {
			return "token", nil
		}
		authMock.NewClaimsFunc = func() auth.Claims {
			return auth.Claims{
				RegisteredClaims: jwt.RegisteredClaims{
					Subject: "user-id",
				},
			}
		}

		svc := service{
			repo:    repoMock,
			authSvc: authMock,
			req: LoginRequest{
				Password: "password",
				UserName: "test@gmail.com",
			},
		}

		resp, err := svc.Login(ctx)
		require.Error(t, err)
		require.ErrorIs(t, err, httputil.ErrUnauthorize)
		require.Nil(t, resp)
	}
	{ // login failed, check password error
		repoMock := repomock.NewUserAccountMock()
		repoMock.GetUserAccountByUserNameFunc = func() (*pgmodel.UserAccount, error) {
			return &pgmodel.UserAccount{
				ID: "user-id",
			}, nil
		}

		authMock := authmock.NewMockService()
		authMock.ComparePasswordAndHashFunc = func() (match bool, err error) {
			return false, errors.New("check pass error")
		}
		authMock.CreateTokenFunc = func() (string, error) {
			return "token", nil
		}
		authMock.NewClaimsFunc = func() auth.Claims {
			return auth.Claims{
				RegisteredClaims: jwt.RegisteredClaims{
					Subject: "user-id",
				},
			}
		}

		svc := service{
			repo:    repoMock,
			authSvc: authMock,
			req: LoginRequest{
				Password: "password",
				UserName: "test@gmail.com",
			},
		}

		resp, err := svc.Login(ctx)
		require.Error(t, err)
		require.ErrorIs(t, err, httputil.ErrUnauthorize)
		require.Nil(t, resp)
	}
	{ // login failed, create token failed
		repoMock := repomock.NewUserAccountMock()
		repoMock.GetUserAccountByUserNameFunc = func() (*pgmodel.UserAccount, error) {
			return &pgmodel.UserAccount{
				ID: "user-id",
			}, nil
		}

		authMock := authmock.NewMockService()
		authMock.ComparePasswordAndHashFunc = func() (match bool, err error) {
			return true, nil
		}
		authMock.CreateTokenFunc = func() (string, error) {
			return "", errors.New("create token failed")
		}
		authMock.NewClaimsFunc = func() auth.Claims {
			return auth.Claims{
				RegisteredClaims: jwt.RegisteredClaims{
					Subject: "user-id",
				},
			}
		}

		svc := service{
			repo:    repoMock,
			authSvc: authMock,
			req: LoginRequest{
				Password: "password",
				UserName: "test@gmail.com",
			},
		}

		resp, err := svc.Login(ctx)
		require.Error(t, err)
		require.Nil(t, resp)
	}

	{ // login failed, request body wrong
		repoMock := repomock.NewUserAccountMock()
		repoMock.GetUserAccountByUserNameFunc = func() (*pgmodel.UserAccount, error) {
			return &pgmodel.UserAccount{
				ID: "user-id",
			}, nil
		}

		authMock := authmock.NewMockService()
		authMock.ComparePasswordAndHashFunc = func() (match bool, err error) {
			return true, nil
		}
		authMock.CreateTokenFunc = func() (string, error) {
			return "", errors.New("create token failed")
		}
		authMock.NewClaimsFunc = func() auth.Claims {
			return auth.Claims{
				RegisteredClaims: jwt.RegisteredClaims{
					Subject: "user-id",
				},
			}
		}

		svc := service{
			repo:    repoMock,
			authSvc: authMock,
			req: LoginRequest{
				Password: "",
				UserName: "test",
			},
		}

		resp, err := svc.Login(ctx)
		require.Error(t, err)
		require.ErrorIs(t, err, httputil.ErrInvalidRequest)
		require.Nil(t, resp)
	}
}

func TestNewParams(t *testing.T) {
	require.NoError(t, database.SetupDatabase())
	{ // success
		result, err := NewParams(model.LoginRequest{
			Password: "password",
			UserName: "userName",
		})
		require.NoError(t, err)
		require.NotNil(t, result)
		require.Equal(t, "password", result.Password)
		require.Equal(t, "userName", result.UserName)
	}
}
