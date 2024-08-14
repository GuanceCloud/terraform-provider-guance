package blacklist

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// blackListResourceModel maps the resource schema data.
type blackListResourceModel struct {
	UUID      types.String `tfsdk:"uuid"`
	CreatedAt types.String `tfsdk:"created_at"`
	Source    types.String `tfsdk:"source"`
	Type      types.String `tfsdk:"type"`
	Filters   []*filter    `tfsdk:"filters"`
}

// filter maps the resource schema data.
type filter struct {
	Name      types.String   `tfsdk:"name"`
	Operation types.String   `tfsdk:"operation"`
	Condition types.String   `tfsdk:"condition"`
	Values    []types.String `tfsdk:"values"`
}
