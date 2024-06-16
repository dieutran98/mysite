package database

import (
	"context"
	"log/slog"
	"mysite/pkgs/env"
	"mysite/pkgs/logger"
	"mysite/utils"
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
	var cancel context.CancelFunc
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
			tx.Rollback()
		}
	}()

	// execute queries
	if err := fn(ctx, tx); err != nil {
		return errors.Wrap(err, "failed to execute queries")
	}

	// rollback if running as test
	if utils.CheckIsRunAsTest(ctx) {
		if err := tx.Rollback(); err != nil {
			return errors.Wrap(err, "failed to roll back")
		}
	}

	if err := tx.Commit(); err != nil {
		return errors.Wrap(err, "failed to commit")
	}

	return nil
}
