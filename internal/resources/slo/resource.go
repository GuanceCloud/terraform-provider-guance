package slo

import (
	"context"
	_ "embed"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/GuanceCloud/terraform-provider-guance/internal/api"
	"github.com/GuanceCloud/terraform-provider-guance/internal/consts"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

//go:embed README.md
var resourceDocument string

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &sloResource{}
	_ resource.ResourceWithConfigure   = &sloResource{}
	_ resource.ResourceWithImportState = &sloResource{}
)

// NewSloResource creates a new SLO resource
func NewSloResource() resource.Resource {
	return &sloResource{}
}

// sloResource is the resource implementation.
type sloResource struct {
	client *api.Client
}

// Schema defines the schema for the data source.
func (r *sloResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = resourceSchema
}

// Configure adds the provider configured client to the data source.
func (r *sloResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.client = req.ProviderData.(*api.Client)
}

// Metadata returns the data source type name.
func (r *sloResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_slo"
}

// Create creates the resource and sets the initial Terraform state.
func (r *sloResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan sloResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	item := r.getSloFromPlan(&plan)
	rb, _ := json.Marshal(item)
	tflog.Debug(ctx, fmt.Sprintf("============= body: %s", string(rb)))
	content := &api.SLOContent{}
	err := r.client.Create(consts.TypeNameSLO, item, content)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating SLO",
			"Could not create SLO, unexpected error: "+err.Error(),
		)
		return
	}

	plan.UUID = types.StringValue(content.UUID)
	plan.WorkspaceUUID = types.StringValue(content.WorkspaceUUID)
	plan.CreateAt = types.Int64Value(int64(content.CreateAt))
	plan.UpdateAt = types.Int64Value(int64(content.UpdateAt))
	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *sloResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state sloResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	content := &api.SLOContent{}
	err := r.client.Read(consts.TypeNameSLO, state.UUID.ValueString(), content)
	if err != nil {
		if err == api.Error404 {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error reading SLO",
			"Could not read SLO, unexpected error: "+err.Error(),
		)
		return
	}

	state.UUID = types.StringValue(content.UUID)
	state.Name = types.StringValue(content.Name)
	state.Interval = types.StringValue(content.Config.Interval)
	state.Goal = types.Float64Value(content.Config.Goal)
	state.MinGoal = types.Float64Value(content.Config.MinGoal)
	state.Describe = types.StringValue(content.Config.Describe)

	// Map SLI UUIDs
	sliUUIDs := make([]types.String, len(content.Config.SliInfos))
	for i, sli := range content.Config.SliInfos {
		sliUUIDs[i] = types.StringValue(sli.ID)
	}
	state.SliUUIDs = sliUUIDs

	// Map alert policy UUIDs
	alertPolicyUUIDs := make([]types.String, 0)
	if content.AlertPolicyInfos != nil {
		for _, info := range content.AlertPolicyInfos {
			if infoMap, ok := info.(map[string]interface{}); ok {
				if uuid, ok := infoMap["uuid"].(string); ok {
					alertPolicyUUIDs = append(alertPolicyUUIDs, types.StringValue(uuid))
				}
			}
		}
	}
	state.AlertPolicyUUIDs = alertPolicyUUIDs

	state.CreateAt = types.Int64Value(int64(content.CreateAt))
	state.UpdateAt = types.Int64Value(int64(content.UpdateAt))
	state.WorkspaceUUID = types.StringValue(content.WorkspaceUUID)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *sloResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan sloResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	item := r.getSloFromPlan(&plan)
	content := &api.SLOContent{}
	err := r.client.Update(consts.TypeNameSLO, plan.UUID.ValueString(), item, content)

	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating SLO",
			"Could not update SLO, unexpected error: "+err.Error(),
		)
		return
	}

	plan.UpdateAt = types.Int64Value(int64(content.UpdateAt))
	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *sloResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state sloResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete existing SLO
	if err := r.client.Delete(consts.TypeNameSLO, state.UUID.ValueString()); err != nil {
		resp.Diagnostics.AddError(
			"Error deleting SLO",
			"Could not delete SLO, unexpected error: "+err.Error(),
		)
		return
	}
}

func (r *sloResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("uuid"), req, resp)
}

func (r *sloResource) getSloFromPlan(plan *sloResourceModel) *api.SLO {

	s := &api.SLO{
		Name:     plan.Name.ValueString(),
		Interval: plan.Interval.ValueString(),
		Goal:     plan.Goal.ValueFloat64(),
		MinGoal:  plan.MinGoal.ValueFloat64(),
	}

	if !plan.Describe.IsNull() {
		s.Describe = plan.Describe.ValueString()
	}

	if len(plan.SliUUIDs) > 0 {
		sliUUIDs := make([]string, len(plan.SliUUIDs))
		for i, sli := range plan.SliUUIDs {
			sliUUIDs[i] = sli.ValueString()
		}
		s.SliUUIDs = sliUUIDs
	}

	if len(plan.AlertPolicyUUIDs) > 0 {
		alertPolicyUUIDs := make([]string, len(plan.AlertPolicyUUIDs))
		for i, policy := range plan.AlertPolicyUUIDs {
			alertPolicyUUIDs[i] = policy.ValueString()
		}
		s.AlertPolicyUUIDs = alertPolicyUUIDs
	}

	if len(plan.Tags) > 0 {
		tags := make([]string, len(plan.Tags))
		for i, tag := range plan.Tags {
			tags[i] = tag.ValueString()
		}
		s.Tags = tags
	}

	return s
}
