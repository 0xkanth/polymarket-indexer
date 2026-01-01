// Package util provides initialization utilities for logger and configuration.
package util

import (
	"os"
	"strings"

	"github.com/knadh/koanf/parsers/toml"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
	"github.com/rs/zerolog"
)

// InitLogger initializes and returns a zerolog logger based on configuration.
// It supports both JSON (production) and pretty console (development) output.
func InitLogger() *zerolog.Logger {
	// Default to info level
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	// Create logger with timestamp
	var logger zerolog.Logger

	// Check if we're in a terminal for pretty output
	if isTerminal() {
		// Pretty console output for development
		logger = zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout}).
			With().
			Timestamp().
			Caller().
			Logger()
	} else {
		// JSON output for production
		logger = zerolog.New(os.Stdout).
			With().
			Timestamp().
			Str("service", "polymarket-indexer").
			Logger()
	}

	return &logger
}

// InitConfig initializes and returns a koanf configuration instance.
// It loads configuration from the TOML file and allows environment variable overrides.
func InitConfig(logger *zerolog.Logger, configPath string) *koanf.Koanf {
	ko := koanf.New(".")

	// Load configuration from TOML file
	if err := ko.Load(file.Provider(configPath), toml.Parser()); err != nil {
		logger.Fatal().
			Err(err).
			Str("path", configPath).
			Msg("failed to load config file")
	}

	// Load environment variables with prefix handling
	// Environment variables like CHAIN_RPC_ENDPOINT will override chain.rpc_endpoint
	if err := ko.Load(env.Provider("", ".", func(s string) string {
		// Convert CHAIN_RPC_ENDPOINT to chain.rpc_endpoint
		return strings.Replace(strings.ToLower(s), "_", ".", -1)
	}), nil); err != nil {
		logger.Warn().
			Err(err).
			Msg("failed to load environment variables")
	}

	logger.Info().
		Str("config_file", configPath).
		Msg("configuration loaded successfully")

	return ko
}

// UpdateLogLevel updates the global log level based on configuration.
func UpdateLogLevel(ko *koanf.Koanf, logger *zerolog.Logger) {
	levelStr := ko.String("logging.level")
	if levelStr == "" {
		levelStr = "info"
	}

	var level zerolog.Level
	switch strings.ToLower(levelStr) {
	case "debug":
		level = zerolog.DebugLevel
	case "info":
		level = zerolog.InfoLevel
	case "warn", "warning":
		level = zerolog.WarnLevel
	case "error":
		level = zerolog.ErrorLevel
	default:
		level = zerolog.InfoLevel
		logger.Warn().
			Str("configured_level", levelStr).
			Str("using_level", "info").
			Msg("unknown log level, defaulting to info")
	}

	zerolog.SetGlobalLevel(level)
	logger.Info().
		Str("level", level.String()).
		Msg("log level set")
}

// isTerminal checks if stdout is a terminal (for pretty console output).
func isTerminal() bool {
	fileInfo, _ := os.Stdout.Stat()
	return (fileInfo.Mode() & os.ModeCharDevice) != 0
}
