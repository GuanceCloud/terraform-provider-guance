package monitor

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/require"

	"github.com/GuanceCloud/terraform-provider-guance/internal/api"
)

func TestStateFromMonitorContentUsesMonitorNameAndCanonicalJSON(t *testing.T) {
	state := &monitorDataSourceModel{}
	content := &api.MonitorContent{
		UUID:              "rul_xxx",
		Type:              "trigger",
		Status:            0,
		AlertPolicyUUIDs:  []string{"altpl_xxx"},
		DashboardUUID:     "dsbd_xxx",
		Tags:              []string{"terraform", "alert"},
		Secret:            "secret_xxx",
		OpenPermissionSet: true,
		PermissionSet:     []string{"wsAdmin"},
		CreateAt:          1710000000,
		UpdateAt:          1710000300,
		WorkspaceUUID:     "wksp_xxx",
		MonitorUUID:       "monitor_xxx",
		MonitorName:       "Terraform Monitor",
		Extend: map[string]any{
			"isNeedCreateIssue": false,
		},
		JsonScript: map[string]any{
			"title": "JSON Script Title",
			"type":  "simpleCheck",
		},
	}

	require.NoError(t, stateFromMonitorContent(state, content))

	require.Equal(t, "rul_xxx", state.UUID.ValueString())
	require.Equal(t, "Terraform Monitor", state.Name.ValueString())
	require.Equal(t, "trigger", state.Type.ValueString())
	require.Equal(t, int64(0), state.Status.ValueInt64())
	require.Equal(t, []types.String{types.StringValue("altpl_xxx")}, state.AlertPolicyUUIDs)
	require.Equal(t, "dsbd_xxx", state.DashboardUUID.ValueString())
	require.Equal(t, []types.String{types.StringValue("terraform"), types.StringValue("alert")}, state.Tags)
	require.Equal(t, "secret_xxx", state.Secret.ValueString())
	require.True(t, state.OpenPermissionSet.ValueBool())
	require.Equal(t, []types.String{types.StringValue("wsAdmin")}, state.PermissionSet)
	require.Equal(t, `{"isNeedCreateIssue":false}`, state.Extend.ValueString())
	require.Equal(t, `{"title":"JSON Script Title","type":"simpleCheck"}`, state.JsonScript.ValueString())
}

func TestMonitorNameFallsBackToJSONScript(t *testing.T) {
	require.Equal(t, "Title Name", monitorName(api.MonitorContent{
		JsonScript: map[string]any{
			"title": "Title Name",
			"name":  "Legacy Name",
		},
	}))
	require.Equal(t, "Legacy Name", monitorName(api.MonitorContent{
		JsonScript: map[string]any{
			"name": "Legacy Name",
		},
	}))
}
