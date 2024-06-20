package internal

import (
	"context"
	"mysite/models/model"
	"mysite/models/pgmodel"
	"mysite/pkgs/auth"
	"mysite/pkgs/database"
	"mysite/pkgs/env"
	"mysite/pkgs/validate"
	"mysite/repositories/useraccountrepo"
	"mysite/utils/httputil"
	"time"

	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

type service struct {
	repo    useraccountrepo.UserAccountRepo
	req     RefreshRequest
	authSvc auth.AuthService
}

type RefreshRequest struct {
	RefreshToken string ` validate:"required"`
}

type RefreshResponse struct {
	// AccessToken access token
	AccessToken string
}

func NewService(req RefreshRequest) service {
	return service{
		repo:    useraccountrepo.NewRepo(),
		req:     req,
		authSvc: auth.NewAuthService(),
	}
}

func NewParams(req model.RefreshJSONRequestBody) (*RefreshRequest, error) {
	var result RefreshRequest
	if err := mapstructure.Decode(req, &result); err != nil {
		return nil, errors.Wrap(err, "failed decode")
	}

	return &result, nil
}

func (s service) RefreshToken(ctx context.Context) (*RefreshResponse, error) {
	if err := validateParams(s.req); err != nil {
		return nil, errors.Wrap(err, "failed validate refresh token request body")
	}

	jwtEnv := env.GetEnv().Jwt
	// parse  token
	claims, err := s.authSvc.ParseToken(s.req.RefreshToken, []byte(jwtEnv.RefreshKey))
	if err != nil {
		return nil, errors.Wrapf(httputil.ErrUnauthorize, "failed to parse token: %s", err.Error())
	}

	var pgUserAccount *pgmodel.UserAccount
	// get user by user id
	if err := database.NewBoilerTransaction(ctx, func(ctx context.Context, tx boil.ContextTransactor) error {
		var err error
		pgUserAccount, err = s.repo.GetActiveUserAccountById(ctx, tx, claims.Subject)
		if err != nil {
			return errors.Wrap(err, "failed get userAccount")
		}
		return nil
	}); err != nil {
		return nil, errors.Wrap(httputil.ErrUnauthorize, err.Error())
	}

	// generate access token
	accessToken, err := s.authSvc.CreateToken(s.authSvc.NewClaims(pgUserAccount.ID, time.Now().Add(time.Minute*15)), []byte(jwtEnv.AccessKey))
	if err != nil {
		return nil, errors.Wrap(err, "failed create access token")
	}

	return &RefreshResponse{
		AccessToken: accessToken,
	}, nil
}

func validateParams(params RefreshRequest) error {
	if err := validate.ValidateStruct(params); err != nil {
		return errors.Wrap(httputil.ErrInvalidRequest, err.Error())
	}
	return nil
}
