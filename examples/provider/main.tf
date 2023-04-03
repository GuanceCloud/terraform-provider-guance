terraform {
  required_version = ">=0.12"

  required_providers {
    guance = {
      source = "GuanceCloud/guance"
    }
  }
}

provider "guance" {
  endpoint = "http://127.0.0.1:8080"
}

