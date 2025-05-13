package blacklist

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// blackListResourceModel maps the resource schema data.
type blackListResourceModel struct {
	UUID          types.String   `tfsdk:"uuid"`
	Name          types.String   `tfsdk:"name"`
	Desc          types.String   `tfsdk:"desc"`
	CreateAt      types.String   `tfsdk:"create_at"`
	UpdateAt      types.String   `tfsdk:"update_at"`
	Source        types.String   `tfsdk:"source"`
	Sources       []types.String `tfsdk:"sources"`
	Type          types.String   `tfsdk:"type"`
	Filters       []*filter      `tfsdk:"filters"`
	WorkspaceUUID types.String   `tfsdk:"workspace_uuid"`
}

// filter maps the resource schema data.
type filter struct {
	Name      types.String   `tfsdk:"name"`
	Operation types.String   `tfsdk:"operation"`
	Condition types.String   `tfsdk:"condition"`
	Values    []types.String `tfsdk:"values"`
}
