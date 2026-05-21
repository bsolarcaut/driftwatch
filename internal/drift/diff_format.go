package drift

import (
	"fmt"
	"strings"
)

// Diff represents a single attribute difference between state and live resource.
type Diff struct {
	Attribute  string
	StateValue string
	LiveValue  string
}

// FormatSlack formats a slice of Diffs for a Slack message block.
func FormatSlack(resourceID string, diffs []Diff) string {
	if len(diffs) == 0 {
		return ""
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("*Drift detected in `%s`*\n", resourceID))
	for _, d := range diffs {
		sb.WriteString(fmt.Sprintf("  • `%s`: `%s` → `%s`\n", d.Attribute, d.StateValue, d.LiveValue))
	}
	return sb.String()
}

// FormatPagerDuty formats a slice of Diffs into a PagerDuty alert summary string.
func FormatPagerDuty(resourceID string, diffs []Diff) string {
	if len(diffs) == 0 {
		return ""
	}
	parts := make([]string, 0, len(diffs))
	for _, d := range diffs {
		parts = append(parts, fmt.Sprintf("%s(%s->%s)", d.Attribute, d.StateValue, d.LiveValue))
	}
	return fmt.Sprintf("drift:%s [%s]", resourceID, strings.Join(parts, ", "))
}
