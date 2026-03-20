package alert_policy

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

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &alertPolicyResource{}
	_ resource.ResourceWithConfigure   = &alertPolicyResource{}
	_ resource.ResourceWithImportState = &alertPolicyResource{}
)

// NewAlertPolicyResource creates a new alert policy resource
func NewAlertPolicyResource() resource.Resource {
	return &alertPolicyResource{}
}

// alertPolicyResource is the resource implementation.
type alertPolicyResource struct {
	client *api.Client
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
	r.client = req.ProviderData.(*api.Client)
}

// Metadata returns the data source type name.
func (r *alertPolicyResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_alert_policy"
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

	item := r.getAlertPolicyFromPlan(&plan)
	content := &api.AlertPolicyContent{}
	err := r.client.Create(consts.TypeNameAlertPolicy, item, content)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating alert policy",
			"Could not create alert policy, unexpected error: "+err.Error(),
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
func (r *alertPolicyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state alertPolicyResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	content := &api.AlertPolicyContent{}
	err := r.client.Read(consts.TypeNameAlertPolicy, state.UUID.ValueString(), content)
	if err != nil {
		if err == api.Error404 {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error reading alert policy",
			"Could not read alert policy, unexpected error: "+err.Error(),
		)
		return
	}

	state.UUID = types.StringValue(content.UUID)
	state.Name = types.StringValue(content.Name)
	state.Desc = types.StringValue(content.Desc)
	state.OpenPermissionSet = types.BoolValue(content.OpenPermissionSet)

	if content.PermissionSet != nil {
		state.PermissionSet = make([]types.String, len(content.PermissionSet))
		for i, perm := range content.PermissionSet {
			state.PermissionSet[i] = types.StringValue(perm)
		}
	}

	if content.CheckerUUIDs != nil {
		state.CheckerUUIDs = make([]types.String, len(content.CheckerUUIDs))
		for i, checker := range content.CheckerUUIDs {
			state.CheckerUUIDs[i] = types.StringValue(checker)
		}
	}
	if content.SecurityRuleUUIDs != nil {
		state.SecurityRuleUUIDs = make([]types.String, len(content.SecurityRuleUUIDs))
		for i, rule := range content.SecurityRuleUUIDs {
			state.SecurityRuleUUIDs[i] = types.StringValue(rule)
		}
	}
	state.RuleTimezone = types.StringValue(content.RuleTimezone)
	// TODO: Map alertOpt from content to state
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
func (r *alertPolicyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan alertPolicyResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	item := r.getAlertPolicyFromPlan(&plan)

	content := &api.AlertPolicyContent{}
	err := r.client.Update(consts.TypeNameAlertPolicy, plan.UUID.ValueString(), item, content)

	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating alert policy",
			"Could not update alert policy, unexpected error: "+err.Error(),
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
func (r *alertPolicyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state alertPolicyResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	body := map[string][]string{
		"alertPolicyUUIDs": {state.UUID.ValueString()},
	}

	// Delete existing alert policy
	if err := r.client.DeleteByPost(
		consts.TypeNameAlertPolicy,
		"",
		body,
		nil,
	); err != nil {
		resp.Diagnostics.AddError(
			"Error deleting alert policy",
			"Could not delete alert policy, unexpected error: "+err.Error(),
		)
		return
	}
}

func (r *alertPolicyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("uuid"), req, resp)
}

func (r *alertPolicyResource) getAlertPolicyFromPlan(plan *alertPolicyResourceModel) *api.AlertPolicy {
	ap := &api.AlertPolicy{
		Name:         plan.Name.ValueString(),
		RuleTimezone: plan.RuleTimezone.ValueString(),
	}

	if !plan.Desc.IsNull() {
		ap.Desc = plan.Desc.ValueString()
	}

	if !plan.OpenPermissionSet.IsNull() {
		ap.OpenPermissionSet = plan.OpenPermissionSet.ValueBool()
	}

	if len(plan.PermissionSet) > 0 {
		permissionSet := make([]string, len(plan.PermissionSet))
		for i, perm := range plan.PermissionSet {
			permissionSet[i] = perm.ValueString()
		}
		ap.PermissionSet = permissionSet
	}

	if len(plan.CheckerUUIDs) > 0 {
		checkerUUIDs := make([]string, len(plan.CheckerUUIDs))
		for i, checker := range plan.CheckerUUIDs {
			checkerUUIDs[i] = checker.ValueString()
		}
		ap.CheckerUUIDs = checkerUUIDs
	}

	if len(plan.SecurityRuleUUIDs) > 0 {
		securityRuleUUIDs := make([]string, len(plan.SecurityRuleUUIDs))
		for i, rule := range plan.SecurityRuleUUIDs {
			securityRuleUUIDs[i] = rule.ValueString()
		}
		ap.SecurityRuleUUIDs = securityRuleUUIDs
	}

	if plan.AlertOpt != nil {
		alertOpt := &api.AlertOpt{}

		if !plan.AlertOpt.AggType.IsNull() {
			alertOpt.AggType = plan.AlertOpt.AggType.ValueString()
		}

		if !plan.AlertOpt.IgnoreOK.IsNull() {
			alertOpt.IgnoreOK = plan.AlertOpt.IgnoreOK.ValueBool()
		}

		if !plan.AlertOpt.AlertType.IsNull() {
			alertOpt.AlertType = plan.AlertOpt.AlertType.ValueString()
		}

		if !plan.AlertOpt.SilentTimeout.IsNull() {
			alertOpt.SilentTimeout = int(plan.AlertOpt.SilentTimeout.ValueInt64())
		}

		if !plan.AlertOpt.SilentTimeoutByStatusEnable.IsNull() {
			alertOpt.SilentTimeoutByStatusEnable = plan.AlertOpt.SilentTimeoutByStatusEnable.ValueBool()
		}

		if len(plan.AlertOpt.SilentTimeoutByStatus) > 0 {
			silentTimeoutByStatus := make([]api.SilentTimeoutByStatus, len(plan.AlertOpt.SilentTimeoutByStatus))
			for i, sts := range plan.AlertOpt.SilentTimeoutByStatus {
				silentTimeoutByStatus[i] = api.SilentTimeoutByStatus{
					Status:        sts.Status.ValueString(),
					SilentTimeout: int(sts.SilentTimeout.ValueInt64()),
				}
			}
			alertOpt.SilentTimeoutByStatus = silentTimeoutByStatus
		}

		if len(plan.AlertOpt.AlertTarget) > 0 {
			alertTarget := make([]api.AlertTarget, len(plan.AlertOpt.AlertTarget))
			for i, at := range plan.AlertOpt.AlertTarget {
				alertTarget[i] = api.AlertTarget{
					Name:            at.Name.ValueString(),
					Crontab:         at.Crontab.ValueString(),
					CrontabDuration: int(at.CrontabDuration.ValueInt64()),
					CustomStartTime: at.CustomStartTime.ValueString(),
					CustomDuration:  int(at.CustomDuration.ValueInt64()),
				}

				if len(at.CustomDateUUIDs) > 0 {
					customDateUUIDs := make([]string, len(at.CustomDateUUIDs))
					for j, uuid := range at.CustomDateUUIDs {
						customDateUUIDs[j] = uuid.ValueString()
					}
					alertTarget[i].CustomDateUUIDs = customDateUUIDs
				}

				if len(at.Targets) > 0 {
					targets := make([]api.Target, len(at.Targets))
					for j, t := range at.Targets {
						targets[j] = api.Target{
							Status:       t.Status.ValueString(),
							DfSource:     t.DfSource.ValueString(),
							FilterString: t.FilterString.ValueString(),
						}

						if len(t.To) > 0 {
							to := make([]string, len(t.To))
							for k, recipient := range t.To {
								to[k] = recipient.ValueString()
							}
							targets[j].To = to
						}

						if len(t.UpgradeTargets) > 0 {
							upgradeTargets := make([]api.UpgradeTarget, len(t.UpgradeTargets))
							for k, ut := range t.UpgradeTargets {
								upgradeTargets[k] = api.UpgradeTarget{
									Duration: int(ut.Duration.ValueInt64()),
								}

								if len(ut.To) > 0 {
									to := make([]string, len(ut.To))
									for l, recipient := range ut.To {
										to[l] = recipient.ValueString()
									}
									upgradeTargets[k].To = to
								}

								if len(ut.ToWay) > 0 {
									toWay := make([]string, len(ut.ToWay))
									for l, way := range ut.ToWay {
										toWay[l] = way.ValueString()
									}
									upgradeTargets[k].ToWay = toWay
								}
							}
							targets[j].UpgradeTargets = upgradeTargets
						}

						if len(t.Tags) > 0 {
							targets[j].Tags = t.Tags
						}
					}
					alertTarget[i].Targets = targets
				}

				if len(at.AlertInfo) > 0 {
					alertInfo := make([]api.AlertInfo, len(at.AlertInfo))
					for j, ai := range at.AlertInfo {
						alertInfo[j] = api.AlertInfo{
							Name:         ai.Name.ValueString(),
							FilterString: ai.FilterString.ValueString(),
						}

						if len(ai.MemberInfo) > 0 {
							memberInfo := make([]string, len(ai.MemberInfo))
							for k, member := range ai.MemberInfo {
								memberInfo[k] = member.ValueString()
							}
							alertInfo[j].MemberInfo = memberInfo
						}

						if len(ai.Targets) > 0 {
							targets := make([]api.Target, len(ai.Targets))
							for k, t := range ai.Targets {
								targets[k] = api.Target{
									Status:       t.Status.ValueString(),
									DfSource:     t.DfSource.ValueString(),
									FilterString: t.FilterString.ValueString(),
								}

								if len(t.To) > 0 {
									to := make([]string, len(t.To))
									for l, recipient := range t.To {
										to[l] = recipient.ValueString()
									}
									targets[k].To = to
								}

								if len(t.UpgradeTargets) > 0 {
									upgradeTargets := make([]api.UpgradeTarget, len(t.UpgradeTargets))
									for l, ut := range t.UpgradeTargets {
										upgradeTargets[l] = api.UpgradeTarget{
											Duration: int(ut.Duration.ValueInt64()),
										}

										if len(ut.To) > 0 {
											to := make([]string, len(ut.To))
											for m, recipient := range ut.To {
												to[m] = recipient.ValueString()
											}
											upgradeTargets[l].To = to
										}

										if len(ut.ToWay) > 0 {
											toWay := make([]string, len(ut.ToWay))
											for m, way := range ut.ToWay {
												toWay[m] = way.ValueString()
											}
											upgradeTargets[l].ToWay = toWay
										}
									}
									targets[k].UpgradeTargets = upgradeTargets
								}

								if len(t.Tags) > 0 {
									targets[k].Tags = t.Tags
								}
							}
							alertInfo[j].Targets = targets
						}
					}
					alertTarget[i].AlertInfo = alertInfo
				}
			}
			alertOpt.AlertTarget = alertTarget
		}

		if !plan.AlertOpt.AggInterval.IsNull() {
			alertOpt.AggInterval = int(plan.AlertOpt.AggInterval.ValueInt64())
		}

		if len(plan.AlertOpt.AggFields) > 0 {
			aggFields := make([]string, len(plan.AlertOpt.AggFields))
			for i, field := range plan.AlertOpt.AggFields {
				aggFields[i] = field.ValueString()
			}
			alertOpt.AggFields = aggFields
		}

		if len(plan.AlertOpt.AggLabels) > 0 {
			aggLabels := make([]string, len(plan.AlertOpt.AggLabels))
			for i, label := range plan.AlertOpt.AggLabels {
				aggLabels[i] = label.ValueString()
			}
			alertOpt.AggLabels = aggLabels
		}

		if len(plan.AlertOpt.AggClusterFields) > 0 {
			aggClusterFields := make([]string, len(plan.AlertOpt.AggClusterFields))
			for i, field := range plan.AlertOpt.AggClusterFields {
				aggClusterFields[i] = field.ValueString()
			}
			alertOpt.AggClusterFields = aggClusterFields
		}

		if !plan.AlertOpt.AggSendFirst.IsNull() {
			alertOpt.AggSendFirst = plan.AlertOpt.AggSendFirst.ValueBool()
		}

		ap.AlertOpt = alertOpt
	}

	return ap
}
