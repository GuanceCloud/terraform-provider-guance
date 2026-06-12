package monitor

import "github.com/hashicorp/terraform-plugin-framework/types"

type monitorDataSourceModel struct {
	UUID              types.String   `tfsdk:"uuid"`
	Name              types.String   `tfsdk:"name"`
	Type              types.String   `tfsdk:"type"`
	Status            types.Int64    `tfsdk:"status"`
	Extend            types.String   `tfsdk:"extend"`
	AlertPolicyUUIDs  []types.String `tfsdk:"alert_policy_uuids"`
	DashboardUUID     types.String   `tfsdk:"dashboard_uuid"`
	Tags              []types.String `tfsdk:"tags"`
	Secret            types.String   `tfsdk:"secret"`
	JsonScript        types.String   `tfsdk:"json_script"`
	OpenPermissionSet types.Bool     `tfsdk:"open_permission_set"`
	PermissionSet     []types.String `tfsdk:"permission_set"`
	CreateAt          types.Int64    `tfsdk:"create_at"`
	UpdateAt          types.Int64    `tfsdk:"update_at"`
	WorkspaceUUID     types.String   `tfsdk:"workspace_uuid"`
	MonitorUUID       types.String   `tfsdk:"monitor_uuid"`
	MonitorName       types.String   `tfsdk:"monitor_name"`
}
