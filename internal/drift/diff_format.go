package drift

import (
	"fmt"
	"strings"
)

// FormatSlack formats a slice of Diffs into a Slack-friendly markdown message.
func FormatSlack(diffs []Diff) string {
	if len(diffs) == 0 {
		return ":white_check_mark: No infrastructure drift detected."
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf(":warning: *Infrastructure Drift Detected* — %d difference(s)\n", len(diffs)))
	for _, d := range diffs {
		sb.WriteString(fmt.Sprintf("• `%s.%s` → `%s`: state=`%v` live=`%v`\n",
			d.ResourceType, d.ResourceName, d.Attribute, d.StateValue, d.LiveValue))
	}
	return sb.String()
}

// FormatPagerDuty returns a summary string suitable for a PagerDuty alert.
func FormatPagerDuty(diffs []Diff) (summary string, details string) {
	if len(diffs) == 0 {
		return "", ""
	}
	summary = fmt.Sprintf("Infrastructure drift detected: %d resource(s) differ from Terraform state", len(diffs))

	var sb strings.Builder
	for _, d := range diffs {
		sb.WriteString(d.String())
		sb.WriteString("\n")
	}
	return summary, sb.String()
}
