package membergroup

import (
	"context"
	_ "embed"

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
	_ resource.Resource                = &memberGroupResource{}
	_ resource.ResourceWithConfigure   = &memberGroupResource{}
	_ resource.ResourceWithImportState = &memberGroupResource{}
)

// NewMemberGroupResource is <no value>
func NewMemberGroupResource() resource.Resource {
	return &memberGroupResource{}
}

// memberGroupResource is the resource implementation.
type memberGroupResource struct {
	client *api.Client
}

// Schema defines the schema for the data source.
func (r *memberGroupResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = resourceSchema
}

// Configure adds the provider configured client to the data source.
func (r *memberGroupResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.client = req.ProviderData.(*api.Client)
}

// Metadata returns the data source type name.
func (r *memberGroupResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_membergroup"
}

// Create creates the resource and sets the initial Terraform state.
func (r *memberGroupResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan memberGroupResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	item := r.getMembergroupFromPlan(&plan)
	content := &api.MembergroupContent{}
	err := r.client.Create(consts.TypeNameMemberGroup, item, content)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating membergroup",
			"Could not create membergroup, unexpected error: "+err.Error(),
		)
		return
	}

	plan.UUID = types.StringValue(content.UUID)

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *memberGroupResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state memberGroupResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	content := &api.MembergroupReadContent{}
	if err := r.client.Read(consts.TypeNameMemberGroup, state.UUID.ValueString(), content); err != nil {

		if err == api.Error404 {
			resp.State.RemoveResource(ctx)
			return
		}

		resp.Diagnostics.AddError(
			"Error reading memberGroup",
			"Could not read memberGroup, unexpected error: "+err.Error(),
		)
		return
	}

	state.AccountUUIDs = make([]basetypes.StringValue, 0)
	for _, v := range content.GroupMembers {
		state.AccountUUIDs = append(state.AccountUUIDs, types.StringValue(v.UUID))
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *memberGroupResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan memberGroupResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	item := r.getMembergroupFromPlan(&plan)
	item.UUID = "" // set empty to avoid update failure
	content := false
	if err := r.client.Update(consts.TypeNameMemberGroup, plan.UUID.ValueString(), item, &content); err != nil {
		resp.Diagnostics.AddError(
			"Error updating memberGroup",
			"Could not update memberGroup, unexpected error: "+err.Error(),
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
func (r *memberGroupResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state memberGroupResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete existing memberGroup
	if err := r.client.Delete(consts.TypeNameMemberGroup, state.UUID.ValueString()); err != nil {
		resp.Diagnostics.AddError(
			"Error deleting memberGroup",
			"Could not delete memberGroup, unexpected error: "+err.Error(),
		)
		return
	}
}

func (r *memberGroupResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("uuid"), req, resp)
}

func (r *memberGroupResource) getMembergroupFromPlan(plan *memberGroupResourceModel) *api.Membergroup {

	g := &api.Membergroup{
		UUID: plan.UUID.ValueString(),
		Name: plan.Name.ValueString(),
	}

	g.AccountUUIDs = make([]string, 0)
	for _, v := range plan.AccountUUIDs {
		g.AccountUUIDs = append(g.AccountUUIDs, v.ValueString())
	}

	return g
}
