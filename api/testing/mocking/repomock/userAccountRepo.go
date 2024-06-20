package repomock

import (
	"context"
	"mysite/models/pgmodel"
	"mysite/repositories/useraccount"

	"github.com/volatiletech/sqlboiler/v4/boil"
)

type mockService struct {
	useraccount.UserAccountRepo
	GetUserAccountByUserNameFunc func() (*pgmodel.UserAccount, error)
	ActiveUserFunc               func() error
	InsertFunc                   func() error
}

func NewUserAccountMock() mockService {
	return mockService{}
}

func (m mockService) GetUserAccountByUserName(ctx context.Context, tx boil.ContextTransactor, userName string) (*pgmodel.UserAccount, error) {
	return m.GetUserAccountByUserNameFunc()
}

func (m mockService) Insert(ctx context.Context, tx boil.ContextTransactor, user pgmodel.UserAccount) error {
	return m.InsertFunc()
}

func (m mockService) ActiveUser(ctx context.Context, tx boil.ContextTransactor, pgUser pgmodel.UserAccount) error {
	return m.ActiveUserFunc()
}
