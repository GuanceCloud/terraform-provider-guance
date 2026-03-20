resource "guance_monitor_json" "example" {
  checker_json = file("monitor.json")
  type = "trigger"
}

output "monitor_json_uuid" {
  value       = guance_monitor_json.example.uuid
  description = "The UUID of the created monitor json resource"
}

output "monitor_json_workspace_uuid" {
  value       = guance_monitor_json.example.workspace_uuid
  description = "The UUID of the workspace"
}
