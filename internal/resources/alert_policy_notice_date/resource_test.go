package alert_policy_notice_date

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/require"

	"github.com/GuanceCloud/terraform-provider-guance/internal/api"
)

func TestNoticeDateFromPlan(t *testing.T) {
	plan := &alertPolicyNoticeDateResourceModel{
		Name: types.StringValue("codex-notice-date"),
		NoticeDates: []types.String{
			types.StringValue("2026/06/10"),
			types.StringValue("2026/06/11"),
			types.StringNull(),
			types.StringUnknown(),
		},
	}

	got := noticeDateFromPlan(plan)

	require.Equal(t, "codex-notice-date", got.Name)
	require.Equal(t, []string{"2026/06/10", "2026/06/11"}, got.NoticeDates)
}

func TestApplyContentToState(t *testing.T) {
	state := &alertPolicyNoticeDateResourceModel{}
	content := &api.AlertPolicyNoticeDateContent{
		UUID:          "ndate_xxx",
		Name:          "codex-notice-date",
		Dates:         []string{"2026/06/10", "2026/06/11"},
		CreateAt:      1781165177,
		UpdateAt:      1781165200,
		WorkspaceUUID: "wksp_xxx",
	}

	applyContentToState(state, content)

	require.Equal(t, "ndate_xxx", state.UUID.ValueString())
	require.Equal(t, "codex-notice-date", state.Name.ValueString())
	require.Equal(t, []types.String{types.StringValue("2026/06/10"), types.StringValue("2026/06/11")}, state.NoticeDates)
	require.Equal(t, int64(1781165177), state.CreateAt.ValueInt64())
	require.Equal(t, int64(1781165200), state.UpdateAt.ValueInt64())
	require.Equal(t, "wksp_xxx", state.WorkspaceUUID.ValueString())
}
