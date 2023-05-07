package Monitor_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/GuanceCloud/terraform-provider-guance/internal/provider"
)

func TestAccMonitor(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: provider.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: provider.Config + `
variable "email" {
  type = string
}

data "guance_members" "demo" {
  filters = [
    {
      name   = "email"
      values = [var.email]
    }
  ]
}

resource "guance_membergroup" "demo" {
  name       = "oac-demo"
  member_ids = data.guance_members.demo.items[*].id
}

resource "guance_alertpolicy" "demo" {
  name           = "oac-demo"
  silent_timeout = "1h"

  statuses = [
    "critical",
    "error",
    "warning",
    "info",
    "ok",
    "nodata",
    "nodata_ok",
    "nodata_as_ok",
  ]

  alert_targets = [
    {
      type         = "member_group"
      member_group = {
        id = guance_membergroup.demo.id
      }
    },
  ]
}

resource "guance_monitor" "demo" {
  manifest     = file("${path.module}/monitor.json")
  alert_policy = {
    id = guance_alertpolicy.demo.id
  }
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(),
			},

			// Delete testing automatically occurs in TestCase
		},
	})
}
