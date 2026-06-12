package monitors

import "github.com/hashicorp/terraform-plugin-framework/types"

type monitorsDataSourceModel struct {
	Search          types.String       `tfsdk:"search"`
	Type            types.String       `tfsdk:"type"`
	Status          types.String       `tfsdk:"status"`
	TagsUUID        types.String       `tfsdk:"tags_uuid"`
	AlertPolicyUUID types.String       `tfsdk:"alert_policy_uuid"`
	DashboardUUID   types.String       `tfsdk:"dashboard_uuid"`
	CheckerUUID     types.String       `tfsdk:"checker_uuid"`
	Monitors        []monitorListModel `tfsdk:"monitors"`
}

type monitorListModel struct {
	UUID             types.String   `tfsdk:"uuid"`
	Name             types.String   `tfsdk:"name"`
	Type             types.String   `tfsdk:"type"`
	Status           types.Int64    `tfsdk:"status"`
	AlertPolicyUUIDs []types.String `tfsdk:"alert_policy_uuids"`
	DashboardUUID    types.String   `tfsdk:"dashboard_uuid"`
	Tags             []types.String `tfsdk:"tags"`
	WorkspaceUUID    types.String   `tfsdk:"workspace_uuid"`
	MonitorUUID      types.String   `tfsdk:"monitor_uuid"`
	MonitorName      types.String   `tfsdk:"monitor_name"`
}
