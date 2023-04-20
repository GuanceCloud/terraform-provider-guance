package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"

	"github.com/GuanceCloud/terraform-provider-guance/guance"
)

const (
	// Config is a shared configuration to combine with the actual
	// test configuration so the Guance Cloud client is properly configured.
	// It is also possible to use the GUANCE_ environment variables instead,
	// such as updating the Makefile and running the testing through that tool.
	Config = `
terraform {
	required_version = ">=0.12"

	required_providers {
		guance = {
			source = "GuanceCloud/guance"
		}
	}
}

provider "guance" {
	region = "hangzhou"
	token = ""
}
`
)

var (
	// TestAccProtoV6ProviderFactories are used to instantiate a provider during
	// acceptance testing. The factory function will be invoked for every Terraform
	// CLI command executed to create a provider server to which the CLI can
	// reattach.
	TestAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
		"guance": providerserver.NewProtocol6WithError(guance.New()),
	}
)
