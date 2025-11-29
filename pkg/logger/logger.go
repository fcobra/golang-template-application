package logger

import (
	"io"
	"log/slog"
	"os"
)

const (
	LevelDebug = "debug"
	LevelInfo  = "info"
	LevelWarn  = "warn"
	LevelError = "error"
)

func New(level, destination string) (*slog.Logger, error) {
	var lvl slog.Level
	switch level {
	case LevelDebug:
		lvl = slog.LevelDebug
	case LevelInfo:
		lvl = slog.LevelInfo
	case LevelWarn:
		lvl = slog.LevelWarn
	case LevelError:
		lvl = slog.LevelError
	default:
		lvl = slog.LevelInfo
	}

	var output io.Writer
	var err error

	if destination == "stdout" {
		output = os.Stdout
	} else {
		output, err = os.OpenFile(destination, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			return nil, err
		}
	}

	return slog.New(slog.NewJSONHandler(output, &slog.HandlerOptions{Level: lvl})), nil
}

// NewDiscardLogger returns a logger that discards all messages.
func NewDiscardLogger() *slog.Logger {
	return slog.New(slog.NewJSONHandler(io.Discard, nil))
}
