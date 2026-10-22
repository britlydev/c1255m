package httpapi

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/britlydev/c210cc/internal/config"
	"github.com/britlydev/c210cc/internal/transport/httpapi/oapi"
	"github.com/gofiber/fiber/v3"
)

type Server struct {
	app             *fiber.App
	addr            string
	shutdownTimeout time.Duration
	logger          *slog.Logger
}

func NewServer(cfg config.HTTP, h oapi.StrictServerInterface, logger *slog.Logger) *Server {
	app := fiber.New(fiber.Config{
		AppName:      "Math API",
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		IdleTimeout:  cfg.IdleTimeout,
	})

	oapi.RegisterHandlers(app, oapi.NewStrictHandler(h, nil))

	return &Server{
		app:             app,
		addr:            cfg.Addr,
		shutdownTimeout: cfg.ShutdownTimeout,
		logger:          logger,
	}
}

func (s *Server) Run(ctx context.Context) error {
	listenErr := make(chan error, 1)
	go func() {
		s.logger.Info("http server started", slog.String("addr", s.addr))
		listenErr <- s.app.Listen(s.addr, fiber.ListenConfig{DisableStartupMessage: true})
	}()

	select {
	case err := <-listenErr:
		if err != nil {
			return fmt.Errorf("http listen: %w", err)
		}
		return nil
	case <-ctx.Done():
	}

	s.logger.Info("shutting down http server")

	shutdownCtx, cancel := context.WithTimeout(ctx, s.shutdownTimeout)
	defer cancel()

	if err := s.app.ShutdownWithContext(shutdownCtx); err != nil {
		return fmt.Errorf("http shutdown: %w", err)
	}

	s.logger.Info("http server stopped")

	return nil
}
