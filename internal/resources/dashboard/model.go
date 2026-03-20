package dashboard

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// dashboardResourceModel maps the resource schema data
type dashboardResourceModel struct {
	UUID                 types.String   `tfsdk:"uuid"`
	Name                 types.String   `tfsdk:"name"`
	Desc                 types.String   `tfsdk:"desc"`
	Identifier           types.String   `tfsdk:"identifier"`
	TagNames             []types.String `tfsdk:"tag_names"`
	TemplateInfo         types.String   `tfsdk:"template_info"`
	TemplateInfoExport   types.String   `tfsdk:"template_info_export"`
	SpecifyDashboardUUID types.String   `tfsdk:"specify_dashboard_uuid"`
	IsPublic             types.Int64    `tfsdk:"is_public"`
	PermissionSet        []types.String `tfsdk:"permission_set"`
	ReadPermissionSet    []types.String `tfsdk:"read_permission_set"`
	CreateAt             types.Int64    `tfsdk:"create_at"`
	UpdateAt             types.Int64    `tfsdk:"update_at"`
	WorkspaceUUID        types.String   `tfsdk:"workspace_uuid"`
}
