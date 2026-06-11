package Mute_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/GuanceCloud/terraform-provider-guance/internal/provider"
)

func TestAccMute(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: provider.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: provider.Config + `
resource "guance_notify_object" "demo" {
  type = "simpleHTTPRequest"
  name = "oac-mute-notify-demo"

  opt_set = jsonencode({
    url = "https://example.com/terraform-provider-guance-mute-test"
    headersConfig = {
      isOpen = false
      items  = []
    }
  })
}

resource "guance_alert_policy" "demo" {
  name          = "oac-mute-alert-policy-demo"
  desc          = "acceptance mute alert policy"
  rule_timezone = "Asia/Shanghai"

  alert_opt = {
    alert_type     = "status"
    silent_timeout = 300
    agg_interval   = 60

    alert_target = [{
      name = "default"

      targets = [{
        to     = [guance_notify_object.demo.uuid]
        status = "critical,error,warning"
      }]
    }]
  }
}

resource "guance_mute" "demo" {
  name        = "oac-mute-demo"
  description = "acceptance alert policy mute"
  type        = "alertPolicy"
  timezone    = "Asia/Shanghai"
  enabled     = false

  mute_ranges = [{
    name              = guance_alert_policy.demo.name
    alert_policy_uuid = guance_alert_policy.demo.uuid
  }]

  repeat_time_set = 0
  start_time      = "2026/12/31 10:00:00"
  end_time        = "2026/12/31 11:00:00"

  notify_time_str = "2026/12/31 09:50:00"
  notify_message  = "mute starts soon"

  notify_targets = [{
    type = "notifyObject"
    to   = [guance_notify_object.demo.uuid]
  }]
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("guance_mute.demo", "name", "oac-mute-demo"),
					resource.TestCheckResourceAttr("guance_mute.demo", "type", "alertPolicy"),
					resource.TestCheckResourceAttr("guance_mute.demo", "enabled", "false"),
					resource.TestCheckResourceAttr("guance_mute.demo", "status", "2"),
				),
			},

			// Delete testing automatically occurs in TestCase
		},
	})
}
