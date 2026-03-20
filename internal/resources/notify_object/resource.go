package notify_object

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
	_ resource.Resource                = &notifyObjectResource{}
	_ resource.ResourceWithConfigure   = &notifyObjectResource{}
	_ resource.ResourceWithImportState = &notifyObjectResource{}
)

// NewNotifyObjectResource creates a new notify object resource
func NewNotifyObjectResource() resource.Resource {
	return &notifyObjectResource{}
}

// notifyObjectResource is the resource implementation.
type notifyObjectResource struct {
	client *api.Client
}

// Schema defines the schema for the data source.
func (r *notifyObjectResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = resourceSchema
}

// Configure adds the provider configured client to the data source.
func (r *notifyObjectResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.client = req.ProviderData.(*api.Client)
}

// Metadata returns the data source type name.
func (r *notifyObjectResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_notify_object"
}

// Create creates the resource and sets the initial Terraform state.
func (r *notifyObjectResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan notifyObjectResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	item := r.getNotifyObjectFromPlan(&plan)
	rb, _ := json.Marshal(item)
	tflog.Debug(ctx, fmt.Sprintf("============= body: %s", string(rb)))
	content := &api.NotifyObjectContent{}
	err := r.client.Create(consts.TypeNameNotifyObject, item, content)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating notify object",
			"Could not create notify object, unexpected error: "+err.Error(),
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
func (r *notifyObjectResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state notifyObjectResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	content := &api.NotifyObjectContent{}
	err := r.client.GetNotifyObject(state.UUID.ValueString(), content)
	if err != nil {
		if err == api.Error404 {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error reading notify object",
			"Could not read notify object, unexpected error: "+err.Error(),
		)
		return
	}

	state.UUID = types.StringValue(content.UUID)
	state.Type = types.StringValue(content.Type)
	state.Name = types.StringValue(content.Name)
	if content.OptSet != nil {
		optSetBytes, _ := json.Marshal(content.OptSet)
		state.OptSet = types.StringValue(string(optSetBytes))
	}
	state.OpenPermissionSet = types.BoolValue(content.OpenPermissionSet)
	if len(content.PermissionSet) > 0 {
		permissionSet := make([]types.String, len(content.PermissionSet))
		for i, perm := range content.PermissionSet {
			permissionSet[i] = types.StringValue(perm)
		}
		state.PermissionSet = permissionSet
	}
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
func (r *notifyObjectResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan notifyObjectResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	item := r.getNotifyObjectFromPlan(&plan)
	// For notify object, the UUID is in the body, not the path
	itemWithUUID := map[string]interface{}{
		"notifyObjectUUID":  plan.UUID.ValueString(),
		"name":              item.Name,
		"optSet":            item.OptSet,
		"openPermissionSet": item.OpenPermissionSet,
		"permissionSet":     item.PermissionSet,
	}

	content := &api.NotifyObjectContent{}
	err := r.client.UpdateNotifyObject(itemWithUUID, content)

	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating notify object",
			"Could not update notify object, unexpected error: "+err.Error(),
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
func (r *notifyObjectResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state notifyObjectResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteNotifyObject(state.UUID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting notify object",
			"Could not delete notify object, unexpected error: "+err.Error(),
		)
		return
	}
}

func (r *notifyObjectResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("uuid"), req, resp)
}

func (r *notifyObjectResource) getNotifyObjectFromPlan(plan *notifyObjectResourceModel) *api.NotifyObject {

	n := &api.NotifyObject{
		Type: plan.Type.ValueString(),
		Name: plan.Name.ValueString(),
	}

	if !plan.OptSet.IsNull() {
		var optSet interface{}
		if err := json.Unmarshal([]byte(plan.OptSet.ValueString()), &optSet); err == nil {
			n.OptSet = optSet
		}
	}

	if !plan.OpenPermissionSet.IsNull() {
		n.OpenPermissionSet = plan.OpenPermissionSet.ValueBool()
	}

	if len(plan.PermissionSet) > 0 {
		permissionSet := make([]string, len(plan.PermissionSet))
		for i, perm := range plan.PermissionSet {
			permissionSet[i] = perm.ValueString()
		}
		n.PermissionSet = permissionSet
	}

	return n
}
