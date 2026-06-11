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
		Name:         types.StringValue("codex-status-policy"),
		Desc:         types.StringValue("status mode"),
		RuleTimezone: types.StringValue("Asia/Shanghai"),
		AlertOpt: &alertOptModel{
			AlertType:     types.StringValue("status"),
			AggInterval:   types.Int64Value(60),
			AggFields:     []types.String{types.StringValue("df_monitor_checker_id")},
			SilentTimeout: types.Int64Value(300),
			AlertTarget: []alertTarget{{
				Name:            types.StringValue("default"),
				CustomDateUUIDs: []types.String{types.StringValue("ndate_xxx")},
				CustomStartTime: types.StringValue("09:30:00"),
				CustomDuration:  types.Int64Value(3600),
				Targets: []target{{
					To:     []types.String{types.StringValue("notify_xxx")},
					Status: types.StringValue("critical,error,warning"),
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
	require.Equal(t, "Asia/Shanghai", got.RuleTimezone)
	require.NotNil(t, got.AlertOpt)
	require.Equal(t, "status", got.AlertOpt.AlertType)
	require.Equal(t, 60, got.AlertOpt.AggInterval)
	require.Equal(t, []string{"df_monitor_checker_id"}, got.AlertOpt.AggFields)
	require.Equal(t, 300, got.AlertOpt.SilentTimeout)
	require.Len(t, got.AlertOpt.AlertTarget, 1)
	require.Equal(t, "default", got.AlertOpt.AlertTarget[0].Name)
	require.Equal(t, []string{"ndate_xxx"}, got.AlertOpt.AlertTarget[0].CustomDateUUIDs)
	require.Equal(t, "09:30:00", got.AlertOpt.AlertTarget[0].CustomStartTime)
	require.Equal(t, 3600, got.AlertOpt.AlertTarget[0].CustomDuration)
	require.Len(t, got.AlertOpt.AlertTarget[0].Targets, 1)
	require.Equal(t, []string{"notify_xxx"}, got.AlertOpt.AlertTarget[0].Targets[0].To)
	require.Equal(t, "critical,error,warning", got.AlertOpt.AlertTarget[0].Targets[0].Status)
	require.Len(t, got.AlertOpt.AlertTarget[0].Targets[0].UpgradeTargets, 1)
	require.Equal(t, []string{"notify_yyy"}, got.AlertOpt.AlertTarget[0].Targets[0].UpgradeTargets[0].To)
	require.Equal(t, 300, got.AlertOpt.AlertTarget[0].Targets[0].UpgradeTargets[0].Duration)
	require.Equal(t, []string{"mail"}, got.AlertOpt.AlertTarget[0].Targets[0].UpgradeTargets[0].ToWay)
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
	got := alertOptFromContent(&api.AlertOpt{
		AlertType:     "member",
		AggInterval:   60,
		SilentTimeout: 300,
		AlertTarget: []api.AlertTarget{{
			Name: "member target",
			AlertInfo: []api.AlertInfo{{
				Name:       "member route",
				MemberInfo: []string{"acnt_xxx"},
				Targets: []api.Target{{
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
