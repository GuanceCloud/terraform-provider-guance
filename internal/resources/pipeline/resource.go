package pipeline

import (
	"context"
	_ "embed"
	"encoding/base64"
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

	plan.UUID = types.StringValue(content.UUID)
	plan.CreateAt = types.StringValue(fmt.Sprintf("%d", content.CreateAt))

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
	state.Category = types.StringValue(pl.Category)

	// only set extend when it's not nil
	if pl.Extend != nil && (pl.Extend.AppID != nil || pl.Extend.Measurement != nil) {
		state.Extend = &pipelineExtend{}
		for _, id := range pl.Extend.AppID {
			state.Extend.AppID = append(state.Extend.AppID, types.StringValue(id))
		}
		for _, m := range pl.Extend.Measurement {
			state.Extend.Measurement = append(state.Extend.Measurement, types.StringValue(m))
		}
	}

	content, err := base64.StdEncoding.DecodeString(pl.Content)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading pipeline",
			"Could not read pipeline, unexpected error: "+err.Error(),
		)
		return
	}

	state.Content = types.StringValue(string(content))

	state.CreateAt = types.StringValue(fmt.Sprintf("%d", pl.CreateAt))
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
	var plan pipelineResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	item := r.getPipelineFromPlan(&plan)
	item.UUID = plan.UUID.ValueString()

	content := &api.Pipeline{}
	err := r.client.Update(consts.TypeNamePipeline, item.UUID, item, content)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating pipeline",
			"Could not update pipeline, unexpected error: "+err.Error(),
		)
		return
	}

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
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
		Name:      plan.Name.ValueString(),
		Category:  plan.Category.ValueString(),
		Content:   base64.StdEncoding.EncodeToString([]byte(plan.Content.ValueString())),
		TestData:  base64.StdEncoding.EncodeToString([]byte(plan.TestData.ValueString())),
		Type:      plan.Type.ValueString(),
		UUID:      plan.UUID.ValueString(),
		AsDefault: plan.AsDefault.ValueInt64(),
		IsForce:   plan.IsForce.ValueBool(),
	}

	if len(plan.Source) > 0 {
		source := []string{}
		for _, s := range plan.Source {
			source = append(source, s.ValueString())
		}

		item.Source = source
	}

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
	}
	// set is_force to true when as_default is set to 1
	if item.AsDefault == 1 {
		item.IsForce = true
	}

	return item
}
