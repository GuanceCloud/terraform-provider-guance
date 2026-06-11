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

func TestMuteFromPlanSupportsCheckerTagAndCustomRanges(t *testing.T) {
	cases := []struct {
		name       string
		plan       *muteResourceModel
		assertions func(t *testing.T, got *api.Mute)
	}{
		{
			name: "checker",
			plan: &muteResourceModel{
				Name:        types.StringValue("codex-checker-mute"),
				Description: types.StringValue("checker mute"),
				Type:        types.StringValue("checker"),
				MuteRanges: []muteRange{{
					Name:        types.StringValue("codex checker"),
					Type:        types.StringValue("monitor"),
					CheckerUUID: types.StringValue("rul_xxx"),
					MonitorUUID: types.StringValue("monitor_xxx"),
				}},
				FilterString:  types.StringValue("host:codex-checker"),
				RepeatTimeSet: types.Int64Value(0),
				StartTime:     types.StringValue("2026/12/31 10:00:00"),
				EndTime:       types.StringValue("2026/12/31 11:00:00"),
				Timezone:      types.StringValue("Asia/Shanghai"),
			},
			assertions: func(t *testing.T, got *api.Mute) {
				require.Equal(t, "checker", got.Type)
				require.Equal(t, "rul_xxx", got.MuteRanges[0].CheckerUUID)
				require.Equal(t, "monitor_xxx", got.MuteRanges[0].MonitorUUID)
				require.Equal(t, "host:codex-checker", got.FilterString)
			},
		},
		{
			name: "tag",
			plan: &muteResourceModel{
				Name: types.StringValue("codex-tag-mute"),
				Type: types.StringValue("tag"),
				MuteRanges: []muteRange{{
					Name:    types.StringValue("codex tag"),
					Type:    types.StringValue("tag"),
					TagUUID: types.StringValue("tag_xxx"),
				}},
				RepeatTimeSet: types.Int64Value(0),
				StartTime:     types.StringValue("2026/12/31 12:00:00"),
				EndTime:       types.StringValue("2026/12/31 13:00:00"),
				Timezone:      types.StringValue("Asia/Shanghai"),
			},
			assertions: func(t *testing.T, got *api.Mute) {
				require.Equal(t, "tag", got.Type)
				require.Equal(t, "tag_xxx", got.MuteRanges[0].TagUUID)
			},
		},
		{
			name: "custom",
			plan: &muteResourceModel{
				Name:          types.StringValue("codex-custom-mute"),
				Type:          types.StringValue("custom"),
				MuteRanges:    []muteRange{},
				FilterString:  types.StringValue("host:codex-custom AND service:api"),
				RepeatTimeSet: types.Int64Value(0),
				StartTime:     types.StringValue("2026/12/31 14:00:00"),
				EndTime:       types.StringValue("2026/12/31 15:00:00"),
				Timezone:      types.StringValue("Asia/Shanghai"),
				Declaration:   map[string]string{"source": "terraform"},
			},
			assertions: func(t *testing.T, got *api.Mute) {
				require.Equal(t, "custom", got.Type)
				require.Empty(t, got.MuteRanges)
				require.Equal(t, "host:codex-custom AND service:api", got.FilterString)
				require.Equal(t, "terraform", got.Declaration["source"])
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := muteFromPlan(tc.plan)
			tc.assertions(t, got)
		})
	}
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

func TestApplyContentToStatePreservesCheckerTagAndCustomRanges(t *testing.T) {
	state := &muteResourceModel{
		MuteRanges: []muteRange{{
			Name:        types.StringValue("existing checker"),
			CheckerUUID: types.StringValue("rul_existing"),
		}},
	}
	content := &api.MuteContent{
		UUID: "mute_xxx",
		Name: "codex-custom-mute",
		Type: "custom",
		MuteRanges: []api.MuteRange{
			{
				Name:        "codex checker",
				Type:        "monitor",
				CheckerUUID: "rul_xxx",
				MonitorUUID: "monitor_xxx",
			},
			{
				Name:    "codex tag",
				Type:    "tag",
				TagUUID: "tag_xxx",
			},
			{
				Name:            "codex policy",
				Type:            "alertPolicy",
				AlertPolicyUUID: "altpl_xxx",
			},
		},
		FilterString: "host:codex-custom",
	}

	applyContentToState(state, content)

	require.Equal(t, "mute_xxx", state.UUID.ValueString())
	require.Equal(t, "custom", state.Type.ValueString())
	require.Equal(t, "host:codex-custom", state.FilterString.ValueString())
	require.Len(t, state.MuteRanges, 3)
	require.Equal(t, "rul_xxx", state.MuteRanges[0].CheckerUUID.ValueString())
	require.Equal(t, "monitor_xxx", state.MuteRanges[0].MonitorUUID.ValueString())
	require.Equal(t, "tag_xxx", state.MuteRanges[1].TagUUID.ValueString())
	require.Equal(t, "altpl_xxx", state.MuteRanges[2].AlertPolicyUUID.ValueString())
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
