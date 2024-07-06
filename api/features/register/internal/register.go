package internal

import (
	"context"
	"mysite/constants"
	"mysite/dtos"
	"mysite/entities"
	"mysite/pkgs/auth"
	"mysite/pkgs/database"
	"mysite/pkgs/validate"
	"mysite/repositories/useraccountrepo"
	"mysite/utils/httputil"

	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

type service struct {
	repo    useraccountrepo.UserAccountRepo
	authSvc auth.AuthService
	req     RegisterRequest
}

type RegisterRequest struct {
	// Password password
	Password       string `validate:"required"`
	HashedPassword string

	// UserName email
	UserName string `validate:"email,required"`

	// Email email of user
	Email *string `json:"email,omitempty" validate:"omitempty,email"`

	// gender of user
	Gender *string `json:"gender,omitempty" validate:"omitempty,oneof=male female other"`

	// Name name of user
	Name *string `json:"name,omitempty"`

	// Phone phone number of user
	Phone *string `json:"phone,omitempty" validate:"omitempty,number"`
}

func NewService(req RegisterRequest) service {
	return service{
		repo:    useraccountrepo.NewRepo(),
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

func NewParams(req dtos.RegisterRequest) (*RegisterRequest, error) {
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

	user = &entities.UserAccount{
		UserName:  s.req.UserName,
		Password:  s.req.HashedPassword,
		IsActive:  true,
		IsDeleted: false,
	}
	if err := s.repo.Insert(ctx, tx, user); err != nil {
		return errors.Wrap(err, "failed to insert user")
	}

	if !s.req.hasUserInfo() {
		return nil
	}

	userInfo := entities.UserInfo{
		Name:          null.StringFromPtr(s.req.Name),
		Phone:         null.StringFromPtr(s.req.Phone),
		Email:         null.StringFromPtr(s.req.Email),
		Gender:        null.StringFromPtr(s.req.Gender),
		MembershipID:  null.IntFrom(constants.Bronze),
		UserAccountID: user.ID,
	}
	if err := user.AddUserInfos(ctx, tx, true, &userInfo); err != nil {
		return errors.Wrap(err, "failed to save userInfo")
	}

	return nil
}
func (req RegisterRequest) hasUserInfo() bool {
	return req.Email != nil || req.Gender != nil || req.Phone != nil || req.Name != nil
}
