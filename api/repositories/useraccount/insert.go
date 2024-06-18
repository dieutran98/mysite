package useraccount

import (
	"context"
	"mysite/models/pgmodel"

	"github.com/friendsofgo/errors"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

func (u userAccountRepo) Insert(ctx context.Context, tx boil.ContextTransactor, pgUser pgmodel.UserAccount) error {
	if err := pgUser.Insert(ctx, tx, boil.Infer()); err != nil {
		return errors.Wrap(err, "failed to insert user")
	}

	return nil
}
