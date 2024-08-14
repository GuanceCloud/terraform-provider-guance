resource "guance_blacklist" "demo" {
  type   = "logging"
  source = "mysql"

  filters = [
    {
      name      = "foo1"
      operation = "in"
      condition = "and"
      values    = ["oac-*"]
    }
  ]
}
