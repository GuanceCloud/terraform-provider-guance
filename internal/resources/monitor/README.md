# Monitor Resource

The `guance_monitor` resource manages monitor checker rules in Guance Cloud.

## Example Usage

```terraform
resource "guance_monitor" "example" {
  type                = "trigger"
  status              = 0
  tags                = ["example", "terraform"]
  secret              = "secret_xxxxx"
  open_permission_set = false
  permission_set      = []
  alert_policy_uuids  = []

  extend = jsonencode({
    isNeedCreateIssue = false
    issueLevelUUID    = ""
    needRecoverIssue  = false
  })

  json_script = {
    type                      = "simpleCheck"
    title                     = "SSH Service Exception"
    message                   = ">Content: Host SSH Status Failed\n>Suggestion: Check Host SSH Service Status"
    every                     = "1m"
    interval                  = 300
    recover_need_period_count = 2
    disable_check_end_time    = false
    group_by                  = ["host"]

    targets = [{
      dql   = "M::`ssh`:(count(`ssh_check`)) BY `host`"
      alias = "Result"
      qtype = "dql"
    }]

    checker_opt = {
      info_event = false

      rules = [
        {
          condition_logic = "and"
          status          = "critical"

          conditions = [{
            alias    = "Result"
            operator = ">="
            operands = ["90"]
          }]
        },
        {
          condition_logic = "and"
          status          = "error"

          conditions = [{
            alias    = "Result"
            operator = ">="
            operands = ["0"]
          }]
        },
      ]
    }

    channels            = []
    at_accounts         = []
    at_no_data_accounts = []
  }
}
```
