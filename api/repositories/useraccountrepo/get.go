package useraccountrepo

import (
	"context"
	"database/sql"
	"mysite/models/pgmodel"

	"github.com/friendsofgo/errors"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

func (u userAccountRepo) GetUserAccountByUserName(ctx context.Context, tx boil.ContextTransactor, userName string) (*pgmodel.UserAccount, error) {
	mods := []qm.QueryMod{
		pgmodel.UserAccountWhere.UserName.EQ(userName),
	}
	pgUserAccount, err := pgmodel.UserAccounts(mods...).One(ctx, tx)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, errors.Wrap(err, "failed check user exist")
	}

	return pgUserAccount, nil
}

func (u userAccountRepo) GetActiveUserAccountById(ctx context.Context, tx boil.ContextTransactor, userId int) (*pgmodel.UserAccount, error) {
	mods := []qm.QueryMod{
		pgmodel.UserAccountWhere.ID.EQ(userId),
		pgmodel.UserAccountWhere.IsActive.EQ(true),
		pgmodel.UserAccountWhere.IsDeleted.EQ(false),
	}

	pgUserAccount, err := pgmodel.UserAccounts(mods...).One(ctx, tx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get userAccount")
	}

	return pgUserAccount, nil
}

func (u userAccountRepo) GetActiveUserAccountByName(ctx context.Context, tx boil.ContextTransactor, userName string) (*pgmodel.UserAccount, error) {
	mods := []qm.QueryMod{
		pgmodel.UserAccountWhere.UserName.EQ(userName),
		pgmodel.UserAccountWhere.IsActive.EQ(true),
		pgmodel.UserAccountWhere.IsDeleted.EQ(false),
	}

	pgUserAccount, err := pgmodel.UserAccounts(mods...).One(ctx, tx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get userAccount")
	}

	return pgUserAccount, nil
}
