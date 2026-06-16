terraform {
  required_version = ">= 1.0"

  required_providers {
    guance = {
      source = "GuanceCloud/guance"
    }
  }
}

provider "guance" {
  # You can set your API key here or use the GUANCE_ACCESS_TOKEN environment variable.
  # access_token = "your-api-key"
}
