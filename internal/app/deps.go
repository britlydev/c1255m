package app

import (
	"context"
	"log/slog"
	"os"

	"github.com/britlydev/c210cc/internal/config"
	"github.com/britlydev/c210cc/internal/platform/logger"
)

type deps struct {
	logger *slog.Logger
}

func newDeps(ctx context.Context, cfg config.Config) (*deps, error) {
	logger, err := logger.New(os.Stdout, cfg.Log)
	if err != nil {
		return nil, err
	}

	return &deps{logger: logger}, nil
}

func (d *deps) close() {

}
