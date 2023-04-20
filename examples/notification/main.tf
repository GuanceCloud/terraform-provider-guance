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
