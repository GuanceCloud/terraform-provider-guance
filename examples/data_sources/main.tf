data "guance_members" "all" {}

data "guance_permissions" "all" {}

data "guance_notify_object" "example" {
  count = var.notify_object_name == "" ? 0 : 1
  name  = var.notify_object_name
}

data "guance_alert_policy_notice_date" "example" {
  count = var.alert_policy_notice_date_name == "" ? 0 : 1
  name  = var.alert_policy_notice_date_name
}

data "guance_alert_policy" "example" {
  count = var.alert_policy_name == "" ? 0 : 1
  name  = var.alert_policy_name
}

data "guance_mute" "example" {
  count = var.mute_name == "" ? 0 : 1
  name  = var.mute_name
}

data "guance_monitor" "example" {
  count = var.monitor_name == "" ? 0 : 1
  name  = var.monitor_name
  type  = var.monitor_type
}

data "guance_monitors" "examples" {
  search = var.monitor_search
  type   = var.monitor_type
  status = var.monitor_status
}
