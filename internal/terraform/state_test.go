package terraform_test

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/your-org/driftwatch/internal/terraform"
)

func writeTempState(t *testing.T, v interface{}) string {
	t.Helper()
	data, err := json.Marshal(v)
	if err != nil {
		t.Fatalf("marshal state: %v", err)
	}
	f, err := os.CreateTemp(t.TempDir(), "*.tfstate")
	if err != nil {
		t.Fatalf("create temp file: %v", err)
	}
	if _, err := f.Write(data); err != nil {
		t.Fatalf("write temp file: %v", err)
	}
	f.Close()
	return f.Name()
}

func TestLoadState_ValidState(t *testing.T) {
	raw := map[string]interface{}{
		"version": 4,
		"resources": []interface{}{
			map[string]interface{}{
				"type":     "aws_instance",
				"name":     "web",
				"provider": "provider[\"registry.terraform.io/hashicorp/aws\"]",
				"instances": []interface{}{
					map[string]interface{}{
						"attributes": map[string]interface{}{"id": "i-12345", "instance_type": "t3.micro"},
					},
				},
			},
		},
	}
	path := writeTempState(t, raw)
	state, err := terraform.LoadState(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if state.Version != 4 {
		t.Errorf("expected version 4, got %d", state.Version)
	}
	if len(state.Resources) != 1 {
		t.Fatalf("expected 1 resource, got %d", len(state.Resources))
	}
	if state.Resources[0].Type != "aws_instance" {
		t.Errorf("expected type aws_instance, got %s", state.Resources[0].Type)
	}
	if state.Resources[0].Attributes["id"] != "i-12345" {
		t.Errorf("unexpected id attribute: %v", state.Resources[0].Attributes["id"])
	}
}

func TestLoadState_MissingFile(t *testing.T) {
	_, err := terraform.LoadState("/nonexistent/path/terraform.tfstate")
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}

func TestLoadState_InvalidJSON(t *testing.T) {
	f, _ := os.CreateTemp(t.TempDir(), "*.tfstate")
	f.WriteString("not valid json")
	f.Close()
	_, err := terraform.LoadState(f.Name())
	if err == nil {
		t.Fatal("expected error for invalid JSON, got nil")
	}
}

func TestLoadState_EmptyResources(t *testing.T) {
	raw := map[string]interface{}{"version": 4, "resources": []interface{}{}}
	path := writeTempState(t, raw)
	state, err := terraform.LoadState(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(state.Resources) != 0 {
		t.Errorf("expected 0 resources, got %d", len(state.Resources))
	}
}
