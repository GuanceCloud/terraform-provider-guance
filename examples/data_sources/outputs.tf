output "member_count" {
  description = "Number of workspace members returned by guance_members."
  value       = length(data.guance_members.all.members)
}

output "permissions" {
  description = "Workspace permissions returned by guance_permissions."
  value       = data.guance_permissions.all
}

output "notify_object_uuid" {
  description = "UUID of the looked up notify object."
  value       = try(data.guance_notify_object.example[0].uuid, null)
}

output "alert_policy_notice_date_uuid" {
  description = "UUID of the looked up alert policy notice date."
  value       = try(data.guance_alert_policy_notice_date.example[0].uuid, null)
}

output "alert_policy_uuid" {
  description = "UUID of the looked up alert policy."
  value       = try(data.guance_alert_policy.example[0].uuid, null)
}

output "mute_uuid" {
  description = "UUID of the looked up mute rule."
  value       = try(data.guance_mute.example[0].uuid, null)
}

output "monitor_uuid" {
  description = "UUID of the looked up monitor/checker."
  value       = try(data.guance_monitor.example[0].uuid, null)
}

output "monitor_count" {
  description = "Number of monitors returned by guance_monitors."
  value       = length(data.guance_monitors.examples.monitors)
}
