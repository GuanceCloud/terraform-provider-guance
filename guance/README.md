# Terraform Provider: Guance

The Guance Provider provides resources to manage [Guance Cloud](https://en.guance.com/) resources.

To learn the basics of Terraform using this provider, follow the hands-on get started tutorials.

Interested in the provider's latest features, or want to make sure you're up to date? Check out the changelog for version information and release notes.

## Authenticating to Guance Cloud

Terraform supports a number of different methods for authenticating to Guance Cloud:

* [Workspace Key](https://console.guance.com/workspace/apiManage)

## Usage

```terraform
# We strongly recommend using the required_providers block to set the
# Guance Cloud Provider source and version being used
terraform {
  required_version = ">=0.12"

  required_providers {
    guance = {
      source = "GuanceCloud/guance"
      version = "=0.0.6"
    }
  }
}

// We also recommend use secret environment variables to set the provider,
// Such as GUANCE_ACCESS_TOKEN and GUANCE_REGION
provider "guance" {
  # access_token = "your access token, recommend store in environment variable"
  region = "hangzhou"
  # end_point = "https://openapi.guance.com"
}
```

## More Examples

* [Example Source Code](https://github.com/GuanceCloud/terraform-provider-guance/tree/main/examples)
