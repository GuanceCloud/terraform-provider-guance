resource "guance_slo" "example" {
  name             = "Example SLO"
  interval         = "5m"
  goal             = 99.9
  min_goal         = 99.0
  sli_uuids        = var.sli_uuids
  describe         = "This is an example SLO"
  alert_policy_uuids = var.alert_policy_uuids
  tags             = var.tags
}
