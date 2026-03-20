package alert_policy

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// alertPolicyResourceModel maps the resource schema data
type alertPolicyResourceModel struct {
	UUID               types.String            `tfsdk:"uuid"`
	Name               types.String            `tfsdk:"name"`
	Desc               types.String            `tfsdk:"desc"`
	OpenPermissionSet  types.Bool              `tfsdk:"open_permission_set"`
	PermissionSet      []types.String          `tfsdk:"permission_set"`
	CheckerUUIDs       []types.String          `tfsdk:"checker_uuids"`
	SecurityRuleUUIDs  []types.String          `tfsdk:"security_rule_uuids"`
	RuleTimezone       types.String            `tfsdk:"rule_timezone"`
	AlertOpt           *alertOptModel          `tfsdk:"alert_opt"`
	CreateAt           types.Int64             `tfsdk:"create_at"`
	UpdateAt           types.Int64             `tfsdk:"update_at"`
	WorkspaceUUID      types.String            `tfsdk:"workspace_uuid"`
}

// alertOptModel maps the alertOpt schema data
type alertOptModel struct {
	AggType                     types.String            `tfsdk:"agg_type"`
	IgnoreOK                    types.Bool              `tfsdk:"ignore_ok"`
	AlertType                   types.String            `tfsdk:"alert_type"`
	SilentTimeout               types.Int64             `tfsdk:"silent_timeout"`
	SilentTimeoutByStatusEnable types.Bool              `tfsdk:"silent_timeout_by_status_enable"`
	SilentTimeoutByStatus       []silentTimeoutByStatus `tfsdk:"silent_timeout_by_status"`
	AlertTarget                 []alertTarget           `tfsdk:"alert_target"`
	AggInterval                 types.Int64             `tfsdk:"agg_interval"`
	AggFields                   []types.String          `tfsdk:"agg_fields"`
	AggLabels                   []types.String          `tfsdk:"agg_labels"`
	AggClusterFields            []types.String          `tfsdk:"agg_cluster_fields"`
	AggSendFirst                types.Bool              `tfsdk:"agg_send_first"`
}

// silentTimeoutByStatus maps the silentTimeoutByStatus schema data
type silentTimeoutByStatus struct {
	Status        types.String `tfsdk:"status"`
	SilentTimeout types.Int64  `tfsdk:"silent_timeout"`
}

// alertTarget maps the alertTarget schema data
type alertTarget struct {
	Name             types.String      `tfsdk:"name"`
	Targets          []target          `tfsdk:"targets"`
	Crontab          types.String      `tfsdk:"crontab"`
	CrontabDuration  types.Int64       `tfsdk:"crontab_duration"`
	CustomDateUUIDs  []types.String    `tfsdk:"custom_date_uuids"`
	CustomStartTime  types.String      `tfsdk:"custom_start_time"`
	CustomDuration   types.Int64       `tfsdk:"custom_duration"`
	AlertInfo        []alertInfo       `tfsdk:"alert_info"`
}

// target maps the targets schema data
type target struct {
	To              []types.String      `tfsdk:"to"`
	Status          types.String        `tfsdk:"status"`
	DfSource        types.String        `tfsdk:"df_source"`
	UpgradeTargets  []upgradeTarget     `tfsdk:"upgrade_targets"`
	Tags            map[string][]string `tfsdk:"tags"`
	FilterString    types.String        `tfsdk:"filter_string"`
}

// upgradeTarget maps the upgradeTargets schema data
type upgradeTarget struct {
	To       []types.String `tfsdk:"to"`
	Duration types.Int64    `tfsdk:"duration"`
	ToWay    []types.String `tfsdk:"to_way"`
}

// alertInfo maps the alertInfo schema data
type alertInfo struct {
	Name        types.String `tfsdk:"name"`
	Targets     []target     `tfsdk:"targets"`
	FilterString types.String `tfsdk:"filter_string"`
	MemberInfo  []types.String `tfsdk:"member_info"`
}
