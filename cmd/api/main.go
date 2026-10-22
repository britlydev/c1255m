package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/britlydev/c210cc/internal/app"
	"github.com/britlydev/c210cc/internal/config"
	"github.com/britlydev/c210cc/internal/version"
)

func main() {
	fmt.Println("Math API")
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run() error {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	configPath := os.Getenv("C1255M_CONFIG_PATH")
	if configPath == "" {
		configPath = "configs/config.yaml"
	}
	flag.StringVar(&configPath, "config", configPath, "path to config file")

	var showVersion bool
	flag.BoolVar(&showVersion, "version", false, "print version and exit")

	flag.Parse()

	if showVersion {
		v := version.Get()
		fmt.Printf("version: %s\ncommit:  %s\ndate:    %s\ngo:      %s\n", v.Version, v.Commit, v.Date, v.GoVersion)
		return nil
	}

	cfg, err := config.Load(configPath)
	if err != nil {
		return err
	}

	errCh := make(chan error, 1)
	go func() {
		errCh <- app.RunAPI(ctx, cfg)
	}()

	select {
	case err := <-errCh:
		return err
	case <-ctx.Done():
	}

	select {
	case err := <-errCh:
		return err
	case <-time.After(30 * time.Second):
		return errors.New("graceful shutdown time out")
	}
}
