package Alertpolicy_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/GuanceCloud/terraform-provider-guance/internal/provider"
)

func TestAccAlertpolicy(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: provider.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: provider.Config + `
resource "guance_notify_object" "demo" {
  type = "simpleHTTPRequest"
  name = "oac-alert-policy-demo"

  opt_set = jsonencode({
    url = "https://example.com/terraform-provider-guance-alert-policy-test"
    headersConfig = {
      isOpen = false
      items  = []
    }
  })
}

resource "guance_alert_policy_notice_date" "demo" {
  name = "oac-alert-policy-date-demo"

  notice_dates = [
    "2026/06/10",
    "2026/06/11",
  ]
}

resource "guance_alert_policy" "demo" {
  name          = "oac-alert-policy-demo"
  desc          = "acceptance alert policy"
  rule_timezone = "Asia/Shanghai"

  alert_opt = {
    alert_type     = "status"
    silent_timeout = 300
    agg_interval   = 60
    agg_fields     = ["df_monitor_checker_id"]

    alert_target = [{
      name              = "default"
      custom_date_uuids = [guance_alert_policy_notice_date.demo.uuid]
      custom_start_time = "09:30:00"
      custom_duration   = 3600

      targets = [{
        to     = [guance_notify_object.demo.uuid]
        status = "critical,error,warning"
      }]
    }]
  }
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("guance_alert_policy.demo", "name", "oac-alert-policy-demo"),
					resource.TestCheckResourceAttr("guance_alert_policy.demo", "alert_opt.alert_type", "status"),
					resource.TestCheckResourceAttr("guance_alert_policy.demo", "alert_opt.agg_interval", "60"),
				),
			},

			// Delete testing automatically occurs in TestCase
		},
	})
}
