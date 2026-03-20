# Dashboard Example

This example demonstrates how to use the `guance_dashboard` resource to create and manage dashboards in Guance Cloud.

## Requirements

- Terraform 1.0+
- Guance Cloud API key

## Usage

1. Configure your Guance Cloud API key in the `provider.tf` file.
2. Modify the `main.tf` file to customize your dashboard configuration.
3. Run `terraform init` to initialize the provider.
4. Run `terraform plan` to preview the changes.
5. Run `terraform apply` to create the dashboard.
6. Run `terraform destroy` to delete the dashboard.

## Example Configuration

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

## Notes

- The `template_info` field expects a JSON string that defines the dashboard structure.
- The `extend` and `mapping` fields also expect JSON strings.
- The `specify_dashboard_uuid` field must follow the pattern `dsbd_custom_` followed by 32 lowercase alphanumeric characters.
