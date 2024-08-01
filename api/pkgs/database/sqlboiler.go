package database

import (
	"context"
	"log/slog"
	"mysite/constants"
	"mysite/pkgs/env"
	"mysite/pkgs/logger"
	"time"

	"github.com/friendsofgo/errors"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

type ExecuteQueriesFunc func(ctx context.Context, tx boil.ContextTransactor) error

func newBoilerDb() error {
	boil.SetDB(getDb().db)
	return nil
}

func NewBoilerTransaction(ctx context.Context, fn ExecuteQueriesFunc) error {
	var cancel context.CancelFunc = func() {
		// do nothing
	}
	// // set timeout
	if _, hasDeadline := ctx.Deadline(); !hasDeadline {
		dbEnv := env.GetEnv().Database
		ctx, cancel = context.WithTimeout(ctx, time.Duration(dbEnv.TransactionTimeout)*time.Second)
	}
	defer cancel()

	ctx = boil.WithDebug(ctx, true)
	ctx = boil.WithDebugWriter(ctx, logger.NewBoilerLogger(ctx))

	tx, err := boil.BeginTx(ctx, nil)
	if err != nil {
		return errors.Wrap(err, "failed to start transaction")
	}

	// add recovery on panic
	defer func() {
		if r := recover(); r != nil {
			slog.Error("panic occur in transaction", slog.Any("panic", r))
			rollback(tx)
		}
	}()

	// execute queries
	if err := fn(ctx, tx); err != nil {
		rollback(tx)
		return errors.Wrap(err, "failed to execute queries")
	}

	// rollback if running as test
	if checkIsRunAsTest(ctx) {
		rollback(tx)
		return nil
	}

	if err := tx.Commit(); err != nil {
		return errors.Wrap(err, "failed to commit")
	}

	return nil
}

func rollback(tx boil.ContextTransactor) {
	if err := tx.Rollback(); err != nil {
		slog.Error("rollback error", logger.AttrError(errors.Wrap(err, "failed rollback")))
	} else {
		slog.Debug("rollback!")
	}
}

func checkIsRunAsTest(ctx context.Context) bool {
	isTesting, found := ctx.Value(constants.Testing).(bool)
	return found && isTesting
}
