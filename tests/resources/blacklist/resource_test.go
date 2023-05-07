package Blacklist_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/GuanceCloud/terraform-provider-guance/internal/provider"
)

func TestAccBlacklist(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: provider.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: provider.Config + `
resource "guance_blacklist" "demo" {
  source = {
    type = "logging"
    name = "nginx"
  }

  filter_rules = [
    {
      name      = "foo"
      operation = "in"
      condition = "and"
      values    = ["oac-*"]
    }
  ]
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(),
			},

			// Delete testing automatically occurs in TestCase
		},
	})
}
