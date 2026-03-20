package custom_region

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/GuanceCloud/terraform-provider-guance/internal/api"
	"github.com/GuanceCloud/terraform-provider-guance/internal/consts"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &customRegionResource{}
	_ resource.ResourceWithConfigure   = &customRegionResource{}
	_ resource.ResourceWithImportState = &customRegionResource{}
)

// NewCustomRegionResource creates a new custom region resource
func NewCustomRegionResource() resource.Resource {
	return &customRegionResource{}
}

// customRegionResource is the resource implementation.
type customRegionResource struct {
	client *api.Client
}

// Schema defines the schema for the data source.
func (r *customRegionResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = resourceSchema
}

// Configure adds the provider configured client to the data source.
func (r *customRegionResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.client = req.ProviderData.(*api.Client)
}

// Metadata returns the data source type name.
func (r *customRegionResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_custom_region"
}

// Create creates the resource and sets the initial Terraform state.
func (r *customRegionResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan customRegionResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	item := r.getCustomRegionFromPlan(&plan)
	response := &api.CustomRegionResponse{}
	tflog.Debug(ctx, fmt.Sprintf("Creating custom region with item: %v", item))
	err := r.client.Create(consts.TypeNameCustomRegion, item, response)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating custom region",
			"Could not create custom region, unexpected error: "+err.Error(),
		)
		return
	}
	tflog.Debug(ctx, fmt.Sprintf("Create custom region response: %v", response))
	tflog.Debug(ctx, fmt.Sprintf("Create custom region response Content: %v", response.Content))
	tflog.Debug(ctx, fmt.Sprintf("Create custom region response Content.AK: %v", response.Content.AK))
	tflog.Debug(ctx, fmt.Sprintf("Create custom region response Content.RegionInfo: %v", response.Content.RegionInfo))

	// Set state from response
	if response.Content.RegionInfo != nil {
		plan.UUID = types.StringValue(response.Content.RegionInfo.UUID)
		plan.CreateAt = types.Int64Value(response.Content.RegionInfo.CreateAt)
		plan.ExternalID = types.StringValue(response.Content.RegionInfo.ExternalID)
		plan.Heartbeat = types.Int64Value(response.Content.RegionInfo.Heartbeat)
		plan.Owner = types.StringValue(response.Content.RegionInfo.Owner)
		plan.ParentAK = types.StringValue(response.Content.RegionInfo.ParentAK)
		plan.Region = types.StringValue(response.Content.RegionInfo.Region)
		plan.Status = types.StringValue(response.Content.RegionInfo.Status)
		plan.NameEn = types.StringValue(response.Content.RegionInfo.NameEn)
		plan.ExtendInfo = types.StringValue(response.Content.RegionInfo.ExtendInfo)
	} else {
		// Set default values for computed attributes
		plan.UUID = types.StringValue("")
		plan.CreateAt = types.Int64Value(0)
		plan.ExternalID = types.StringValue("")
		plan.Heartbeat = types.Int64Value(0)
		plan.Owner = types.StringValue("")
		plan.ParentAK = types.StringValue("")
		plan.Region = types.StringValue("")
		plan.Status = types.StringValue("")
		plan.NameEn = types.StringValue("")
		plan.ExtendInfo = types.StringValue("")
	}

	// Set server and declaration info
	if response.Content.Server != "" {
		plan.Server = types.StringValue(response.Content.Server)
	} else {
		plan.Server = types.StringValue("")
	}

	if response.Content.Declaration != nil {
		declarationMap, diags := types.MapValueFrom(ctx, types.StringType, response.Content.Declaration)
		if diags.HasError() {
			resp.Diagnostics.Append(diags...)
			return
		}
		plan.Declaration = declarationMap
	} else {
		// Create empty declaration map
		declarationMap, diags := types.MapValueFrom(ctx, types.StringType, map[string]string{})
		if diags.HasError() {
			resp.Diagnostics.Append(diags...)
			return
		}
		plan.Declaration = declarationMap
	}

	// Set AK/SK info
	if response.Content.AK != nil {
		// Create AK object using ObjectValue
		akObject, diags := types.ObjectValue(map[string]attr.Type{
			"ak":          types.StringType,
			"sk":          types.StringType,
			"external_id": types.StringType,
			"owner":       types.StringType,
			"parent_ak":   types.StringType,
			"status":      types.StringType,
			"update_at":   types.Int64Type,
		}, map[string]attr.Value{
			"ak":          types.StringValue(response.Content.AK.AK),
			"sk":          types.StringValue(response.Content.AK.SK),
			"external_id": types.StringValue(response.Content.AK.ExternalID),
			"owner":       types.StringValue(response.Content.AK.Owner),
			"parent_ak":   types.StringValue(response.Content.AK.ParentAK),
			"status":      types.StringValue(response.Content.AK.Status),
			"update_at":   types.Int64Value(response.Content.AK.UpdateAt),
		})
		if diags.HasError() {
			resp.Diagnostics.Append(diags...)
			return
		}
		plan.AK = akObject
	} else {
		// Create empty AK object
		akObject, diags := types.ObjectValue(map[string]attr.Type{
			"ak":          types.StringType,
			"sk":          types.StringType,
			"external_id": types.StringType,
			"owner":       types.StringType,
			"parent_ak":   types.StringType,
			"status":      types.StringType,
			"update_at":   types.Int64Type,
		}, map[string]attr.Value{
			"ak":          types.StringValue(""),
			"sk":          types.StringValue(""),
			"external_id": types.StringValue(""),
			"owner":       types.StringValue(""),
			"parent_ak":   types.StringValue(""),
			"status":      types.StringValue(""),
			"update_at":   types.Int64Value(0),
		})
		if diags.HasError() {
			resp.Diagnostics.Append(diags...)
			return
		}
		plan.AK = akObject
	}

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *customRegionResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state customRegionResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	response := &api.CustomRegionResponse{}
	tflog.Debug(ctx, fmt.Sprintf("Reading custom region with UUID: '%s'", state.UUID.ValueString()))
	if state.UUID.IsUnknown() || state.UUID.ValueString() == "" {
		// Resource hasn't been created yet, return without error
		return
	}
	err := r.client.Read(consts.TypeNameCustomRegion, state.UUID.ValueString(), response)
	if err != nil {
		if err == api.Error404 {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error reading custom region",
			"Could not read custom region, unexpected error: "+err.Error(),
		)
		return
	}

	// Update state from response
	if response.Content.RegionInfo != nil {
		state.UUID = types.StringValue(response.Content.RegionInfo.UUID)
		state.Internal = types.BoolValue(response.Content.RegionInfo.Internal)
		state.ISP = types.StringValue(response.Content.RegionInfo.ISP)
		state.Country = types.StringValue(response.Content.RegionInfo.Country)
		state.Province = types.StringValue(response.Content.RegionInfo.Province)
		state.City = types.StringValue(response.Content.RegionInfo.City)
		state.Name = types.StringValue(response.Content.RegionInfo.Name)
		state.NameEn = types.StringValue(response.Content.RegionInfo.NameEn)
		state.Company = types.StringValue(response.Content.RegionInfo.Company)
		state.Keycode = types.StringValue(response.Content.RegionInfo.Keycode)
		state.CreateAt = types.Int64Value(response.Content.RegionInfo.CreateAt)
		state.ExtendInfo = types.StringValue(response.Content.RegionInfo.ExtendInfo)
		state.ExternalID = types.StringValue(response.Content.RegionInfo.ExternalID)
		state.Heartbeat = types.Int64Value(response.Content.RegionInfo.Heartbeat)
		state.Owner = types.StringValue(response.Content.RegionInfo.Owner)
		state.ParentAK = types.StringValue(response.Content.RegionInfo.ParentAK)
		state.Region = types.StringValue(response.Content.RegionInfo.Region)
		state.Status = types.StringValue(response.Content.RegionInfo.Status)
	}

	// Update server and declaration info
	if response.Content.Server != "" {
		state.Server = types.StringValue(response.Content.Server)
	} else {
		state.Server = types.StringValue("")
	}

	if response.Content.Declaration != nil {
		declarationMap, diags := types.MapValueFrom(ctx, types.StringType, response.Content.Declaration)
		if diags.HasError() {
			resp.Diagnostics.Append(diags...)
			return
		}
		state.Declaration = declarationMap
	} else {
		// Create empty declaration map
		declarationMap, diags := types.MapValueFrom(ctx, types.StringType, map[string]string{})
		if diags.HasError() {
			resp.Diagnostics.Append(diags...)
			return
		}
		state.Declaration = declarationMap
	}

	// Update AK/SK info
	if response.Content.AK != nil {
		// Create AK object using ObjectValue
		akObject, diags := types.ObjectValue(map[string]attr.Type{
			"ak":          types.StringType,
			"sk":          types.StringType,
			"external_id": types.StringType,
			"owner":       types.StringType,
			"parent_ak":   types.StringType,
			"status":      types.StringType,
			"update_at":   types.Int64Type,
		}, map[string]attr.Value{
			"ak":          types.StringValue(response.Content.AK.AK),
			"sk":          types.StringValue(response.Content.AK.SK),
			"external_id": types.StringValue(response.Content.AK.ExternalID),
			"owner":       types.StringValue(response.Content.AK.Owner),
			"parent_ak":   types.StringValue(response.Content.AK.ParentAK),
			"status":      types.StringValue(response.Content.AK.Status),
			"update_at":   types.Int64Value(response.Content.AK.UpdateAt),
		})
		if diags.HasError() {
			resp.Diagnostics.Append(diags...)
			return
		}
		state.AK = akObject
	} else {
		// Create empty AK object
		akObject, diags := types.ObjectValue(map[string]attr.Type{
			"ak":          types.StringType,
			"sk":          types.StringType,
			"external_id": types.StringType,
			"owner":       types.StringType,
			"parent_ak":   types.StringType,
			"status":      types.StringType,
			"update_at":   types.Int64Type,
		}, map[string]attr.Value{
			"ak":          types.StringValue(""),
			"sk":          types.StringValue(""),
			"external_id": types.StringValue(""),
			"owner":       types.StringValue(""),
			"parent_ak":   types.StringValue(""),
			"status":      types.StringValue(""),
			"update_at":   types.Int64Value(0),
		})
		if diags.HasError() {
			resp.Diagnostics.Append(diags...)
			return
		}
		state.AK = akObject
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *customRegionResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Custom region does not support update operation
	resp.Diagnostics.AddError(
		"Error updating custom region",
		"Custom region does not support update operation",
	)
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *customRegionResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state customRegionResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete existing custom region using DeleteByPost
	// Some APIs require POST for delete operations
	if err := r.client.DeleteByPost(consts.TypeNameCustomRegion, state.UUID.ValueString(), nil, nil); err != nil {
		resp.Diagnostics.AddError(
			"Error deleting custom region",
			"Could not delete custom region, unexpected error: "+err.Error(),
		)
		return
	}
}

func (r *customRegionResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("uuid"), req, resp)
}

func (r *customRegionResource) getCustomRegionFromPlan(plan *customRegionResourceModel) *api.CustomRegion {
	return &api.CustomRegion{
		Internal: plan.Internal.ValueBool(),
		ISP:      plan.ISP.ValueString(),
		Country:  plan.Country.ValueString(),
		Province: plan.Province.ValueString(),
		City:     plan.City.ValueString(),
		Name:     plan.Name.ValueString(),
		Company:  plan.Company.ValueString(),
		Keycode:  plan.Keycode.ValueString(),
	}
}
