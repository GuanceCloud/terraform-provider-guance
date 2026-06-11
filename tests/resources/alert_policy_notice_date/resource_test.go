package AlertPolicyNoticeDate_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/GuanceCloud/terraform-provider-guance/internal/provider"
)

func TestAccAlertPolicyNoticeDate(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: provider.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: provider.Config + `
resource "guance_alert_policy_notice_date" "demo" {
  name                     = "oac-alert-policy-notice-date-demo"
  skip_ref_check_on_delete = false

  notice_dates = [
    "2026/06/10",
    "2026/06/11",
  ]
}

data "guance_alert_policy_notice_date" "demo" {
  name = guance_alert_policy_notice_date.demo.name

  depends_on = [guance_alert_policy_notice_date.demo]
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("guance_alert_policy_notice_date.demo", "name", "oac-alert-policy-notice-date-demo"),
					resource.TestCheckResourceAttr("guance_alert_policy_notice_date.demo", "notice_dates.0", "2026/06/10"),
					resource.TestCheckResourceAttr("guance_alert_policy_notice_date.demo", "notice_dates.1", "2026/06/11"),
					resource.TestCheckResourceAttr("guance_alert_policy_notice_date.demo", "skip_ref_check_on_delete", "false"),
					resource.TestCheckResourceAttrSet("data.guance_alert_policy_notice_date.demo", "uuid"),
				),
			},
		},
	})
}
