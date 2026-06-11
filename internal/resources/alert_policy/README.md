# Alert Policy Resource

The `guance_alert_policy` resource manages Guance alert policies. Alert policies control alert aggregation, silence behavior, and notification targets.

## Example Usage

### Basic Alert Policy

```hcl
resource "guance_alert_policy" "example" {
  name          = "High CPU Alert"
  desc          = "Alert when CPU usage exceeds threshold"
  rule_timezone = "Asia/Shanghai"

  alert_opt = {
    alert_type     = "status"
    agg_interval   = 60
    agg_fields     = ["df_monitor_checker_id"]
    silent_timeout = 300

    alert_target = [{
      name = "Default notification"

      targets = [{
        to     = ["notify_xxx"]
        status = "critical,error"
      }]
    }]
  }
}
```

### Alert Policy with Escalation

```hcl
resource "guance_alert_policy" "escalation_example" {
  name          = "Database Alert"
  desc          = "Alert on database connectivity issues"
  rule_timezone = "Asia/Shanghai"

  alert_opt = {
    alert_type = "status"

    alert_target = [{
      name = "Database Alert Target"

      targets = [{
        to     = ["notify_xxx"]
        status = "critical"

        upgrade_targets = [{
          to       = ["notify_yyy"]
          duration = 300
        }]
      }]
    }]
  }
}
```

### Alert Policy with Custom Notice Dates

```hcl
resource "guance_alert_policy_notice_date" "holiday" {
  name                     = "Holiday notice dates"
  skip_ref_check_on_delete = false

  notice_dates = [
    "2026/01/01",
    "2026/05/01",
  ]
}

resource "guance_alert_policy" "custom_date_example" {
  name          = "Holiday Alert"
  rule_timezone = "Asia/Shanghai"

  alert_opt = {
    alert_type = "status"

    alert_target = [{
      name              = "Holiday notification"
      custom_date_uuids = [guance_alert_policy_notice_date.holiday.uuid]
      custom_start_time = "09:30:00"
      custom_duration   = 3600

      targets = [{
        to     = ["notify_xxx"]
        status = "critical,error"
      }]
    }]
  }
}
```

### Full Alert Chain

```hcl
resource "guance_notify_object" "ops" {
  name                = "Ops Webhook"
  type                = "http"
  opt_set             = jsonencode({ url = "https://example.com/alert" })
  open_permission_set = false
}

resource "guance_alert_policy_notice_date" "holiday" {
  name                     = "Holiday notice dates"
  skip_ref_check_on_delete = false

  notice_dates = [
    "2026/01/01",
    "2026/05/01",
  ]
}

resource "guance_alert_policy" "ops" {
  name          = "Ops Alert Policy"
  desc          = "Route severe production alerts"
  rule_timezone = "Asia/Shanghai"

  alert_opt = {
    alert_type                      = "status"
    agg_type                        = "byFields"
    agg_interval                    = 120
    agg_fields                      = ["df_monitor_checker_id", "df_label"]
    agg_labels                      = ["service"]
    agg_send_first                  = true
    ignore_ok                       = true
    silent_timeout                  = 300
    silent_timeout_by_status_enable = true

    silent_timeout_by_status = [{
      status         = "critical,error"
      silent_timeout = 600
    }]

    alert_target = [{
      name              = "Business hours"
      custom_date_uuids = [guance_alert_policy_notice_date.holiday.uuid]
      custom_start_time = "09:00:00"
      custom_duration   = 28800

      targets = [{
        to            = [guance_notify_object.ops.uuid]
        status        = "critical,error"
        filter_string = "`service` IN ['checkout']"

        upgrade_targets = [{
          to       = [guance_notify_object.ops.uuid]
          duration = 900
        }]
      }]
    }]
  }
}

resource "guance_mute" "maintenance" {
  name        = "Ops Alert Maintenance"
  type        = "alertPolicy"
  timezone    = "Asia/Shanghai"
  start_time  = "2026/06/12 01:00:00"
  end_time    = "2026/06/12 02:00:00"

  mute_ranges = [{
    alert_policy_uuid = guance_alert_policy.ops.uuid
  }]
}
```

### Alert Policy by Member

```hcl
data "guance_members" "all" {}

resource "guance_alert_policy" "member_example" {
  name          = "Member Alert"
  desc          = "Route alerts by member"
  rule_timezone = "Asia/Shanghai"

  alert_opt = {
    alert_type     = "member"
    agg_interval   = 60
    silent_timeout = 300

    alert_target = [{
      name = "Member notification"

      alert_info = [{
        name        = "Owner route"
        member_info = [data.guance_members.all.members[0].uuid]

        targets = [{
          to     = ["notify_xxx"]
          status = "critical,error,warning"
        }]
      }]
    }]
  }
}
```

## Argument Reference

| Name | Type | Required | Default | Description |
|------|------|----------|---------|-------------|
| `name` | string | Yes | - | Alert policy name. |
| `desc` | string | No | `null` | Alert policy description. |
| `open_permission_set` | bool | No | `false` | Whether to enable custom operation permissions. |
| `permission_set` | list(string) | No | `[]` | Role, member, or team UUIDs allowed to operate this policy. |
| `checker_uuids` | list(string) | No | `null` | Monitor, smart monitor, smart inspection, or SLO UUIDs associated with this policy. |
| `security_rule_uuids` | list(string) | No | `null` | Security rule UUIDs associated with this policy. |
| `rule_timezone` | string | Yes | - | Timezone used by the alert policy, for example `Asia/Shanghai`. |
| `alert_opt` | object | No | `null` | Alert delivery, aggregation, and silence settings. Use object assignment syntax: `alert_opt = { ... }`. |

## `alert_opt` Object

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `agg_type` | string | No | Alert aggregation type: `byFields`, `byCluster`, `byAI`, or `byCustom`. |
| `ignore_ok` | bool | No | Whether normal level only generates events and skips notifications. |
| `alert_type` | string | No | Notification type: `status` or `member`. |
| `silent_timeout` | number | No | Minimum repeated alert interval in seconds. |
| `silent_timeout_by_status_enable` | bool | No | Whether to enable status-specific silence intervals. |
| `silent_timeout_by_status` | list(object) | No | Status-specific silence interval list. |
| `alert_target` | list(object) | No | Notification target groups. |
| `agg_interval` | number | No | Alert aggregation interval in seconds. |
| `agg_fields` | list(string) | No | Aggregation field list. |
| `agg_labels` | list(string) | No | Aggregation labels. |
| `agg_cluster_fields` | list(string) | No | Smart aggregation field list. |
| `agg_send_first` | bool | No | Whether to send the first alert directly before aggregation. |

For `alert_type = "status"`, configure notification recipients under `alert_target.targets`.
For `alert_type = "member"`, configure member routing under `alert_target.alert_info`; each `alert_info.member_info` entry is a member UUID and its `targets` list defines recipients for that member route. In real OpenAPI validation, `agg_interval` is required for member mode.

The Forethought UI exposes alert policy enable/disable through `/alert_policy/set_disable`, but that route is not exported in the OpenAPI alert policy module used by this provider. Terraform currently manages create/read/update/delete and the exported v2 alert option fields.

## Data Source

The `guance_alert_policy` data source reads an existing alert policy by `uuid` or exact `name`.

Lookup by name:

```hcl
data "guance_alert_policy" "example" {
  name = "High CPU Alert"
}

output "alert_policy_uuid" {
  value = data.guance_alert_policy.example.uuid
}

output "first_notify_object" {
  value = data.guance_alert_policy.example.alert_opt.alert_target[0].targets[0].to[0]
}
```

Lookup by UUID:

```hcl
data "guance_alert_policy" "example" {
  uuid = "altpl_xxx"
}
```

Name lookup must match exactly one alert policy. The data source exports:

| Name | Type | Description |
|------|------|-------------|
| `uuid` | string | Alert policy UUID. |
| `name` | string | Alert policy name. |
| `desc` | string | Alert policy description. |
| `open_permission_set` | bool | Whether custom operation permissions are enabled. |
| `permission_set` | list(string) | Role, member, or team UUIDs allowed to operate this policy. |
| `checker_uuids` | list(string) | Monitor, smart monitor, smart inspection, or SLO UUIDs associated with this policy. |
| `security_rule_uuids` | list(string) | Security rule UUIDs associated with this policy. |
| `rule_timezone` | string | Timezone used by the alert policy. |
| `alert_opt` | object | Alert delivery, aggregation, silence, and target settings. |
| `create_at` | number | Created timestamp in seconds. |
| `update_at` | number | Updated timestamp in seconds. |
| `workspace_uuid` | string | Workspace UUID. |

## Import

```bash
terraform import guance_alert_policy.example altpl_xxx
```
