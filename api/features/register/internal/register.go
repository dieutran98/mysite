package internal

import (
	"context"
	"mysite/models/model"
	"mysite/models/pgmodel"
	"mysite/pkgs/auth"
	"mysite/pkgs/database"
	"mysite/pkgs/validate"
	useraccount "mysite/repositories/useraccount"
	"mysite/utils/httputil"

	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

type service struct {
	repo    useraccount.UserAccountRepo
	authSvc auth.AuthService
	req     RegisterRequest
}

type RegisterRequest struct {
	// Password password
	Password       string `validate:"required"`
	HashedPassword string

	// UserName email
	UserName string `validate:"email,required"`
}

func NewService(req RegisterRequest) service {
	return service{
		repo:    useraccount.NewRepo(),
		authSvc: auth.NewAuthService(),
		req:     req,
	}
}

func (s service) Register(ctx context.Context) error {
	// validate data
	if err := validateRequest(s.req); err != nil {
		return errors.Wrap(httputil.ErrInvalidRequest, err.Error())
	}

	// hash password
	hash, err := s.authSvc.HashPassword(s.req.Password)
	if err != nil {
		return errors.Wrap(err, "failed to hash password")
	}
	s.req.HashedPassword = hash

	// save userName and password

	if err := database.NewBoilerTransaction(ctx, s.registerUser); err != nil {
		return errors.Wrap(err, "failed insert user")
	}

	return nil
}

func NewParams(req model.RegisterRequest) (*RegisterRequest, error) {
	var result RegisterRequest
	if err := mapstructure.Decode(req, &result); err != nil {
		return nil, errors.Wrap(err, "failed mapping struct")
	}
	return &result, nil
}

func validateRequest(req RegisterRequest) error {
	if err := validate.ValidateStruct(req); err != nil {
		return errors.Wrap(err, "failed to validate RegisterRequest")
	}
	return nil
}

func (s service) registerUser(ctx context.Context, tx boil.ContextTransactor) error {
	// checking user is exist
	user, err := s.repo.GetUserAccountByUserName(ctx, tx, s.req.UserName)
	if err != nil {
		return errors.Wrap(err, "failed checking IsUserNameExist")
	}

	if user != nil {
		switch {
		case !user.IsActive, user.IsDeleted:
			if err := s.repo.ActiveUser(ctx, tx, *user); err != nil {
				return errors.Wrap(err, "failed to active user")
			}
		default:
			return errors.Wrap(httputil.ErrInvalidRequest, "user existed")
		}
		return nil
	}

	user = &pgmodel.UserAccount{
		UserName:  s.req.UserName,
		Password:  s.req.HashedPassword,
		IsActive:  true,
		IsDeleted: false,
	}
	if err := s.repo.Insert(ctx, tx, *user); err != nil {
		return errors.Wrap(err, "failed to insert user")
	}

	return nil
}
