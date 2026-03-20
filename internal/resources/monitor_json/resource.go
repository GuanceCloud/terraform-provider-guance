package monitor_json

import (
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/GuanceCloud/terraform-provider-guance/internal/api"
	"github.com/GuanceCloud/terraform-provider-guance/internal/consts"
)

//go:embed README.md
var resourceDocument string

var (
	_ resource.Resource                = &monitorJsonResource{}
	_ resource.ResourceWithConfigure   = &monitorJsonResource{}
	_ resource.ResourceWithImportState = &monitorJsonResource{}
)

func NewMonitorJsonResource() resource.Resource {
	return &monitorJsonResource{}
}

type monitorJsonResource struct {
	client *api.Client
}

func (r *monitorJsonResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = resourceSchema
}

func (r *monitorJsonResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.client = req.ProviderData.(*api.Client)
}

func (r *monitorJsonResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_monitor_json"
}

func (r *monitorJsonResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan monitorJsonResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if plan.CheckerJson.IsNull() {
		resp.Diagnostics.AddError(
			"Missing checker_json",
			"checker_json is required for creating a monitor",
		)
		return
	}

	// Parse the checker_json as a single object (not an array)
	checker, err := r.parseCheckerJson(plan.CheckerJson.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Invalid checker_json format",
			"Could not parse checker_json, please ensure it's a valid JSON object: "+err.Error(),
		)
		return
	}

	// For Import API, we need to wrap the single checker in an array
	checkers := []interface{}{checker}

	importRequest := &api.MonitorJsonImportRequest{
		Checkers:             checkers,
		SkipRepeatNameCreate: true,
	}

	if !plan.Type.IsNull() {
		importRequest.Type = plan.Type.ValueString()
	}

	importContent := &api.MonitorJsonImportContent{}
	err = r.client.Import(consts.TypeNameMonitorJson, "", importRequest, importContent)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error importing monitor json",
			"Could not import monitor json, unexpected error: "+err.Error(),
		)
		return
	}

	// Use the first monitor's UUID and timestamps if available
	if len(importContent.Rule) == 1 {
		firstMonitor := importContent.Rule[0]
		plan.UUID = types.StringValue(firstMonitor.UUID)
		plan.CreateAt = types.Int64Value(int64(firstMonitor.CreateAt))
		plan.UpdateAt = types.Int64Value(int64(firstMonitor.UpdateAt))
		plan.WorkspaceUUID = types.StringValue(firstMonitor.WorkspaceUUID)
	} else if len(importContent.Rule) == 0 {
		resp.Diagnostics.AddError(
			"No monitors imported",
			"No monitors were imported. This may be due to a name conflict, invalid monitor configuration, or missing required fields.",
		)
		return
	} else {
		resp.Diagnostics.AddError(
			"Multiple monitors in configuration",
			"The checker_json must contain exactly one monitor configuration, but got "+fmt.Sprintf("%d", len(importContent.Rule)),
		)
		return
	}

	// Export the monitor to get the latest configuration
	checkerJson, err := r.exportMonitorJson(ctx, plan.UUID.ValueString())
	if err == nil {
		// Set the exported checker json
		plan.CheckerJsonExport = types.StringValue(checkerJson)
	} else {
		resp.Diagnostics.AddError(
			"Error exporting monitor configuration",
			"Could not export the created monitor to get the latest configuration: "+err.Error(),
		)
		return
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *monitorJsonResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state monitorJsonResourceModel
	var err error
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	checkerJson, err := r.exportMonitorJson(ctx, state.UUID.ValueString())
	if err == nil {
		if checkerJson != state.CheckerJsonExport.ValueString() {
			// Set the exported checker json
			state.CheckerJsonExport = types.StringValue(checkerJson)
			state.CheckerJson = types.StringValue(checkerJson)
		}
	} else {
		resp.State.RemoveResource(ctx)
		return
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *monitorJsonResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan monitorJsonResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state monitorJsonResourceModel
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	checkerJsonChanged := !plan.CheckerJson.IsNull() && (state.CheckerJson.IsNull() || plan.CheckerJson.ValueString() != state.CheckerJson.ValueString())

	if checkerJsonChanged {
		// Parse the checker_json as a single object (not an array)
		checker, err := r.parseCheckerJson(plan.CheckerJson.ValueString())
		if err != nil {
			resp.Diagnostics.AddError(
				"Invalid checker_json format",
				"Could not parse checker_json, please ensure it's a valid JSON object: "+err.Error(),
			)
			return
		}

		// For Import API, we need to wrap the single checker in an array
		checkers := []interface{}{checker}

		importRequest := &api.MonitorJsonImportRequest{
			Checkers: checkers,
		}

		if !plan.Type.IsNull() {
			importRequest.Type = plan.Type.ValueString()
		}

		importContent := &api.MonitorJsonImportContent{}
		// Use Replace API for update
		replaceRequest := &api.MonitorJsonReplaceRequest{
			Checker: checker,
		}

		err = r.client.Replace(consts.TypeNameMonitorJson, state.UUID.ValueString(), replaceRequest, importContent)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error updating monitor",
				"Could not update the monitor using the Replace API: "+err.Error(),
			)
			return
		}

		plan.UpdateAt = types.Int64Value(time.Now().Unix())

		// Use the first monitor's update timestamp if available
		if len(importContent.Rule) > 0 {
			firstMonitor := importContent.Rule[0]
			if firstMonitor.UpdateAt > 0 {
				plan.UpdateAt = types.Int64Value(int64(firstMonitor.UpdateAt))
			}
		}
	}

	// Export the latest configuration to update checker_json and checker_json_export
	checkerJson, err := r.exportMonitorJson(ctx, state.UUID.ValueString())
	if err == nil {
		plan.CheckerJsonExport = types.StringValue(checkerJson)
	} else {
		resp.Diagnostics.AddError(
			"Error exporting monitor configuration",
			"Could not export the monitor to get the latest configuration: "+err.Error(),
		)
		return
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *monitorJsonResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state monitorJsonResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Use the UUID from the state directly
	if !state.UUID.IsNull() {
		uuid := state.UUID.ValueString()
		if err := r.client.DeleteMonitor(uuid); err != nil {
			resp.Diagnostics.AddError(
				"Error deleting monitor",
				"Could not delete monitor, unexpected error: "+err.Error(),
			)
			return
		}
	}
}

func (r *monitorJsonResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("uuid"), req, resp)
}

// exportMonitorJson exports the monitor and returns the checker JSON
func (r *monitorJsonResource) exportMonitorJson(ctx context.Context, uuid string) (string, error) {
	uuids := []string{uuid}
	exportRequest := &api.MonitorJsonExportRequest{
		Checkers: uuids,
	}
	exportContent := &api.MonitorJsonContent{}
	err := r.client.ExportWithBody(consts.TypeNameMonitorJson, exportRequest, exportContent)
	if err != nil {
		return "", err
	}

	if len(exportContent.Checkers) != 1 {
		return "", fmt.Errorf("expected 1 monitor, got %d", len(exportContent.Checkers))
	}

	checkerJsonBytes, err := json.Marshal(exportContent.Checkers[0])
	if err != nil {
		return "", err
	}

	return string(checkerJsonBytes), nil
}

// parseCheckerJson parses the checker JSON and returns the checker object
func (r *monitorJsonResource) parseCheckerJson(checkerJson string) (interface{}, error) {
	var checker interface{}
	if err := json.Unmarshal([]byte(checkerJson), &checker); err != nil {
		return nil, err
	}
	return checker, nil
}
