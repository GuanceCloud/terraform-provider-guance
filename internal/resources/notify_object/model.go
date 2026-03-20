package notify_object

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// notifyObjectResourceModel maps the resource schema data
type notifyObjectResourceModel struct {
	UUID                types.String   `tfsdk:"uuid"`
	Type                types.String   `tfsdk:"type"`
	Name                types.String   `tfsdk:"name"`
	OptSet              types.String   `tfsdk:"opt_set"`
	OpenPermissionSet   types.Bool     `tfsdk:"open_permission_set"`
	PermissionSet       []types.String `tfsdk:"permission_set"`
	CreateAt            types.Int64    `tfsdk:"create_at"`
	UpdateAt            types.Int64    `tfsdk:"update_at"`
	WorkspaceUUID       types.String   `tfsdk:"workspace_uuid"`
}
