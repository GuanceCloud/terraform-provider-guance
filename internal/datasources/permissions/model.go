package permissions

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type permissionSubResourceModel struct {
	Desc                types.String `tfsdk:"desc"`
	Disabled            types.Int64  `tfsdk:"disabled"`
	IsSupportCustomRole types.Int64  `tfsdk:"is_support_custom_role"`
	IsSupportGeneral    types.Int64  `tfsdk:"is_support_general"`
	IsSupportOwner      types.Int64  `tfsdk:"is_support_owner"`
	IsSupportReadOnly   types.Int64  `tfsdk:"is_support_read_only"`
	IsSupportWsAdmin    types.Int64  `tfsdk:"is_support_ws_admin"`
	Key                 types.String `tfsdk:"key"`
	Name                types.String `tfsdk:"name"`
}

// permissionResourceModel maps the resource schema data.
type permissionResourceModel struct {
	Desc                types.String                  `tfsdk:"desc"`
	Disabled            types.Int64                   `tfsdk:"disabled"`
	IsSupportCustomRole types.Int64                   `tfsdk:"is_support_custom_role"`
	IsSupportGeneral    types.Int64                   `tfsdk:"is_support_general"`
	IsSupportOwner      types.Int64                   `tfsdk:"is_support_owner"`
	IsSupportReadOnly   types.Int64                   `tfsdk:"is_support_read_only"`
	IsSupportWsAdmin    types.Int64                   `tfsdk:"is_support_ws_admin"`
	Key                 types.String                  `tfsdk:"key"`
	Name                types.String                  `tfsdk:"name"`
	Subs                []*permissionSubResourceModel `tfsdk:"subs"`
}

// permissionDataSourceModel maps the resource schema data.
type permissionDataSourceModel struct {
	Permissions         []permissionResourceModel `tfsdk:"permissions"`
	IsSupportCustomRole types.Bool                `tfsdk:"is_support_custom_role"`
}
