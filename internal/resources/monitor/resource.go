package monitor

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
	_ resource.Resource                = &monitorResource{}
	_ resource.ResourceWithConfigure   = &monitorResource{}
	_ resource.ResourceWithImportState = &monitorResource{}
)

// NewMonitorResource creates a new monitor resource
func NewMonitorResource() resource.Resource {
	return &monitorResource{}
}

// monitorResource is the resource implementation.
type monitorResource struct {
	client *api.Client
}

// Schema defines the schema for the data source.
func (r *monitorResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = resourceSchema
}

// Configure adds the provider configured client to the data source.
func (r *monitorResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.client = req.ProviderData.(*api.Client)
}

// Metadata returns the data source type name.
func (r *monitorResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_monitor"
}

// Create creates the resource and sets the initial Terraform state.
func (r *monitorResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan monitorResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	item := r.getMonitorFromPlan(&plan)
	rb, _ := json.Marshal(item)
	tflog.Debug(ctx, fmt.Sprintf("============= body: %s", string(rb)))
	content := &api.MonitorContent{}
	err := r.client.Create(consts.TypeNameMonitor, item, content)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating monitor",
			"Could not create monitor, unexpected error: "+err.Error(),
		)
		return
	}

	plan.UUID = types.StringValue(content.UUID)
	plan.WorkspaceUUID = types.StringValue(content.WorkspaceUUID)
	plan.MonitorUUID = types.StringValue(content.MonitorUUID)
	plan.MonitorName = types.StringValue(content.MonitorName)
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
func (r *monitorResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state monitorResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	content := &api.MonitorContent{}
	err := r.client.Read(consts.TypeNameMonitor, state.UUID.ValueString(), content)
	if err != nil {
		if err == api.Error404 {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error reading monitor",
			"Could not read monitor, unexpected error: "+err.Error(),
		)
		return
	}

	state.UUID = types.StringValue(content.UUID)
	state.Type = types.StringValue(content.Type)
	state.Status = types.Int64Value(int64(content.Status))
	state.MonitorUUID = types.StringValue(content.MonitorUUID)
	state.MonitorName = types.StringValue(content.MonitorName)
	state.WorkspaceUUID = types.StringValue(content.WorkspaceUUID)
	state.CreateAt = types.Int64Value(int64(content.CreateAt))
	state.UpdateAt = types.Int64Value(int64(content.UpdateAt))

	// Map JsonScript from API response
	if content.JsonScript != nil {
		if jsMap, ok := content.JsonScript.(map[string]interface{}); ok {
			// Map top-level fields
			if typeVal, ok := jsMap["type"].(string); ok {
				state.JsonScript.Type = types.StringValue(typeVal)
			}
			if titleVal, ok := jsMap["title"].(string); ok {
				state.JsonScript.Title = types.StringValue(titleVal)
			}
			if messageVal, ok := jsMap["message"].(string); ok {
				state.JsonScript.Message = types.StringValue(messageVal)
			}
			if everyVal, ok := jsMap["every"].(string); ok {
				state.JsonScript.Every = types.StringValue(everyVal)
			}
			if intervalVal, ok := jsMap["interval"].(float64); ok {
				state.JsonScript.Interval = types.Int64Value(int64(intervalVal))
			}
			if intervalVal, ok := jsMap["interval"].(int); ok {
				state.JsonScript.Interval = types.Int64Value(int64(intervalVal))
			}
			if recoverVal, ok := jsMap["recoverNeedPeriodCount"].(float64); ok {
				state.JsonScript.RecoverNeedPeriodCount = types.Int64Value(int64(recoverVal))
			}
			if recoverVal, ok := jsMap["recoverNeedPeriodCount"].(int); ok {
				state.JsonScript.RecoverNeedPeriodCount = types.Int64Value(int64(recoverVal))
			}
			if disableVal, ok := jsMap["disableCheckEndTime"].(bool); ok {
				state.JsonScript.DisableCheckEndTime = types.BoolValue(disableVal)
			}

			// Map groupBy
			if groupByVal, ok := jsMap["groupBy"].([]interface{}); ok {
				groupBy := make([]types.String, len(groupByVal))
				for i, gb := range groupByVal {
					if gbStr, ok := gb.(string); ok {
						groupBy[i] = types.StringValue(gbStr)
					}
				}
				state.JsonScript.GroupBy = groupBy
			}

			// Map targets
			if targetsVal, ok := jsMap["targets"].([]interface{}); ok {
				targets := make([]Target, len(targetsVal))
				for i, t := range targetsVal {
					if targetMap, ok := t.(map[string]interface{}); ok {
						target := Target{}
						if dqlVal, ok := targetMap["dql"].(string); ok {
							target.Dql = types.StringValue(dqlVal)
						}
						if aliasVal, ok := targetMap["alias"].(string); ok {
							target.Alias = types.StringValue(aliasVal)
						}
						if qtypeVal, ok := targetMap["qtype"].(string); ok {
							target.Qtype = types.StringValue(qtypeVal)
						}
						targets[i] = target
					}
				}
				state.JsonScript.Targets = targets
			}

			// Map checkerOpt
			if checkerOptVal, ok := jsMap["checkerOpt"].(map[string]interface{}); ok {
				if infoEventVal, ok := checkerOptVal["infoEvent"].(bool); ok {
					state.JsonScript.CheckerOpt.InfoEvent = types.BoolValue(infoEventVal)
				}

				// Map rules
				if rulesVal, ok := checkerOptVal["rules"].([]interface{}); ok {
					rules := make([]Rule, len(rulesVal))
					for i, r := range rulesVal {
						if ruleMap, ok := r.(map[string]interface{}); ok {
							rule := Rule{}
							if conditionLogicVal, ok := ruleMap["conditionLogic"].(string); ok {
								rule.ConditionLogic = types.StringValue(conditionLogicVal)
							}
							if statusVal, ok := ruleMap["status"].(string); ok {
								rule.Status = types.StringValue(statusVal)
							}

							// Map conditions
							if conditionsVal, ok := ruleMap["conditions"].([]interface{}); ok {
								conditions := make([]Condition, len(conditionsVal))
								for j, c := range conditionsVal {
									if conditionMap, ok := c.(map[string]interface{}); ok {
										condition := Condition{}
										if aliasVal, ok := conditionMap["alias"].(string); ok {
											condition.Alias = types.StringValue(aliasVal)
										}
										if operatorVal, ok := conditionMap["operator"].(string); ok {
											condition.Operator = types.StringValue(operatorVal)
										}
										if operandsVal, ok := conditionMap["operands"].([]interface{}); ok {
											operands := make([]types.String, len(operandsVal))
											for k, op := range operandsVal {
												if opStr, ok := op.(string); ok {
													operands[k] = types.StringValue(opStr)
												}
											}
											condition.Operands = operands
										}
										conditions[j] = condition
									}
								}
								rule.Conditions = conditions
							}
							rules[i] = rule
						}
					}
					state.JsonScript.CheckerOpt.Rules = rules
				}
			}

			// Map channels
			if channelsVal, ok := jsMap["channels"].([]interface{}); ok {
				channels := make([]types.String, len(channelsVal))
				for i, ch := range channelsVal {
					if chStr, ok := ch.(string); ok {
						channels[i] = types.StringValue(chStr)
					}
				}
				state.JsonScript.Channels = channels
			}

			// Map atAccounts
			if atAccountsVal, ok := jsMap["atAccounts"].([]interface{}); ok {
				atAccounts := make([]types.String, len(atAccountsVal))
				for i, acct := range atAccountsVal {
					if acctStr, ok := acct.(string); ok {
						atAccounts[i] = types.StringValue(acctStr)
					}
				}
				state.JsonScript.AtAccounts = atAccounts
			}

			// Map atNoDataAccounts
			if atNoDataAccountsVal, ok := jsMap["atNoDataAccounts"].([]interface{}); ok {
				atNoDataAccounts := make([]types.String, len(atNoDataAccountsVal))
				for i, acct := range atNoDataAccountsVal {
					if acctStr, ok := acct.(string); ok {
						atNoDataAccounts[i] = types.StringValue(acctStr)
					}
				}
				state.JsonScript.AtNoDataAccounts = atNoDataAccounts
			}
		}
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *monitorResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan monitorResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	item := r.getMonitorFromPlan(&plan)
	// Type is not updatable, so we need to set it to empty string to avoid error
	item.Type = ""
	content := &api.MonitorContent{}
	err := r.client.Update(consts.TypeNameMonitor, plan.UUID.ValueString(), item, content)

	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating monitor",
			"Could not update monitor, unexpected error: "+err.Error(),
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
func (r *monitorResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state monitorResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete existing monitor
	if err := r.client.DeleteMonitor(
		state.UUID.ValueString(),
	); err != nil {
		resp.Diagnostics.AddError(
			"Error deleting monitor",
			"Could not delete monitor, unexpected error: "+err.Error(),
		)
		return
	}
}

func (r *monitorResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("uuid"), req, resp)
}

func (r *monitorResource) getMonitorFromPlan(plan *monitorResourceModel) *api.Monitor {

	m := &api.Monitor{}

	if !plan.Type.IsNull() {
		m.Type = plan.Type.ValueString()
	}

	if !plan.Status.IsNull() {
		m.Status = int(plan.Status.ValueInt64())
	}

	if !plan.Extend.IsNull() {
		var extend interface{}
		if err := json.Unmarshal([]byte(plan.Extend.ValueString()), &extend); err == nil {
			m.Extend = extend
		}
	}

	if len(plan.AlertPolicyUUIDs) > 0 {
		alertPolicyUUIDs := make([]string, len(plan.AlertPolicyUUIDs))
		for i, uuid := range plan.AlertPolicyUUIDs {
			alertPolicyUUIDs[i] = uuid.ValueString()
		}
		m.AlertPolicyUUIDs = alertPolicyUUIDs
	}

	if !plan.DashboardUUID.IsNull() {
		m.DashboardUUID = plan.DashboardUUID.ValueString()
	}

	if len(plan.Tags) > 0 {
		tags := make([]string, len(plan.Tags))
		for i, tag := range plan.Tags {
			tags[i] = tag.ValueString()
		}
		m.Tags = tags
	}

	if !plan.Secret.IsNull() {
		m.Secret = plan.Secret.ValueString()
	}

	// Build jsonScript from nested structure
	jsonScript := make(map[string]interface{})

	if !plan.JsonScript.Type.IsNull() {
		jsonScript["type"] = plan.JsonScript.Type.ValueString()
	}

	if !plan.JsonScript.Title.IsNull() {
		jsonScript["title"] = plan.JsonScript.Title.ValueString()
	}

	if !plan.JsonScript.Message.IsNull() {
		jsonScript["message"] = plan.JsonScript.Message.ValueString()
	}

	if !plan.JsonScript.Every.IsNull() {
		jsonScript["every"] = plan.JsonScript.Every.ValueString()
	}

	if !plan.JsonScript.Interval.IsNull() {
		jsonScript["interval"] = plan.JsonScript.Interval.ValueInt64()
	}

	if !plan.JsonScript.RecoverNeedPeriodCount.IsNull() {
		jsonScript["recoverNeedPeriodCount"] = plan.JsonScript.RecoverNeedPeriodCount.ValueInt64()
	}

	if !plan.JsonScript.DisableCheckEndTime.IsNull() {
		jsonScript["disableCheckEndTime"] = plan.JsonScript.DisableCheckEndTime.ValueBool()
	}

	if len(plan.JsonScript.GroupBy) > 0 {
		groupBy := make([]string, len(plan.JsonScript.GroupBy))
		for i, gb := range plan.JsonScript.GroupBy {
			groupBy[i] = gb.ValueString()
		}
		jsonScript["groupBy"] = groupBy
	}

	if len(plan.JsonScript.Targets) > 0 {
		targets := make([]map[string]interface{}, len(plan.JsonScript.Targets))
		for i, target := range plan.JsonScript.Targets {
			t := make(map[string]interface{})
			if !target.Dql.IsNull() {
				t["dql"] = target.Dql.ValueString()
			}
			if !target.Alias.IsNull() {
				t["alias"] = target.Alias.ValueString()
			}
			if !target.Qtype.IsNull() {
				t["qtype"] = target.Qtype.ValueString()
			}
			targets[i] = t
		}
		jsonScript["targets"] = targets
	}

	// Build checkerOpt
	checkerOpt := make(map[string]interface{})
	if !plan.JsonScript.CheckerOpt.InfoEvent.IsNull() {
		checkerOpt["infoEvent"] = plan.JsonScript.CheckerOpt.InfoEvent.ValueBool()
	}

	if len(plan.JsonScript.CheckerOpt.Rules) > 0 {
		rules := make([]map[string]interface{}, len(plan.JsonScript.CheckerOpt.Rules))
		for i, rule := range plan.JsonScript.CheckerOpt.Rules {
			r := make(map[string]interface{})
			if !rule.ConditionLogic.IsNull() {
				r["conditionLogic"] = rule.ConditionLogic.ValueString()
			}
			if !rule.Status.IsNull() {
				r["status"] = rule.Status.ValueString()
			}

			if len(rule.Conditions) > 0 {
				conditions := make([]map[string]interface{}, len(rule.Conditions))
				for j, cond := range rule.Conditions {
					c := make(map[string]interface{})
					if !cond.Alias.IsNull() {
						c["alias"] = cond.Alias.ValueString()
					}
					if !cond.Operator.IsNull() {
						c["operator"] = cond.Operator.ValueString()
					}
					if len(cond.Operands) > 0 {
						operands := make([]string, len(cond.Operands))
						for k, op := range cond.Operands {
							operands[k] = op.ValueString()
						}
						c["operands"] = operands
					}
					conditions[j] = c
				}
				r["conditions"] = conditions
			}
			rules[i] = r
		}
		checkerOpt["rules"] = rules
	}

	if len(checkerOpt) > 0 {
		jsonScript["checkerOpt"] = checkerOpt
	}

	if len(plan.JsonScript.Channels) > 0 {
		channels := make([]string, len(plan.JsonScript.Channels))
		for i, ch := range plan.JsonScript.Channels {
			channels[i] = ch.ValueString()
		}
		jsonScript["channels"] = channels
	}

	if len(plan.JsonScript.AtAccounts) > 0 {
		atAccounts := make([]string, len(plan.JsonScript.AtAccounts))
		for i, acct := range plan.JsonScript.AtAccounts {
			atAccounts[i] = acct.ValueString()
		}
		jsonScript["atAccounts"] = atAccounts
	}

	if len(plan.JsonScript.AtNoDataAccounts) > 0 {
		atNoDataAccounts := make([]string, len(plan.JsonScript.AtNoDataAccounts))
		for i, acct := range plan.JsonScript.AtNoDataAccounts {
			atNoDataAccounts[i] = acct.ValueString()
		}
		jsonScript["atNoDataAccounts"] = atNoDataAccounts
	}

	if len(jsonScript) > 0 {
		m.JsonScript = jsonScript
	}

	if !plan.OpenPermissionSet.IsNull() {
		m.OpenPermissionSet = plan.OpenPermissionSet.ValueBool()
	}

	if len(plan.PermissionSet) > 0 {
		permissionSet := make([]string, len(plan.PermissionSet))
		for i, perm := range plan.PermissionSet {
			permissionSet[i] = perm.ValueString()
		}
		m.PermissionSet = permissionSet
	}

	return m
}
