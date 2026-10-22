package app

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/britlydev/c210cc/internal/config"
	"github.com/britlydev/c210cc/internal/service/calc"
	"github.com/britlydev/c210cc/internal/transport/httpapi"
	"github.com/britlydev/c210cc/internal/transport/httpapi/handler"
	"golang.org/x/sync/errgroup"
)

func RunAPI(ctx context.Context, cfg config.Config) error {
	deps, err := newDeps(ctx, cfg)
	if err != nil {
		return nil
	}
	defer deps.close()

	deps.logger.Info("api started", slog.String("env", cfg.Env))

	h := &handler.Handler{
		HealthHandler: handler.NewHealthHandler(cfg.Env),
		CalcHandler:   handler.NewCalcHandler(calc.New()),
	}

	srv := httpapi.NewServer(cfg.HTTP, h, deps.logger)

	g, gctx := errgroup.WithContext(ctx)
	g.Go(func() error {
		return srv.Run(gctx)
	})

	if err := g.Wait(); err != nil {
		return fmt.Errorf("run api: %w", err)
	}

	deps.logger.Info("api stopped")

	return nil
}
