package role

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// roleResourceModel maps the resource schema data.
type roleResourceModel struct {
	UUID types.String   `tfsdk:"uuid"`
	Name types.String   `tfsdk:"name"`
	Desc types.String   `tfsdk:"desc"`
	Keys []types.String `tfsdk:"keys"`
}
