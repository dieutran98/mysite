package internal

import (
	"context"
	"fmt"
	"mysite/dtos"
	"mysite/entities"
	"mysite/pkgs/auth"
	"mysite/pkgs/database"
	"mysite/testing/dbtest"
	"mysite/testing/mocking/pkgmock"
	"mysite/testing/mocking/repomock"
	"testing"

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
	req := RefreshRequest{
		RefreshToken: "refresh-token",
	}

	{ // refresh success
		repoMock := &repomock.UserAccountRepoMock{}
		repoMock.GetActiveUserAccountByIdFunc = func(ctx context.Context, tx boil.ContextTransactor, userId int) (*entities.UserAccount, error) {
			return &entities.UserAccount{
				ID: 1,
			}, nil
		}

		jwtMock := &pkgmock.JwtHandlerMock{}
		jwtMock.ParseTokenFunc = func(tokenString string, claims auth.Claims) error {
			claims.(*auth.CustomClaims[any]).Subject = "1"
			return nil
		}
		jwtMock.CreateTokenFunc = func() (string, error) { return "token", nil }
		jwtMock.WithClaimsFunc = func(claims auth.Claims) auth.JwtHandler { return jwtMock }

		svc := service{
			repo:       repoMock,
			jwtHandler: jwtMock,
		}

		resp, err := svc.RefreshToken(ctx, req)
		require.NoError(t, err)
		require.NotNil(t, resp)
		require.Equal(t, "token", resp.AccessToken)
	}
	{ // refresh failed, validate failed

		req := RefreshRequest{
			RefreshToken: "",
		}
		svc := service{}

		resp, err := svc.RefreshToken(ctx, req)
		require.Error(t, err)
		require.Nil(t, resp)
	}
	{ // refresh failed, parse token failed
		repoMock := &repomock.UserAccountRepoMock{}
		repoMock.GetActiveUserAccountByIdFunc = func(ctx context.Context, tx boil.ContextTransactor, userId int) (*entities.UserAccount, error) {
			return &entities.UserAccount{
				ID: 1,
			}, nil
		}

		jwtMock := &pkgmock.JwtHandlerMock{}
		jwtMock.ParseTokenFunc = func(tokenString string, claims auth.Claims) error {
			return errors.New("parse token failed")
		}

		svc := service{
			repo:       repoMock,
			jwtHandler: jwtMock,
		}

		resp, err := svc.RefreshToken(ctx, req)
		require.Error(t, err)
		require.Nil(t, resp)
	}
	{ // refresh failed, get user failed
		repoMock := &repomock.UserAccountRepoMock{}
		repoMock.GetActiveUserAccountByIdFunc = func(ctx context.Context, tx boil.ContextTransactor, userId int) (*entities.UserAccount, error) {
			return nil, errors.New("get user account failed")
		}

		jwtMock := &pkgmock.JwtHandlerMock{}
		jwtMock.ParseTokenFunc = func(tokenString string, claims auth.Claims) error {
			return nil
		}

		svc := service{
			repo:       repoMock,
			jwtHandler: jwtMock,
		}

		resp, err := svc.RefreshToken(ctx, req)
		require.Error(t, err)
		require.Nil(t, resp)
	}
	{ // refresh failed, createToken failed
		repoMock := &repomock.UserAccountRepoMock{}
		repoMock.GetActiveUserAccountByIdFunc = func(ctx context.Context, tx boil.ContextTransactor, userId int) (*entities.UserAccount, error) {
			return &entities.UserAccount{
				ID: 1,
			}, nil
		}

		jwtMock := &pkgmock.JwtHandlerMock{}
		jwtMock.ParseTokenFunc = func(tokenString string, claims auth.Claims) error {
			return nil
		}
		jwtMock.CreateTokenFunc = func() (string, error) {
			return "", errors.New("create token failed")
		}
		jwtMock.WithClaimsFunc = func(claims auth.Claims) auth.JwtHandler { return jwtMock }

		svc := service{
			repo:       repoMock,
			jwtHandler: jwtMock,
		}

		resp, err := svc.RefreshToken(ctx, req)
		require.Error(t, err)
		require.Nil(t, resp)
	}

}

func TestNewParams(t *testing.T) {
	{ // create params success
		testReq := dtos.RefreshJSONRequestBody{
			RefreshToken: "token",
		}
		req, err := NewParams(testReq)
		require.NoError(t, err)
		require.NotNil(t, req)
		require.Equal(t, "token", req.RefreshToken)

	}
}
