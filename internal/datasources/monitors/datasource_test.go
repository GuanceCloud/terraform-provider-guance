package monitors

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/require"

	"github.com/GuanceCloud/terraform-provider-guance/internal/api"
)

func TestMonitorFromContent(t *testing.T) {
	got := monitorFromContent(api.MonitorContent{
		UUID:             "rul_xxx",
		Type:             "trigger",
		Status:           2,
		AlertPolicyUUIDs: []string{"altpl_xxx"},
		DashboardUUID:    "dsbd_xxx",
		Tags:             []string{"terraform"},
		WorkspaceUUID:    "wksp_xxx",
		MonitorUUID:      "monitor_xxx",
		JsonScript: map[string]any{
			"title": "Monitor From JSON",
		},
	})

	require.Equal(t, "rul_xxx", got.UUID.ValueString())
	require.Equal(t, "Monitor From JSON", got.Name.ValueString())
	require.Equal(t, "trigger", got.Type.ValueString())
	require.Equal(t, int64(2), got.Status.ValueInt64())
	require.Equal(t, []types.String{types.StringValue("altpl_xxx")}, got.AlertPolicyUUIDs)
	require.Equal(t, "dsbd_xxx", got.DashboardUUID.ValueString())
	require.Equal(t, []types.String{types.StringValue("terraform")}, got.Tags)
	require.Equal(t, "wksp_xxx", got.WorkspaceUUID.ValueString())
	require.Equal(t, "monitor_xxx", got.MonitorUUID.ValueString())
}
