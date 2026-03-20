# Alert Policy Resource

## Overview

The `guance_alert_policy` resource manages alert policies in Guance Cloud. Alert policies define how alerts are generated, aggregated, and delivered to notification targets.

## Core Functionality

- Create, read, update, and delete alert policies
- Configure alert aggregation settings
- Define notification targets and escalation rules
- Associate with monitors and security rules
- Support for both status-based and member-based alert types

## Applicable Scenarios

- Setting up monitoring and alerting for infrastructure
- Configuring escalation paths for critical alerts
- Managing alert policies across multiple environments
- Integrating with existing monitoring resources

## Configuration Parameters

| Parameter | Type | Required | Default | Description |
|-----------|------|----------|---------|-------------|
| `name` | string | Yes | - | The name of the alert policy |
| `desc` | string | No | "" | The description of the alert policy |
| `open_permission_set` | bool | No | false | Whether to open custom permission configuration |
| `permission_set` | list(string) | No | [] | Operation permission configuration |
| `checker_uuids` | list(string) | No | [] | Monitor/smart monitor/smart inspection/slo uuid |
| `security_rule_uuids` | list(string) | No | [] | Security monitoring (cspm, siem) uuid |
| `rule_timezone` | string | Yes | - | Alert policy timezone |
| `alert_opt` | block | No | - | Alert settings |
| `alert_opt.agg_type` | string | No | null | Alert aggregation type |
| `alert_opt.ignore_ok` | bool | No | false | Whether to ignore OK level alerts |
| `alert_opt.alert_type` | string | No | "status" | Alert policy notification type |
| `alert_opt.silent_timeout` | number | No | 0 | Minimum alert interval in seconds |
| `alert_opt.silent_timeout_by_status_enable` | bool | No | false | Whether to enable status-based silent timeout |
| `alert_opt.silent_timeout_by_status` | list(object) | No | [] | Status-based silent timeout configuration |
| `alert_opt.alert_target` | list(object) | No | [] | Notification target configuration |
| `alert_opt.agg_interval` | number | No | 0 | Alert aggregation interval in seconds |
| `alert_opt.agg_fields` | list(string) | No | [] | Aggregation field list |
| `alert_opt.agg_labels` | list(string) | No | [] | Label list for aggregation |
| `alert_opt.agg_cluster_fields` | list(string) | No | [] | Field list for smart aggregation |
| `alert_opt.agg_send_first` | bool | No | false | Whether to send first alert directly |

## Example Usage

### Basic Alert Policy

```hcl
resource "guance_alert_policy" "example" {
  name          = "High CPU Alert"
  desc          = "Alert when CPU usage exceeds threshold"
  rule_timezone = "Asia/Shanghai"

  alert_opt {
    agg_interval = 60
    agg_fields   = ["df_monitor_checker_id"]
    alert_target {
      targets {
        to     = ["notify_xxx"]
        status = "critical,error"
      }
    }
  }
}
```

### Alert Policy with Escalation

```hcl
resource "guance_alert_policy" "escalation_example" {
  name          = "Database Alert"
  desc          = "Alert on database connectivity issues"
  rule_timezone = "Asia/Shanghai"

  alert_opt {
    alert_type = "status"
    alert_target {
      name = "Database Alert Target"
      targets {
        to       = ["acnt_xxx", "group_xxx"]
        status   = "critical"
        upgrade_targets {
          to       = ["acnt_yyy"]
          duration = 300
        }
      }
    }
  }
}
```

## API Call Examples

### Create Alert Policy

```bash
curl 'https://openapi.example.com/api/v1/alert_policy/add_v2' \
  -H 'DF-API-KEY: <DF-API-KEY>' \
  -H 'Content-Type: application/json' \
  --data-raw '{"name":"High CPU Alert","ruleTimezone":"Asia/Shanghai","alertOpt":{"aggInterval":60,"aggFields":["df_monitor_checker_id"],"alertTarget":[{"targets":[{"to":["notify_xxx"],"status":"critical,error"}]}]}}'
```

### Update Alert Policy

```bash
curl 'https://openapi.example.com/api/v1/alert_policy/{uuid}/modify_v2' \
  -H 'DF-API-KEY: <DF-API-KEY>' \
  -H 'Content-Type: application/json' \
  --data-raw '{"name":"Updated Alert","ruleTimezone":"Asia/Shanghai","alertOpt":{"aggInterval":30,"alertTarget":[{"targets":[{"to":["notify_xxx","group_xxx"],"status":"critical"}]}]}}'
```

### Delete Alert Policy

```bash
curl 'https://openapi.example.com/api/v1/alert_policy/delete' \
  -H 'DF-API-KEY: <DF-API-KEY>' \
  -H 'Content-Type: application/json' \
  --data-raw '{"alertPolicyUUIDs":["altpl_xxx"]}'
```

## Common Issues and Solutions

### Resource Creation Failure

**Symptom**: Terraform apply fails with API error

**Possible Causes**:
- Invalid `rule_timezone` format
- Missing required parameters
- Permission denied for API key

**Solution**:
- Verify timezone format (e.g., "Asia/Shanghai")
- Check all required parameters are provided
- Ensure API key has sufficient permissions

### Permission Issues

**Symptom**: "Permission denied" errors

**Solution**:
- Enable `open_permission_set` and configure `permission_set`
- Include necessary roles, teams, or members in permission set
- Verify API key has appropriate workspace permissions

### Aggregation Configuration

**Symptom**: Alerts not aggregating as expected

**Solution**:
- Set `agg_interval` to a non-zero value
- Ensure `agg_fields` includes relevant fields
- For smart aggregation, include "CLUSTER" in `agg_fields`

### Notification Delivery

**Symptom**: Alerts not being delivered

**Solution**:
- Verify `alert_target` configuration
- Check notification target UUIDs are correct
- Ensure notification channels are properly configured in Guance Cloud

## Importing Existing Alert Policies

To import an existing alert policy into Terraform:

```bash
terraform import guance_alert_policy.example altpl_xxx
```

Replace `altpl_xxx` with the actual UUID of the alert policy.