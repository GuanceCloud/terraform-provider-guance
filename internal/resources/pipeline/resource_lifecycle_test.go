package pipeline

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/stretchr/testify/require"

	"github.com/GuanceCloud/terraform-provider-guance/internal/api"
)

// mockPipelineAPI serves GET /pipeline/<uuid>/get with a settable extend and
// records the body of every POST /pipeline/<uuid>/modify. The handler only
// records; all assertions run in the test goroutine.
type mockPipelineAPI struct {
	mu         sync.Mutex
	liveExtend map[string]any
	modify     [][]byte
	unexpected []string
}

func (m *mockPipelineAPI) setLiveExtend(extend map[string]any) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.liveExtend = extend
}

func (m *mockPipelineAPI) modifyBodies() [][]byte {
	m.mu.Lock()
	defer m.mu.Unlock()
	return append([][]byte(nil), m.modify...)
}

func (m *mockPipelineAPI) unexpectedRequests() []string {
	m.mu.Lock()
	defer m.mu.Unlock()
	return append([]string(nil), m.unexpected...)
}

func (m *mockPipelineAPI) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	m.mu.Lock()
	defer m.mu.Unlock()

	content := map[string]any{
		"uuid":          "pl_xxx",
		"workspaceUUID": "wksp_xxx",
		"name":          "demo",
		"category":      "logging",
		"type":          "local",
		"source":        []string{"nginx"},
		"content":       base64.StdEncoding.EncodeToString([]byte("add_key(a, 1)")),
		"testData":      base64.StdEncoding.EncodeToString([]byte("x")),
	}
	if m.liveExtend != nil {
		content["extend"] = m.liveExtend
	}

	switch {
	case r.Method == http.MethodGet && r.URL.Path == "/api/v1/pipeline/pl_xxx/get":
	case r.Method == http.MethodPost && r.URL.Path == "/api/v1/pipeline/pl_xxx/modify":
		body, _ := io.ReadAll(r.Body)
		m.modify = append(m.modify, body)
	default:
		m.unexpected = append(m.unexpected, r.Method+" "+r.URL.Path)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	_ = json.NewEncoder(w).Encode(map[string]any{"code": 200, "success": true, "content": content})
}

func newMockedResource(t *testing.T) (*pipelineResource, *mockPipelineAPI) {
	t.Helper()
	mock := &mockPipelineAPI{}
	server := httptest.NewServer(mock)
	t.Cleanup(server.Close)
	t.Cleanup(func() { require.Empty(t, mock.unexpectedRequests()) })

	return &pipelineResource{client: &api.Client{
		EndPoint:   server.URL,
		HTTPClient: server.Client(),
	}}, mock
}

// configModel is the pipeline as written in configuration, with the given
// extend and pipeline script.
func configModel(extend *pipelineExtend, script string) *pipelineResourceModel {
	return &pipelineResourceModel{
		UUID:          types.StringValue("pl_xxx"),
		WorkspaceUUID: types.StringValue("wksp_xxx"),
		CreateAt:      types.StringValue("0"),
		UpdateAt:      types.StringValue("0"),
		Name:          types.StringValue("demo"),
		Source:        []types.String{types.StringValue("nginx")},
		Content:       types.StringValue(script),
		TestData:      types.StringValue("x"),
		IsForce:       types.BoolValue(false),
		IsDisabled:    types.BoolValue(false),
		Category:      types.StringValue("logging"),
		AsDefault:     types.Int64Null(),
		Status:        types.Int64Value(0),
		Type:          types.StringValue("local"),
		DataType:      types.StringValue(""),
		Extend:        extend,
	}
}

func newState(t *testing.T, m *pipelineResourceModel) tfsdk.State {
	t.Helper()
	ctx := context.Background()
	st := tfsdk.State{Schema: resourceSchema, Raw: tftypes.NewValue(resourceSchema.Type().TerraformType(ctx), nil)}
	require.False(t, st.Set(ctx, m).HasError())
	return st
}

func newPlan(t *testing.T, m *pipelineResourceModel) tfsdk.Plan {
	t.Helper()
	ctx := context.Background()
	pl := tfsdk.Plan{Schema: resourceSchema, Raw: tftypes.NewValue(resourceSchema.Type().TerraformType(ctx), nil)}
	require.False(t, pl.Set(ctx, m).HasError())
	return pl
}

// refresh runs the resource Read against the mock and returns the new state.
func refresh(t *testing.T, r *pipelineResource, prior tfsdk.State) tfsdk.State {
	t.Helper()
	resp := &resource.ReadResponse{State: tfsdk.State{Schema: prior.Schema, Raw: prior.Raw.Copy()}}
	r.Read(context.Background(), resource.ReadRequest{State: prior}, resp)
	require.False(t, resp.Diagnostics.HasError(), "%v", resp.Diagnostics)
	return resp.State
}

// apply runs the resource Update against the mock and returns the new state.
func apply(t *testing.T, r *pipelineResource, prior tfsdk.State, plan tfsdk.Plan) tfsdk.State {
	t.Helper()
	resp := &resource.UpdateResponse{State: tfsdk.State{Schema: prior.Schema, Raw: prior.Raw.Copy()}}
	r.Update(context.Background(), resource.UpdateRequest{State: prior, Plan: plan}, resp)
	require.False(t, resp.Diagnostics.HasError(), "%v", resp.Diagnostics)
	return resp.State
}

func extendAttr(t *testing.T, st tfsdk.State) types.Object {
	t.Helper()
	var extend types.Object
	require.False(t, st.GetAttribute(context.Background(), path.Root("extend"), &extend).HasError())
	return extend
}

func sentExtend(t *testing.T, body []byte) string {
	t.Helper()
	var sent map[string]json.RawMessage
	require.NoError(t, json.Unmarshal(body, &sent))
	return string(sent["extend"])
}

// Review P2: a script-only update must not drop a live "loggingIndex": ""
// that the plan does not show as removed.
func TestScriptOnlyUpdateKeepsLiveEmptyLoggingIndex(t *testing.T) {
	r, mock := newMockedResource(t)

	// 1. applied with extend = {}; the API first returns an empty extend
	mock.setLiveExtend(map[string]any{})
	state := refresh(t, r, newState(t, configModel(&pipelineExtend{}, "add_key(a, 1)")))

	// 2. a console edit writes loggingIndex "" and an unmanaged key
	mock.setLiveExtend(map[string]any{"loggingIndex": "", "ifAiCreate": true})
	state = refresh(t, r, state)

	// 3. only the script changes; extend in the plan equals extend in state
	plannedExtend := extendAttr(t, newState(t, configModel(&pipelineExtend{}, "add_key(b, 2)")))
	require.True(t, extendAttr(t, state).Equal(plannedExtend), "plan shows no change to extend")

	// 4. apply: the modify body must keep loggingIndex
	apply(t, r, state, newPlan(t, configModel(&pipelineExtend{}, "add_key(b, 2)")))

	bodies := mock.modifyBodies()
	require.Len(t, bodies, 1)
	require.JSONEq(t, `{"loggingIndex":"","ifAiCreate":true}`, sentExtend(t, bodies[0]))
}

// A logging_index that the plan removes is still removed.
func TestPlannedLoggingIndexRemovalIsSent(t *testing.T) {
	r, mock := newMockedResource(t)

	mock.setLiveExtend(map[string]any{"loggingIndex": "default", "ifAiCreate": true})
	state := refresh(t, r, newState(t, configModel(
		&pipelineExtend{LoggingIndex: types.StringValue("default")}, "add_key(a, 1)")))

	apply(t, r, state, newPlan(t, configModel(&pipelineExtend{}, "add_key(a, 1)")))

	bodies := mock.modifyBodies()
	require.Len(t, bodies, 1)
	require.JSONEq(t, `{"ifAiCreate":true}`, sentExtend(t, bodies[0]))
}

// Review P2: configured empty lists must survive refresh when the API leaves
// the list fields out, or every plan proposes adding them again.
func TestConfiguredEmptyListsSurviveRefreshAndApply(t *testing.T) {
	for name, live := range map[string]map[string]any{
		"default loggingIndex": {"loggingIndex": ""},
		"empty extend":         {},
	} {
		t.Run(name, func(t *testing.T) {
			r, mock := newMockedResource(t)
			mock.setLiveExtend(live)

			config := configModel(&pipelineExtend{
				AppID:       []types.String{},
				Measurement: []types.String{},
			}, "add_key(a, 1)")
			want := extendAttr(t, newState(t, config))

			// refresh after create: no diff against configuration
			state := refresh(t, r, newState(t, config))
			require.True(t, extendAttr(t, state).Equal(want), "refresh: got %s, want %s", extendAttr(t, state), want)

			// apply a plan, refresh again: still no diff
			state = apply(t, r, state, newPlan(t, config))
			state = refresh(t, r, state)
			require.True(t, extendAttr(t, state).Equal(want), "after apply: got %s, want %s", extendAttr(t, state), want)
		})
	}
}

// Lists that are not configured stay null when the API returns them empty.
func TestUnsetListsStayNullWhenRemoteListsAreEmpty(t *testing.T) {
	r, mock := newMockedResource(t)
	mock.setLiveExtend(map[string]any{"loggingIndex": "default", "appID": []string{}, "measurement": []string{}})

	config := configModel(&pipelineExtend{LoggingIndex: types.StringValue("default")}, "add_key(a, 1)")
	state := refresh(t, r, newState(t, config))

	require.True(t, extendAttr(t, state).Equal(extendAttr(t, newState(t, config))))
}

// A remote list that was emptied while configuration has items still shows
// up as a diff.
func TestRemoteEmptiedListIsDetected(t *testing.T) {
	r, mock := newMockedResource(t)
	mock.setLiveExtend(map[string]any{"loggingIndex": ""})

	config := configModel(&pipelineExtend{
		AppID:        []types.String{types.StringValue("app_xxx")},
		LoggingIndex: types.StringValue(""),
	}, "add_key(a, 1)")
	state := refresh(t, r, newState(t, config))

	require.False(t, extendAttr(t, state).Equal(extendAttr(t, newState(t, config))))
}
