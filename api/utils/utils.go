package utils

import (
	"context"
	"mysite/constants"
)

func CheckIsRunAsTest(ctx context.Context) bool {
	isTesting, found := ctx.Value(constants.Testing).(bool)
	return found && isTesting
}
