package useraccountrepo

import (
	"context"
	"mysite/entities"

	"github.com/volatiletech/sqlboiler/v4/boil"
)

type Get interface {
	GetUserAccountByUserName(ctx context.Context, tx boil.ContextTransactor, userName string) (*entities.UserAccount, error)
	GetActiveUserAccountById(ctx context.Context, tx boil.ContextTransactor, userId int) (*entities.UserAccount, error)
	GetActiveUserAccountByName(ctx context.Context, tx boil.ContextTransactor, userName string) (*entities.UserAccount, error)
}

type Insert interface {
	Insert(ctx context.Context, tx boil.ContextTransactor, user *entities.UserAccount) error
}

type Update interface {
	ActiveUser(ctx context.Context, tx boil.ContextTransactor, pgUser entities.UserAccount) error
}

type Delete interface{}

//go:generate moq -pkg repomock -out ../../testing/mocking/repomock/useraccountmock.go . UserAccountRepo
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
