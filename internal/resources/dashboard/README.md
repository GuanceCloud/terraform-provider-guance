# Dashboard Resource

The `guance_dashboard` resource manages dashboards in Guance Cloud. Dashboards are visualizations that display metrics, logs, traces, and other observability data from your systems.

## Example Usage

```hcl
resource "guance_dashboard" "example" {
  name     = "example-dashboard"
  desc     = "An example dashboard created with Terraform"
  is_public = 0
  
  tag_names = [
    "example",
    "terraform"
  ]
  
  template_info = jsonencode({
    title = "Example Dashboard"
    main = {
      charts = [
        {
          name = "Request Count"
          type = "sequence"
          queries = [
            {
              checked = true
              datasource = "dataflux"
              qtype = "dql"
              query = {
                q = "T::re(`.*`):(count(`trace_id`)) [::auto]"
              }
            }
          ]
        }
      ]
    }
  })
}
```

## Argument Reference

- `name` - (Required) The name of the dashboard.
- `desc` - (Optional) The description of the dashboard.
- `identifier` - (Optional) The identifier of the dashboard.
- `tag_names` - (Optional) List of associated tag names.
- `template_info` - (Optional) Dashboard template data in JSON format.
- `specify_dashboard_uuid` - (Optional) Specified dashboard UUID.
- `is_public` - (Optional) Whether the dashboard is public.
- `permission_set` - (Optional) Custom operation permissions.
- `read_permission_set` - (Optional) Custom read permissions.

## Attribute Reference

- `uuid` - The UUID of the dashboard.
- `create_at` - The timestamp when the dashboard was created.
- `update_at` - The timestamp when the dashboard was last updated.
- `workspace_uuid` - The UUID of the workspace.
- `template_info` - The dashboard template data in JSON format.
- `template_info_export` - The exported dashboard template data in JSON format.
