package notify_object

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/require"

	"github.com/GuanceCloud/terraform-provider-guance/internal/api"
)

func TestGetNotifyObjectFromPlan(t *testing.T) {
	resource := &notifyObjectResource{}
	plan := &notifyObjectResourceModel{
		Type:              types.StringValue("simpleHTTPRequest"),
		Name:              types.StringValue("codex-notify-object"),
		OptSet:            types.StringValue(`{ "url" : "https://example.com/hook", "headersConfig" : { "items" : [], "isOpen" : false } }`),
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
	require.Equal(t, `{"headersConfig":{"isOpen":false,"items":[]},"url":"https://example.com/hook"}`, plan.OptSet.ValueString())
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

func TestStateFromNotifyObjectContentCanonicalizesOptSet(t *testing.T) {
	state := &notifyObjectDataSourceModel{}
	content := &api.NotifyObjectContent{
		UUID: "notify_xxx",
		Type: "simpleHTTPRequest",
		Name: "codex-notify-object",
		OptSet: map[string]any{
			"url": "https://example.com/hook",
			"headersConfig": map[string]any{
				"items":  []any{},
				"isOpen": false,
			},
		},
	}

	err := stateFromNotifyObjectContent(state, content)

	require.NoError(t, err)
	require.Equal(t, `{"headersConfig":{"isOpen":false,"items":[]},"url":"https://example.com/hook"}`, state.OptSet.ValueString())
}
