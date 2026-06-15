variable "name_prefix" {
  description = "Prefix used for all alert example resource names."
  type        = string
  default     = "terraform-alert-example"
}

variable "timezone" {
  description = "Timezone used by the alert policy and mute rule."
  type        = string
  default     = "Asia/Shanghai"
}

variable "webhook_url" {
  description = "Webhook URL used by the simple HTTP request notify object."
  type        = string
  default     = "https://example.com/guance-alert-example"
}

variable "notice_dates" {
  description = "Custom notice dates in YYYY/MM/DD format."
  type        = list(string)
  default = [
    "2026/12/31",
    "2027/01/01",
  ]
}

variable "alert_filter" {
  description = "DQL-style alert target filter expression."
  type        = string
  default     = null
}

variable "mute_start_time" {
  description = "One-time mute start time in YYYY/MM/DD HH:mm:ss format."
  type        = string
  default     = "2026/12/31 10:00:00"
}

variable "mute_end_time" {
  description = "One-time mute end time in YYYY/MM/DD HH:mm:ss format."
  type        = string
  default     = "2026/12/31 11:00:00"
}

variable "mute_notify_time" {
  description = "Mute pre-notification time in YYYY/MM/DD HH:mm:ss format."
  type        = string
  default     = "2026/12/31 09:50:00"
}
