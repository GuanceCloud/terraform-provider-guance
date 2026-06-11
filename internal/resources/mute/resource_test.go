package mute

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/require"

	"github.com/GuanceCloud/terraform-provider-guance/internal/api"
)

func TestMuteFromPlanRepeatedMuteWithNotify(t *testing.T) {
	plan := &muteResourceModel{
		Name:        types.StringValue("codex-mute-weekly"),
		Description: types.StringValue("weekly mute"),
		Type:        types.StringValue("alertPolicy"),
		MuteRanges: []muteRange{{
			Name:            types.StringValue("codex-policy"),
			AlertPolicyUUID: types.StringValue("altpl_xxx"),
		}},
		NotifyTargets: []notifyTarget{{
			Type: types.StringValue("notifyObject"),
			To:   []types.String{types.StringValue("notify_xxx")},
		}},
		NotifyMessage: types.StringValue("mute starts soon"),
		NotifyTimeStr: types.StringValue("2026/12/31 11:50:00"),
		RepeatTimeSet: types.Int64Value(1),
		RepeatCrontabSet: &repeatCrontabSet{
			Min:   types.StringValue("0"),
			Hour:  types.StringValue("0"),
			Day:   types.StringValue("*"),
			Month: types.StringValue("*"),
			Week:  types.StringValue("1,2,3,4,5"),
		},
		CrontabDuration:  types.Int64Value(3600),
		RepeatExpireTime: types.StringValue("0"),
		Timezone:         types.StringValue("Asia/Shanghai"),
		Tags: map[string][]string{
			"service": {"codex-weekly"},
		},
	}

	got := muteFromPlan(plan)

	require.Equal(t, "codex-mute-weekly", got.Name)
	require.Equal(t, "weekly mute", got.Description)
	require.Equal(t, "alertPolicy", got.Type)
	require.Len(t, got.MuteRanges, 1)
	require.Equal(t, "codex-policy", got.MuteRanges[0].Name)
	require.Equal(t, "altpl_xxx", got.MuteRanges[0].AlertPolicyUUID)
	require.Len(t, got.NotifyTargets, 1)
	require.Equal(t, "notifyObject", got.NotifyTargets[0].Type)
	require.Equal(t, []string{"notify_xxx"}, got.NotifyTargets[0].To)
	require.Equal(t, "mute starts soon", got.NotifyMessage)
	require.Equal(t, "2026/12/31 11:50:00", got.NotifyTimeStr)
	require.Equal(t, 1, got.RepeatTimeSet)
	require.NotNil(t, got.RepeatCrontabSet)
	require.Equal(t, "0", got.RepeatCrontabSet.Min)
	require.Equal(t, "1,2,3,4,5", got.RepeatCrontabSet.Week)
	require.Equal(t, 3600, got.CrontabDuration)
	require.Equal(t, "0", got.RepeatExpireTime)
	require.Equal(t, "Asia/Shanghai", got.Timezone)
	require.Equal(t, []string{"codex-weekly"}, got.Tags["service"])
}

func TestApplyContentToStateInfersRepeatedMuteAndPreservesUnconfiguredWindow(t *testing.T) {
	state := &muteResourceModel{
		RepeatTimeSet: types.Int64Value(1),
		StartTime:     types.StringNull(),
		EndTime:       types.StringNull(),
		Timezone:      types.StringValue("Asia/Shanghai"),
	}
	content := &api.MuteContent{
		UUID:          "mute_xxx",
		Name:          "codex-mute-weekly",
		Type:          "alertPolicy",
		StartTime:     "2026/06/11 00:00:00",
		EndTime:       "2026/06/11 01:00:00",
		RepeatTimeSet: 0,
		RepeatCrontabSet: &api.RepeatCrontabSet{
			Min:   "0",
			Hour:  "0",
			Day:   "*",
			Month: "*",
			Week:  "1,2,3,4,5",
		},
		CrontabDuration:  3600,
		RepeatExpireTime: "-1",
		Timezone:         "Asia/Shanghai",
		Status:           0,
		WorkspaceUUID:    "wksp_xxx",
	}

	applyContentToState(state, content)

	require.Equal(t, "mute_xxx", state.UUID.ValueString())
	require.Equal(t, int64(1), state.RepeatTimeSet.ValueInt64())
	require.True(t, state.StartTime.IsNull())
	require.True(t, state.EndTime.IsNull())
	require.NotNil(t, state.RepeatCrontabSet)
	require.Equal(t, "0", state.RepeatCrontabSet.Min.ValueString())
	require.Equal(t, "1,2,3,4,5", state.RepeatCrontabSet.Week.ValueString())
	require.Equal(t, int64(3600), state.CrontabDuration.ValueInt64())
	require.True(t, state.RepeatExpireTime.IsNull())
	require.Equal(t, "wksp_xxx", state.WorkspaceUUID.ValueString())
}

func TestApplyContentToStateOneTimeMuteKeepsReturnedWindowAndNotifyTargets(t *testing.T) {
	state := &muteResourceModel{
		RepeatTimeSet: types.Int64Value(0),
		StartTime:     types.StringValue("2026/12/31 10:00:00"),
		EndTime:       types.StringValue("2026/12/31 11:00:00"),
	}
	content := &api.MuteContent{
		UUID:          "mute_xxx",
		Name:          "codex-mute-notify",
		Type:          "alertPolicy",
		StartTime:     "2026/12/31 12:00:00",
		EndTime:       "2026/12/31 13:00:00",
		NotifyTimeStr: "2026/12/31 11:50:00",
		NotifyMessage: "mute starts soon",
		NotifyTargets: []api.MuteNotifyTarget{{
			Type: "notifyObject",
			To:   []string{"notify_xxx"},
		}},
	}

	applyContentToState(state, content)

	require.Equal(t, int64(0), state.RepeatTimeSet.ValueInt64())
	require.Equal(t, "2026/12/31 12:00:00", state.StartTime.ValueString())
	require.Equal(t, "2026/12/31 13:00:00", state.EndTime.ValueString())
	require.Equal(t, "2026/12/31 11:50:00", state.NotifyTimeStr.ValueString())
	require.Equal(t, "mute starts soon", state.NotifyMessage.ValueString())
	require.Len(t, state.NotifyTargets, 1)
	require.Equal(t, "notifyObject", state.NotifyTargets[0].Type.ValueString())
	require.Equal(t, []types.String{types.StringValue("notify_xxx")}, state.NotifyTargets[0].To)
}
