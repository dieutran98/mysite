package internal

import (
	"context"
	"fmt"
	"testing"
	"time"

	"mysite/dtos"
	"mysite/entities"
	"mysite/pkgs/auth"
	"mysite/pkgs/database"
	"mysite/testing/dbtest"
	"mysite/testing/mocking/pkgmock"
	"mysite/testing/mocking/repomock"
	"mysite/utils/httputil"

	"github.com/golang-jwt/jwt/v5"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
	"github.com/volatiletech/sqlboiler/v4/boil"
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
		repoMock := &repomock.UserAccountRepoMock{}
		repoMock.GetActiveUserAccountByNameFunc = func(ctx context.Context, tx boil.ContextTransactor, userName string) (*entities.UserAccount, error) {
			return &entities.UserAccount{
				ID: 1,
			}, nil
		}

		authMock := &pkgmock.AuthServiceMock{}
		authMock.ComparePasswordAndHashFunc = func(password, encodedHash string) (bool, error) { return true, nil }
		authMock.CreateTokenFunc = func(claims auth.Claims, signingKey []byte) (string, error) { return "token", nil }
		authMock.NewClaimsFunc = func(userId int, expiredTime time.Time) auth.Claims {
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
		repoMock := &repomock.UserAccountRepoMock{}
		repoMock.GetActiveUserAccountByNameFunc = func(ctx context.Context, tx boil.ContextTransactor, userName string) (*entities.UserAccount, error) {
			return nil, errors.New("failed get user")
		}

		authMock := &pkgmock.AuthServiceMock{}
		authMock.ComparePasswordAndHashFunc = func(password, encodedHash string) (bool, error) { return true, nil }
		authMock.CreateTokenFunc = func(claims auth.Claims, signingKey []byte) (string, error) { return "token", nil }
		authMock.NewClaimsFunc = func(userId int, expiredTime time.Time) auth.Claims {
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
		repoMock := &repomock.UserAccountRepoMock{}
		repoMock.GetActiveUserAccountByNameFunc = func(ctx context.Context, tx boil.ContextTransactor, userName string) (*entities.UserAccount, error) {
			return nil, nil
		}

		authMock := &pkgmock.AuthServiceMock{}
		authMock.ComparePasswordAndHashFunc = func(password, encodedHash string) (bool, error) { return true, nil }
		authMock.CreateTokenFunc = func(claims auth.Claims, signingKey []byte) (string, error) { return "token", nil }
		authMock.NewClaimsFunc = func(userId int, expiredTime time.Time) auth.Claims {
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
		repoMock := &repomock.UserAccountRepoMock{}
		repoMock.GetActiveUserAccountByNameFunc = func(ctx context.Context, tx boil.ContextTransactor, userName string) (*entities.UserAccount, error) {
			return &entities.UserAccount{
				ID: 1,
			}, nil
		}

		authMock := &pkgmock.AuthServiceMock{}
		authMock.ComparePasswordAndHashFunc = func(password, encodedHash string) (bool, error) { return false, nil }
		authMock.CreateTokenFunc = func(claims auth.Claims, signingKey []byte) (string, error) { return "token", nil }
		authMock.NewClaimsFunc = func(userId int, expiredTime time.Time) auth.Claims {
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
		repoMock := &repomock.UserAccountRepoMock{}
		repoMock.GetActiveUserAccountByNameFunc = func(ctx context.Context, tx boil.ContextTransactor, userName string) (*entities.UserAccount, error) {
			return &entities.UserAccount{
				ID: 1,
			}, nil
		}

		authMock := &pkgmock.AuthServiceMock{}
		authMock.ComparePasswordAndHashFunc = func(password, encodedHash string) (bool, error) {
			return false, errors.New("check pass error")
		}
		authMock.CreateTokenFunc = func(claims auth.Claims, signingKey []byte) (string, error) { return "token", nil }
		authMock.NewClaimsFunc = func(userId int, expiredTime time.Time) auth.Claims {
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
		repoMock := &repomock.UserAccountRepoMock{}
		repoMock.GetActiveUserAccountByNameFunc = func(ctx context.Context, tx boil.ContextTransactor, userName string) (*entities.UserAccount, error) {
			return &entities.UserAccount{
				ID: 1,
			}, nil
		}

		authMock := &pkgmock.AuthServiceMock{}
		authMock.ComparePasswordAndHashFunc = func(password, encodedHash string) (bool, error) { return true, nil }
		authMock.CreateTokenFunc = func(claims auth.Claims, signingKey []byte) (string, error) {
			return "", errors.New("create token failed")
		}
		authMock.NewClaimsFunc = func(userId int, expiredTime time.Time) auth.Claims {
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
		repoMock := &repomock.UserAccountRepoMock{}
		repoMock.GetActiveUserAccountByNameFunc = func(ctx context.Context, tx boil.ContextTransactor, userName string) (*entities.UserAccount, error) {
			return &entities.UserAccount{
				ID: 1,
			}, nil
		}

		authMock := &pkgmock.AuthServiceMock{}
		authMock.ComparePasswordAndHashFunc = func(password, encodedHash string) (bool, error) { return true, nil }
		authMock.CreateTokenFunc = func(claims auth.Claims, signingKey []byte) (string, error) {
			return "", errors.New("create token failed")
		}
		authMock.NewClaimsFunc = func(userId int, expiredTime time.Time) auth.Claims {
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
		result, err := NewParams(dtos.LoginRequest{
			Password: "password",
			UserName: "userName",
		})
		require.NoError(t, err)
		require.NotNil(t, result)
		require.Equal(t, "password", result.Password)
		require.Equal(t, "userName", result.UserName)
	}
}
