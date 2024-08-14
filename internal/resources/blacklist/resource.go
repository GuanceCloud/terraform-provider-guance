package blacklist

import (
	"context"
	_ "embed"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/GuanceCloud/terraform-provider-guance/internal/api"
	"github.com/GuanceCloud/terraform-provider-guance/internal/consts"
)

//go:embed README.md
var resourceDocument string

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &blackListResource{}
	_ resource.ResourceWithConfigure   = &blackListResource{}
	_ resource.ResourceWithImportState = &blackListResource{}
)

// NewBlackListResource is <no value>
func NewBlackListResource() resource.Resource {
	return &blackListResource{}
}

// blackListResource is the resource implementation.
type blackListResource struct {
	client *api.Client
}

// Schema defines the schema for the data source.
func (r *blackListResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = resourceSchema
}

// Configure adds the provider configured client to the data source.
func (r *blackListResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.client = req.ProviderData.(*api.Client)
}

// Metadata returns the data source type name.
func (r *blackListResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_blacklist"
}

// Create creates the resource and sets the initial Terraform state.
func (r *blackListResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan blackListResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	item := r.getBlacklistFromPlan(&plan)
	content := &api.Blacklist{}
	err := r.client.Create(consts.TypeNameBlackList, item, content)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating blackList",
			"Could not create blackList, unexpected error: "+err.Error(),
		)
		return
	}

	plan.UUID = types.StringValue(content.UUID)
	plan.CreatedAt = types.StringValue(fmt.Sprintf("%d", content.CreateAt))

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *blackListResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state blackListResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	content := &api.Blacklist{}
	err := r.client.Read(consts.TypeNameBlackList, state.UUID.ValueString(), content)
	if err != nil {
		if err == api.Error404 {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error reading blacklist",
			"Could not read blacklist, unexpected error: "+err.Error(),
		)
		return
	}

	state.CreatedAt = types.StringValue(fmt.Sprintf("%d", content.CreateAt))
	state.UUID = types.StringValue(content.UUID)
	state.Source = types.StringValue(content.Source)
	state.Type = types.StringValue(content.Type)
	state.Filters = make([]*filter, 0)
	for _, item := range content.Filters {
		f := &filter{
			Name:      types.StringValue(item.Name),
			Condition: types.StringValue(item.Condition),
			Operation: types.StringValue(item.Operation),
		}

		for _, value := range item.Value {
			f.Values = append(f.Values, types.StringValue(value))
		}

		state.Filters = append(state.Filters, f)
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *blackListResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan blackListResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	item := r.getBlacklistFromPlan(&plan)
	item.UUID = plan.UUID.ValueString()
	body := &api.Blacklist{}
	err := r.client.Update(consts.TypeNameBlackList, item.UUID, item, body)

	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating blacklist",
			"Could not update blacklist, unexpected error: "+err.Error(),
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
func (r *blackListResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state blackListResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	//Delete existing blackList
	if err := r.client.Delete(consts.TypeNameBlackList, state.UUID.ValueString()); err != nil {
		resp.Diagnostics.AddError(
			"Error deleting blackList",
			"Could not delete blackList, unexpected error: "+err.Error(),
		)
		return
	}
}

func (r *blackListResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("uuid"), req, resp)
}

func (r *blackListResource) getBlacklistFromPlan(plan *blackListResourceModel) *api.Blacklist {
	item := &api.Blacklist{
		Type:   plan.Type.ValueString(),
		Source: plan.Source.ValueString(),
	}

	if len(plan.Filters) > 0 {
		filters := []api.Filter{}
		for _, filter := range plan.Filters {
			f := api.Filter{
				Name:      filter.Name.ValueString(),
				Condition: filter.Condition.ValueString(),
				Operation: filter.Operation.ValueString(),
			}

			for _, value := range filter.Values {
				f.Value = append(f.Value, value.ValueString())
			}

			filters = append(filters, f)
		}
		item.Filters = filters
	}

	return item
}
