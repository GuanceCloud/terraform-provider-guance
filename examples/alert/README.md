# Alert Example

This example creates a complete alert delivery chain:

- `guance_notify_object` for a webhook notification target
- `guance_alert_policy_notice_date` for a reusable custom notice date calendar
- `guance_alert_policy` for status-based alert routing
- `guance_mute` for a one-time maintenance mute
- Matching data sources to read the created resources back by exact name

## Prerequisites

1. Terraform 1.0+ installed
2. Guance Cloud account with API access
3. `GUANCE_ACCESS_TOKEN` configured in your environment, or `access_token` configured in `provider.tf`

## Usage

```bash
terraform init
terraform plan
terraform apply
terraform destroy
```

## Configuration

Override the default names if you want to run the example more than once in the same workspace:

```bash
terraform apply \
  -var='name_prefix=my-alert-example'
```

The default notification object uses `simpleHTTPRequest` and posts to `https://example.com/guance-alert-example`. Replace `var.webhook_url` with a real endpoint before using this as a production alert route.

