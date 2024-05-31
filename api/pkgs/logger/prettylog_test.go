package logger

import (
	"bytes"
	"context"
	"testing"
	"time"

	"log/slog"

	"github.com/stretchr/testify/require"
)

func TestNewPrettyHandler(t *testing.T) {
	var buf bytes.Buffer
	levelVar := new(slog.LevelVar)
	levelVar.Set(slog.LevelInfo)
	opts := &Option{
		TimeFormat: time.RFC3339,
		Level:      levelVar,
	}

	handler := NewPrettyHandler(&buf, opts)

	require.NotNil(t, handler)
	require.Equal(t, opts.TimeFormat, handler.option.TimeFormat)
	require.Equal(t, opts.Level, handler.option.Level)
	require.NotNil(t, handler.writer)
	require.NotNil(t, handler.mutex)
}

func TestHandle(t *testing.T) {
	{
		var buf bytes.Buffer
		levelVar := new(slog.LevelVar)
		levelVar.Set(slog.LevelInfo)
		opts := &Option{
			TimeFormat: time.RFC3339,
			Level:      levelVar,
		}
		handler := NewPrettyHandler(&buf, opts)

		record := slog.Record{
			Time:  time.Now(),
			Level: slog.LevelInfo,
		}
		record.AddAttrs(slog.Attr{Key: "test", Value: slog.AnyValue("Test message")})

		err := handler.Handle(context.Background(), record)
		require.NoError(t, err)
		require.Contains(t, buf.String(), "Test message")
	}
	{
		var buf bytes.Buffer
		levelVar := new(slog.LevelVar)
		levelVar.Set(slog.LevelInfo)
		opts := &Option{
			TimeFormat: time.RFC3339,
			Level:      levelVar,
		}
		handler := NewPrettyHandler(&buf, opts)

		record := slog.Record{
			Time:  time.Now(),
			Level: slog.LevelError,
		}
		record.AddAttrs(slog.Attr{Key: "test", Value: slog.AnyValue("Test message")})

		err := handler.Handle(context.Background(), record)
		require.NoError(t, err)
		require.Contains(t, buf.String(), "Test message")
	}
	{
		var buf bytes.Buffer
		levelVar := new(slog.LevelVar)
		levelVar.Set(slog.LevelInfo)
		opts := &Option{
			TimeFormat: time.RFC3339,
			Level:      levelVar,
		}
		handler := NewPrettyHandler(&buf, opts)

		record := slog.Record{
			Time:  time.Now(),
			Level: slog.LevelWarn,
		}
		record.AddAttrs(slog.Attr{Key: "test", Value: slog.AnyValue("Test message")})

		err := handler.Handle(context.Background(), record)
		require.NoError(t, err)
		require.Contains(t, buf.String(), "Test message")
	}
	{
		var buf bytes.Buffer
		levelVar := new(slog.LevelVar)
		levelVar.Set(slog.LevelInfo)
		opts := &Option{
			TimeFormat: time.RFC3339,
			Level:      levelVar,
		}
		handler := NewPrettyHandler(&buf, opts)

		record := slog.Record{
			Time:  time.Now(),
			Level: slog.LevelDebug,
		}
		record.AddAttrs(slog.Attr{Key: "test", Value: slog.AnyValue("Test message")})

		err := handler.Handle(context.Background(), record)
		require.NoError(t, err)
		require.Contains(t, buf.String(), "Test message")
	}

}
func TestEnabled(t *testing.T) {
	var buf bytes.Buffer
	levelVar := new(slog.LevelVar)
	levelVar.Set(slog.LevelInfo)
	opts := &Option{
		TimeFormat: time.RFC3339,
		Level:      levelVar,
	}
	handler := NewPrettyHandler(&buf, opts)

	require.True(t, handler.Enabled(context.Background(), slog.LevelInfo))
	require.False(t, handler.Enabled(context.Background(), slog.LevelDebug))
}

func TestWithAttrs(t *testing.T) {
	var buf bytes.Buffer
	levelVar := new(slog.LevelVar)
	levelVar.Set(slog.LevelInfo)
	opts := &Option{
		TimeFormat: time.RFC3339,
		Level:      levelVar,
	}
	handler := NewPrettyHandler(&buf, opts)

	attrs := []slog.Attr{
		{Key: "key1", Value: slog.AnyValue("value1")},
	}

	newHandler := handler.WithAttrs(attrs).(*prettyHandler)

	require.NotNil(t, newHandler)
	require.Len(t, newHandler.attrs, 1)
	require.Equal(t, "key1", newHandler.attrs[0].Key)
	require.Equal(t, "value1", newHandler.attrs[0].Value.String())
}

func TestWithGroup(t *testing.T) {
	var buf bytes.Buffer
	levelVar := new(slog.LevelVar)
	levelVar.Set(slog.LevelInfo)
	opts := &Option{
		TimeFormat: time.RFC3339,
		Level:      levelVar,
	}
	handler := NewPrettyHandler(&buf, opts)

	groupName := "testGroup"
	newHandler := handler.WithGroup(groupName).(*prettyHandler)

	require.NotNil(t, newHandler)
	require.Len(t, newHandler.group, 1)
	require.Equal(t, groupName, newHandler.group[0])
}

func TestPrintAttributes(t *testing.T) {
	var buf bytes.Buffer
	levelVar := new(slog.LevelVar)
	levelVar.Set(slog.LevelInfo)
	opts := &Option{
		TimeFormat: time.RFC3339,
		Level:      levelVar,
	}
	handler := NewPrettyHandler(&buf, opts)

	record := slog.Record{
		Time:  time.Now(),
		Level: slog.LevelInfo,
	}
	record.AddAttrs(slog.Attr{Key: "test", Value: slog.AnyValue("Test message")})

	attrStr, err := handler.printAttributes(context.Background(), record)
	require.NoError(t, err)
	require.Contains(t, attrStr, "Test message")
}
