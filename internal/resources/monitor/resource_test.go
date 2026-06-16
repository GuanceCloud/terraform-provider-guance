package monitor

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/require"
)

func TestGetMonitorFromPlanParsesExtendJSON(t *testing.T) {
	r := &monitorResource{}
	plan := &monitorResourceModel{
		Type:   types.StringValue("trigger"),
		Extend: types.StringValue(`{"isNeedCreateIssue":false,"issueLevelUUID":"","needRecoverIssue":false}`),
		JsonScript: JsonScript{
			Type: types.StringValue("simpleCheck"),
		},
	}

	got, err := r.getMonitorFromPlan(plan)

	require.NoError(t, err)
	extend, ok := got.Extend.(map[string]any)
	require.True(t, ok)
	require.Equal(t, false, extend["isNeedCreateIssue"])
	require.Equal(t, "", extend["issueLevelUUID"])
	require.Equal(t, false, extend["needRecoverIssue"])
}

func TestGetMonitorFromPlanRejectsInvalidExtendJSON(t *testing.T) {
	r := &monitorResource{}
	plan := &monitorResourceModel{
		Type:   types.StringValue("trigger"),
		Extend: types.StringValue(`{"isNeedCreateIssue":`),
		JsonScript: JsonScript{
			Type: types.StringValue("simpleCheck"),
		},
	}

	got, err := r.getMonitorFromPlan(plan)

	require.Nil(t, got)
	require.Error(t, err)
}

func TestOptionalStringFromContentClearsEmptyRemoteValues(t *testing.T) {
	require.True(t, optionalStringFromContent("").IsNull())
	require.Equal(t, "dashboard_xxx", optionalStringFromContent("dashboard_xxx").ValueString())
}

func TestApplyExtendFromContentPreservesBackendExpandedValue(t *testing.T) {
	state := &monitorResourceModel{
		Extend: types.StringValue(`{"isNeedCreateIssue":false}`),
	}

	err := applyExtendFromContent(state, map[string]any{
		"isNeedCreateIssue": false,
		"querylist": []any{
			map[string]any{"qtype": "dql"},
		},
	})

	require.NoError(t, err)
	require.Equal(t, `{"isNeedCreateIssue":false}`, state.Extend.ValueString())
}

func TestApplyExtendFromContentOverwritesConfiguredDrift(t *testing.T) {
	state := &monitorResourceModel{
		Extend: types.StringValue(`{"isNeedCreateIssue":false}`),
	}

	err := applyExtendFromContent(state, map[string]any{
		"isNeedCreateIssue": true,
		"querylist": []any{
			map[string]any{"qtype": "dql"},
		},
	})

	require.NoError(t, err)
	require.Equal(t, `{"isNeedCreateIssue":true,"querylist":[{"qtype":"dql"}]}`, state.Extend.ValueString())
}
