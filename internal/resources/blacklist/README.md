# Black List

Guance Cloud supports filtering data that meets the conditions by setting a blacklist.

After configuring the blacklist, the data that meets the conditions will no longer be reported to the Guance Cloud
workspace, helping you save data storage costs.

## Create

The first let me create a resource. We will send the create operation to the resource management service

```terraform
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
```
