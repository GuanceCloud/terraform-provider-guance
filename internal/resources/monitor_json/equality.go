package monitor_json

import (
	"context"
	"encoding/json"
	"reflect"
	"sort"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
)

// unorderedListKeys are the top-level checker fields whose arrays are sets:
// the platform returns their items in a different order on each export, and
// the order carries no meaning. Every other array (rules, conditions, query
// lists, group-by fields ...) keeps its order in the comparison.
var unorderedListKeys = []string{
	"alertPolicyNames",
	"tagInfo",
}

// checkerJsonEqual reports whether two checker JSON documents describe the same
// monitor: equal after decoding (so whitespace and object key order are
// ignored), with the arrays listed in unorderedListKeys compared as multisets.
// Invalid JSON is only equal to the identical string.
func checkerJsonEqual(a, b string) bool {
	if a == b {
		return true
	}
	va, err := normalizeCheckerJson(a)
	if err != nil {
		return false
	}
	vb, err := normalizeCheckerJson(b)
	if err != nil {
		return false
	}
	return reflect.DeepEqual(va, vb)
}

func normalizeCheckerJson(s string) (interface{}, error) {
	var v interface{}
	if err := json.Unmarshal([]byte(s), &v); err != nil {
		return nil, err
	}
	obj, ok := v.(map[string]interface{})
	if !ok {
		return v, nil
	}
	for _, key := range unorderedListKeys {
		list, ok := obj[key].([]interface{})
		if !ok || len(list) < 2 {
			continue
		}
		sorted, err := sortByEncoding(list)
		if err != nil {
			return nil, err
		}
		obj[key] = sorted
	}
	return obj, nil
}

// sortByEncoding returns a copy of list ordered by each item's canonical JSON
// encoding (encoding/json sorts map keys), so it works for strings and objects.
func sortByEncoding(list []interface{}) ([]interface{}, error) {
	type item struct {
		key string
		val interface{}
	}
	items := make([]item, len(list))
	for i, v := range list {
		b, err := json.Marshal(v)
		if err != nil {
			return nil, err
		}
		items[i] = item{string(b), v}
	}
	sort.SliceStable(items, func(i, j int) bool { return items[i].key < items[j].key })
	out := make([]interface{}, len(items))
	for i, it := range items {
		out[i] = it.val
	}
	return out, nil
}

// checkerJsonOrderInsensitive keeps the prior state value in the plan when the
// configured checker_json differs from it only in ways checkerJsonEqual ignores,
// so a reordered tag or alert-policy list does not plan an update.
func checkerJsonOrderInsensitive() planmodifier.String {
	return checkerJsonOrderInsensitiveModifier{}
}

type checkerJsonOrderInsensitiveModifier struct{}

func (m checkerJsonOrderInsensitiveModifier) Description(_ context.Context) string {
	return "Ignores differences in JSON formatting, object key order and the order of tagInfo and alertPolicyNames."
}

func (m checkerJsonOrderInsensitiveModifier) MarkdownDescription(ctx context.Context) string {
	return m.Description(ctx)
}

func (m checkerJsonOrderInsensitiveModifier) PlanModifyString(_ context.Context, req planmodifier.StringRequest, resp *planmodifier.StringResponse) {
	if req.StateValue.IsNull() || req.StateValue.IsUnknown() ||
		req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}
	if checkerJsonEqual(req.ConfigValue.ValueString(), req.StateValue.ValueString()) {
		resp.PlanValue = req.StateValue
	}
}

// ModifyPlan keeps the computed attributes from state when the plan leaves
// checker_json and type as they are in state. The framework marks unset
// computed attributes unknown whenever the configuration differs from state,
// which happens before checkerJsonOrderInsensitive restores the state value,
// so without this a reordered list would still plan an update.
func (r *monitorJsonResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	if req.State.Raw.IsNull() || req.Plan.Raw.IsNull() {
		return
	}
	var plan, state monitorJsonResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if !plan.CheckerJson.Equal(state.CheckerJson) || !plan.Type.Equal(state.Type) {
		return
	}
	plan.CheckerJsonExport = state.CheckerJsonExport
	plan.UpdateAt = state.UpdateAt
	resp.Diagnostics.Append(resp.Plan.Set(ctx, plan)...)
}
