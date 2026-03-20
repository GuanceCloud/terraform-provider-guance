# SLO Example

This example demonstrates how to create an SLO (Service Level Objective) resource using the Guance Cloud Terraform provider.

## Prerequisites

1. [Terraform](https://www.terraform.io/downloads.html) 1.0+ installed
2. Guance Cloud account with API access
3. API key configured in your environment

## Configuration

1. Update the `main.tf` file with your specific values:
   - `name`: The name of your SLO
   - `interval`: Detection frequency (`5m` or `10m`)
   - `goal`: SLO expected goal (0-100)
   - `min_goal`: SLO minimum goal (0-100, must be less than goal)
   - `sli_uuids`: List of SLI UUIDs
   - `describe`: Optional description
   - `alert_policy_uuids`: Optional list of alert policy UUIDs
   - `tags`: Optional tags

2. Update the `variables.tf` file with your preferred default values

## Usage

1. Initialize the Terraform working directory:

```bash
terraform init
```

2. Review the execution plan:

```bash
terraform plan
```

3. Apply the configuration:

```bash
terraform apply
```

4. Destroy the resources when no longer needed:

```bash
terraform destroy
```

## Example Output

After applying the configuration, you will see output similar to:

```
Apply complete! Resources: 1 added, 0 changed, 0 destroyed.

Outputs:

SLO_name = "Example SLO"
SLO_uuid = "monitor_123456"
SLO_workspace_uuid = "wksp_123456"
```
