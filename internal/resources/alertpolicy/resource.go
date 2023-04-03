package alertpolicy

import (
	"context"
	_ "embed"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"

	"github.com/GuanceCloud/terraform-provider-guance/internal/sdk"
)

//go:embed README.md
var resourceDocument string

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &alertPolicyResource{}
	_ resource.ResourceWithConfigure   = &alertPolicyResource{}
	_ resource.ResourceWithImportState = &alertPolicyResource{}
)

// NewAlertPolicyResource is <no value>
func NewAlertPolicyResource() resource.Resource {
	return &alertPolicyResource{}
}

// alertPolicyResource is the resource implementation.
type alertPolicyResource struct {
	client *sdk.Client[*alertPolicyResourceModel]
}

// Schema defines the schema for the data source.
func (r *alertPolicyResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = resourceSchema
}

// Configure adds the provider configured client to the data source.
func (r *alertPolicyResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.client = req.ProviderData.(*sdk.Client[*alertPolicyResourceModel])
}

// Metadata returns the data source type name.
func (r *alertPolicyResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_alertpolicy"
}

// Create creates the resource and sets the initial Terraform state.
func (r *alertPolicyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan alertPolicyResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.Create(ctx, &plan); err != nil {
		resp.Diagnostics.AddError(
			"Error creating alertPolicy",
			"Could not create alertPolicy, unexpected error: "+err.Error(),
		)
		return
	}

	if err := r.client.Read(ctx, &plan); err != nil {
		resp.Diagnostics.AddError(
			"Error creating alertPolicy",
			"Could not create alertPolicy, unexpected error: "+err.Error(),
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

// Read refreshes the Terraform state with the latest data.
func (r *alertPolicyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state alertPolicyResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.Read(ctx, &state); err != nil {
		resp.Diagnostics.AddError(
			"Error reading alertPolicy",
			"Could not read alertPolicy, unexpected error: "+err.Error(),
		)
		return
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *alertPolicyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan alertPolicyResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *alertPolicyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state alertPolicyResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete existing alertPolicy
	if err := r.client.Delete(ctx, &state); err != nil {
		resp.Diagnostics.AddError(
			"Error deleting alertPolicy",
			"Could not delete alertPolicy, unexpected error: "+err.Error(),
		)
		return
	}
}

func (r *alertPolicyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
