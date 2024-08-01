package internal

import (
	"context"
	"log/slog"
	"mysite/dtos"
	"mysite/entities"
	"mysite/pkgs/auth"
	"mysite/pkgs/database"
	"mysite/pkgs/logger"
	"mysite/pkgs/validate"
	"mysite/repositories/useraccountrepo"
	"mysite/utils/httputil"
	"strconv"
	"time"

	"github.com/friendsofgo/errors"
	"github.com/mitchellh/mapstructure"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

type service struct {
	repo       useraccountrepo.UserAccountRepo
	req        LoginRequest
	authSvc    auth.AuthService
	jwtHandler auth.JwtHandler
}

type LoginRequest struct {
	// Password password
	Password string `json:"password" validate:"required"`

	// UserName email
	UserName string `json:"userName" validate:"email,required"`
}

type LoginResponse struct {
	// AccessToken access token
	AccessToken string

	// RefreshToken refresh token
	RefreshToken string
}

func NewService(req LoginRequest) *service {
	return &service{
		repo:       useraccountrepo.NewRepo(),
		req:        req,
		authSvc:    auth.NewAuthService(),
		jwtHandler: auth.NewJwtHandler(),
	}
}

func NewParams(req dtos.LoginJSONRequestBody) (*LoginRequest, error) {
	var result LoginRequest
	if err := mapstructure.Decode(req, &result); err != nil {
		return nil, errors.Wrap(err, "failed decode")
	}

	return &result, nil
}

func (s *service) Login(ctx context.Context) (*LoginResponse, error) {
	// validate params
	if err := validateParams(s.req); err != nil {
		return nil, errors.Wrap(err, "failed validate login request")
	}

	var user *entities.UserAccount
	// get password from db by userName
	if err := database.NewBoilerTransaction(ctx, func(ctx context.Context, tx boil.ContextTransactor) error {
		var err error
		user, err = s.repo.GetActiveUserAccountByName(ctx, tx, s.req.UserName)
		if err != nil || user == nil {
			return errors.Wrap(httputil.ErrUnauthorize, "login failed at step 1")
		}

		return nil
	}); err != nil {
		return nil, err
	}

	// check password hash
	match, err := s.authSvc.ComparePasswordAndHash(s.req.Password, user.Password)
	if err != nil || !match {
		return nil, errors.Wrap(httputil.ErrUnauthorize, "login failed at step 2")
	}

	// generate access token and refresh token
	accessToken, refreshToken, err := s.generateToken(user.ID)
	if err != nil {
		slog.Error("failed generate token", logger.AttrError(err))
		return nil, errors.Wrap(httputil.ErrUnauthorize, "login failed at step 3")
	}

	return &LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func validateParams(params LoginRequest) error {
	if err := validate.ValidateStruct(params); err != nil {
		return errors.Wrap(httputil.ErrInvalidRequest, err.Error())
	}
	return nil
}

func (s *service) generateToken(userId int) (accessToken string, refreshToken string, err error) {
	claims := auth.NewCustomClaims[any]()
	claims.Subject = strconv.Itoa(userId)
	accessClaims := claims.Clone().WithExpireAt(time.Now().Add(15 * time.Minute))
	accessClaims.KeyType = auth.AccessKey

	accessToken, err = s.jwtHandler.WithClaims(accessClaims).CreateToken()
	if err != nil {
		return "", "", errors.Wrap(err, "failed create accessKey")
	}

	refreshClaims := claims.Clone().WithExpireAt(time.Now().Add(72 * time.Hour))
	refreshClaims.KeyType = auth.RefreshKey
	refreshToken, err = s.jwtHandler.WithClaims(refreshClaims).CreateToken()
	if err != nil {
		return "", "", errors.Wrap(err, "failed create refreshKey")
	}

	return
}
