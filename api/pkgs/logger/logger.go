package logger

import (
	"io"
	"log/slog"
)

var programLevel = new(slog.LevelVar) // default is info

func SetLogger(wr io.Writer) {
	logger := slog.New(NewPrettyHandler(wr, &Option{
		Level: programLevel,
	}))

	slog.SetDefault(logger)

}

func SetLogLevel(level slog.Level) {
	programLevel.Set(level)
}
