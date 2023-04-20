package Intelligentinspection_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/GuanceCloud/terraform-provider-guance/internal/provider"
)

func TestAccIntelligentinspection(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: provider.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: provider.Config + `
resource "guance_intelligentinspection" "demo" {
	name        = "oac-demo"
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(),
			},

			// Create and Read testing
			{
				Config: provider.Config + `
resource "guance_intelligentinspection" "demo" {
  name = "oac-demo-complete"
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(),
			},

			// Delete testing automatically occurs in TestCase
		},
	})
}
