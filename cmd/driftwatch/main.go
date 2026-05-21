package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"

	"github.com/yourorg/driftwatch/internal/config"
)

const defaultConfigPath = "driftwatch.yaml"

func main() {
	os.Exit(run())
}

func run() int {
	configPath := flag.String("config", defaultConfigPath, "path to driftwatch config file")
	logLevel := flag.String("log-level", "info", "log level: debug, info, warn, error")
	flag.Parse()

	level, err := parseLogLevel(*logLevel)
	if err != nil {
		fmt.Fprintf(os.Stderr, "invalid log level %q: %v\n", *logLevel, err)
		return 1
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: level}))
	slog.SetDefault(logger)

	cfg, err := config.Load(*configPath)
	if err != nil {
		slog.Error("failed to load configuration", "path", *configPath, "error", err)
		return 1
	}

	slog.Info("driftwatch starting",
		"state_path", cfg.Terraform.StatePath,
		"workspace", cfg.Terraform.Workspace,
		"scan_interval", cfg.Scan.Interval.String(),
		"slack_enabled", cfg.Slack.Enabled,
		"pagerduty_enabled", cfg.PagerDuty.Enabled,
	)

	// TODO: wire up drift scanner and notifiers.
	slog.Info("startup complete — drift scanner not yet implemented")
	return 0
}

func parseLogLevel(s string) (slog.Level, error) {
	switch s {
	case "debug":
		return slog.LevelDebug, nil
	case "info":
		return slog.LevelInfo, nil
	case "warn", "warning":
		return slog.LevelWarn, nil
	case "error":
		return slog.LevelError, nil
	default:
		return slog.LevelInfo, fmt.Errorf("unknown level %q", s)
	}
}
