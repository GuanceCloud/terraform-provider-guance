package members

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// memberResourceModel maps the resource schema data.
type memberResourceModel struct {
	UUID     types.String `tfsdk:"uuid"`
	CreateAt types.String `tfsdk:"create_at"`
	Email    types.String `tfsdk:"email"`
	Roles    []roleModel  `tfsdk:"roles"`
	Name     types.String `tfsdk:"name"`
}

type roleModel struct {
	Name types.String `tfsdk:"name"`
	UUID types.String `tfsdk:"uuid"`
}

// memberDataSourceModel maps the resource schema data.
type memberDataSourceModel struct {
	Members []memberResourceModel `tfsdk:"members"`
	Search  types.String          `tfsdk:"search"`
}
