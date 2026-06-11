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

func TestAccAlertpolicyComplexStatus(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: provider.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: provider.Config + `
resource "guance_notify_object" "demo" {
  type = "simpleHTTPRequest"
  name = "oac-alert-policy-complex-demo"

  opt_set = jsonencode({
    url = "https://example.com/terraform-provider-guance-alert-policy-complex-test"
    headersConfig = {
      isOpen = false
      items  = []
    }
  })
}

resource "guance_alert_policy_notice_date" "demo" {
  name = "oac-alert-policy-complex-date-demo"

  notice_dates = [
    "2026/06/10",
    "2026/06/11",
  ]
}

resource "guance_alert_policy" "demo" {
  name          = "oac-alert-policy-complex-demo"
  desc          = "acceptance complex status alert policy"
  rule_timezone = "Asia/Shanghai"

  alert_opt = {
    agg_type                        = "byFields"
    ignore_ok                       = true
    alert_type                      = "status"
    silent_timeout                  = 300
    silent_timeout_by_status_enable = true
    silent_timeout_by_status = [{
      status         = "critical"
      silent_timeout = 120
    }]
    agg_interval       = 60
    agg_fields         = ["df_monitor_checker_id", "df_label"]
    agg_labels         = ["service"]
    agg_cluster_fields = ["df_title"]
    agg_send_first     = true

    alert_target = [{
      name              = "complex route"
      custom_date_uuids = [guance_alert_policy_notice_date.demo.uuid]
      custom_start_time = "09:30:00"
      custom_duration   = 3600

      targets = [
        {
          to            = [guance_notify_object.demo.uuid]
          status        = "critical,error"
          filter_string = "host:oac-alert-policy-complex service:terraform"
          tags = {
            service = ["oac-alert-policy-complex"]
          }
          upgrade_targets = [{
            to       = [guance_notify_object.demo.uuid]
            duration = 300
            to_way   = ["mail"]
          }]
        },
        {
          to        = [guance_notify_object.demo.uuid]
          status    = "critical"
          df_source = "security"
        },
      ]
    }]
  }
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("guance_alert_policy.demo", "alert_opt.agg_type", "byFields"),
					resource.TestCheckResourceAttr("guance_alert_policy.demo", "alert_opt.ignore_ok", "true"),
					resource.TestCheckResourceAttr("guance_alert_policy.demo", "alert_opt.silent_timeout_by_status_enable", "true"),
					resource.TestCheckResourceAttr("guance_alert_policy.demo", "alert_opt.alert_target.0.targets.0.upgrade_targets.0.duration", "300"),
					resource.TestCheckResourceAttr("guance_alert_policy.demo", "alert_opt.alert_target.0.targets.1.df_source", "security"),
				),
			},
		},
	})
}

func TestAccAlertpolicyMemberMode(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: provider.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: provider.Config + `
data "guance_members" "demo" {}

resource "guance_notify_object" "demo" {
  type = "simpleHTTPRequest"
  name = "oac-alert-policy-member-demo"

  opt_set = jsonencode({
    url = "https://example.com/terraform-provider-guance-alert-policy-member-test"
    headersConfig = {
      isOpen = false
      items  = []
    }
  })
}

resource "guance_alert_policy" "demo" {
  name          = "oac-alert-policy-member-demo"
  desc          = "acceptance member alert policy"
  rule_timezone = "Asia/Shanghai"

  alert_opt = {
    alert_type     = "member"
    silent_timeout = 300
    agg_interval   = 60

    alert_target = [{
      name = "member notification target"

      alert_info = [{
        name        = "member route"
        member_info = [data.guance_members.demo.members[0].uuid]

        targets = [{
          to     = [guance_notify_object.demo.uuid]
          status = "critical,error,warning"
        }]
      }]
    }]
  }
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("guance_alert_policy.demo", "alert_opt.alert_type", "member"),
					resource.TestCheckResourceAttr("guance_alert_policy.demo", "alert_opt.alert_target.0.alert_info.0.name", "member route"),
					resource.TestCheckResourceAttrSet("guance_alert_policy.demo", "alert_opt.alert_target.0.alert_info.0.member_info.0"),
				),
			},
		},
	})
}
