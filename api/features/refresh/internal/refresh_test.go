package internal

import (
	"context"
	"fmt"
	"mysite/models/model"
	"mysite/models/pgmodel"
	"mysite/pkgs/auth"
	"mysite/pkgs/database"
	"mysite/testing/dbtest"
	"mysite/testing/mocking/pkgmock"
	"mysite/testing/mocking/repomock"
	"testing"
	"time"

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

func TestRefreshToken(t *testing.T) {
	t.Parallel()
	require.NoError(t, database.SetupDatabase())
	ctx := dbtest.SetTestTransactionCtx(context.Background())

	{ // refresh success
		repoMock := &repomock.UserAccountRepoMock{}
		repoMock.GetActiveUserAccountByIdFunc = func(ctx context.Context, tx boil.ContextTransactor, userId int) (*pgmodel.UserAccount, error) {
			return &pgmodel.UserAccount{
				ID: 1,
			}, nil
		}

		authMock := &pkgmock.AuthServiceMock{}
		authMock.ParseTokenFunc = func(tokenStr string, signingKey []byte) (*auth.Claims, error) {
			return &auth.Claims{
				RegisteredClaims: jwt.RegisteredClaims{
					Subject: "1",
				},
			}, nil
		}
		authMock.CreateTokenFunc = func(claims auth.Claims, signingKey []byte) (string, error) { return "token", nil }
		authMock.NewClaimsFunc = func(userId int, expiredTime time.Time) auth.Claims {
			return auth.Claims{
				RegisteredClaims: jwt.RegisteredClaims{
					Subject: "1",
				},
			}
		}

		svc := service{
			repo:    repoMock,
			authSvc: authMock,
			req: RefreshRequest{
				RefreshToken: "refresh-token",
			},
		}

		resp, err := svc.RefreshToken(ctx)
		require.NoError(t, err)
		require.NotNil(t, resp)
		require.Equal(t, "token", resp.AccessToken)
	}
	{ // refresh failed, validate failed
		svc := service{
			req: RefreshRequest{
				RefreshToken: "",
			},
		}

		resp, err := svc.RefreshToken(ctx)
		require.Error(t, err)
		require.Nil(t, resp)
	}
	{ // refresh failed, parse token failed
		repoMock := &repomock.UserAccountRepoMock{}
		repoMock.GetActiveUserAccountByIdFunc = func(ctx context.Context, tx boil.ContextTransactor, userId int) (*pgmodel.UserAccount, error) {
			return &pgmodel.UserAccount{
				ID: 1,
			}, nil
		}

		authMock := &pkgmock.AuthServiceMock{}
		authMock.ParseTokenFunc = func(tokenStr string, signingKey []byte) (*auth.Claims, error) {
			return nil, errors.New("parse token failed")
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
			req: RefreshRequest{
				RefreshToken: "refresh-token",
			},
		}

		resp, err := svc.RefreshToken(ctx)
		require.Error(t, err)
		require.Nil(t, resp)
	}
	{ // refresh failed, get user failed
		repoMock := &repomock.UserAccountRepoMock{}
		repoMock.GetActiveUserAccountByIdFunc = func(ctx context.Context, tx boil.ContextTransactor, userId int) (*pgmodel.UserAccount, error) {
			return nil, errors.New("get user account failed")
		}
		authMock := &pkgmock.AuthServiceMock{}
		authMock.ParseTokenFunc = func(tokenStr string, signingKey []byte) (*auth.Claims, error) {
			return &auth.Claims{
				RegisteredClaims: jwt.RegisteredClaims{
					Subject: "1",
				},
			}, nil
		}
		authMock.CreateTokenFunc = func(claims auth.Claims, signingKey []byte) (string, error) {
			return "token", nil
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
			req: RefreshRequest{
				RefreshToken: "refresh-token",
			},
		}

		resp, err := svc.RefreshToken(ctx)
		require.Error(t, err)
		require.Nil(t, resp)
	}
	{ // refresh failed, createToken failed
		repoMock := &repomock.UserAccountRepoMock{}
		repoMock.GetActiveUserAccountByIdFunc = func(ctx context.Context, tx boil.ContextTransactor, userId int) (*pgmodel.UserAccount, error) {
			return &pgmodel.UserAccount{
				ID: 1,
			}, nil
		}
		authMock := &pkgmock.AuthServiceMock{}
		authMock.ParseTokenFunc = func(tokenStr string, signingKey []byte) (*auth.Claims, error) {
			return &auth.Claims{
				RegisteredClaims: jwt.RegisteredClaims{
					Subject: "1",
				},
			}, nil
		}
		authMock.CreateTokenFunc = func(claims auth.Claims, signingKey []byte) (string, error) {
			return "", errors.New("create token failed")
		}
		authMock.NewClaimsFunc = func(userId int, expiredTime time.Time) auth.Claims {
			return auth.Claims{
				RegisteredClaims: jwt.RegisteredClaims{
					Subject: "1",
				},
			}
		}

		svc := service{
			repo:    repoMock,
			authSvc: authMock,
			req: RefreshRequest{
				RefreshToken: "refresh-token",
			},
		}

		resp, err := svc.RefreshToken(ctx)
		require.Error(t, err)
		require.Nil(t, resp)
	}

}

func TestNewParams(t *testing.T) {
	{ // create params success
		testReq := model.RefreshJSONRequestBody{
			RefreshToken: "token",
		}
		req, err := NewParams(testReq)
		require.NoError(t, err)
		require.NotNil(t, req)
		require.Equal(t, "token", req.RefreshToken)

	}
}
