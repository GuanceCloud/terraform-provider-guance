variable "region" {
  description = "The region of Guance Cloud."
  type        = string
  default     = "hangzhou"
}
variable "alert_policy_uuids" {
  description = "List of alert policy UUIDs."
  type        = list(string)
  default     = []
}

variable "tags" {
  description = "List of tag names for filtering."
  type        = list(string)
  default     = ["example", "terraform"]
}

variable "secret" {
  description = "Unique identifier secret for the middle section of the Webhook address."
  type        = string
  default     = "secret_xxxxx"
}

variable "open_permission_set" {
  description = "Whether to enable custom permission configuration."
  type        = bool
  default     = false
}

variable "permission_set" {
  description = "Operation permission configuration."
  type        = list(string)
  default     = ["wsAdmin"]
}

variable "sli_uuids" {
  description = "SLI UUIDs."
  type        = list(string)
  default     = ["rul-aaaaaa", "rul-bbbbbb"]
}
