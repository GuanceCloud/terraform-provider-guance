package mute

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
	_ resource.Resource                = &muteResource{}
	_ resource.ResourceWithConfigure   = &muteResource{}
	_ resource.ResourceWithImportState = &muteResource{}
)

func NewMuteResource() resource.Resource {
	return &muteResource{}
}

type muteResource struct {
	client *api.Client
}

func (r *muteResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = resourceSchema
}

func (r *muteResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.client = req.ProviderData.(*api.Client)
}

func (r *muteResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_mute"
}

func (r *muteResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan muteResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	content := &api.MuteContent{}
	if err := r.client.Create(consts.TypeNameMute, muteFromPlan(&plan), content); err != nil {
		resp.Diagnostics.AddError(
			"Error creating mute rule",
			"Could not create mute rule, unexpected error: "+err.Error(),
		)
		return
	}

	applyContentToState(&plan, content)
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *muteResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state muteResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	content := &api.MuteContent{}
	if err := r.client.GetMute(state.UUID.ValueString(), content); err != nil {
		if err == api.Error404 {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error reading mute rule",
			"Could not read mute rule, unexpected error: "+err.Error(),
		)
		return
	}

	applyContentToState(&state, content)
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

func (r *muteResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan muteResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	content := &api.MuteContent{}
	if err := r.client.Update(consts.TypeNameMute, plan.UUID.ValueString(), muteFromPlan(&plan), content); err != nil {
		resp.Diagnostics.AddError(
			"Error updating mute rule",
			"Could not update mute rule, unexpected error: "+err.Error(),
		)
		return
	}

	applyContentToState(&plan, content)
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *muteResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state muteResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.DeleteByPost(consts.TypeNameMute, state.UUID.ValueString(), nil, nil); err != nil {
		resp.Diagnostics.AddError(
			"Error deleting mute rule",
			"Could not delete mute rule, unexpected error: "+err.Error(),
		)
	}
}

func (r *muteResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("uuid"), req, resp)
}

func muteFromPlan(plan *muteResourceModel) *api.Mute {
	item := &api.Mute{
		Name:          plan.Name.ValueString(),
		Type:          plan.Type.ValueString(),
		MuteRanges:    muteRangesFromPlan(plan.MuteRanges),
		RepeatTimeSet: int(plan.RepeatTimeSet.ValueInt64()),
		Timezone:      plan.Timezone.ValueString(),
	}

	if !plan.Description.IsNull() {
		item.Description = plan.Description.ValueString()
	}
	if plan.Tags != nil {
		item.Tags = plan.Tags
	}
	if !plan.FilterString.IsNull() {
		item.FilterString = plan.FilterString.ValueString()
	}
	if len(plan.NotifyTargets) > 0 {
		item.NotifyTargets = notifyTargetsFromPlan(plan.NotifyTargets)
	}
	if !plan.NotifyMessage.IsNull() {
		item.NotifyMessage = plan.NotifyMessage.ValueString()
	}
	if !plan.NotifyTimeStr.IsNull() {
		item.NotifyTimeStr = plan.NotifyTimeStr.ValueString()
	}
	if !plan.StartTime.IsNull() {
		item.StartTime = plan.StartTime.ValueString()
	}
	if !plan.EndTime.IsNull() {
		item.EndTime = plan.EndTime.ValueString()
	}
	if plan.RepeatCrontabSet != nil {
		item.RepeatCrontabSet = repeatCrontabSetFromPlan(plan.RepeatCrontabSet)
	}
	if !plan.CrontabDuration.IsNull() {
		item.CrontabDuration = int(plan.CrontabDuration.ValueInt64())
	}
	if !plan.RepeatExpireTime.IsNull() {
		item.RepeatExpireTime = plan.RepeatExpireTime.ValueString()
	}
	if plan.Declaration != nil {
		item.Declaration = plan.Declaration
	}

	return item
}

func applyContentToState(state *muteResourceModel, content *api.MuteContent) {
	state.UUID = types.StringValue(content.UUID)
	state.Name = types.StringValue(content.Name)
	state.Description = stringValueOrExisting(content.Description, state.Description)
	state.Type = types.StringValue(content.Type)
	state.MuteRanges = muteRangesFromContent(content.MuteRanges, state.MuteRanges)
	if len(content.Tags) > 0 {
		state.Tags = content.Tags
	}
	state.FilterString = stringValueOrExisting(content.FilterString, state.FilterString)
	state.NotifyTargets = notifyTargetsFromContent(content.NotifyTargets, state.NotifyTargets)
	state.NotifyMessage = stringValueOrExisting(content.NotifyMessage, state.NotifyMessage)
	state.NotifyTimeStr = stringValueOrExisting(content.NotifyTimeStr, state.NotifyTimeStr)
	repeatTimeSet := repeatTimeSetFromContent(content)
	if repeatTimeSet != 0 || state.RepeatTimeSet.IsNull() || state.RepeatTimeSet.IsUnknown() {
		state.RepeatTimeSet = types.Int64Value(int64(repeatTimeSet))
	}
	if state.RepeatTimeSet.ValueInt64() == 1 {
		state.StartTime = stringValueOrConfigured(content.StartTime, state.StartTime)
		state.EndTime = stringValueOrConfigured(content.EndTime, state.EndTime)
	} else {
		state.StartTime = stringValueOrExisting(content.StartTime, state.StartTime)
		state.EndTime = stringValueOrExisting(content.EndTime, state.EndTime)
	}
	if content.RepeatCrontabSet != nil {
		state.RepeatCrontabSet = repeatCrontabSetFromContent(content.RepeatCrontabSet)
	}
	if content.CrontabDuration != 0 || !state.CrontabDuration.IsNull() {
		state.CrontabDuration = types.Int64Value(int64(content.CrontabDuration))
	}
	state.RepeatExpireTime = repeatExpireTimeValueOrExisting(content.RepeatExpireTime, state.RepeatExpireTime)
	state.Timezone = stringValueOrExisting(content.Timezone, state.Timezone)
	if len(content.Declaration) > 0 && state.Declaration != nil {
		state.Declaration = declarationFromContent(content.Declaration, state.Declaration)
	}
	state.Status = types.Int64Value(int64(content.Status))
	state.CreateAt = types.Int64Value(int64(content.CreateAt))
	state.UpdateAt = types.Int64Value(int64(content.UpdateAt))
	state.WorkspaceUUID = types.StringValue(content.WorkspaceUUID)
}

func stringValueOrExisting(value string, existing types.String) types.String {
	if value == "" {
		return existing
	}
	return types.StringValue(value)
}

func stringValueOrConfigured(value string, existing types.String) types.String {
	if existing.IsNull() || existing.IsUnknown() {
		return existing
	}
	return stringValueOrExisting(value, existing)
}

func repeatExpireTimeValueOrExisting(value string, existing types.String) types.String {
	if value == "" || value == "-1" {
		return existing
	}
	return types.StringValue(value)
}

func muteRangesFromPlan(values []muteRange) []api.MuteRange {
	result := make([]api.MuteRange, 0, len(values))
	for _, value := range values {
		result = append(result, api.MuteRange{
			Name:            value.Name.ValueString(),
			Type:            value.Type.ValueString(),
			CheckerUUID:     value.CheckerUUID.ValueString(),
			MonitorUUID:     value.MonitorUUID.ValueString(),
			SLOUUID:         value.SLOUUID.ValueString(),
			AlertPolicyUUID: value.AlertPolicyUUID.ValueString(),
			TagUUID:         value.TagUUID.ValueString(),
		})
	}
	return result
}

func muteRangesFromContent(values []api.MuteRange, existing []muteRange) []muteRange {
	if len(values) == 0 && len(existing) > 0 {
		return existing
	}
	result := make([]muteRange, 0, len(values))
	for _, value := range values {
		next := muteRange{
			Name:            types.StringValue(value.Name),
			Type:            stringPointerValue(value.Type),
			CheckerUUID:     stringPointerValue(value.CheckerUUID),
			MonitorUUID:     stringPointerValue(value.MonitorUUID),
			SLOUUID:         stringPointerValue(value.SLOUUID),
			AlertPolicyUUID: stringPointerValue(value.AlertPolicyUUID),
			TagUUID:         stringPointerValue(value.TagUUID),
		}
		if len(existing) > len(result) {
			next.Name = stringValueOrExisting(value.Name, existing[len(result)].Name)
			next.Type = stringValueOrExisting(value.Type, existing[len(result)].Type)
			next.CheckerUUID = stringValueOrExisting(value.CheckerUUID, existing[len(result)].CheckerUUID)
			next.MonitorUUID = stringValueOrExisting(value.MonitorUUID, existing[len(result)].MonitorUUID)
			next.SLOUUID = stringValueOrExisting(value.SLOUUID, existing[len(result)].SLOUUID)
			next.AlertPolicyUUID = stringValueOrExisting(value.AlertPolicyUUID, existing[len(result)].AlertPolicyUUID)
			next.TagUUID = stringValueOrExisting(value.TagUUID, existing[len(result)].TagUUID)
		}
		result = append(result, next)
	}
	return result
}

func stringPointerValue(value string) types.String {
	if value == "" {
		return types.StringNull()
	}
	return types.StringValue(value)
}

func notifyTargetsFromPlan(values []notifyTarget) []api.MuteNotifyTarget {
	result := make([]api.MuteNotifyTarget, 0, len(values))
	for _, value := range values {
		result = append(result, api.MuteNotifyTarget{
			To:   stringsFromTypes(value.To),
			Type: value.Type.ValueString(),
		})
	}
	return result
}

func notifyTargetsFromContent(values []api.MuteNotifyTarget, existing []notifyTarget) []notifyTarget {
	if len(values) == 0 {
		return existing
	}
	result := make([]notifyTarget, 0, len(values))
	for _, value := range values {
		result = append(result, notifyTarget{
			To:   typesFromStrings(value.To),
			Type: types.StringValue(value.Type),
		})
	}
	return result
}

func repeatCrontabSetFromPlan(value *repeatCrontabSet) *api.RepeatCrontabSet {
	return &api.RepeatCrontabSet{
		Min:   value.Min.ValueString(),
		Hour:  value.Hour.ValueString(),
		Day:   value.Day.ValueString(),
		Month: value.Month.ValueString(),
		Week:  value.Week.ValueString(),
	}
}

func repeatCrontabSetFromContent(value *api.RepeatCrontabSet) *repeatCrontabSet {
	return &repeatCrontabSet{
		Min:   types.StringValue(value.Min),
		Hour:  types.StringValue(value.Hour),
		Day:   types.StringValue(value.Day),
		Month: types.StringValue(value.Month),
		Week:  types.StringValue(value.Week),
	}
}

func repeatTimeSetFromContent(content *api.MuteContent) int {
	if content.RepeatTimeSet != 0 {
		return content.RepeatTimeSet
	}
	if content.RepeatCrontabSet != nil {
		return 1
	}
	return 0
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

func declarationFromContent(values map[string]any, existing map[string]string) map[string]string {
	result := make(map[string]string, len(values))
	for key, value := range values {
		stringValue, ok := value.(string)
		if !ok {
			continue
		}
		result[key] = stringValue
	}
	if len(result) == 0 {
		return existing
	}
	return result
}
