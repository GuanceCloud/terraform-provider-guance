resource "guance_blacklist" "demo" {
  source = {
    type = "logging"
    name = "nginx"
  }

  filter_rules = [
    {
      name      = "foo"
      operation = "in"
      condition = "and"
      values    = ["oac-*"]
    }
  ]
}
