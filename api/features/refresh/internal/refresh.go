package internal

import (
	"context"
	"mysite/dtos"
	"mysite/pkgs/auth"
	"mysite/pkgs/database"
	"mysite/pkgs/validate"
	"mysite/repositories/useraccountrepo"
	"mysite/utils/httputil"
	"strconv"
	"time"

	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

type service struct {
	repo       useraccountrepo.UserAccountRepo
	jwtHandler auth.JwtHandler
}

type RefreshRequest struct {
	RefreshToken string `validate:"required"`
}

type RefreshResponse struct {
	// AccessToken access token
	AccessToken string
}

func NewService() service {
	return service{
		repo:       useraccountrepo.NewRepo(),
		jwtHandler: auth.NewJwtHandler(),
	}
}

func NewParams(req dtos.RefreshJSONRequestBody) (*RefreshRequest, error) {
	var result RefreshRequest
	if err := mapstructure.Decode(req, &result); err != nil {
		return nil, errors.Wrap(err, "failed decode")
	}

	return &result, nil
}

func (s service) RefreshToken(ctx context.Context, req RefreshRequest) (*RefreshResponse, error) {
	if err := validateParams(req); err != nil {
		return nil, errors.Wrap(err, "failed validate refresh token request body")
	}

	// parse  token
	var claims auth.CustomClaims[any]
	claims.KeyType = auth.AccessKey

	err := s.jwtHandler.ParseToken(req.RefreshToken, &claims)
	if err != nil {
		return nil, errors.Wrapf(httputil.ErrUnauthorize, "failed to parse token: %s", err.Error())
	}

	// get userId from claims
	userId, err := strconv.Atoi(claims.Subject)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get user Id")
	}

	// get user by user id
	if err := database.NewBoilerTransaction(ctx, func(ctx context.Context, tx boil.ContextTransactor) error {
		var err error
		pgUserAccount, err := s.repo.GetActiveUserAccountById(ctx, tx, userId)
		if err != nil {
			return errors.Wrap(err, "failed get userAccount")
		}
		if pgUserAccount == nil {
			return errors.New("failed get userAccount")
		}
		return nil
	}); err != nil {
		return nil, errors.Wrap(httputil.ErrUnauthorize, err.Error())
	}

	// generate access token
	accessClaims := claims.Clone().WithExpireAt(time.Now().Add(time.Minute * 15))
	claims.KeyType = auth.AccessKey
	accessToken, err := s.jwtHandler.WithClaims(accessClaims).CreateToken()
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
