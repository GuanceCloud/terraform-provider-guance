package pipeline

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/require"

	"github.com/GuanceCloud/terraform-provider-guance/internal/api"
)

func stringPtr(s string) *string {
	return &s
}

func TestGetPipelineFromPlanSendsEmptyLoggingIndexOnly(t *testing.T) {
	r := &pipelineResource{}
	plan := &pipelineResourceModel{
		Name:   types.StringValue("demo"),
		Extend: &pipelineExtend{LoggingIndex: types.StringValue("")},
	}

	item := r.getPipelineFromPlan(plan)
	body, err := json.Marshal(item.Extend)

	require.NoError(t, err)
	require.JSONEq(t, `{"loggingIndex":""}`, string(body))
}

func TestGetPipelineFromPlanOmitsUnsetLoggingIndex(t *testing.T) {
	r := &pipelineResource{}
	plan := &pipelineResourceModel{
		Name: types.StringValue("demo"),
		Extend: &pipelineExtend{
			AppID: []types.String{types.StringValue("app_xxx")},
		},
	}

	item := r.getPipelineFromPlan(plan)
	body, err := json.Marshal(item.Extend)

	require.NoError(t, err)
	require.JSONEq(t, `{"appID":["app_xxx"]}`, string(body))
}

func TestExtendFromContentKeepsEmptyLoggingIndexAfterImport(t *testing.T) {
	got := extendFromContent(nil, &api.PipelineExtend{LoggingIndex: stringPtr("")})

	require.NotNil(t, got)
	require.False(t, got.LoggingIndex.IsNull())
	require.Equal(t, "", got.LoggingIndex.ValueString())
	require.Nil(t, got.AppID)
	require.Nil(t, got.Measurement)
}

func TestExtendFromContentClearsStateWhenRemoteExtendIsEmpty(t *testing.T) {
	prior := &pipelineExtend{LoggingIndex: types.StringValue("")}

	require.Nil(t, extendFromContent(prior, &api.PipelineExtend{}))
	require.Nil(t, extendFromContent(prior, nil))
}

func TestExtendFromContentKeepsEmptyConfiguredBlock(t *testing.T) {
	prior := &pipelineExtend{}

	require.Equal(t, prior, extendFromContent(prior, &api.PipelineExtend{}))
}

func TestExtendFromContentKeepsUnsetLoggingIndexNull(t *testing.T) {
	prior := &pipelineExtend{
		AppID: []types.String{types.StringValue("app_xxx")},
	}

	got := extendFromContent(prior, &api.PipelineExtend{
		AppID:        []string{"app_xxx"},
		LoggingIndex: stringPtr(""),
	})

	require.NotNil(t, got)
	require.True(t, got.LoggingIndex.IsNull())
	require.Equal(t, []types.String{types.StringValue("app_xxx")}, got.AppID)
}

func TestExtendFromContentDetectsRemoteLoggingIndexChange(t *testing.T) {
	prior := &pipelineExtend{LoggingIndex: types.StringValue("")}

	got := extendFromContent(prior, &api.PipelineExtend{LoggingIndex: stringPtr("default")})

	require.Equal(t, "default", got.LoggingIndex.ValueString())
}

func TestMergeExtendKeepsUnmanagedKeys(t *testing.T) {
	live := &api.PipelineExtend{
		LoggingIndex: stringPtr("old"),
		Extra: map[string]json.RawMessage{
			"ifAiCreate": json.RawMessage(`true`),
		},
	}

	got := mergeExtend(live, &api.PipelineExtend{LoggingIndex: stringPtr("default")}, true)
	body, err := json.Marshal(got)

	require.NoError(t, err)
	require.JSONEq(t, `{"ifAiCreate":true,"loggingIndex":"default"}`, string(body))

	got = mergeExtend(live, nil, false)
	body, err = json.Marshal(got)

	require.NoError(t, err)
	require.JSONEq(t, `{"ifAiCreate":true}`, string(body))
}

func TestMergeExtendKeepsLiveLoggingIndexOnlyWhenAsked(t *testing.T) {
	live := &api.PipelineExtend{LoggingIndex: stringPtr("")}

	body, err := json.Marshal(mergeExtend(live, &api.PipelineExtend{}, true))
	require.NoError(t, err)
	require.JSONEq(t, `{"loggingIndex":""}`, string(body))

	body, err = json.Marshal(mergeExtend(live, &api.PipelineExtend{}, false))
	require.NoError(t, err)
	require.JSONEq(t, `{}`, string(body))

	// a planned loggingIndex always wins
	body, err = json.Marshal(mergeExtend(live, &api.PipelineExtend{LoggingIndex: stringPtr("default")}, true))
	require.NoError(t, err)
	require.JSONEq(t, `{"loggingIndex":"default"}`, string(body))
}

func TestGetPipelineForUpdateKeepsLiveExtendKeys(t *testing.T) {
	var (
		mu          sync.Mutex
		gotRequests []string
	)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// record only; the handler runs outside the test goroutine
		mu.Lock()
		gotRequests = append(gotRequests, r.Method+" "+r.URL.Path)
		mu.Unlock()

		_ = json.NewEncoder(w).Encode(map[string]any{
			"code":    200,
			"success": true,
			"content": map[string]any{
				"uuid": "pl_xxx",
				"extend": map[string]any{
					"loggingIndex":      "",
					"ifAiCreate":        true,
					"exampleType":       "json",
					"aiTargetStructure": map[string]any{"fields": []string{"a"}},
				},
			},
		})
	}))
	defer server.Close()

	r := &pipelineResource{client: &api.Client{
		EndPoint:   server.URL,
		HTTPClient: server.Client(),
	}}
	plan := &pipelineResourceModel{
		UUID:    types.StringValue("pl_xxx"),
		Name:    types.StringValue("demo"),
		Content: types.StringValue("add_key(a, 1)"),
		Extend:  &pipelineExtend{LoggingIndex: types.StringValue("")},
	}

	item, err := r.getPipelineForUpdate(plan, &pipelineResourceModel{Extend: plan.Extend})
	require.NoError(t, err)

	mu.Lock()
	require.Equal(t, []string{"GET /api/v1/pipeline/pl_xxx/get"}, gotRequests)
	mu.Unlock()

	body, err := json.Marshal(item)
	require.NoError(t, err)

	var sent map[string]json.RawMessage
	require.NoError(t, json.Unmarshal(body, &sent))
	require.JSONEq(t,
		`{"loggingIndex":"","ifAiCreate":true,"exampleType":"json","aiTargetStructure":{"fields":["a"]}}`,
		string(sent["extend"]),
	)
}

func TestGetPipelineForUpdateFailsWhenLiveReadFails(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	r := &pipelineResource{client: &api.Client{
		EndPoint:   server.URL,
		HTTPClient: server.Client(),
	}}

	item, err := r.getPipelineForUpdate(&pipelineResourceModel{UUID: types.StringValue("pl_xxx")}, nil)

	require.Nil(t, item)
	require.Error(t, err)
}
