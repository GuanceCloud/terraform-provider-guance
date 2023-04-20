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
variable "ding_talk_webhook" {
  type = string
}

variable "ding_talk_secret" {
  type = string
}

variable "email" {
  type = string
}

data "guance_members" "demo" {
  emails = [
    "liyufei906@guance.com"
  ]
}

resource "guance_membergroup" "demo" {
  name       = "oac-demo"
  member_ids = data.guance_members.demo.items[*].id
}

resource "guance_mute" "demo" {
  name = "oac-demo"

  // mute ranges
  mute_ranges {
    type = "monitor"

    monitor {
      id = ""
    }
  }

  mute_ranges {
    type = "alert_policy"

    alert_policy {
      id = ""
    }
  }

  // notify options
  notify {
    message = <<EOF
      Muted
    EOF

    before_time = "15m"
  }

  notify_targets {
    type = "member_group"

    member_group {
      id = guance_membergroup.demo.id
    }
  }

  notify_targets {
    type = "notification"

    notification {
      id = ""
    }
  }

  // ont-time options
  onetime {
    start = "2022-08-04T12:00:00Z"
    end   = "2023-12-31T12:00:00Z"
  }

  // cron options
  repeat {
    crontab_duration = 30 // 30s
    start            = "05:00:00Z"
    end              = "10:00:00Z"
    expire           = "2023-12-31T12:00:00Z"
    crontab {
      min   = "0"
      hour  = "0"
      day   = "*"
      month = "*"
      week  = "*"
    }
  }
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(),
			},

			// Delete testing automatically occurs in TestCase
		},
	})
}
