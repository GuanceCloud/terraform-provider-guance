resource "guance_blacklist" "demo" {
  name = "blacklist-demo"
  type   = "logging"
  sources = ["mysql", "oracle"]
  desc = "this is a demo"

  filters = [
    {
      name      = "foo1"
      operation = "in"
      condition = "and"
      values    = ["oac-*"]
    }
  ]
}
