package mute

import "github.com/hashicorp/terraform-plugin-framework/types"

type muteResourceModel struct {
	UUID             types.String        `tfsdk:"uuid"`
	Name             types.String        `tfsdk:"name"`
	Description      types.String        `tfsdk:"description"`
	Type             types.String        `tfsdk:"type"`
	MuteRanges       []muteRange         `tfsdk:"mute_ranges"`
	Tags             map[string][]string `tfsdk:"tags"`
	FilterString     types.String        `tfsdk:"filter_string"`
	NotifyTargets    []notifyTarget      `tfsdk:"notify_targets"`
	NotifyMessage    types.String        `tfsdk:"notify_message"`
	NotifyTimeStr    types.String        `tfsdk:"notify_time_str"`
	StartTime        types.String        `tfsdk:"start_time"`
	EndTime          types.String        `tfsdk:"end_time"`
	RepeatTimeSet    types.Int64         `tfsdk:"repeat_time_set"`
	RepeatCrontabSet *repeatCrontabSet   `tfsdk:"repeat_crontab_set"`
	CrontabDuration  types.Int64         `tfsdk:"crontab_duration"`
	RepeatExpireTime types.String        `tfsdk:"repeat_expire_time"`
	Timezone         types.String        `tfsdk:"timezone"`
	Declaration      map[string]string   `tfsdk:"declaration"`
	Enabled          types.Bool          `tfsdk:"enabled"`
	Status           types.Int64         `tfsdk:"status"`
	CreateAt         types.Int64         `tfsdk:"create_at"`
	UpdateAt         types.Int64         `tfsdk:"update_at"`
	WorkspaceUUID    types.String        `tfsdk:"workspace_uuid"`
}

type muteRange struct {
	Name            types.String `tfsdk:"name"`
	Type            types.String `tfsdk:"type"`
	CheckerUUID     types.String `tfsdk:"checker_uuid"`
	MonitorUUID     types.String `tfsdk:"monitor_uuid"`
	SLOUUID         types.String `tfsdk:"slo_uuid"`
	AlertPolicyUUID types.String `tfsdk:"alert_policy_uuid"`
	TagUUID         types.String `tfsdk:"tag_uuid"`
}

type notifyTarget struct {
	To   []types.String `tfsdk:"to"`
	Type types.String   `tfsdk:"type"`
}

type repeatCrontabSet struct {
	Min   types.String `tfsdk:"min"`
	Hour  types.String `tfsdk:"hour"`
	Day   types.String `tfsdk:"day"`
	Month types.String `tfsdk:"month"`
	Week  types.String `tfsdk:"week"`
}
