# Monitor Json Resource

The `guance_monitor_json` resource manages monitors in Guance Cloud using JSON import/export functionality. This resource allows you to import and export monitor configurations using the checker JSON format.

## Example Usage

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
      title = "Host ${host} SSH error"
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

## Argument Reference

- `checker_json` - (Optional) The checker JSON configuration for import. This should be a single JSON object of monitor configuration.
- `type` - (Optional) Monitor type. Valid values are `trigger` (normal monitor) or `smartMonitor` (smart monitor). Default is `trigger`.

## Attribute Reference

- `uuid` - The UUID of the resource.
- `checker_json` - The checker JSON configuration.
- `checker_json_export` - The exported checker JSON configuration containing the created monitors.
- `create_at` - The timestamp when the resource was created.
- `update_at` - The timestamp when the resource was last updated.
- `workspace_uuid` - The UUID of the workspace.

## Import

You can import a monitor json resource using its UUID:

```sh
terraform import guance_monitor_json.example <uuid>
```

## Notes

- The `checker_json` field accepts a single JSON object of monitor configuration (not an array).
- The `checker_json_export` field contains the actual monitor configurations after import, including the generated UUIDs.
- When updating the resource, the Replace API is used to update the existing monitor.
- The import operation uses the Guance Cloud checker import API endpoints.
- You can export monitors from Guance Cloud and use the exported JSON as the `checker_json` input.
