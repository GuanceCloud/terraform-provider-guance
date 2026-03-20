package monitor_json

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// monitorJsonResourceModel maps the resource schema data
type monitorJsonResourceModel struct {
	UUID            types.String `tfsdk:"uuid"`
	CheckerJson     types.String `tfsdk:"checker_json"`
	CheckerJsonExport types.String `tfsdk:"checker_json_export"`
	Type            types.String `tfsdk:"type"`
	CreateAt        types.Int64  `tfsdk:"create_at"`
	UpdateAt        types.Int64  `tfsdk:"update_at"`
	WorkspaceUUID   types.String `tfsdk:"workspace_uuid"`
}
