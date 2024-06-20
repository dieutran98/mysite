package logger

import (
	"context"
	"io"
	"log/slog"
	"strings"

	"github.com/go-chi/chi/v5/middleware"
)

const (
	boilerQueries   int = 0
	boilerArgs      int = 1
	boilerLoggedAll int = 2
)

type boilerLogger struct {
	writeCount    int
	TransactionId string
	Query         string
	Args          string
}

func NewBoilerLogger(ctx context.Context) io.Writer {
	bLogger := boilerLogger{
		TransactionId: middleware.GetReqID(ctx),
	}
	if bLogger.TransactionId == "" {
		bLogger.TransactionId = "unknown"
	}
	return &bLogger
}

func (bl *boilerLogger) Write(b []byte) (n int, err error) {
	bl.write(b)
	bl.log()
	return
}

func (bl *boilerLogger) write(b []byte) {
	switch bl.writeCount {
	case boilerQueries:
		bl.Query = strings.TrimSpace((string(b)))
	case boilerArgs:
		bl.Args = strings.TrimSpace((string(b)))
	}
	bl.writeCount++
}

func (bl *boilerLogger) log() {
	// check if all are logged
	if bl.writeCount < boilerLoggedAll {
		return
	}

	slog.Debug("boiler log",
		slog.String("transactionId", bl.TransactionId),
		slog.String("query", bl.Query),
		slog.String("args", bl.Args),
	)
	// reset counter
	bl.writeCount = 0
}
