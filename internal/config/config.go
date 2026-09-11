package config

import (
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	Environment         string
	HTTPAddr            string
	DatabaseURL         string
	LogLevel            slog.Level
	SessionCookieSecure bool
}

func Load() (Config, error) {
	return load(os.LookupEnv)
}

func load(lookup func(string) (string, bool)) (Config, error) {
	cfg := Config{HTTPAddr: ":8080", LogLevel: slog.LevelInfo}

	var ok bool
	cfg.Environment, ok = lookup("APP_ENV")
	if !ok || strings.TrimSpace(cfg.Environment) == "" {
		return Config{}, fmt.Errorf("APP_ENV is required")
	}
	if cfg.Environment != "development" && cfg.Environment != "test" && cfg.Environment != "production" {
		return Config{}, fmt.Errorf("APP_ENV must be development, test or production")
	}

	if value, exists := lookup("HTTP_ADDR"); exists && strings.TrimSpace(value) != "" {
		cfg.HTTPAddr = value
	}

	cfg.DatabaseURL, ok = lookup("DATABASE_URL")
	if !ok || strings.TrimSpace(cfg.DatabaseURL) == "" {
		return Config{}, fmt.Errorf("DATABASE_URL is required")
	}

	if value, exists := lookup("LOG_LEVEL"); exists {
		level, err := parseLogLevel(value)
		if err != nil {
			return Config{}, err
		}
		cfg.LogLevel = level
	}

	defaultSecure := cfg.Environment == "production"
	cfg.SessionCookieSecure = defaultSecure
	if value, exists := lookup("SESSION_COOKIE_SECURE"); exists {
		secure, err := strconv.ParseBool(value)
		if err != nil {
			return Config{}, fmt.Errorf("SESSION_COOKIE_SECURE must be true or false")
		}
		cfg.SessionCookieSecure = secure
	}
	if cfg.Environment == "production" && !cfg.SessionCookieSecure {
		return Config{}, fmt.Errorf("SESSION_COOKIE_SECURE must be true in production")
	}

	return cfg, nil
}

func parseLogLevel(value string) (slog.Level, error) {
	switch strings.ToLower(value) {
	case "debug":
		return slog.LevelDebug, nil
	case "info", "":
		return slog.LevelInfo, nil
	case "warn":
		return slog.LevelWarn, nil
	case "error":
		return slog.LevelError, nil
	default:
		return 0, fmt.Errorf("LOG_LEVEL must be debug, info, warn or error")
	}
}
