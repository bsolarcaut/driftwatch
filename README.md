# driftwatch

> Detect infrastructure drift between Terraform state and live cloud resources — report diffs to Slack or PagerDuty.

---

## Installation

```bash
go install github.com/yourorg/driftwatch@latest
```

Or download a pre-built binary from the [releases page](https://github.com/yourorg/driftwatch/releases).

---

## Usage

Point `driftwatch` at your Terraform state and let it compare against your live cloud environment:

```bash
driftwatch scan \
  --state s3://my-bucket/terraform.tfstate \
  --provider aws \
  --region us-east-1 \
  --notify slack \
  --slack-webhook https://hooks.slack.com/services/XXX/YYY/ZZZ
```

To report drift to PagerDuty instead:

```bash
driftwatch scan \
  --state s3://my-bucket/terraform.tfstate \
  --provider aws \
  --notify pagerduty \
  --pagerduty-key YOUR_INTEGRATION_KEY
```

Run on a schedule (e.g. via cron or CI) to get continuous drift alerts before they become incidents.

### Configuration file

You can also use a config file instead of flags:

```yaml
# driftwatch.yaml
state: s3://my-bucket/terraform.tfstate
provider: aws
region: us-east-1
notify: slack
slack_webhook: https://hooks.slack.com/services/XXX/YYY/ZZZ
```

```bash
driftwatch scan --config driftwatch.yaml
```

---

## License

MIT © yourorg