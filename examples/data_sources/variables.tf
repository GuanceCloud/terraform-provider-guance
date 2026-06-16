variable "notify_object_name" {
  description = "Existing notify object name to look up."
  type        = string
  default     = ""
}

variable "alert_policy_notice_date_name" {
  description = "Existing alert policy notice date name to look up."
  type        = string
  default     = ""
}

variable "alert_policy_name" {
  description = "Existing alert policy name to look up."
  type        = string
  default     = ""
}

variable "mute_name" {
  description = "Existing mute rule name to look up."
  type        = string
  default     = ""
}

variable "monitor_name" {
  description = "Existing monitor/checker name to look up."
  type        = string
  default     = ""
}

variable "monitor_type" {
  description = "Optional checker type filter for guance_monitors, such as simpleCheck. Leave empty to list all monitor/checker types."
  type        = string
  default     = ""
}

variable "monitor_search" {
  description = "Optional monitor search keyword for guance_monitors."
  type        = string
  default     = ""
}

variable "monitor_status" {
  description = "Optional monitor status filter for guance_monitors, such as 0 or 2."
  type        = string
  default     = ""
}
