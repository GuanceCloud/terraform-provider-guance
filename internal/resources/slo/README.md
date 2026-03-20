# SLO Resource

The `guance_slo` resource allows you to manage SLO (Service Level Objective) resources in Guance Cloud.

## Example Usage

```hcl
resource "guance_slo" "example" {
  name             = "Example SLO"
  interval         = "5m"
  goal             = 99.9
  min_goal         = 99.0
  sli_uuids        = ["rul-aaaaaa", "rul-bbbbbb"]
  describe         = "This is an example SLO"
  alert_policy_uuids = ["altpl-xxxxxx"]
  tags             = ["example", "terraform"]
}
```

## Argument Reference

The following arguments are supported:

### Required

- `name` - (String) SLO name. Maximum length: 256 characters.
- `interval` - (String) Detection frequency. Valid values: `5m`, `10m`.
- `goal` - (Float) SLO expected goal. Range: 0-100.
- `min_goal` - (Float) SLO minimum goal. Range: 0-100, must be less than goal.
- `sli_uuids` - (List of String) SLI UUID list.

### Optional

- `describe` - (String) SLO description. Maximum length: 3000 characters.
- `alert_policy_uuids` - (List of String) Alert policy UUIDs.
- `tags` - (List of String) Tag names for filtering.

## Attribute Reference

The following attributes are exported:

- `uuid` - (String) The UUID of the SLO.
- `name` - (String) The name of the SLO.
- `interval` - (String) The detection frequency of the SLO.
- `goal` - (Float) The expected goal of the SLO.
- `min_goal` - (Float) The minimum goal of the SLO.
- `sli_uuids` - (List of String) The SLI UUID list of the SLO.
- `describe` - (String) The description of the SLO.
- `alert_policy_uuids` - (List of String) The alert policy UUIDs of the SLO.
- `tags` - (List of String) The tags of the SLO.
- `create_at` - (Int64) The timestamp seconds of the resource created at.
- `update_at` - (Int64) The timestamp seconds of the resource updated at.
- `workspace_uuid` - (String) The UUID of the workspace.

## API Endpoints

The following API endpoints are used by this resource:

- **Create**: `POST /api/v1/slo/add`
- **Read**: `GET /api/v1/slo/{slo_uuid}/get`
- **Update**: `POST /api/v1/slo/{slo_uuid}/modify`
- **Delete**: `GET /api/v1/slo/{slo_uuid}/delete`

## Import

SLOs can be imported using their UUID:

```bash
terraform import guance_slo.example monitor_123456
```

## Common Issues and Solutions

### Error: Could not create SLO

**Possible causes:**
- Invalid SLI UUIDs
- Goal value out of range
- MinGoal value greater than or equal to Goal
- Insufficient permissions

**Solutions:**
- Verify that the SLI UUIDs are valid and exist
- Ensure Goal is between 0 and 100
- Ensure MinGoal is between 0 and 100 and less than Goal
- Check that your API key has sufficient permissions

### Error: Could not update SLO

**Possible causes:**
- SLO UUID not found
- Invalid update parameters
- Insufficient permissions

**Solutions:**
- Verify that the SLO UUID exists
- Check that all update parameters are valid
- Ensure your API key has sufficient permissions

### Error: Could not delete SLO

**Possible causes:**
- SLO UUID not found
- Insufficient permissions

**Solutions:**
- Verify that the SLO UUID exists
- Ensure your API key has sufficient permissions
