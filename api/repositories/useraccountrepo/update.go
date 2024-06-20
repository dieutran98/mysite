package useraccountrepo

import (
	"context"
	"mysite/models/pgmodel"

	"github.com/friendsofgo/errors"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

func (u userAccountRepo) ActiveUser(ctx context.Context, tx boil.ContextTransactor, pgUser pgmodel.UserAccount) error {
	pgUser.IsActive = true
	pgUser.IsDeleted = false
	pgUser.DeletedAt = null.Time{}
	rowEffected, err := pgUser.Update(ctx, tx, boil.Infer())
	if err != nil || rowEffected == 0 {
		return errors.Wrap(err, "failed to update UserAccount")
	}
	return nil
}
