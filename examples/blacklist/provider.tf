
terraform {
  required_version = ">= 1.0"

  required_providers {
    guance = {
      source = "GuanceCloud/guance"
    }
  }
}

provider "guance" {
  region = "hangzhou"
}
