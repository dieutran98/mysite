package useraccountrepo

import (
	"context"
	"mysite/models/pgmodel"

	"github.com/volatiletech/sqlboiler/v4/boil"
)

type Get interface {
	GetUserAccountByUserName(ctx context.Context, tx boil.ContextTransactor, userName string) (*pgmodel.UserAccount, error)
	GetActiveUserAccountById(ctx context.Context, tx boil.ContextTransactor, userId string) (*pgmodel.UserAccount, error)
	GetActiveUserAccountByName(ctx context.Context, tx boil.ContextTransactor, userName string) (*pgmodel.UserAccount, error)
}

type Insert interface {
	Insert(ctx context.Context, tx boil.ContextTransactor, user pgmodel.UserAccount) error
}

type Update interface {
	ActiveUser(ctx context.Context, tx boil.ContextTransactor, pgUser pgmodel.UserAccount) error
}

type Delete interface{}

type UserAccountRepo interface {
	Get
	Insert
	Update
	Delete
}

type userAccountRepo struct {
}

func NewRepo() UserAccountRepo {
	return &userAccountRepo{}
}
