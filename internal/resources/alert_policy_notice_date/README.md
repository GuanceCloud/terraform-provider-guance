# Alert Policy Notice Date Resource

The `guance_alert_policy_notice_date` resource manages custom notice dates for Guance alert policies. These dates can be referenced by `guance_alert_policy.alert_opt.alert_target.custom_date_uuids`.

## Example Usage

```hcl
resource "guance_alert_policy_notice_date" "example" {
  name                     = "Holiday notice dates"
  skip_ref_check_on_delete = false

  notice_dates = [
    "2026/01/01",
    "2026/05/01",
  ]
}
```

## Argument Reference

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `name` | string | Yes | The custom notice date name. The backend stores up to 64 characters. |
| `notice_dates` | list(string) | Yes | Custom notice dates. Each value must use `YYYY/MM/DD` format. Up to 366 dates are allowed. |
| `skip_ref_check_on_delete` | bool | No | Whether deletion bypasses backend reference checks. Defaults to `true` for compatibility. Set to `false` to let Guance reject deletion while the date is referenced by an alert policy. |

## Attribute Reference

| Name | Type | Description |
|------|------|-------------|
| `uuid` | string | The notice date UUID. |
| `create_at` | number | Created timestamp in seconds. |
| `update_at` | number | Updated timestamp in seconds. |
| `workspace_uuid` | string | Workspace UUID. |

## Data Source

The `guance_alert_policy_notice_date` data source reads an existing custom notice date by `uuid` or exact `name`.

Lookup by name:

```hcl
data "guance_alert_policy_notice_date" "holiday" {
  name = "Holiday notice dates"
}

resource "guance_alert_policy" "example" {
  name          = "Holiday Alert"
  rule_timezone = "Asia/Shanghai"

  alert_opt = {
    alert_type = "status"

    alert_target = [{
      custom_date_uuids = [data.guance_alert_policy_notice_date.holiday.uuid]

      targets = [{
        to     = ["notify_xxx"]
        status = "critical,error"
      }]
    }]
  }
}
```

Lookup by UUID:

```hcl
data "guance_alert_policy_notice_date" "holiday" {
  uuid = "ndate_xxx"
}
```

Name lookup must match exactly one notice date. The data source exports `uuid`, `name`, `notice_dates`, `create_at`, `update_at`, and `workspace_uuid`.

## Import

```bash
terraform import guance_alert_policy_notice_date.example ndate_xxx
```
