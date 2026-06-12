output "notify_object_uuid" {
  description = "The UUID of the created notify object."
  value       = guance_notify_object.example.uuid
}

output "notice_date_uuid" {
  description = "The UUID of the created alert policy notice date."
  value       = guance_alert_policy_notice_date.example.uuid
}

output "alert_policy_uuid" {
  description = "The UUID of the created alert policy."
  value       = guance_alert_policy.example.uuid
}

output "mute_uuid" {
  description = "The UUID of the created mute rule."
  value       = guance_mute.maintenance.uuid
}

output "data_source_alert_type" {
  description = "Alert type read back from the alert policy data source."
  value       = data.guance_alert_policy.example.alert_opt.alert_type
}

output "data_source_mute_alert_policy_uuid" {
  description = "Alert policy UUID read back from the mute data source."
  value       = data.guance_mute.maintenance.mute_ranges[0].alert_policy_uuid
}

