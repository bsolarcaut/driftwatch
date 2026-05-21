package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/yourusername/driftwatch/internal/drift"
)

// SlackPayload represents the JSON body sent to a Slack webhook.
type SlackPayload struct {
	Text string `json:"text"`
}

// PagerDutyPayload represents a simplified PagerDuty Events API v2 payload.
type PagerDutyPayload struct {
	RoutingKey  string            `json:"routing_key"`
	EventAction string            `json:"event_action"`
	Payload     PagerDutyDetail   `json:"payload"`
}

type PagerDutyDetail struct {
	Summary  string `json:"summary"`
	Severity string `json:"severity"`
	Source   string `json:"source"`
}

// Notifier sends drift reports to configured destinations.
type Notifier struct {
	client         *http.Client
	slackURL       string
	pagerDutyKey   string
	pagerDutyURL   string
}

// New creates a Notifier with the given Slack webhook URL and PagerDuty routing key.
func New(slackURL, pagerDutyKey string) *Notifier {
	return &Notifier{
		client:       &http.Client{Timeout: 10 * time.Second},
		slackURL:     slackURL,
		pagerDutyKey: pagerDutyKey,
		pagerDutyURL: "https://events.pagerduty.com/v2/enqueue",
	}
}

// NotifySlack sends a drift report to Slack if a webhook URL is configured.
func (n *Notifier) NotifySlack(diffs []drift.Diff) error {
	if n.slackURL == "" {
		return nil
	}
	message := drift.FormatSlack(diffs)
	payload := SlackPayload{Text: message}
	return n.post(n.slackURL, payload)
}

// NotifyPagerDuty sends a drift report to PagerDuty if a routing key is configured.
func (n *Notifier) NotifyPagerDuty(diffs []drift.Diff) error {
	if n.pagerDutyKey == "" {
		return nil
	}
	summary := drift.FormatPagerDuty(diffs)
	payload := PagerDutyPayload{
		RoutingKey:  n.pagerDutyKey,
		EventAction: "trigger",
		Payload: PagerDutyDetail{
			Summary:  summary,
			Severity: "warning",
			Source:   "driftwatch",
		},
	}
	return n.post(n.pagerDutyURL, payload)
}

func (n *Notifier) post(url string, payload any) error {
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("notify: marshal payload: %w", err)
	}
	resp, err := n.client.Post(url, "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("notify: http post: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		return fmt.Errorf("notify: unexpected status %d from %s", resp.StatusCode, url)
	}
	return nil
}
