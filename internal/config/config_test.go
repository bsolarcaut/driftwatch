package config

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func writeTempConfig(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "driftwatch.yaml")
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatalf("writing temp config: %v", err)
	}
	return path
}

func TestLoad_ValidConfig(t *testing.T) {
	content := `
terraform:
  state_path: ./terraform.tfstate
  workspace: default
  backend_type: local
slack:
  enabled: true
  webhook_url: https://hooks.slack.com/services/TEST
  channel: "#alerts"
pagerduty:
  enabled: false
  integration_key: ""
  severity: warning
scan:
  interval: 10m
  resource_types:
    - aws_instance
    - aws_s3_bucket
`
	path := writeTempConfig(t, content)
	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Terraform.StatePath != "./terraform.tfstate" {
		t.Errorf("expected state_path './terraform.tfstate', got %q", cfg.Terraform.StatePath)
	}
	if cfg.Scan.Interval != 10*time.Minute {
		t.Errorf("expected interval 10m, got %v", cfg.Scan.Interval)
	}
	if len(cfg.Scan.ResourceTypes) != 2 {
		t.Errorf("expected 2 resource types, got %d", len(cfg.Scan.ResourceTypes))
	}
}

func TestLoad_DefaultInterval(t *testing.T) {
	content := `
terraform:
  state_path: ./terraform.tfstate
`
	path := writeTempConfig(t, content)
	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Scan.Interval != 5*time.Minute {
		t.Errorf("expected default interval 5m, got %v", cfg.Scan.Interval)
	}
}

func TestLoad_MissingStatePath(t *testing.T) {
	content := `terraform: {}`
	path := writeTempConfig(t, content)
	_, err := Load(path)
	if err == nil {
		t.Fatal("expected error for missing state_path, got nil")
	}
}

func TestLoad_SlackEnabledMissingURL(t *testing.T) {
	content := `
terraform:
  state_path: ./terraform.tfstate
slack:
  enabled: true
`
	path := writeTempConfig(t, content)
	_, err := Load(path)
	if err == nil {
		t.Fatal("expected error for missing slack webhook_url, got nil")
	}
}

func TestLoad_FileNotFound(t *testing.T) {
	_, err := Load("/nonexistent/path/config.yaml")
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}
