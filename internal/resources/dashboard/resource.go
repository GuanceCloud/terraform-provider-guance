package dashboard

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
	_ resource.Resource                = &dashboardResource{}
	_ resource.ResourceWithConfigure   = &dashboardResource{}
	_ resource.ResourceWithImportState = &dashboardResource{}
)

// NewDashboardResource creates a new dashboard resource
func NewDashboardResource() resource.Resource {
	return &dashboardResource{}
}

// dashboardResource is the resource implementation.
type dashboardResource struct {
	client *api.Client
}

// Schema defines the schema for the data source.
func (r *dashboardResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = resourceSchema
}

// Configure adds the provider configured client to the data source.
func (r *dashboardResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.client = req.ProviderData.(*api.Client)
}

// Metadata returns the data source type name.
func (r *dashboardResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_dashboard"
}

// Create creates the resource and sets the initial Terraform state.
func (r *dashboardResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan dashboardResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	item := r.getDashboardFromPlan(&plan)
	rb, _ := json.Marshal(item)
	tflog.Debug(ctx, fmt.Sprintf("============= body: %s", string(rb)))
	content := &api.DashboardContent{}
	err := r.client.Create(consts.TypeNameDashboard, item, content)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating dashboard",
			"Could not create dashboard, unexpected error: "+err.Error(),
		)
		return
	}

	plan.UUID = types.StringValue(content.UUID)
	plan.WorkspaceUUID = types.StringValue(content.WorkspaceUUID)
	plan.CreateAt = types.Int64Value(int64(content.CreateAt))
	plan.UpdateAt = types.Int64Value(int64(content.UpdateAt))

	// Export dashboard to get template_info and template_info_export
	templateInfoExportStr, err := r.exportDashboard(ctx, plan.UUID.ValueString())
	if err == nil {
		plan.TemplateInfoExport = types.StringValue(templateInfoExportStr)
	}

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *dashboardResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state dashboardResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	content := &api.DashboardContent{}
	err := r.client.Read(consts.TypeNameDashboard, state.UUID.ValueString(), content)
	if err != nil {
		if err == api.Error404 {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error reading dashboard",
			"Could not read dashboard, unexpected error: "+err.Error(),
		)
		return
	}

	state.UUID = types.StringValue(content.UUID)
	state.Name = types.StringValue(content.Name)
	state.Desc = types.StringValue(content.Desc)
	state.Identifier = types.StringValue(content.Identifier)
	state.CreateAt = types.Int64Value(int64(content.CreateAt))
	state.UpdateAt = types.Int64Value(int64(content.UpdateAt))
	state.WorkspaceUUID = types.StringValue(content.WorkspaceUUID)
	state.IsPublic = types.Int64Value(int64(content.IsPublic))

	// Export dashboard to get template_info_export
	templateInfoExportStr, err := r.exportDashboard(ctx, state.UUID.ValueString())
	if err == nil {
		// If template_info_export has changed, update both template_info and template_info_export
		if templateInfoExportStr != state.TemplateInfoExport.ValueString() {
			state.TemplateInfo = types.StringValue(templateInfoExportStr)
		}
		// Update template_info_export
		state.TemplateInfoExport = types.StringValue(templateInfoExportStr)
	} else {
		if err == api.Error404 {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error exporting dashboard",
			"Could not export dashboard, unexpected error: "+err.Error(),
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
func (r *dashboardResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan and state
	var plan dashboardResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state dashboardResourceModel
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Check if template_info has changed
	templateInfoChanged := !plan.TemplateInfo.IsNull() && (state.TemplateInfo.IsNull() || plan.TemplateInfo.ValueString() != state.TemplateInfo.ValueString())

	// If template_info has changed, import the template
	if templateInfoChanged {
		// Import the template
		var templateInfo interface{}
		if err := json.Unmarshal([]byte(plan.TemplateInfo.ValueString()), &templateInfo); err == nil {
			importRequest := map[string]interface{}{
				"templateInfo": templateInfo,
			}

			if !plan.SpecifyDashboardUUID.IsNull() {
				importRequest["specifyDashboardUUID"] = plan.SpecifyDashboardUUID.ValueString()
			}

			rb, _ := json.Marshal(importRequest)
			tflog.Debug(ctx, fmt.Sprintf("============= import body: %s", string(rb)))

			importContent := &api.DashboardContent{}
			err := r.client.Import(consts.TypeNameDashboard, plan.UUID.ValueString(), importRequest, importContent)
			if err != nil {
				resp.Diagnostics.AddError(
					"Error importing dashboard template",
					"Could not import dashboard template, unexpected error: "+err.Error(),
				)
				return
			}

			plan.UpdateAt = types.Int64Value(int64(importContent.UpdateAt))
		}
	}

	// Always update dashboard normally first
	item := r.getDashboardFromPlan(&plan)
	item.TemplateInfo = nil

	rb, _ := json.Marshal(item)
	tflog.Debug(ctx, fmt.Sprintf("============= body: %s", string(rb)))

	content := &api.DashboardContent{}
	err := r.client.Update(consts.TypeNameDashboard, plan.UUID.ValueString(), item, content)

	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating dashboard",
			"Could not update dashboard, unexpected error: "+err.Error(),
		)
		return
	}

	plan.UpdateAt = types.Int64Value(int64(content.UpdateAt))

	// Export dashboard to get template_info_export
	templateInfoExportStr, err := r.exportDashboard(ctx, plan.UUID.ValueString())
	if err == nil {
		tflog.Debug(ctx, fmt.Sprintf("============= export body: %s", templateInfoExportStr))
		// Only update template_info_export, keep template_info as user set
		plan.TemplateInfoExport = types.StringValue(templateInfoExportStr)
	} else {
		resp.Diagnostics.AddError(
			"Error exporting dashboard after import",
			"Could not export dashboard after import, unexpected error: "+err.Error(),
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
func (r *dashboardResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state dashboardResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete existing dashboard
	if err := r.client.DeleteByPost(
		consts.TypeNameDashboard,
		state.UUID.ValueString(),
		nil,
		nil,
	); err != nil {
		resp.Diagnostics.AddError(
			"Error deleting dashboard",
			"Could not delete dashboard, unexpected error: "+err.Error(),
		)
		return
	}
}

func (r *dashboardResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("uuid"), req, resp)
}

// exportDashboard exports the dashboard and updates the template_info_export field
func (r *dashboardResource) exportDashboard(ctx context.Context, uuid string) (string, error) {
	var exportContent interface{}
	err := r.client.Export(consts.TypeNameDashboard, uuid, &exportContent)
	if err != nil {
		return "", err
	}

	templateInfoBytes, err := json.Marshal(exportContent)
	if err != nil {
		return "", err
	}

	return string(templateInfoBytes), nil
}

func (r *dashboardResource) getDashboardFromPlan(plan *dashboardResourceModel) *api.Dashboard {

	d := &api.Dashboard{
		Name: plan.Name.ValueString(),
	}

	if !plan.Desc.IsNull() {
		d.Desc = plan.Desc.ValueString()
	}

	if !plan.Identifier.IsNull() {
		d.Identifier = plan.Identifier.ValueString()
	}

	if len(plan.TagNames) > 0 {
		tagNames := make([]string, len(plan.TagNames))
		for i, tag := range plan.TagNames {
			tagNames[i] = tag.ValueString()
		}
		d.TagNames = tagNames
	}

	if !plan.TemplateInfo.IsNull() {
		var templateInfo interface{}
		if err := json.Unmarshal([]byte(plan.TemplateInfo.ValueString()), &templateInfo); err == nil {
			d.TemplateInfo = templateInfo
		}
	}

	if !plan.SpecifyDashboardUUID.IsNull() {
		d.SpecifyDashboardUUID = plan.SpecifyDashboardUUID.ValueString()
	}

	if !plan.IsPublic.IsNull() {
		d.IsPublic = int(plan.IsPublic.ValueInt64())
		if d.IsPublic == -1 { // custom permission
			if len(plan.PermissionSet) > 0 {
				permissionSet := make([]string, len(plan.PermissionSet))
				for i, perm := range plan.PermissionSet {
					permissionSet[i] = perm.ValueString()
				}
				d.PermissionSet = permissionSet
			}

			if len(plan.ReadPermissionSet) > 0 {
				readPermissionSet := make([]string, len(plan.ReadPermissionSet))
				for i, perm := range plan.ReadPermissionSet {
					readPermissionSet[i] = perm.ValueString()
				}
				d.ReadPermissionSet = readPermissionSet
			}
		}
	}

	return d
}
