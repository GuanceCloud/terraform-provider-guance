# Monitor Example

This example demonstrates how to create a Guance monitor/checker with the `guance_monitor` resource.

Use `guance_monitor` when you want Terraform to manage the checker through structured fields such as `json_script`, `status`, `tags`, alert policy bindings, and operation permissions. Use `guance_monitor_json` instead when you want to manage an exported checker JSON document directly.

## Prerequisites

1. Terraform 1.0+ installed.
2. Guance Cloud account with OpenAPI access.
3. `GUANCE_ACCESS_TOKEN` configured in your environment, or `access_token` configured in `provider.tf`.

## Usage

```bash
terraform init
terraform plan
terraform apply
terraform destroy
```

## Notes

- `json_script` uses Terraform object syntax and maps to the Forethought OpenAPI `jsonScript` payload.
- `extend` is a JSON string used by Forethought for issue-related and frontend echo fields.
- The backend may add frontend echo fields to `extend`; Terraform keeps the configured subset stable while still detecting changes to fields you configured.
- `alert_policy_uuids` can bind the monitor/checker to alert policies.
- `secret` should be unique in the workspace when used. Clearing an existing non-empty `secret` with `secret = ""` currently depends on a pending Forethought OpenAPI adjustment; avoid using empty string as a clear operation for now.
