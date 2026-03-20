package slo

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// sloResourceModel maps the resource schema data
type sloResourceModel struct {
	UUID              types.String   `tfsdk:"uuid"`
	Name              types.String   `tfsdk:"name"`
	Interval          types.String   `tfsdk:"interval"`
	Goal              types.Float64  `tfsdk:"goal"`
	MinGoal           types.Float64  `tfsdk:"min_goal"`
	SliUUIDs          []types.String `tfsdk:"sli_uuids"`
	Describe          types.String   `tfsdk:"describe"`
	AlertPolicyUUIDs  []types.String `tfsdk:"alert_policy_uuids"`
	Tags              []types.String `tfsdk:"tags"`
	CreateAt          types.Int64    `tfsdk:"create_at"`
	UpdateAt          types.Int64    `tfsdk:"update_at"`
	WorkspaceUUID     types.String   `tfsdk:"workspace_uuid"`
}