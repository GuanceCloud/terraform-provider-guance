
terraform {
	required_version = ">=0.12"

	required_providers {
		guance = {
			source = "GuanceCloud/guance"
		}
	}
}

provider "guance" {
}
