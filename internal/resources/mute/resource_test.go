package mute

import (
	"encoding/json"
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

func TestMuteUpdateBodyPreservesClearableZeroValues(t *testing.T) {
	got := muteUpdateBody(&api.Mute{
		Name:          "codex-mute",
		Description:   "",
		Type:          "custom",
		MuteRanges:    []api.MuteRange{},
		RepeatTimeSet: 0,
		Timezone:      "Asia/Shanghai",
	})

	require.Equal(t, "codex-mute", got["name"])
	require.Equal(t, "", got["description"])
	require.Equal(t, "custom", got["type"])
	require.Equal(t, []api.MuteRange{}, got["muteRanges"])
	require.Equal(t, map[string][]string{}, got["tags"])
	require.Equal(t, "", got["filterString"])
	require.Equal(t, []api.MuteNotifyTarget{}, got["notifyTargets"])
	require.Equal(t, "", got["notifyMessage"])
	require.Equal(t, "", got["notifyTimeStr"])
	require.NotContains(t, got, "startTime")
	require.NotContains(t, got, "endTime")
	require.Equal(t, 0, got["repeatTimeSet"])
	require.Nil(t, got["repeatCrontabSet"])
	require.Equal(t, 0, got["crontabDuration"])
	require.Equal(t, "", got["repeatExpireTime"])
	require.Equal(t, "Asia/Shanghai", got["timezone"])
	require.Equal(t, map[string]string{}, got["declaration"])
}

func TestMuteUpdateBodyIncludesConfiguredTimes(t *testing.T) {
	got := muteUpdateBody(&api.Mute{
		Name:          "codex-mute",
		Type:          "alertPolicy",
		MuteRanges:    []api.MuteRange{},
		RepeatTimeSet: 1,
		StartTime:     "2026/12/31 10:00:00",
		EndTime:       "2026/12/31 11:00:00",
		Timezone:      "Asia/Shanghai",
	})

	require.Equal(t, "2026/12/31 10:00:00", got["startTime"])
	require.Equal(t, "2026/12/31 11:00:00", got["endTime"])
}

func TestApplyContentToStateInfersRepeatedMuteAndPreservesUnconfiguredWindow(t *testing.T) {
	state := &muteResourceModel{
		RepeatTimeSet: types.Int64Value(1),
		StartTime:     types.StringNull(),
		EndTime:       types.StringNull(),
		Timezone:      types.StringValue("Asia/Shanghai"),
		Enabled:       types.BoolValue(true),
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
	require.True(t, state.Enabled.ValueBool())
	require.Equal(t, "wksp_xxx", state.WorkspaceUUID.ValueString())
}

func TestApplyContentToStatePreservesConfiguredRepeatedWindowWhenAPIReturnsEmpty(t *testing.T) {
	state := &muteResourceModel{
		RepeatTimeSet: types.Int64Value(1),
		StartTime:     types.StringValue("2026/12/31 10:00:00"),
		EndTime:       types.StringValue("2026/12/31 11:00:00"),
		Timezone:      types.StringValue("Asia/Shanghai"),
		Enabled:       types.BoolValue(true),
	}
	content := muteContentFromJSON(t, `{
		"uuid": "mute_xxx",
		"name": "codex-mute-weekly",
		"type": "alertPolicy",
		"startTime": "2026/12/31 10:00:00",
		"endTime": "",
		"repeatTimeSet": 1,
		"repeatCrontabSet": {
			"min": "0",
			"hour": "10",
			"day": "*",
			"month": "*",
			"week": "*"
		},
		"crontabDuration": 3600,
		"repeatExpireTime": "2027/01/01 00:00:00",
		"timezone": "Asia/Shanghai",
		"status": 0
	}`)

	applyContentToState(state, content)

	require.Equal(t, "2026/12/31 10:00:00", state.StartTime.ValueString())
	require.Equal(t, "2026/12/31 11:00:00", state.EndTime.ValueString())
	require.NotNil(t, state.RepeatCrontabSet)
	require.Equal(t, int64(3600), state.CrontabDuration.ValueInt64())
}

func TestApplyContentToStateMapsDisabledStatusToEnabled(t *testing.T) {
	state := &muteResourceModel{
		Enabled: types.BoolValue(true),
	}
	content := &api.MuteContent{
		UUID:   "mute_xxx",
		Name:   "codex-disabled-mute",
		Type:   "custom",
		Status: 2,
	}

	applyContentToState(state, content)

	require.Equal(t, int64(2), state.Status.ValueInt64())
	require.False(t, state.Enabled.ValueBool())
}

func TestMuteEnabledFromStatusPreservesUnknownStatuses(t *testing.T) {
	require.True(t, muteEnabledFromStatus(1, types.BoolValue(true)).ValueBool())
	require.False(t, muteEnabledFromStatus(1, types.BoolValue(false)).ValueBool())
	require.True(t, muteEnabledFromStatus(0, types.BoolValue(false)).ValueBool())
	require.False(t, muteEnabledFromStatus(2, types.BoolValue(true)).ValueBool())
}

func TestMuteStatusNeedsEnabledChange(t *testing.T) {
	require.False(t, muteStatusNeedsEnabledChange(0, true))
	require.False(t, muteStatusNeedsEnabledChange(2, false))
	require.False(t, muteStatusNeedsEnabledChange(1, true))
	require.True(t, muteStatusNeedsEnabledChange(2, true))
	require.True(t, muteStatusNeedsEnabledChange(0, false))
	require.True(t, muteStatusNeedsEnabledChange(1, false))
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

func TestApplyContentToStateAppliesRemoteClears(t *testing.T) {
	state := &muteResourceModel{
		Description: types.StringValue("existing description"),
		Type:        types.StringValue("alertPolicy"),
		MuteRanges: []muteRange{{
			Name: types.StringValue("existing policy"),
		}},
		Tags: map[string][]string{
			"service": {"api"},
		},
		FilterString: types.StringValue("host:old"),
		NotifyTargets: []notifyTarget{{
			Type: types.StringValue("notifyObject"),
			To:   []types.String{types.StringValue("notify_xxx")},
		}},
		NotifyMessage: types.StringValue("old message"),
		NotifyTimeStr: types.StringValue("2026/12/31 11:50:00"),
		StartTime:     types.StringValue("2026/12/31 12:00:00"),
		EndTime:       types.StringValue("2026/12/31 13:00:00"),
		RepeatTimeSet: types.Int64Value(1),
		RepeatCrontabSet: &repeatCrontabSet{
			Min:   types.StringValue("0"),
			Hour:  types.StringValue("9"),
			Day:   types.StringValue("*"),
			Month: types.StringValue("*"),
			Week:  types.StringValue("1"),
		},
		CrontabDuration:  types.Int64Value(3600),
		RepeatExpireTime: types.StringValue("2026/12/31 23:59:59"),
		Timezone:         types.StringValue("Asia/Shanghai"),
		Declaration: map[string]string{
			"reason": "maintenance",
		},
		Enabled: types.BoolValue(true),
	}
	content := muteContentFromJSON(t, `{
		"uuid": "mute_xxx",
		"name": "codex-mute-clear",
		"description": "",
		"type": "alertPolicy",
		"muteRanges": [],
		"tags": {},
		"filterString": "",
		"notifyTargets": [],
		"notifyMessage": "",
		"notifyTimeStr": "",
		"startTime": "",
		"endTime": "",
		"repeatTimeSet": 0,
		"repeatCrontabSet": null,
		"crontabDuration": 0,
		"repeatExpireTime": "",
		"timezone": "Asia/Shanghai",
		"declaration": {},
		"status": 0
	}`)

	applyContentToState(state, content)

	require.Equal(t, "mute_xxx", state.UUID.ValueString())
	require.Equal(t, "", state.Description.ValueString())
	require.Empty(t, state.MuteRanges)
	require.Empty(t, state.Tags)
	require.Equal(t, "", state.FilterString.ValueString())
	require.Empty(t, state.NotifyTargets)
	require.Equal(t, "", state.NotifyMessage.ValueString())
	require.Equal(t, "", state.NotifyTimeStr.ValueString())
	require.Equal(t, int64(0), state.RepeatTimeSet.ValueInt64())
	require.Equal(t, "", state.StartTime.ValueString())
	require.Equal(t, "", state.EndTime.ValueString())
	require.Nil(t, state.RepeatCrontabSet)
	require.Equal(t, int64(0), state.CrontabDuration.ValueInt64())
	require.Equal(t, "", state.RepeatExpireTime.ValueString())
	require.Empty(t, state.Declaration)
}

func TestApplyContentToStateKeepsUnconfiguredOptionalNulls(t *testing.T) {
	state := &muteResourceModel{
		FilterString:     types.StringNull(),
		NotifyMessage:    types.StringNull(),
		NotifyTimeStr:    types.StringNull(),
		CrontabDuration:  types.Int64Null(),
		RepeatExpireTime: types.StringNull(),
	}
	content := muteContentFromJSON(t, `{
		"uuid": "mute_xxx",
		"name": "codex-mute-null",
		"type": "alertPolicy",
		"muteRanges": [],
		"tags": {},
		"filterString": "",
		"notifyTargets": [],
		"notifyMessage": "",
		"notifyTimeStr": "",
		"crontabDuration": 0,
		"repeatExpireTime": "",
		"declaration": {}
	}`)

	applyContentToState(state, content)

	require.Empty(t, state.MuteRanges)
	require.Nil(t, state.Tags)
	require.True(t, state.FilterString.IsNull())
	require.Empty(t, state.NotifyTargets)
	require.True(t, state.NotifyMessage.IsNull())
	require.True(t, state.NotifyTimeStr.IsNull())
	require.True(t, state.CrontabDuration.IsNull())
	require.True(t, state.RepeatExpireTime.IsNull())
	require.Nil(t, state.Declaration)
}

func TestMuteContentFieldPresentDistinguishesExplicitClears(t *testing.T) {
	content := muteContentFromJSON(t, `{
		"filterString": "",
		"notifyTargets": [],
		"repeatCrontabSet": null,
		"declaration": {}
	}`)

	require.True(t, content.FieldPresent("filterString"))
	require.True(t, content.FieldPresent("notifyTargets"))
	require.True(t, content.FieldPresent("repeatCrontabSet"))
	require.True(t, content.FieldPresent("declaration"))
	require.False(t, content.FieldPresent("notifyMessage"))
	require.False(t, content.FieldPresent("tags"))
}

func muteContentFromJSON(t *testing.T, value string) *api.MuteContent {
	t.Helper()
	var content api.MuteContent
	require.NoError(t, json.Unmarshal([]byte(value), &content))
	return &content
}
