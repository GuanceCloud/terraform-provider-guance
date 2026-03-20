variable "region" {
  description = "The region of Guance Cloud."
  type        = string
  default     = "hangzhou"
}

variable "sli_uuids" {
  description = "SLI UUIDs."
  type        = list(string)
  default     = ["rul-aaaaaa", "rul-bbbbbb"]
}

variable "alert_policy_uuids" {
  description = "Alert policy UUIDs."
  type        = list(string)
  default     = ["altpl-xxxxxx"]
}

variable "tags" {
  description = "Tags for the SLO."
  type        = list(string)
  default     = ["example", "terraform"]
}
