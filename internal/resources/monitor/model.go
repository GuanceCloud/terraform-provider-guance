package monitor

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// monitorResourceModel maps the resource schema data
type monitorResourceModel struct {
	UUID              types.String   `tfsdk:"uuid"`
	Type              types.String   `tfsdk:"type"`
	Status            types.Int64    `tfsdk:"status"`
	Extend            types.String   `tfsdk:"extend"`
	AlertPolicyUUIDs  []types.String `tfsdk:"alert_policy_uuids"`
	DashboardUUID     types.String   `tfsdk:"dashboard_uuid"`
	Tags              []types.String `tfsdk:"tags"`
	Secret            types.String   `tfsdk:"secret"`
	JsonScript        JsonScript     `tfsdk:"json_script"`
	OpenPermissionSet types.Bool     `tfsdk:"open_permission_set"`
	PermissionSet     []types.String `tfsdk:"permission_set"`
	CreateAt          types.Int64    `tfsdk:"create_at"`
	UpdateAt          types.Int64    `tfsdk:"update_at"`
	WorkspaceUUID     types.String   `tfsdk:"workspace_uuid"`
	MonitorUUID       types.String   `tfsdk:"monitor_uuid"`
	MonitorName       types.String   `tfsdk:"monitor_name"`
}

// JsonScript represents the json_script nested structure
type JsonScript struct {
	Type                 types.String   `tfsdk:"type"`
	Title                types.String   `tfsdk:"title"`
	Message              types.String   `tfsdk:"message"`
	Every                types.String   `tfsdk:"every"`
	Interval             types.Int64    `tfsdk:"interval"`
	RecoverNeedPeriodCount types.Int64    `tfsdk:"recover_need_period_count"`
	DisableCheckEndTime  types.Bool     `tfsdk:"disable_check_end_time"`
	GroupBy              []types.String `tfsdk:"group_by"`
	Targets              []Target       `tfsdk:"targets"`
	CheckerOpt           CheckerOpt     `tfsdk:"checker_opt"`
	Channels             []types.String `tfsdk:"channels"`
	AtAccounts           []types.String `tfsdk:"at_accounts"`
	AtNoDataAccounts     []types.String `tfsdk:"at_no_data_accounts"`
}

// Target represents a target in the targets list
type Target struct {
	Dql   types.String `tfsdk:"dql"`
	Alias types.String `tfsdk:"alias"`
	Qtype types.String `tfsdk:"qtype"`
}

// CheckerOpt represents checker options
type CheckerOpt struct {
	InfoEvent types.Bool `tfsdk:"info_event"`
	Rules     []Rule     `tfsdk:"rules"`
}

// Rule represents a rule in the rules list
type Rule struct {
	ConditionLogic types.String    `tfsdk:"condition_logic"`
	Conditions     []Condition     `tfsdk:"conditions"`
	Status         types.String    `tfsdk:"status"`
}

// Condition represents a condition in the conditions list
type Condition struct {
	Alias     types.String   `tfsdk:"alias"`
	Operands  []types.String `tfsdk:"operands"`
	Operator  types.String   `tfsdk:"operator"`
}
