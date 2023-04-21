# Notification

Notification is a feature that allows you to receive alerts when an event occurs.

Guance Cloud supports setting the notification object for the alert event. When the alert event occurs, the alert
information will be sent to the notification object.

Guance Cloud supports multiple different notification methods. You can refer to the documentation to learn more.

Relationships:

```mermaid
graph LR

A[Monitor] --> B[Alert Policy] --> C[Notification]
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

resource "guance_notification" "demo" {
  name            = "oac-demo"
  type            = "ding_talk_robot"
  ding_talk_robot = {
    webhook = var.ding_talk_webhook
    secret  = var.ding_talk_secret
  }
}
```
