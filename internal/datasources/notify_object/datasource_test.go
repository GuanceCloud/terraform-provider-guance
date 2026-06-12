package notify_object

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/GuanceCloud/terraform-provider-guance/internal/api"
)

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
