package logger

import (
	"fmt"
	"io"
	"log/slog"

	"github.com/britlydev/c210cc/internal/config"
)

func New(w io.Writer, cfg config.Log) (*slog.Logger, error) {
	level, err := cfg.SlogLevel()
	if err != nil {
		return nil, err
	}

	opts := &slog.HandlerOptions{Level: level}

	var handler slog.Handler
	switch cfg.Format {
	case "json":
		handler = slog.NewJSONHandler(w, opts)
	case "text":
		handler = slog.NewTextHandler(w, opts)
	default:
		return nil, fmt.Errorf("unknown log format %q", cfg.Format)
	}

	return slog.New(handler), nil
}
