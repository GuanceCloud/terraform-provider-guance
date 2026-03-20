resource "guance_monitor" "example" {
  type = "trigger"
  status = 0
  alert_policy_uuids = []
  dashboard_uuid = []
  tags = ["example", "terraform"]
  secret = ["TF-MONITOR-SECRET"]
  open_permission_set = []
  permission_set = []
  
  extend = jsonencode({
    isNeedCreateIssue = false
    issueLevelUUID = ""
    needRecoverIssue = false
  })
  
  json_script = {
    type = "simpleCheck"
    title = "Terraform Monitor Example"
    message = ">Level：{{status}}  \n>Host：{{host}}  \n>Content：Host SSH Status {{ Result |  to_fixed(2) }}%  \n>Suggestion：Check Host SSH Service Status"
    every = "1m"
    interval = 300
    recover_need_period_count = 2
    disable_check_end_time = false
    group_by = ["host"]
    targets = [
      {
        dql = "M::`ssh`:(count(`ssh_check`)) BY `host`"
        alias = "Result"
        qtype = "dql"
      }
    ]
    checker_opt = {
      rules = [
        {
          condition_logic = "and"
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
          condition_logic = "and"
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
      info_event = false
    }
    channels = []
    at_accounts = []
    at_no_data_accounts = []
  }
}
