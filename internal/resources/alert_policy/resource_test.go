package alert_policy

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/require"

	"github.com/GuanceCloud/terraform-provider-guance/internal/api"
)

func TestGetAlertPolicyFromPlanStatusMode(t *testing.T) {
	resource := &alertPolicyResource{}
	plan := &alertPolicyResourceModel{
		Name:              types.StringValue("codex-status-policy"),
		Desc:              types.StringValue("status mode"),
		OpenPermissionSet: types.BoolValue(true),
		PermissionSet:     []types.String{types.StringValue("wsAdmin")},
		CheckerUUIDs:      []types.String{types.StringValue("rul_xxx")},
		RuleTimezone:      types.StringValue("Asia/Shanghai"),
		AlertOpt: &alertOptModel{
			AggType:                     types.StringValue("byFields"),
			IgnoreOK:                    types.BoolValue(true),
			AlertType:                   types.StringValue("status"),
			AggInterval:                 types.Int64Value(60),
			AggFields:                   []types.String{types.StringValue("df_monitor_checker_id"), types.StringValue("df_label")},
			AggLabels:                   []types.String{types.StringValue("service")},
			AggClusterFields:            []types.String{types.StringValue("df_title")},
			AggSendFirst:                types.BoolValue(true),
			SilentTimeout:               types.Int64Value(300),
			SilentTimeoutByStatusEnable: types.BoolValue(true),
			SilentTimeoutByStatus: []silentTimeoutByStatus{{
				Status:        types.StringValue("critical"),
				SilentTimeout: types.Int64Value(120),
			}},
			AlertTarget: []alertTarget{{
				Name:            types.StringValue("default"),
				CustomDateUUIDs: []types.String{types.StringValue("ndate_xxx")},
				CustomStartTime: types.StringValue("09:30:00"),
				CustomDuration:  types.Int64Value(3600),
				Targets: []target{{
					To:           []types.String{types.StringValue("notify_xxx")},
					Status:       types.StringValue("critical,error,warning"),
					DfSource:     types.StringValue("security"),
					FilterString: types.StringValue("host:codex-alert-test"),
					Tags: map[string][]string{
						"service": {"codex"},
					},
					UpgradeTargets: []upgradeTarget{{
						To:       []types.String{types.StringValue("notify_yyy")},
						Duration: types.Int64Value(300),
						ToWay:    []types.String{types.StringValue("mail")},
					}},
				}},
			}},
		},
	}

	got := resource.getAlertPolicyFromPlan(plan)

	require.Equal(t, "codex-status-policy", got.Name)
	require.Equal(t, "status mode", got.Desc)
	require.True(t, got.OpenPermissionSet)
	require.Equal(t, []string{"wsAdmin"}, got.PermissionSet)
	require.Equal(t, []string{"rul_xxx"}, got.CheckerUUIDs)
	require.Equal(t, "Asia/Shanghai", got.RuleTimezone)
	require.NotNil(t, got.AlertOpt)
	require.Equal(t, "byFields", got.AlertOpt.AggType)
	require.True(t, got.AlertOpt.IgnoreOK)
	require.Equal(t, "status", got.AlertOpt.AlertType)
	require.Equal(t, 60, got.AlertOpt.AggInterval)
	require.Equal(t, []string{"df_monitor_checker_id", "df_label"}, got.AlertOpt.AggFields)
	require.Equal(t, []string{"service"}, got.AlertOpt.AggLabels)
	require.Equal(t, []string{"df_title"}, got.AlertOpt.AggClusterFields)
	require.True(t, got.AlertOpt.AggSendFirst)
	require.Equal(t, 300, got.AlertOpt.SilentTimeout)
	require.True(t, got.AlertOpt.SilentTimeoutByStatusEnable)
	require.Len(t, got.AlertOpt.SilentTimeoutByStatus, 1)
	require.Equal(t, "critical", got.AlertOpt.SilentTimeoutByStatus[0].Status)
	require.Equal(t, 120, got.AlertOpt.SilentTimeoutByStatus[0].SilentTimeout)
	require.Len(t, got.AlertOpt.AlertTarget, 1)
	require.Equal(t, "default", got.AlertOpt.AlertTarget[0].Name)
	require.Equal(t, []string{"ndate_xxx"}, got.AlertOpt.AlertTarget[0].CustomDateUUIDs)
	require.Equal(t, "09:30:00", got.AlertOpt.AlertTarget[0].CustomStartTime)
	require.Equal(t, 3600, got.AlertOpt.AlertTarget[0].CustomDuration)
	require.Len(t, got.AlertOpt.AlertTarget[0].Targets, 1)
	require.Equal(t, []string{"notify_xxx"}, got.AlertOpt.AlertTarget[0].Targets[0].To)
	require.Equal(t, "critical,error,warning", got.AlertOpt.AlertTarget[0].Targets[0].Status)
	require.Equal(t, "security", got.AlertOpt.AlertTarget[0].Targets[0].DfSource)
	require.Equal(t, "host:codex-alert-test", got.AlertOpt.AlertTarget[0].Targets[0].FilterString)
	require.Equal(t, []string{"codex"}, got.AlertOpt.AlertTarget[0].Targets[0].Tags["service"])
	require.Len(t, got.AlertOpt.AlertTarget[0].Targets[0].UpgradeTargets, 1)
	require.Equal(t, []string{"notify_yyy"}, got.AlertOpt.AlertTarget[0].Targets[0].UpgradeTargets[0].To)
	require.Equal(t, 300, got.AlertOpt.AlertTarget[0].Targets[0].UpgradeTargets[0].Duration)
	require.Equal(t, []string{"mail"}, got.AlertOpt.AlertTarget[0].Targets[0].UpgradeTargets[0].ToWay)
}

func TestAlertPolicyUpdateBodyPreservesPermissionZeroValues(t *testing.T) {
	got := alertPolicyUpdateBody(&api.AlertPolicy{
		Name:              "codex-status-policy",
		RuleTimezone:      "Asia/Shanghai",
		OpenPermissionSet: false,
		PermissionSet:     nil,
		CheckerUUIDs:      nil,
		SecurityRuleUUIDs: nil,
		AlertOpt: &api.AlertOpt{
			AlertType:                   "status",
			IgnoreOK:                    false,
			SilentTimeout:               0,
			SilentTimeoutByStatusEnable: false,
			AggInterval:                 0,
			AggSendFirst:                false,
			AlertTarget: []api.AlertTarget{{
				Name: "clear schedule",
				Targets: []api.Target{{
					Status: "critical",
				}},
			}},
		},
	})

	require.Equal(t, "codex-status-policy", got["name"])
	require.Equal(t, "", got["desc"])
	require.Equal(t, "Asia/Shanghai", got["ruleTimezone"])
	require.Equal(t, false, got["openPermissionSet"])
	require.Equal(t, []string{}, got["permissionSet"])
	require.Equal(t, []string{}, got["checkerUUIDs"])
	require.Equal(t, []string{}, got["securityRuleUUIDs"])
	alertOpt, ok := got["alertOpt"].(map[string]any)
	require.True(t, ok)
	require.Equal(t, "status", alertOpt["alertType"])
	require.Equal(t, false, alertOpt["ignoreOK"])
	require.Equal(t, 0, alertOpt["silentTimeout"])
	require.Equal(t, false, alertOpt["silentTimeoutByStatusEnable"])
	require.Equal(t, []map[string]any{}, alertOpt["silentTimeoutByStatus"])
	require.Equal(t, []map[string]any{{
		"name": "clear schedule",
		"targets": []map[string]any{{
			"to":             []string{},
			"status":         "critical",
			"upgradeTargets": []map[string]any{},
			"tags":           map[string][]string{},
			"filterString":   "",
		}},
		"crontab":         "",
		"crontabDuration": 0,
		"customDateUUIDs": []string{},
		"customStartTime": "",
		"customDuration":  0,
	}}, alertOpt["alertTarget"])
	require.Equal(t, 0, alertOpt["aggInterval"])
	require.Equal(t, []string{}, alertOpt["aggFields"])
	require.Equal(t, []string{}, alertOpt["aggLabels"])
	require.Equal(t, []string{}, alertOpt["aggClusterFields"])
	require.Equal(t, false, alertOpt["aggSendFirst"])
}

func TestTargetsUpdateBodyHandlesDfSource(t *testing.T) {
	got := targetsUpdateBody([]api.Target{{
		To:       []string{"notify_xxx"},
		Status:   "critical",
		DfSource: "security",
	}, {
		To:     []string{"notify_yyy"},
		Status: "error",
	}})

	require.Equal(t, []map[string]any{{
		"to":             []string{"notify_xxx"},
		"status":         "critical",
		"df_source":      "security",
		"upgradeTargets": []map[string]any{},
		"tags":           map[string][]string{},
		"filterString":   "",
	}, {
		"to":             []string{"notify_yyy"},
		"status":         "error",
		"upgradeTargets": []map[string]any{},
		"tags":           map[string][]string{},
		"filterString":   "",
	}}, got)
}

func TestGetAlertPolicyFromPlanMemberMode(t *testing.T) {
	resource := &alertPolicyResource{}
	plan := &alertPolicyResourceModel{
		Name:         types.StringValue("codex-member-policy"),
		RuleTimezone: types.StringValue("Asia/Shanghai"),
		AlertOpt: &alertOptModel{
			AlertType:     types.StringValue("member"),
			AggInterval:   types.Int64Value(60),
			SilentTimeout: types.Int64Value(300),
			AlertTarget: []alertTarget{{
				Name: types.StringValue("member target"),
				AlertInfo: []alertInfo{{
					Name:       types.StringValue("member route"),
					MemberInfo: []types.String{types.StringValue("acnt_xxx")},
					Targets: []target{{
						To:     []types.String{types.StringValue("notify_xxx")},
						Status: types.StringValue("critical,error,warning"),
					}},
				}},
			}},
		},
	}

	got := resource.getAlertPolicyFromPlan(plan)

	require.Equal(t, "member", got.AlertOpt.AlertType)
	require.Equal(t, 60, got.AlertOpt.AggInterval)
	require.Len(t, got.AlertOpt.AlertTarget, 1)
	require.Empty(t, got.AlertOpt.AlertTarget[0].Targets)
	require.Len(t, got.AlertOpt.AlertTarget[0].AlertInfo, 1)
	require.Equal(t, "member route", got.AlertOpt.AlertTarget[0].AlertInfo[0].Name)
	require.Equal(t, []string{"acnt_xxx"}, got.AlertOpt.AlertTarget[0].AlertInfo[0].MemberInfo)
	require.Len(t, got.AlertOpt.AlertTarget[0].AlertInfo[0].Targets, 1)
	require.Equal(t, []string{"notify_xxx"}, got.AlertOpt.AlertTarget[0].AlertInfo[0].Targets[0].To)
	require.Equal(t, "critical,error,warning", got.AlertOpt.AlertTarget[0].AlertInfo[0].Targets[0].Status)
}

func TestAlertOptFromContentMemberMode(t *testing.T) {
	got := alertOptFromContent(&api.AlertOptContent{
		AlertType:     "member",
		AggInterval:   intPtr(60),
		SilentTimeout: intPtr(300),
		AlertTarget: []api.AlertTargetContent{{
			Name: "member target",
			AlertInfo: []api.AlertInfoContent{{
				Name:       "member route",
				MemberInfo: []string{"acnt_xxx"},
				Targets: []api.TargetContent{{
					To:     []string{"notify_xxx"},
					Status: "critical,error,warning",
				}},
			}},
		}},
	}, nil)

	require.NotNil(t, got)
	require.Equal(t, "member", got.AlertType.ValueString())
	require.Equal(t, int64(60), got.AggInterval.ValueInt64())
	require.Equal(t, int64(300), got.SilentTimeout.ValueInt64())
	require.Len(t, got.AlertTarget, 1)
	require.Equal(t, "member target", got.AlertTarget[0].Name.ValueString())
	require.Len(t, got.AlertTarget[0].AlertInfo, 1)
	require.Equal(t, "member route", got.AlertTarget[0].AlertInfo[0].Name.ValueString())
	require.Equal(t, []types.String{types.StringValue("acnt_xxx")}, got.AlertTarget[0].AlertInfo[0].MemberInfo)
	require.Equal(t, []types.String{types.StringValue("notify_xxx")}, got.AlertTarget[0].AlertInfo[0].Targets[0].To)
}

func TestAlertOptFromContentComplexStatusMode(t *testing.T) {
	got := alertOptFromContent(&api.AlertOptContent{
		AggType:                     "byFields",
		IgnoreOK:                    boolPtr(true),
		AlertType:                   "status",
		SilentTimeout:               intPtr(300),
		SilentTimeoutByStatusEnable: boolPtr(true),
		SilentTimeoutByStatus: []api.SilentTimeoutByStatus{{
			Status:        "critical",
			SilentTimeout: 120,
		}},
		AlertTarget: []api.AlertTargetContent{{
			Name:            "status route",
			CustomDateUUIDs: []string{"ndate_xxx"},
			Targets: []api.TargetContent{{
				To:           []string{"notify_xxx"},
				Status:       "critical,error",
				DfSource:     "security",
				FilterString: "host:codex-alert-test",
				Tags: map[string][]string{
					"service": {"codex"},
				},
				UpgradeTargets: []api.UpgradeTargetContent{{
					To:       []string{"notify_yyy"},
					Duration: intPtr(300),
					ToWay:    []string{"mail"},
				}},
			}},
			Crontab:         "0 9 * * 1",
			CrontabDuration: intPtr(3600),
			CustomStartTime: "09:30:00",
			CustomDuration:  intPtr(1800),
		}},
		AggInterval:      intPtr(60),
		AggFields:        []string{"df_monitor_checker_id", "df_label"},
		AggLabels:        []string{"service"},
		AggClusterFields: []string{"df_title"},
		AggSendFirst:     boolPtr(true),
	}, nil)

	require.Equal(t, "byFields", got.AggType.ValueString())
	require.True(t, got.IgnoreOK.ValueBool())
	require.Equal(t, "status", got.AlertType.ValueString())
	require.Equal(t, int64(300), got.SilentTimeout.ValueInt64())
	require.True(t, got.SilentTimeoutByStatusEnable.ValueBool())
	require.Len(t, got.SilentTimeoutByStatus, 1)
	require.Equal(t, "critical", got.SilentTimeoutByStatus[0].Status.ValueString())
	require.Equal(t, int64(120), got.SilentTimeoutByStatus[0].SilentTimeout.ValueInt64())
	require.Equal(t, int64(60), got.AggInterval.ValueInt64())
	require.Equal(t, []types.String{types.StringValue("df_monitor_checker_id"), types.StringValue("df_label")}, got.AggFields)
	require.Equal(t, []types.String{types.StringValue("service")}, got.AggLabels)
	require.Equal(t, []types.String{types.StringValue("df_title")}, got.AggClusterFields)
	require.True(t, got.AggSendFirst.ValueBool())
	require.Len(t, got.AlertTarget, 1)
	require.Equal(t, "status route", got.AlertTarget[0].Name.ValueString())
	require.Equal(t, []types.String{types.StringValue("ndate_xxx")}, got.AlertTarget[0].CustomDateUUIDs)
	require.Equal(t, "0 9 * * 1", got.AlertTarget[0].Crontab.ValueString())
	require.Equal(t, int64(3600), got.AlertTarget[0].CrontabDuration.ValueInt64())
	require.Equal(t, "09:30:00", got.AlertTarget[0].CustomStartTime.ValueString())
	require.Equal(t, int64(1800), got.AlertTarget[0].CustomDuration.ValueInt64())
	require.Equal(t, []types.String{types.StringValue("notify_xxx")}, got.AlertTarget[0].Targets[0].To)
	require.Equal(t, "security", got.AlertTarget[0].Targets[0].DfSource.ValueString())
	require.Equal(t, "host:codex-alert-test", got.AlertTarget[0].Targets[0].FilterString.ValueString())
	require.Equal(t, []string{"codex"}, got.AlertTarget[0].Targets[0].Tags["service"])
	require.Equal(t, []types.String{types.StringValue("notify_yyy")}, got.AlertTarget[0].Targets[0].UpgradeTargets[0].To)
	require.Equal(t, int64(300), got.AlertTarget[0].Targets[0].UpgradeTargets[0].Duration.ValueInt64())
	require.Equal(t, []types.String{types.StringValue("mail")}, got.AlertTarget[0].Targets[0].UpgradeTargets[0].ToWay)
}

func TestAlertOptFromContentPreservesPriorZeroValuesForResourceRead(t *testing.T) {
	prior := &alertOptModel{
		IgnoreOK:                    types.BoolValue(true),
		SilentTimeout:               types.Int64Value(300),
		SilentTimeoutByStatusEnable: types.BoolValue(true),
		AggInterval:                 types.Int64Value(60),
		AggSendFirst:                types.BoolValue(true),
	}

	got := alertOptFromContent(&api.AlertOptContent{}, prior)

	require.True(t, got.IgnoreOK.ValueBool())
	require.Equal(t, int64(300), got.SilentTimeout.ValueInt64())
	require.True(t, got.SilentTimeoutByStatusEnable.ValueBool())
	require.Equal(t, int64(60), got.AggInterval.ValueInt64())
	require.True(t, got.AggSendFirst.ValueBool())
}

func TestAlertOptFromContentAppliesRemoteZeroValuesForResourceRead(t *testing.T) {
	prior := &alertOptModel{
		IgnoreOK:                    types.BoolValue(true),
		SilentTimeout:               types.Int64Value(300),
		SilentTimeoutByStatusEnable: types.BoolValue(true),
		AggInterval:                 types.Int64Value(60),
		AggSendFirst:                types.BoolValue(true),
		AggFields:                   []types.String{types.StringValue("df_monitor_checker_id")},
		AggLabels:                   []types.String{types.StringValue("service")},
		AggClusterFields:            []types.String{types.StringValue("df_title")},
		SilentTimeoutByStatus: []silentTimeoutByStatus{{
			Status:        types.StringValue("critical"),
			SilentTimeout: types.Int64Value(120),
		}},
	}

	got := alertOptFromContent(&api.AlertOptContent{
		IgnoreOK:                    boolPtr(false),
		SilentTimeout:               intPtr(0),
		SilentTimeoutByStatusEnable: boolPtr(false),
		SilentTimeoutByStatus:       []api.SilentTimeoutByStatus{},
		AggInterval:                 intPtr(0),
		AggFields:                   []string{},
		AggLabels:                   []string{},
		AggClusterFields:            []string{},
		AggSendFirst:                boolPtr(false),
		AlertTarget: []api.AlertTargetContent{{
			Name:            "zero durations",
			CrontabDuration: intPtr(0),
			CustomDuration:  intPtr(0),
			Targets: []api.TargetContent{{
				Status: "critical",
				UpgradeTargets: []api.UpgradeTargetContent{{
					Duration: intPtr(0),
				}},
			}},
		}},
	}, prior)

	require.False(t, got.IgnoreOK.IsNull())
	require.False(t, got.IgnoreOK.ValueBool())
	require.False(t, got.SilentTimeout.IsNull())
	require.Equal(t, int64(0), got.SilentTimeout.ValueInt64())
	require.False(t, got.SilentTimeoutByStatusEnable.IsNull())
	require.False(t, got.SilentTimeoutByStatusEnable.ValueBool())
	require.False(t, got.AggInterval.IsNull())
	require.Equal(t, int64(0), got.AggInterval.ValueInt64())
	require.False(t, got.AggSendFirst.IsNull())
	require.False(t, got.AggSendFirst.ValueBool())
	require.Empty(t, got.AggFields)
	require.Empty(t, got.AggLabels)
	require.Empty(t, got.AggClusterFields)
	require.Empty(t, got.SilentTimeoutByStatus)
	require.Len(t, got.AlertTarget, 1)
	require.False(t, got.AlertTarget[0].CrontabDuration.IsNull())
	require.Equal(t, int64(0), got.AlertTarget[0].CrontabDuration.ValueInt64())
	require.False(t, got.AlertTarget[0].CustomDuration.IsNull())
	require.Equal(t, int64(0), got.AlertTarget[0].CustomDuration.ValueInt64())
	require.Len(t, got.AlertTarget[0].Targets, 1)
	require.Len(t, got.AlertTarget[0].Targets[0].UpgradeTargets, 1)
	require.False(t, got.AlertTarget[0].Targets[0].UpgradeTargets[0].Duration.IsNull())
	require.Equal(t, int64(0), got.AlertTarget[0].Targets[0].UpgradeTargets[0].Duration.ValueInt64())
}

func TestAlertOptFromContentPreservesPriorNestedEmptyValuesForResourceRead(t *testing.T) {
	prior := &alertOptModel{
		AlertTarget: []alertTarget{{
			Name:            types.StringValue("default"),
			CrontabDuration: types.Int64Value(0),
			CustomDateUUIDs: []types.String{},
			CustomDuration:  types.Int64Value(0),
			Targets: []target{{
				To:             []types.String{types.StringValue("notify_xxx")},
				Status:         types.StringValue("critical,error,warning"),
				Tags:           map[string][]string{},
				UpgradeTargets: []upgradeTarget{},
			}},
		}},
	}

	got := alertOptFromContent(&api.AlertOptContent{
		AlertTarget: []api.AlertTargetContent{{
			Name: "default",
			Targets: []api.TargetContent{{
				To:     []string{"notify_xxx"},
				Status: "critical,error,warning",
			}},
		}},
	}, prior)

	require.Len(t, got.AlertTarget, 1)
	require.False(t, got.AlertTarget[0].CrontabDuration.IsNull())
	require.Equal(t, int64(0), got.AlertTarget[0].CrontabDuration.ValueInt64())
	require.Empty(t, got.AlertTarget[0].CustomDateUUIDs)
	require.False(t, got.AlertTarget[0].CustomDuration.IsNull())
	require.Equal(t, int64(0), got.AlertTarget[0].CustomDuration.ValueInt64())
	require.Len(t, got.AlertTarget[0].Targets, 1)
	require.Empty(t, got.AlertTarget[0].Targets[0].Tags)
	require.Empty(t, got.AlertTarget[0].Targets[0].UpgradeTargets)
}

func TestAlertOptFromContentForDataSourceIncludesZeroValues(t *testing.T) {
	got := alertOptFromContentForDataSource(&api.AlertOptContent{})

	require.False(t, got.IgnoreOK.IsNull())
	require.False(t, got.IgnoreOK.ValueBool())
	require.False(t, got.SilentTimeout.IsNull())
	require.Equal(t, int64(0), got.SilentTimeout.ValueInt64())
	require.False(t, got.SilentTimeoutByStatusEnable.IsNull())
	require.False(t, got.SilentTimeoutByStatusEnable.ValueBool())
	require.False(t, got.AggInterval.IsNull())
	require.Equal(t, int64(0), got.AggInterval.ValueInt64())
	require.False(t, got.AggSendFirst.IsNull())
	require.False(t, got.AggSendFirst.ValueBool())
}

func boolPtr(value bool) *bool {
	return &value
}

func intPtr(value int) *int {
	return &value
}
