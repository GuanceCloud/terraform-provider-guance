package monitor_json

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/require"
)

const baseChecker = `{"alertPolicyNames":["policy-a","policy-b"],"extend":{"rules":[{"status":"critical"},{"status":"warning"}]},"jsonScript":{"groupBy":["host","service"],"message":"m1","title":"t"},"tagInfo":[{"name":"tag-1"},{"name":"tag-2"}]}`

func TestCheckerJsonEqual(t *testing.T) {
	cases := []struct {
		name  string
		other string
		equal bool
	}{
		{"identical", baseChecker, true},
		{"alert policies reordered",
			`{"alertPolicyNames":["policy-b","policy-a"],"extend":{"rules":[{"status":"critical"},{"status":"warning"}]},"jsonScript":{"groupBy":["host","service"],"message":"m1","title":"t"},"tagInfo":[{"name":"tag-1"},{"name":"tag-2"}]}`, true},
		{"tags reordered",
			`{"alertPolicyNames":["policy-a","policy-b"],"extend":{"rules":[{"status":"critical"},{"status":"warning"}]},"jsonScript":{"groupBy":["host","service"],"message":"m1","title":"t"},"tagInfo":[{"name":"tag-2"},{"name":"tag-1"}]}`, true},
		{"key order and whitespace",
			`{ "tagInfo": [{"name":"tag-2"}, {"name":"tag-1"}], "jsonScript": {"title":"t","message":"m1","groupBy":["host","service"]}, "extend": {"rules":[{"status":"critical"},{"status":"warning"}]}, "alertPolicyNames": ["policy-b","policy-a"] }`, true},
		{"message changed",
			`{"alertPolicyNames":["policy-a","policy-b"],"extend":{"rules":[{"status":"critical"},{"status":"warning"}]},"jsonScript":{"groupBy":["host","service"],"message":"m2","title":"t"},"tagInfo":[{"name":"tag-1"},{"name":"tag-2"}]}`, false},
		{"alert policy replaced",
			`{"alertPolicyNames":["policy-a","policy-c"],"extend":{"rules":[{"status":"critical"},{"status":"warning"}]},"jsonScript":{"groupBy":["host","service"],"message":"m1","title":"t"},"tagInfo":[{"name":"tag-1"},{"name":"tag-2"}]}`, false},
		{"alert policy removed",
			`{"alertPolicyNames":["policy-a"],"extend":{"rules":[{"status":"critical"},{"status":"warning"}]},"jsonScript":{"groupBy":["host","service"],"message":"m1","title":"t"},"tagInfo":[{"name":"tag-1"},{"name":"tag-2"}]}`, false},
		{"tag duplicated instead of second tag",
			`{"alertPolicyNames":["policy-a","policy-b"],"extend":{"rules":[{"status":"critical"},{"status":"warning"}]},"jsonScript":{"groupBy":["host","service"],"message":"m1","title":"t"},"tagInfo":[{"name":"tag-1"},{"name":"tag-1"}]}`, false},
		{"rules reordered (order matters)",
			`{"alertPolicyNames":["policy-a","policy-b"],"extend":{"rules":[{"status":"warning"},{"status":"critical"}]},"jsonScript":{"groupBy":["host","service"],"message":"m1","title":"t"},"tagInfo":[{"name":"tag-1"},{"name":"tag-2"}]}`, false},
		{"groupBy reordered (order matters)",
			`{"alertPolicyNames":["policy-a","policy-b"],"extend":{"rules":[{"status":"critical"},{"status":"warning"}]},"jsonScript":{"groupBy":["service","host"],"message":"m1","title":"t"},"tagInfo":[{"name":"tag-1"},{"name":"tag-2"}]}`, false},
		{"invalid json", `{"alertPolicyNames":`, false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			require.Equal(t, tc.equal, checkerJsonEqual(baseChecker, tc.other))
			require.Equal(t, tc.equal, checkerJsonEqual(tc.other, baseChecker))
		})
	}
}

func TestCheckerJsonEqualDoesNotMutateInputOrder(t *testing.T) {
	reordered := `{"alertPolicyNames":["b","a"]}`
	require.True(t, checkerJsonEqual(reordered, `{"alertPolicyNames":["a","b"]}`))
	require.Equal(t, `{"alertPolicyNames":["b","a"]}`, reordered)
}

func runModifier(state, config types.String) types.String {
	req := planmodifier.StringRequest{StateValue: state, ConfigValue: config, PlanValue: config}
	resp := &planmodifier.StringResponse{PlanValue: req.PlanValue}
	checkerJsonOrderInsensitive().PlanModifyString(context.Background(), req, resp)
	return resp.PlanValue
}

func TestCheckerJsonPlanModifier(t *testing.T) {
	state := types.StringValue(baseChecker)

	reordered := types.StringValue(`{"alertPolicyNames":["policy-b","policy-a"],"extend":{"rules":[{"status":"critical"},{"status":"warning"}]},"jsonScript":{"groupBy":["host","service"],"message":"m1","title":"t"},"tagInfo":[{"name":"tag-2"},{"name":"tag-1"}]}`)
	require.Equal(t, state, runModifier(state, reordered), "reordered sets keep the state value")

	changed := types.StringValue(`{"alertPolicyNames":["policy-b","policy-a"],"extend":{"rules":[{"status":"critical"},{"status":"warning"}]},"jsonScript":{"groupBy":["host","service"],"message":"m2","title":"t"},"tagInfo":[{"name":"tag-2"},{"name":"tag-1"}]}`)
	require.Equal(t, changed, runModifier(state, changed), "a real change keeps the config value")

	require.Equal(t, reordered, runModifier(types.StringNull(), reordered), "create keeps the config value")
	require.Equal(t, types.StringNull(), runModifier(state, types.StringNull()), "removed config is not suppressed")
	require.Equal(t, types.StringUnknown(), runModifier(state, types.StringUnknown()), "unknown config is not suppressed")
}

func modifyPlan(t *testing.T, state, plan monitorJsonResourceModel) monitorJsonResourceModel {
	t.Helper()
	ctx := context.Background()
	st := tfsdk.State{Schema: resourceSchema}
	require.False(t, st.Set(ctx, state).HasError())
	pl := tfsdk.Plan{Schema: resourceSchema}
	require.False(t, pl.Set(ctx, plan).HasError())
	req := resource.ModifyPlanRequest{State: st, Plan: pl}
	resp := &resource.ModifyPlanResponse{Plan: pl}
	(&monitorJsonResource{}).ModifyPlan(ctx, req, resp)
	require.False(t, resp.Diagnostics.HasError())
	var got monitorJsonResourceModel
	require.False(t, resp.Plan.Get(ctx, &got).HasError())
	return got
}

func TestModifyPlanKeepsComputedWhenUnchanged(t *testing.T) {
	state := monitorJsonResourceModel{
		UUID:              types.StringValue("rul_x"),
		CheckerJson:       types.StringValue(baseChecker),
		CheckerJsonExport: types.StringValue(baseChecker),
		Type:              types.StringNull(),
		CreateAt:          types.Int64Value(1),
		UpdateAt:          types.Int64Null(),
		WorkspaceUUID:     types.StringValue("wksp_x"),
	}
	plan := state
	plan.CheckerJsonExport = types.StringUnknown()
	plan.UpdateAt = types.Int64Unknown()

	got := modifyPlan(t, state, plan)
	require.Equal(t, state.CheckerJsonExport, got.CheckerJsonExport)
	require.Equal(t, state.UpdateAt, got.UpdateAt)

	plan.CheckerJson = types.StringValue(`{"alertPolicyNames":["policy-a"]}`)
	got = modifyPlan(t, state, plan)
	require.True(t, got.CheckerJsonExport.IsUnknown(), "a real change leaves the export unknown")
	require.True(t, got.UpdateAt.IsUnknown())
}
