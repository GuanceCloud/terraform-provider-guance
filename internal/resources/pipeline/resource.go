package pipeline

import (
	"context"
	_ "embed"
	"encoding/base64"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	"github.com/GuanceCloud/terraform-provider-guance/internal/api"
	"github.com/GuanceCloud/terraform-provider-guance/internal/consts"
)

//go:embed README.md
var resourceDocument string

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &pipelineResource{}
	_ resource.ResourceWithConfigure   = &pipelineResource{}
	_ resource.ResourceWithImportState = &pipelineResource{}
)

// NewPipelineResource is <no value>
func NewPipelineResource() resource.Resource {
	return &pipelineResource{}
}

// pipelineResource is the resource implementation.
type pipelineResource struct {
	client *api.Client
}

// Schema defines the schema for the data source.
func (r *pipelineResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = resourceSchema
}

// Configure adds the provider configured client to the data source.
func (r *pipelineResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.client = req.ProviderData.(*api.Client)
}

// Metadata returns the data source type name.
func (r *pipelineResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_pipeline"
}

// Create creates the resource and sets the initial Terraform state.
func (r *pipelineResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan pipelineResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	item := r.getPipelineFromPlan(&plan)
	content := &api.Pipeline{}
	err := r.client.Create(consts.TypeNamePipeline, item, content)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating pipeline",
			"Could not create pipeline, unexpected error: "+err.Error(),
		)
		return
	}

	isDisable := plan.IsDisabled.ValueBool()
	if isDisable {
		if err := r.disablePipeline([]string{content.UUID}, isDisable); err != nil {
			resp.Diagnostics.AddError(
				"Error updating pipeline",
				"Could not create pipeline, unexpected error: "+err.Error(),
			)
			return
		}
	}

	plan.UUID = types.StringValue(content.UUID)
	plan.WorkspaceUUID = types.StringValue(content.WorkspaceUUID)
	plan.Status = types.Int64Value(content.Status)
	plan.CreateAt = types.StringValue(fmt.Sprintf("%d", content.CreateAt))
	plan.UpdateAt = types.StringValue(fmt.Sprintf("%0.0f", content.UpdateAt))

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *pipelineResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state pipelineResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	pl := &api.Pipeline{}
	err := r.client.Read(consts.TypeNamePipeline, state.UUID.ValueString(), pl)
	if err != nil {
		if err == api.Error404 {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error reading pipeline",
			"Could not read pipeline, unexpected error: "+err.Error(),
		)
		return
	}

	state.AsDefault = types.Int64Value(pl.AsDefault)
	state.EnableByLogBackup = types.Int64Value(pl.EnableByLogBackup)
	state.Category = types.StringValue(pl.Category)

	state.Extend = extendFromContent(state.Extend, pl.Extend)

	content, err := base64.StdEncoding.DecodeString(pl.Content)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading pipeline",
			"Could not read pipeline, unexpected error: "+err.Error(),
		)
		return
	}

	state.Content = types.StringValue(string(content))
	state.WorkspaceUUID = types.StringValue(pl.WorkspaceUUID)
	state.CreateAt = types.StringValue(fmt.Sprintf("%d", pl.CreateAt))
	state.UpdateAt = types.StringValue(fmt.Sprintf("%d", int64(pl.UpdateAt)))
	state.IsForce = types.BoolValue(pl.IsForce)

	testData, err := base64.StdEncoding.DecodeString(pl.TestData)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading pipeline",
			"Could not read pipeline, unexpected error: "+err.Error(),
		)
		return
	}
	state.TestData = types.StringValue(string(testData))

	state.Name = types.StringValue(pl.Name)
	state.Source = make([]basetypes.StringValue, 0)
	for _, s := range pl.Source {
		state.Source = append(state.Source, types.StringValue(s))
	}

	state.Type = types.StringValue(pl.Type)
	state.Status = types.Int64Value(pl.Status)
	state.IsDisabled = types.BoolValue(pl.Status != 0)
	state.DataType = types.StringValue(pl.DataType)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *pipelineResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan, state pipelineResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	itemID := plan.UUID.ValueString()
	item, err := r.getPipelineForUpdate(&plan, &state)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating pipeline",
			"Could not read pipeline before update, unexpected error: "+err.Error(),
		)
		return
	}

	// disable or enable
	if state.IsDisabled != plan.IsDisabled {
		if err := r.disablePipeline([]string{itemID}, plan.IsDisabled.ValueBool()); err != nil {
			resp.Diagnostics.AddError(
				"Error updating pipeline",
				"Could not update pipeline, unexpected error: "+err.Error(),
			)
			return
		}
	}

	// update pipeline, omit uuid field
	item.UUID = ""
	content := &api.Pipeline{}
	err = r.client.Update(consts.TypeNamePipeline, itemID, item, content)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating pipeline",
			"Could not update pipeline, unexpected error: "+err.Error(),
		)
		return
	}

	plan.Status = types.Int64Value(content.Status)
	plan.UpdateAt = types.StringValue(fmt.Sprintf("%.0f", content.UpdateAt))

	// Set state to fully populated data
	diags := resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *pipelineResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state pipelineResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.Delete(consts.TypeNamePipeline, state.UUID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting pipeline",
			"Could not delete pipeline, unexpected error: "+err.Error(),
		)
		return
	}

}

func (r *pipelineResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("uuid"), req, resp)
}

func (r *pipelineResource) getPipelineFromPlan(plan *pipelineResourceModel) *api.Pipeline {
	item := &api.Pipeline{
		Name:              plan.Name.ValueString(),
		Category:          plan.Category.ValueString(),
		Content:           base64.StdEncoding.EncodeToString([]byte(plan.Content.ValueString())),
		TestData:          base64.StdEncoding.EncodeToString([]byte(plan.TestData.ValueString())),
		Type:              plan.Type.ValueString(),
		UUID:              plan.UUID.ValueString(),
		AsDefault:         plan.AsDefault.ValueInt64(),
		EnableByLogBackup: plan.EnableByLogBackup.ValueInt64(),
		IsForce:           plan.IsForce.ValueBool(),
		DataType:          plan.DataType.ValueString(),
	}

	source := []string{}
	for _, s := range plan.Source {
		source = append(source, s.ValueString())
	}

	item.Source = source

	if plan.Extend != nil {
		item.Extend = &api.PipelineExtend{}
		if len(plan.Extend.AppID) > 0 {
			appID := []string{}
			for _, s := range plan.Extend.AppID {
				appID = append(appID, s.ValueString())
			}
			item.Extend.AppID = appID
		}

		if len(plan.Extend.Measurement) > 0 {
			measurement := []string{}
			for _, s := range plan.Extend.Measurement {
				measurement = append(measurement, s.ValueString())
			}
			item.Extend.Measurement = measurement
		}

		// "" is a valid value (the console stores "loggingIndex": ""), so
		// send it whenever logging_index is set, and leave it out when unset.
		if !plan.Extend.LoggingIndex.IsNull() && !plan.Extend.LoggingIndex.IsUnknown() {
			loggingIndex := plan.Extend.LoggingIndex.ValueString()
			item.Extend.LoggingIndex = &loggingIndex
		}
	}
	// set is_force to true when as_default is set to 1
	if item.AsDefault == 1 {
		item.IsForce = true
	}

	return item
}

// getPipelineForUpdate builds the modify body from the plan. The modify API
// replaces the whole object, extend included, so the live extend is read first
// and its keys that this resource does not manage are sent back unchanged.
//
// A live loggingIndex is also sent back when neither the prior state nor the
// plan has logging_index: refresh hides a remote "" in that case (see
// extendFromContent), so the plan never showed its removal.
func (r *pipelineResource) getPipelineForUpdate(plan, state *pipelineResourceModel) (*api.Pipeline, error) {
	item := r.getPipelineFromPlan(plan)

	live := &api.Pipeline{}
	if err := r.client.Read(consts.TypeNamePipeline, plan.UUID.ValueString(), live); err != nil {
		return nil, err
	}
	keepLiveLoggingIndex := state == nil || state.Extend == nil || state.Extend.LoggingIndex.IsNull()
	item.Extend = mergeExtend(live.Extend, item.Extend, keepLiveLoggingIndex)

	return item, nil
}

// mergeExtend returns the extend to send on modify: appID, measurement and
// loggingIndex come from the plan, every other key comes from the live object.
// When keepLiveLoggingIndex is true and the plan has no loggingIndex, the live
// loggingIndex is kept as well.
func mergeExtend(live, planned *api.PipelineExtend, keepLiveLoggingIndex bool) *api.PipelineExtend {
	if live == nil {
		return planned
	}
	keepIndex := keepLiveLoggingIndex && live.LoggingIndex != nil &&
		(planned == nil || planned.LoggingIndex == nil)
	if len(live.Extra) == 0 && !keepIndex {
		return planned
	}

	merged := &api.PipelineExtend{}
	if len(live.Extra) > 0 {
		merged.Extra = make(map[string]json.RawMessage, len(live.Extra))
		for k, v := range live.Extra {
			merged.Extra[k] = v
		}
	}
	if planned != nil {
		merged.AppID = planned.AppID
		merged.Measurement = planned.Measurement
		merged.LoggingIndex = planned.LoggingIndex
	}
	if keepIndex {
		loggingIndex := *live.LoggingIndex
		merged.LoggingIndex = &loggingIndex
	}
	return merged
}

// extendFromContent builds the extend state from the remote object.
//
// It returns nil when the remote object has none of the managed keys, so that
// a remote extend that was emptied shows up in the plan instead of the prior
// state being kept. A remote "loggingIndex": "" is kept as "" (the console
// writes it, and dropping it is a real change), except when the prior state
// has extend without logging_index, where "" is the API's default.
//
// The API leaves empty appID and measurement lists out, so a list that is
// empty in the prior state stays empty (not null) when the remote list is
// absent or empty; a list that is null in the prior state stays null.
func extendFromContent(prior *pipelineExtend, remote *api.PipelineExtend) *pipelineExtend {
	if remote == nil || (remote.AppID == nil && remote.Measurement == nil && remote.LoggingIndex == nil) {
		// an extend block with no attributes set, or only empty lists, is
		// equivalent to an empty remote extend
		if prior != nil && len(prior.AppID) == 0 && len(prior.Measurement) == 0 && prior.LoggingIndex.IsNull() {
			return prior
		}
		return nil
	}

	extend := &pipelineExtend{LoggingIndex: types.StringNull()}
	var priorAppID, priorMeasurement []types.String
	if prior != nil {
		priorAppID, priorMeasurement = prior.AppID, prior.Measurement
	}
	extend.AppID = listFromContent(priorAppID, remote.AppID)
	extend.Measurement = listFromContent(priorMeasurement, remote.Measurement)
	if remote.LoggingIndex != nil {
		extend.LoggingIndex = types.StringValue(*remote.LoggingIndex)
		if *remote.LoggingIndex == "" && prior != nil && prior.LoggingIndex.IsNull() {
			extend.LoggingIndex = types.StringNull()
		}
	}
	return extend
}

// listFromContent converts a remote list. An absent or empty remote list is
// an empty list when the prior state has an empty list, and null otherwise.
func listFromContent(prior []types.String, remote []string) []types.String {
	if len(remote) == 0 {
		if prior != nil && len(prior) == 0 {
			return []types.String{}
		}
		return nil
	}
	list := make([]types.String, 0, len(remote))
	for _, v := range remote {
		list = append(list, types.StringValue(v))
	}
	return list
}

func (r *pipelineResource) disablePipeline(itemIDs []string, isDisable bool) error {
	isSuccess := false
	if err := r.client.DisablePipeline(itemIDs, isDisable, &isSuccess); err != nil {
		return fmt.Errorf("disable failed: %w", err)
	}

	if !isSuccess {
		return fmt.Errorf("disabled failed")
	}
	return nil
}
