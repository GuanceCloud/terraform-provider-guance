resource "guance_dashboard" "demo" {
  name     = "oac-demo"
  manifest = file("${path.module}/dashboard.json")
}
