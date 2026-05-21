package config

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

// Config holds all driftwatch configuration.
type Config struct {
	Terraform TerraformConfig `yaml:"terraform"`
	Slack     SlackConfig     `yaml:"slack"`
	PagerDuty PagerDutyConfig `yaml:"pagerduty"`
	Scan      ScanConfig      `yaml:"scan"`
}

// TerraformConfig holds Terraform-related settings.
type TerraformConfig struct {
	StatePath   string `yaml:"state_path"`
	Workspace   string `yaml:"workspace"`
	BackendType string `yaml:"backend_type"` // local, s3, gcs
}

// SlackConfig holds Slack notification settings.
type SlackConfig struct {
	Enabled    bool   `yaml:"enabled"`
	WebhookURL string `yaml:"webhook_url"`
	Channel    string `yaml:"channel"`
}

// PagerDutyConfig holds PagerDuty notification settings.
type PagerDutyConfig struct {
	Enabled        bool   `yaml:"enabled"`
	IntegrationKey string `yaml:"integration_key"`
	Severity       string `yaml:"severity"` // critical, error, warning, info
}

// ScanConfig holds drift scan behaviour settings.
type ScanConfig struct {
	Interval        time.Duration `yaml:"interval"`
	ResourceTypes   []string      `yaml:"resource_types"`
	IgnoreResources []string      `yaml:"ignore_resources"`
}

// Load reads a YAML config file from the given path.
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading config file %q: %w", path, err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parsing config file %q: %w", path, err)
	}

	if err := cfg.validate(); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	return &cfg, nil
}

// validate performs basic sanity checks on the loaded config.
func (c *Config) validate() error {
	if c.Terraform.StatePath == "" {
		return fmt.Errorf("terraform.state_path must not be empty")
	}
	if c.PagerDuty.Enabled && c.PagerDuty.IntegrationKey == "" {
		return fmt.Errorf("pagerduty.integration_key required when pagerduty is enabled")
	}
	if c.Slack.Enabled && c.Slack.WebhookURL == "" {
		return fmt.Errorf("slack.webhook_url required when slack is enabled")
	}
	if c.Scan.Interval == 0 {
		c.Scan.Interval = 5 * time.Minute
	}
	return nil
}
