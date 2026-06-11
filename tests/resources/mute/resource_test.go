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

func TestAccMuteRepeatedAndCustom(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: provider.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: provider.Config + `
resource "guance_notify_object" "demo" {
  type = "simpleHTTPRequest"
  name = "oac-mute-scenarios-notify-demo"

  opt_set = jsonencode({
    url = "https://example.com/terraform-provider-guance-mute-scenarios-test"
    headersConfig = {
      isOpen = false
      items  = []
    }
  })
}

resource "guance_alert_policy" "demo" {
  name          = "oac-mute-scenarios-alert-policy-demo"
  desc          = "acceptance mute scenarios alert policy"
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

resource "guance_mute" "weekly" {
  name        = "oac-mute-weekly-demo"
  description = "acceptance repeated alert policy mute"
  type        = "alertPolicy"
  timezone    = "Asia/Shanghai"

  mute_ranges = [{
    name              = guance_alert_policy.demo.name
    alert_policy_uuid = guance_alert_policy.demo.uuid
  }]

  repeat_time_set = 1
  repeat_crontab_set = {
    min   = "0"
    hour  = "0"
    day   = "*"
    month = "*"
    week  = "1,2,3,4,5"
  }
  crontab_duration   = 3600
  repeat_expire_time = "0"

  tags = {
    service = ["oac-mute-weekly-demo"]
  }
}

resource "guance_mute" "custom" {
  name        = "oac-mute-custom-demo"
  description = "acceptance custom mute"
  type        = "custom"
  timezone    = "Asia/Shanghai"
  enabled     = false

  mute_ranges = []

  repeat_time_set = 0
  start_time      = "2026/12/31 14:00:00"
  end_time        = "2026/12/31 15:00:00"
  filter_string   = "host:oac-mute-custom service:terraform"

  declaration = {
    source = "terraform-provider-guance"
  }
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("guance_mute.weekly", "repeat_time_set", "1"),
					resource.TestCheckResourceAttr("guance_mute.weekly", "repeat_crontab_set.hour", "0"),
					resource.TestCheckResourceAttr("guance_mute.weekly", "crontab_duration", "3600"),
					resource.TestCheckResourceAttr("guance_mute.custom", "type", "custom"),
					resource.TestCheckResourceAttr("guance_mute.custom", "enabled", "false"),
					resource.TestCheckResourceAttr("guance_mute.custom", "status", "2"),
					resource.TestCheckResourceAttr("guance_mute.custom", "filter_string", "host:oac-mute-custom service:terraform"),
				),
			},
		},
	})
}
