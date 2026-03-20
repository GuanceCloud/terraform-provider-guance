# Monitor Json Example

This example demonstrates how to use `guance_monitor_json` resource to import and manage monitors in Guance Cloud using JSON configuration.

## Requirements

- Terraform 1.0+
- Guance Cloud API key

## Usage

1. Configure your Guance Cloud API key in the `provider.tf` file.
2. Modify the `main.tf` or `main_from_file.tf` file to customize your monitor configuration.
3. Run `terraform init` to initialize the provider.
4. Run `terraform plan` to preview the changes.
5. Run `terraform apply` to import the monitors.
6. Run `terraform destroy` to delete the monitors.

## Example Configuration

### Using jsonencode

```hcl
resource "guance_monitor_json" "example" {
  checker_json = jsonencode({
    extend = {
      funcName = ""
      isNeedCreateIssue = false
      issueLevelUUID = ""
      needRecoverIssue = false
      querylist = [
        {
          datasource = "dataflux"
          qtype = "dql"
          query = {
            alias = ""
            code = "Result"
            dataSource = "ssh"
            field = "ssh_check"
            fieldFunc = "count"
            fieldType = "float"
            funcList = []
            groupBy = ["host"]
            groupByTime = ""
            namespace = "metric"
            q = "M::`ssh`:(count(`ssh_check`)) BY `host`"
            type = "simple"
          }
          uuid = "aada629a-672e-46f9-9503-8fd61065c382"
        }
      ]
      rules = [
        {
          conditionLogic = "and"
          conditions = [
            {
              alias = "Result"
              operands = ["90"]
              operator = ">="
            }
          ]
          status = "critical"
        }
      ]
    }
    is_disable = false
    jsonScript = {
      atAccounts = []
      atNoDataAccounts = []
      channels = []
      checkerOpt = {
        infoEvent = false
        rules = [
          {
            conditionLogic = "and"
            conditions = [
              {
                alias = "Result"
                operands = ["90"]
                operator = ">="
              }
            ]
            status = "critical"
          }
        ]
      }
      disableCheckEndTime = false
      every = "1m"
      groupBy = ["host"]
      interval = 300
      message = "message1"
      noDataMessage = ""
      noDataTitle = ""
      recoverNeedPeriodCount = 2
      targets = [
        {
          alias = "Result"
          dql = "M::`ssh`:(count(`ssh_check`)) BY `host`"
          qtype = "dql"
        }
      ]
      title = "Host {{ host }} SSH error"
      type = "simpleCheck"
    }
    monitorName = "default"
    secret = ""
    tagInfo = []
    type = "trigger"
  })
  
  type = "trigger"
}
```

### Using file

```hcl
resource "guance_monitor_json" "example_from_file" {
  checker_json = file("monitor.json")
  
  type = "trigger"
}
```

## Notes

- The `checker_json` field accepts a single JSON object of monitor configuration (not an array).
- When updating the resource, the Replace API is used to update the existing monitor.
- The import operation uses the Guance Cloud checker import API endpoints.
- You can export monitors from Guance Cloud and use the exported JSON as the `checker_json` input.
- The `type` field supports two values: `trigger` (normal monitor) and `smartMonitor` (smart monitor).

