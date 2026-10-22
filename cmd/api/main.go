package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/britlydev/c210cc/internal/config"
)

func main() {
	fmt.Println("Math API")
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run() error {
	configPath := os.Getenv("C1255M_CONFIG_PATH")
	if configPath == "" {
		configPath = "configs/config.yaml"
	}
	flag.StringVar(&configPath, "config", configPath, "path to config file")
	flag.Parse()

	cfg, err := config.Load(configPath)
	if err != nil {
		return err
	}

	_ = cfg

	return nil
}
