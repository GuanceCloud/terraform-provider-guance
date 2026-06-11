package alert_policy_notice_date

import "github.com/hashicorp/terraform-plugin-framework/types"

type alertPolicyNoticeDateResourceModel struct {
	UUID                 types.String   `tfsdk:"uuid"`
	Name                 types.String   `tfsdk:"name"`
	NoticeDates          []types.String `tfsdk:"notice_dates"`
	SkipRefCheckOnDelete types.Bool     `tfsdk:"skip_ref_check_on_delete"`
	CreateAt             types.Int64    `tfsdk:"create_at"`
	UpdateAt             types.Int64    `tfsdk:"update_at"`
	WorkspaceUUID        types.String   `tfsdk:"workspace_uuid"`
}
