package drift

import (
	"strings"
	"testing"
)

func TestFormatSlack_NoDiffs(t *testing.T) {
	out := FormatSlack("i-1234", nil)
	if out != "" {
		t.Errorf("expected empty string, got %q", out)
	}
}

func TestFormatSlack_WithDiffs(t *testing.T) {
	diffs := []Diff{
		{Attribute: "instance_type", StateValue: "t2.micro", LiveValue: "t3.small"},
		{Attribute: "ami", StateValue: "ami-aaa", LiveValue: "ami-bbb"},
	}
	out := FormatSlack("i-1234", diffs)
	if !strings.Contains(out, "i-1234") {
		t.Error("expected resource ID in output")
	}
	if !strings.Contains(out, "instance_type") {
		t.Error("expected attribute name in output")
	}
	if !strings.Contains(out, "t2.micro") || !strings.Contains(out, "t3.small") {
		t.Error("expected state and live values in output")
	}
	if !strings.Contains(out, "ami-aaa") || !strings.Contains(out, "ami-bbb") {
		t.Error("expected ami diff in output")
	}
}

func TestFormatPagerDuty_NoDiffs(t *testing.T) {
	out := FormatPagerDuty("i-1234", []Diff{})
	if out != "" {
		t.Errorf("expected empty string, got %q", out)
	}
}

func TestFormatPagerDuty_WithDiffs(t *testing.T) {
	diffs := []Diff{
		{Attribute: "instance_type", StateValue: "t2.micro", LiveValue: "t3.small"},
	}
	out := FormatPagerDuty("i-1234", diffs)
	if !strings.HasPrefix(out, "drift:i-1234") {
		t.Errorf("expected prefix 'drift:i-1234', got %q", out)
	}
	if !strings.Contains(out, "instance_type(t2.micro->t3.small)") {
		t.Errorf("expected formatted diff in output, got %q", out)
	}
}

func TestFormatPagerDuty_MultipleDiffs(t *testing.T) {
	diffs := []Diff{
		{Attribute: "ami", StateValue: "ami-old", LiveValue: "ami-new"},
		{Attribute: "tags", StateValue: "env=prod", LiveValue: "env=staging"},
	}
	out := FormatPagerDuty("i-5678", diffs)
	if !strings.Contains(out, "ami(ami-old->ami-new)") {
		t.Errorf("missing ami diff in output: %q", out)
	}
	if !strings.Contains(out, "tags(env=prod->env=staging)") {
		t.Errorf("missing tags diff in output: %q", out)
	}
}
