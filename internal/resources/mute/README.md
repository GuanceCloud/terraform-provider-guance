# guance_mute

The `guance_mute` resource manages Guance mute rules. Mute rules can silence monitors, alert policies, monitor tags, or custom resource ranges.

## Example Usage

Mute an alert policy for a one-time window:

```hcl
resource "guance_mute" "alert_policy" {
  name        = "terraform-alert-policy-mute"
  description = "Managed by Terraform"
  type        = "alertPolicy"
  timezone    = "Asia/Shanghai"

  mute_ranges = [{
    name              = guance_alert_policy.example.name
    alert_policy_uuid = guance_alert_policy.example.uuid
  }]

  repeat_time_set = 0
  start_time      = "2026/06/11 10:00:00"
  end_time        = "2026/06/11 11:00:00"

  tags = {
    host = ["web001"]
  }
}
```

Repeated mute:

```hcl
resource "guance_mute" "weekly" {
  name     = "terraform-weekly-mute"
  type     = "alertPolicy"
  timezone = "Asia/Shanghai"

  mute_ranges = [{
    alert_policy_uuid = guance_alert_policy.example.uuid
  }]

  repeat_time_set = 1
  repeat_crontab_set = {
    min   = "0"
    hour  = "0"
    day   = "*"
    month = "*"
    week  = "1,2"
  }
  crontab_duration  = 18000
  repeat_expire_time = "0"
}
```

Notify before a mute starts:

```hcl
resource "guance_mute" "with_notify" {
  name     = "terraform-mute-with-notify"
  type     = "alertPolicy"
  timezone = "Asia/Shanghai"

  mute_ranges = [{
    alert_policy_uuid = guance_alert_policy.example.uuid
  }]

  repeat_time_set = 0
  start_time      = "2026/06/11 10:00:00"
  end_time        = "2026/06/11 11:00:00"
  notify_time_str = "2026/06/11 09:50:00"
  notify_message  = "Alert policy mute will start soon."

  notify_targets = [{
    type = "notifyObject"
    to   = [guance_notify_object.example.uuid]
  }]
}
```

## Argument Reference

* `name` - (Required) Mute rule name.
* `type` - (Required) Mute rule type. Valid values are `checker`, `alertPolicy`, `tag`, and `custom`.
* `mute_ranges` - (Required) Mute ranges. An empty list means all resources for the selected type.
* `description` - (Optional) Mute rule description.
* `tags` - (Optional) Event attribute filters. Prefix a key with `-` for negative matching.
* `filter_string` - (Optional) Event attribute filter expression. This has higher priority than `tags`.
* `notify_targets` - (Optional) Notification targets.
* `notify_message` - (Optional) Notification message.
* `notify_time_str` - (Optional) Notification time in `YYYY/MM/DD HH:mm:ss`.
* `start_time` - (Optional) One-time mute start time in `YYYY/MM/DD HH:mm:ss`.
* `end_time` - (Optional) One-time mute end time in `YYYY/MM/DD HH:mm:ss`.
* `repeat_time_set` - (Optional) `0` for one-time mute, `1` for repeated mute. Defaults to `0`.
* `repeat_crontab_set` - (Optional) Repeated mute crontab fields.
* `crontab_duration` - (Optional) Repeated mute duration in seconds.
* `repeat_expire_time` - (Optional) Repeated mute expiration time in `YYYY/MM/DD HH:mm:ss`, or `0` for never expires.
* `timezone` - (Optional) Mute rule timezone. Defaults to `Asia/Shanghai`.
* `declaration` - (Optional) Custom declaration information.
* `enabled` - (Optional) Whether the mute rule is enabled. Defaults to `true`.

### `mute_ranges`

* `name` - (Optional) Display name of the muted resource.
* `type` - (Optional) Resource type returned by the API.
* `checker_uuid` - (Optional) Monitor/checker UUID.
* `monitor_uuid` - (Optional) Monitor UUID.
* `slo_uuid` - (Optional) SLO UUID.
* `alert_policy_uuid` - (Optional) Alert policy UUID.
* `tag_uuid` - (Optional) Monitor tag UUID.

### `notify_targets`

* `type` - (Required) Notification target type, such as `mail` or `notifyObject`.
* `to` - (Required) Notification target UUIDs.

## Attribute Reference

* `uuid` - Mute rule UUID.
* `status` - Mute rule status returned by the API. Status `0` maps to `enabled = true`; status `2` maps to `enabled = false`.
* `create_at` - Creation timestamp in seconds.
* `update_at` - Last update timestamp in seconds.
* `workspace_uuid` - Workspace UUID.

## Data Source

The `guance_mute` data source reads an existing mute rule by `uuid` or exact `name`.

Lookup by name:

```hcl
data "guance_mute" "example" {
  name = "terraform-alert-policy-mute"
}

output "muted_alert_policy_uuid" {
  value = data.guance_mute.example.mute_ranges[0].alert_policy_uuid
}
```

Lookup by UUID:

```hcl
data "guance_mute" "example" {
  uuid = "mute_xxx"
}
```

Name lookup must match exactly one mute rule. The data source exports all resource attributes as read-only values, including `mute_ranges`, `tags`, `filter_string`, `notify_targets`, `notify_message`, `notify_time_str`, `start_time`, `end_time`, `repeat_time_set`, `repeat_crontab_set`, `crontab_duration`, `repeat_expire_time`, `timezone`, `declaration`, `enabled`, `status`, `create_at`, `update_at`, and `workspace_uuid`.

## Import

```shell
terraform import guance_mute.example mute_xxx
```
