package alert_policy_notice_date

import (
	"context"
	_ "embed"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/GuanceCloud/terraform-provider-guance/internal/api"
	"github.com/GuanceCloud/terraform-provider-guance/internal/consts"
)

//go:embed README.md
var resourceDocument string

var (
	_ resource.Resource                = &alertPolicyNoticeDateResource{}
	_ resource.ResourceWithConfigure   = &alertPolicyNoticeDateResource{}
	_ resource.ResourceWithImportState = &alertPolicyNoticeDateResource{}
)

func NewAlertPolicyNoticeDateResource() resource.Resource {
	return &alertPolicyNoticeDateResource{}
}

type alertPolicyNoticeDateResource struct {
	client *api.Client
}

func (r *alertPolicyNoticeDateResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = resourceSchema
}

func (r *alertPolicyNoticeDateResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.client = req.ProviderData.(*api.Client)
}

func (r *alertPolicyNoticeDateResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_alert_policy_notice_date"
}

func (r *alertPolicyNoticeDateResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan alertPolicyNoticeDateResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	content := &api.AlertPolicyNoticeDateContent{}
	if err := r.client.Create(consts.TypeNameAlertPolicyNoticeDate, noticeDateFromPlan(&plan), content); err != nil {
		resp.Diagnostics.AddError(
			"Error creating alert policy notice date",
			"Could not create alert policy notice date, unexpected error: "+err.Error(),
		)
		return
	}

	applyContentToState(&plan, content)
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *alertPolicyNoticeDateResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state alertPolicyNoticeDateResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	content := &api.AlertPolicyNoticeDateContent{}
	if err := r.client.Read(consts.TypeNameAlertPolicyNoticeDate, state.UUID.ValueString(), content); err != nil {
		if err == api.Error404 {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error reading alert policy notice date",
			"Could not read alert policy notice date, unexpected error: "+err.Error(),
		)
		return
	}

	applyContentToState(&state, content)
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

func (r *alertPolicyNoticeDateResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan alertPolicyNoticeDateResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	content := &api.AlertPolicyNoticeDateContent{}
	if err := r.client.Update(consts.TypeNameAlertPolicyNoticeDate, plan.UUID.ValueString(), noticeDateFromPlan(&plan), content); err != nil {
		resp.Diagnostics.AddError(
			"Error updating alert policy notice date",
			"Could not update alert policy notice date, unexpected error: "+err.Error(),
		)
		return
	}

	applyContentToState(&plan, content)
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *alertPolicyNoticeDateResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state alertPolicyNoticeDateResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	body := noticeDateDeleteBody(&state)
	if err := r.client.DeleteByPost(consts.TypeNameAlertPolicyNoticeDate, "", body, nil); err != nil {
		resp.Diagnostics.AddError(
			"Error deleting alert policy notice date",
			"Could not delete alert policy notice date, unexpected error: "+err.Error(),
		)
	}
}

func noticeDateDeleteBody(state *alertPolicyNoticeDateResourceModel) map[string]any {
	return map[string]any{
		"noticeDatesUUIDs": []string{state.UUID.ValueString()},
		"skipRefCheck":     boolValueOrDefault(state.SkipRefCheckOnDelete, true),
	}
}

func (r *alertPolicyNoticeDateResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("uuid"), req, resp)
}

func noticeDateFromPlan(plan *alertPolicyNoticeDateResourceModel) *api.AlertPolicyNoticeDate {
	return &api.AlertPolicyNoticeDate{
		Name:        plan.Name.ValueString(),
		NoticeDates: stringsFromTypes(plan.NoticeDates),
	}
}

func applyContentToState(state *alertPolicyNoticeDateResourceModel, content *api.AlertPolicyNoticeDateContent) {
	state.UUID = types.StringValue(content.UUID)
	state.Name = types.StringValue(content.Name)
	state.NoticeDates = typesFromStrings(content.Dates)
	state.CreateAt = types.Int64Value(int64(content.CreateAt))
	state.UpdateAt = types.Int64Value(int64(content.UpdateAt))
	state.WorkspaceUUID = types.StringValue(content.WorkspaceUUID)
}

func stringsFromTypes(values []types.String) []string {
	result := make([]string, 0, len(values))
	for _, value := range values {
		if value.IsNull() || value.IsUnknown() {
			continue
		}
		result = append(result, value.ValueString())
	}
	return result
}

func typesFromStrings(values []string) []types.String {
	result := make([]types.String, 0, len(values))
	for _, value := range values {
		result = append(result, types.StringValue(value))
	}
	return result
}

func boolValueOrDefault(value types.Bool, defaultValue bool) bool {
	if value.IsNull() || value.IsUnknown() {
		return defaultValue
	}
	return value.ValueBool()
}
