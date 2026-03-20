# Monitor Resource

The `guance_monitor` resource manages monitors in Guance Cloud. Monitors are used to detect anomalies and trigger alerts based on defined rules and conditions.

## Example Usage

```terraform
resource "guance_monitor" "example" {
  type = "trigger"
  status = 0
  alert_policy_uuids = ["altpl_xxxx32"]
  dashboard_uuid = "dsbd_xxxx32"
  tags = ["example", "terraform"]
  secret = "secret_xxxxx"
  open_permission_set = false
  permission_set = ["wsAdmin", "acnt_xxxx", "group_yyyy"]
  
  extend = jsonencode({
    isNeedCreateIssue = false
    issueLevelUUID = ""
    needRecoverIssue = false
  })
  
  json_script = jsonencode({
    type = "simpleCheck"
    title = "SSH Service Exception"
    message = ">Content：Host SSH Status Failed  \n>Suggestion：Check Host SSH Service Status"
    every = "1m"
    interval = 300
    recoverNeedPeriodCount = 2
    disableCheckEndTime = false
    groupBy = ["host"]
    targets = [
      {
        dql = "M::`ssh`:(count(`ssh_check`)) BY `host`"
        alias = "Result"
        qtype = "dql"
      }
    ]
    checkerOpt = {
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
        },
        {
          conditionLogic = "and"
          conditions = [
            {
              alias = "Result"
              operands = ["0"]
              operator = ">="
            }
          ]
          status = "error"
        }
      ]
      infoEvent = false
    }
    channels = []
    atAccounts = []
    atNoDataAccounts = []
  })
}
```
