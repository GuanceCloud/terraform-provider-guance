package intelligentinspection

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
	_ resource.Resource                = &intelligentInspectionResource{}
	_ resource.ResourceWithConfigure   = &intelligentInspectionResource{}
	_ resource.ResourceWithImportState = &intelligentInspectionResource{}
)

// NewIntelligentInspectionResource is <no value>
func NewIntelligentInspectionResource() resource.Resource {
	return &intelligentInspectionResource{}
}

// intelligentInspectionResource is the resource implementation.
type intelligentInspectionResource struct {
	client *sdk.Client[*intelligentInspectionResourceModel]
}

// Schema defines the schema for the data source.
func (r *intelligentInspectionResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = resourceSchema
}

// Configure adds the provider configured client to the data source.
func (r *intelligentInspectionResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.client = req.ProviderData.(*sdk.Client[*intelligentInspectionResourceModel])
}

// Metadata returns the data source type name.
func (r *intelligentInspectionResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_intelligentinspection"
}

// Create creates the resource and sets the initial Terraform state.
func (r *intelligentInspectionResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan intelligentInspectionResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.Create(ctx, &plan); err != nil {
		resp.Diagnostics.AddError(
			"Error creating intelligentInspection",
			"Could not create intelligentInspection, unexpected error: "+err.Error(),
		)
		return
	}

	if err := r.client.Read(ctx, &plan); err != nil {
		resp.Diagnostics.AddError(
			"Error creating intelligentInspection",
			"Could not create intelligentInspection, unexpected error: "+err.Error(),
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
func (r *intelligentInspectionResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state intelligentInspectionResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.Read(ctx, &state); err != nil {
		resp.Diagnostics.AddError(
			"Error reading intelligentInspection",
			"Could not read intelligentInspection, unexpected error: "+err.Error(),
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
func (r *intelligentInspectionResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan intelligentInspectionResourceModel
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
func (r *intelligentInspectionResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state intelligentInspectionResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete existing intelligentInspection
	if err := r.client.Delete(ctx, &state); err != nil {
		resp.Diagnostics.AddError(
			"Error deleting intelligentInspection",
			"Could not delete intelligentInspection, unexpected error: "+err.Error(),
		)
		return
	}
}

func (r *intelligentInspectionResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
