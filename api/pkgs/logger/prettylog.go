package logger

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"os"
	"sync"
	"time"

	"github.com/pkg/errors"
)

type Font string

const (
	Reset         Font = "\033[0m"
	Black         Font = "\033[30m"
	Red           Font = "\033[31m"
	Green         Font = "\033[32m"
	Yellow        Font = "\033[33m"
	Blue          Font = "\033[34m"
	Magenta       Font = "\033[35m"
	Cyan          Font = "\033[36m"
	White         Font = "\033[37m"
	Gray          Font = "\033[90m"
	BrightBlack   Font = "\033[90m"
	BrightRed     Font = "\033[91m"
	BrightGreen   Font = "\033[92m"
	BrightYellow  Font = "\033[93m"
	BrightBlue    Font = "\033[94m"
	BrightMagenta Font = "\033[95m"
	BrightCyan    Font = "\033[96m"
	BrightWhite   Font = "\033[97m"
	BgBlack       Font = "\033[40m"
	BgRed         Font = "\033[41m"
	BgGreen       Font = "\033[42m"
	BgYellow      Font = "\033[43m"
	BgBlue        Font = "\033[44m"
	BgMagenta     Font = "\033[45m"
	BgCyan        Font = "\033[46m"
	BgWhite       Font = "\033[47m"
	LightGreen    Font = "\033[92m"
	Bold          Font = "\033[1m"
	None          Font = "none"
)

type Option struct {
	TimeFormat string
	Level      *slog.LevelVar
}

type prettyHandler struct {
	writer io.Writer   // destination of the log
	mutex  *sync.Mutex // mutex for handling concurrency
	option Option      // option config for logging
	attrs  []slog.Attr
	group  []string
}

func NewPrettyHandler(writer io.Writer, opts *Option) *prettyHandler {
	handlerOpts := makeOptions(opts)
	return &prettyHandler{
		writer: writer,
		mutex:  &sync.Mutex{},
		option: handlerOpts,
	}
}

func makeOptions(opts *Option) Option {
	handlerOpts := Option{
		TimeFormat: time.RFC3339, // Default time format
	}
	if opts == nil {
		return handlerOpts
	}

	if opts.Level != nil {
		handlerOpts.Level = opts.Level
	}

	if opts.TimeFormat != "" {
		handlerOpts.TimeFormat = opts.TimeFormat
	}
	return handlerOpts
}

func (p *prettyHandler) Handle(ctx context.Context, record slog.Record) error {
	var buf bytes.Buffer

	if len(p.attrs) > 0 {
		record.AddAttrs(p.attrs...)
	}

	attrStr, err := p.printAttributes(ctx, record)
	if err != nil {
		fmt.Println(err)
		return errors.Wrap(err, "failed printAttributes")
	}
	buf.WriteString(attrStr)

	return p.write(buf.Bytes())
}

func (p *prettyHandler) write(data []byte) error {
	if _, err := p.writer.Write(data); err != nil {
		return errors.Wrap(err, "failed write log")
	}
	return nil
}

// Enabled reports whether l is greater than or equal to the
func (p *prettyHandler) Enabled(ctx context.Context, level slog.Level) bool {
	if p.option.Level == nil {
		return false
	}
	return p.option.Level.Level() <= level
}

// WithAttrs returns a new Handler whose attributes consist of
// both the receiver's attributes and the arguments.
// The Handler owns the slice: it may retain, modify or discard it.
func (p *prettyHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	if len(attrs) == 0 {
		return p
	}

	p2 := p.clone()
	p2.attrs = append(p2.attrs, attrs...)
	return &p2
}

// WithGroup returns a new Handler with the given group appended to
// the receiver's existing groups.
func (p *prettyHandler) WithGroup(name string) slog.Handler {
	if name == "" {
		return p
	}

	p2 := p.clone()
	p2.group = append(p2.group, name)
	return &p2
}

func (p *prettyHandler) clone() prettyHandler {
	return prettyHandler{
		writer: p.writer,
		mutex:  &sync.Mutex{},
		option: p.option,
		attrs:  append([]slog.Attr{}, p.attrs...),
		group:  append([]string{}, p.group...),
	}
}

func printLevel(level slog.Level) string {
	return level.String()
}

func fontByLevel(level slog.Level) Font {
	switch level {
	case slog.LevelInfo:
		return LightGreen
	case slog.LevelError:
		return Red
	case slog.LevelWarn:
		return Yellow
	case slog.LevelDebug:
		return Cyan
	default:
		return White
	}
}

func (p prettyHandler) printTime(t time.Time) string {
	timeFormat := p.option.TimeFormat
	if p.option.TimeFormat == "" {
		timeFormat = time.DateTime
	}
	return t.Format(timeFormat)
}

func printFont(font Font, src string) string {
	return fmt.Sprintf("%s%s%s", font.String(), src, Reset.String())
}

func (p *prettyHandler) printAttributes(ctx context.Context, record slog.Record) (string, error) {
	var buf bytes.Buffer
	jsonHandler := slog.NewJSONHandler(&buf, &slog.HandlerOptions{
		Level: p.option.Level.Level(),
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			switch a.Key {
			case "time":
				return slog.String(a.Key, p.printTime(record.Time))
			case "level":
				return slog.String(a.Key, printLevel(record.Level))
			case "function":
				return slog.Attr{}
			default:
				return a
			}
		},
		AddSource: true,
	})
	if len(p.attrs) > 0 {
		record.AddAttrs(p.attrs...)
	}

	if err := jsonHandler.Handle(ctx, record); err != nil {
		return "", errors.Wrap(err, "failed handle json log")
	}

	switch p.writer.(type) {
	case *os.File:
		return jsonFormat(buf.Bytes(), None)
	default:
		return jsonFormat(buf.Bytes(), fontByLevel(record.Level))
	}
}

func (f Font) String() string {
	return string(f)
}

func jsonFormat(data []byte, font Font) (string, error) {
	var out bytes.Buffer
	if err := json.Indent(&out, data, "", "   "); err != nil {
		return "", errors.Wrap(err, "failed to indent json string")
	}

	switch font {
	case None:
		return out.String(), nil
	default:
		return printFont(font, out.String()), nil
	}

}
