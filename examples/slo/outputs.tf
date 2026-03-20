output "slo_uuid" {
  description = "The UUID of the created SLO."
  value       = guance_slo.example.uuid
}

output "slo_name" {
  description = "The name of the created SLO."
  value       = guance_slo.example.name
}

output "slo_workspace_uuid" {
  description = "The workspace UUID of the created SLO."
  value       = guance_slo.example.workspace_uuid
}
