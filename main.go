package main

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"

	"github.com/GuanceCloud/terraform-provider-guance/guance"
)

func main() {
	err := providerserver.Serve(context.Background(), guance.New, providerserver.ServeOpts{
		// NOTE: This is not a typical Terraform Registry provider address,
		// such as registry.terraform.io/hashicorp/hashicups. This specific
		// provider address is used in these tutorials in conjunction with a
		// specific Terraform CLI configuration for manual development testing
		// of this provider.
		Address: "registry.terraform.io/GuanceCloud/guance",
	})
	panic(err)
}
