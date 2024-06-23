package repomock

import (
	"context"
	"mysite/models/pgmodel"

	"github.com/volatiletech/sqlboiler/v4/boil"
)

type mockService struct {
	GetUserAccountByUserNameFunc   func() (*pgmodel.UserAccount, error)
	GetActiveUserAccountByIdFunc   func() (*pgmodel.UserAccount, error)
	GetActiveUserAccountByNameFunc func() (*pgmodel.UserAccount, error)
	ActiveUserFunc                 func() error
	InsertFunc                     func() error
}

func NewUserAccountMock() mockService {
	return mockService{}
}

func (m mockService) GetUserAccountByUserName(ctx context.Context, tx boil.ContextTransactor, userName string) (*pgmodel.UserAccount, error) {
	return m.GetUserAccountByUserNameFunc()
}

func (m mockService) GetActiveUserAccountById(ctx context.Context, tx boil.ContextTransactor, userId int) (*pgmodel.UserAccount, error) {
	return m.GetActiveUserAccountByIdFunc()
}

func (m mockService) GetActiveUserAccountByName(ctx context.Context, tx boil.ContextTransactor, userName string) (*pgmodel.UserAccount, error) {
	return m.GetActiveUserAccountByNameFunc()
}

func (m mockService) Insert(ctx context.Context, tx boil.ContextTransactor, user pgmodel.UserAccount) error {
	return m.InsertFunc()
}

func (m mockService) ActiveUser(ctx context.Context, tx boil.ContextTransactor, pgUser pgmodel.UserAccount) error {
	return m.ActiveUserFunc()
}
