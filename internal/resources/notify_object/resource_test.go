package notify_object

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/require"
)

func TestGetNotifyObjectFromPlan(t *testing.T) {
	resource := &notifyObjectResource{}
	plan := &notifyObjectResourceModel{
		Type:              types.StringValue("simpleHTTPRequest"),
		Name:              types.StringValue("codex-notify-object"),
		OptSet:            types.StringValue(`{"url":"https://example.com/hook","headersConfig":{"isOpen":false,"items":[]}}`),
		OpenPermissionSet: types.BoolValue(true),
		PermissionSet: []types.String{
			types.StringValue("wsAdmin"),
			types.StringValue("acnt_xxx"),
		},
	}

	got, err := resource.getNotifyObjectFromPlan(plan)

	require.NoError(t, err)
	require.Equal(t, "simpleHTTPRequest", got.Type)
	require.Equal(t, "codex-notify-object", got.Name)
	require.True(t, got.OpenPermissionSet)
	require.Equal(t, []string{"wsAdmin", "acnt_xxx"}, got.PermissionSet)

	optSet, ok := got.OptSet.(map[string]any)
	require.True(t, ok)
	require.Equal(t, "https://example.com/hook", optSet["url"])
	headersConfig, ok := optSet["headersConfig"].(map[string]any)
	require.True(t, ok)
	require.Equal(t, false, headersConfig["isOpen"])
	require.Empty(t, headersConfig["items"])
}

func TestGetNotifyObjectFromPlanRejectsInvalidOptSetJSON(t *testing.T) {
	resource := &notifyObjectResource{}
	plan := &notifyObjectResourceModel{
		Type:   types.StringValue("simpleHTTPRequest"),
		Name:   types.StringValue("codex-notify-object"),
		OptSet: types.StringValue(`{"url":`),
	}

	got, err := resource.getNotifyObjectFromPlan(plan)

	require.Nil(t, got)
	require.Error(t, err)
}
