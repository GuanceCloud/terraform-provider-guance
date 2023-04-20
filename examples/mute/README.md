# Mute Rule

Mute rule is a feature that allows you to temporarily stop receiving notifications for a specific alert. You can use
mute rules to temporarily silence alerts that are not relevant to you, or to silence alerts that you are already aware
of.

Guance Cloud supports the management of all mute rules in the current workspace. It supports muting different monitors,
smart inspections, self-built inspections, SLOs, and alert policies, so that the muted objects do not send any alert
notifications to any alert notification objects during the mute time.

Relationships:

```mermaid
graph LR
    A[Mute Rule] --> B[Alert Policy]
```

## Create

The first let me create a resource. We will send the create operation to the resource management service

```terraform
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
```
