# Black List

Guance Cloud supports filtering data that meets the conditions by setting a blacklist.

After configuring the blacklist, the data that meets the conditions will no longer be reported to the Guance Cloud
workspace, helping you save data storage costs.

## Example Usage

```terraform
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
```