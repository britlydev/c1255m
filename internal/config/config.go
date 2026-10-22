package config

import (
	"errors"
	"fmt"
	"log/slog"
	"slices"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

const (
	EnvDevelopment = "development"
	EnvStaging     = "staging"
	EnvProduction  = "production"
)

type Config struct {
	Env  string `yaml:"env" env:"C1255M_ENV" env-default:"development"`
	Log  Log    `yaml:"log"`
	HTTP HTTP   `yaml:"http"`
}

type Log struct {
	Level  string `yaml:"level" env:"C1255M_LOG_LEVEL" env-default:"info"`
	Format string `yaml:"format" env:"C1255M_LOG_FORMAT" env-default:"json"`
}

type HTTP struct {
	Addr            string        `yaml:"addr" env:"C1255M_HTTP_ADDR" env-default:"127.0.0.1:4000"`
	ReadTimeout     time.Duration `yaml:"read_timeout" env-default:"10s"`
	WriteTimeout    time.Duration `yaml:"write_timeout" env-default:"15s"`
	IdleTimeout     time.Duration `yaml:"idle_timeout" env-default:"60s"`
	ShutdownTimeout time.Duration `yaml:"shutdown_timeout" env-default:"15s"`
}

func Load(path string) (Config, error) {
	var cfg Config
	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		return Config{}, fmt.Errorf("read config %q: %w", path, err)
	}

	return cfg, nil
}

func (c Config) Validate() error {
	var errs []error

	// env
	if !slices.Contains([]string{EnvDevelopment, EnvProduction, EnvStaging}, c.Env) {
		errs = append(errs, fmt.Errorf("env: unknown value %q", c.Env))
	}

	// log
	if _, err := c.Log.SlogLevel(); err != nil {
		errs = append(errs, err)
	}
	if !slices.Contains([]string{"text", "json"}, c.Log.Level) {
		errs = append(errs, fmt.Errorf("log.format: unknown value %q", c.Log.Format))
	}

	// http
	if c.HTTP.Addr == "" {
		errs = append(errs, fmt.Errorf("http.add: must not be empty"))
	}
	for name, d := range map[string]time.Duration{
		"http.read_timeout":     c.HTTP.ReadTimeout,
		"http.write_timeout":    c.HTTP.WriteTimeout,
		"http.idle_timeout":     c.HTTP.IdleTimeout,
		"http.shutdown_timeout": c.HTTP.ShutdownTimeout,
	} {
		if d <= 0 {
			errs = append(errs, fmt.Errorf("%s: must be positive", name))
		}
	}

	return errors.Join(errs...)
}

func (l Log) SlogLevel() (slog.Level, error) {
	var level slog.Level
	if err := level.UnmarshalText([]byte(l.Level)); err != nil {
		return 0, fmt.Errorf("log.level: %w", err)
	}

	return level, nil
}
