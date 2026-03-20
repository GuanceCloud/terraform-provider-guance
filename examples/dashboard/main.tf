resource "guance_dashboard" "example" {
  name     = "TF-dashboard"
  desc     = "An example dashboard created with Terraform"
  is_public = 1
  identifier = "identifier-example"
  permission_set = []
  read_permission_set = ["*"]
  tag_names = [
    "example",
    "terraform",
  ]
  
  template_info = file("dashboard.json")
}

output "dashboard_uuid" {
  value       = guance_dashboard.example.uuid
  description = "The UUID of the created dashboard"
}

output "dashboard_name" {
  value       = guance_dashboard.example.name
  description = "The name of the created dashboard"
}