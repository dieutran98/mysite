package useraccountrepo

import (
	"context"
	"database/sql"
	"mysite/entities"

	"github.com/friendsofgo/errors"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

func (u userAccountRepo) GetUserAccountByUserName(ctx context.Context, tx boil.ContextTransactor, userName string) (*entities.UserAccount, error) {
	mods := []qm.QueryMod{
		entities.UserAccountWhere.UserName.EQ(userName),
	}
	pgUserAccount, err := entities.UserAccounts(mods...).One(ctx, tx)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, errors.Wrap(err, "failed check user exist")
	}

	return pgUserAccount, nil
}

func (u userAccountRepo) GetActiveUserAccountById(ctx context.Context, tx boil.ContextTransactor, userId int) (*entities.UserAccount, error) {
	mods := []qm.QueryMod{
		entities.UserAccountWhere.ID.EQ(userId),
		entities.UserAccountWhere.IsActive.EQ(true),
		entities.UserAccountWhere.IsDeleted.EQ(false),
	}

	pgUserAccount, err := entities.UserAccounts(mods...).One(ctx, tx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get userAccount")
	}

	return pgUserAccount, nil
}

func (u userAccountRepo) GetActiveUserAccountByName(ctx context.Context, tx boil.ContextTransactor, userName string) (*entities.UserAccount, error) {
	mods := []qm.QueryMod{
		entities.UserAccountWhere.UserName.EQ(userName),
		entities.UserAccountWhere.IsActive.EQ(true),
		entities.UserAccountWhere.IsDeleted.EQ(false),
	}

	pgUserAccount, err := entities.UserAccounts(mods...).One(ctx, tx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get userAccount")
	}

	return pgUserAccount, nil
}
