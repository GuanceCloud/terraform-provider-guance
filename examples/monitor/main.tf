resource "guance_monitor" "example" {
  type                = "trigger"
  status              = 0
  alert_policy_uuids  = var.alert_policy_uuids
  tags                = var.tags
  secret              = var.secret
  open_permission_set = var.open_permission_set
  permission_set      = var.permission_set

  extend = jsonencode({
    isNeedCreateIssue = false
    issueLevelUUID    = ""
    needRecoverIssue  = false
  })

  json_script = {
    type                      = "simpleCheck"
    title                     = "Terraform Monitor Example"
    message                   = ">Level: {{status}}\n>Host: {{host}}\n>Content: Host SSH Status {{ Result | to_fixed(2) }}%\n>Suggestion: Check Host SSH Service Status"
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
