package monitor

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var resourceSchema = schema.Schema{
	Description:         "Monitor",
	MarkdownDescription: resourceDocument,
	Attributes: map[string]schema.Attribute{
		"uuid": schema.StringAttribute{
			Description: "The uuid of resource.",
			Computed:    true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"type": schema.StringAttribute{
			Description: "Monitor type, default trigger, trigger: normal monitor, smartMonitor: smart monitor.",
			Optional:    true,
			Validators: []validator.String{
				stringvalidator.OneOf("trigger", "smartMonitor"),
			},
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.RequiresReplace(),
			},
		},
		"status": schema.Int64Attribute{
			Description: "Monitor status, 0: enabled, 2: disabled, default enabled.",
			Optional:    true,
			Validators:  []validator.Int64{
				// Add int64 validators if needed
			},
		},
		"extend": schema.StringAttribute{
			Description: "Additional information (fields related to exception tracking and some fields used for frontend echo).",
			Optional:    true,
		},
		"alert_policy_uuids": schema.ListAttribute{
			Description: "Alert policy UUIDs.",
			Optional:    true,
			ElementType: types.StringType,
		},
		"dashboard_uuid": schema.StringAttribute{
			Description: "Associated dashboard ID.",
			Optional:    true,
		},
		"tags": schema.ListAttribute{
			Description: "Tag names for filtering.",
			Optional:    true,
			ElementType: types.StringType,
		},
		"secret": schema.StringAttribute{
			Description: "Unique identifier secret for the middle section of the Webhook address (generally take a random uuid to ensure uniqueness within the workspace).",
			Optional:    true,
		},
		"json_script": schema.SingleNestedAttribute{
			Description: "Rule configuration.",
			Required:    true,
			Attributes: map[string]schema.Attribute{
				"type": schema.StringAttribute{
					Description: "Check method type.",
					Required:    true,
				},
				"title": schema.StringAttribute{
					Description: "Generated event title.",
					Optional:    true,
				},
				"message": schema.StringAttribute{
					Description: "Event content.",
					Optional:    true,
				},
				"every": schema.StringAttribute{
					Description: "Check frequency.",
					Optional:    true,
				},
				"interval": schema.Int64Attribute{
					Description: "Query interval, i.e., the time range difference for one query.",
					Optional:    true,
				},
				"recover_need_period_count": schema.Int64Attribute{
					Description: "Specify after how many check cycles an exception generates a recovery event.",
					Optional:    true,
				},
				"disable_check_end_time": schema.BoolAttribute{
					Description: "Whether to disable end time limit.",
					Optional:    true,
				},
				"group_by": schema.ListAttribute{
					Description: "Trigger dimension.",
					Optional:    true,
					ElementType: types.StringType,
				},
				"targets": schema.ListNestedAttribute{
					Description: "Check targets.",
					Optional:    true,
					NestedObject: schema.NestedAttributeObject{
						Attributes: map[string]schema.Attribute{
							"dql": schema.StringAttribute{
								Description: "DQL query statement.",
								Required:    true,
							},
							"alias": schema.StringAttribute{
								Description: "Alias.",
								Required:    true,
							},
							"qtype": schema.StringAttribute{
								Description: "Query type.",
								Optional:    true,
							},
						},
					},
				},
				"checker_opt": schema.SingleNestedAttribute{
					Description: "Check condition settings.",
					Optional:    true,
					Attributes: map[string]schema.Attribute{
						"info_event": schema.BoolAttribute{
							Description: "Whether to generate info events during continuous normal operation.",
							Optional:    true,
						},
						"rules": schema.ListNestedAttribute{
							Description: "Trigger condition list.",
							Optional:    true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"condition_logic": schema.StringAttribute{
										Description: "Logic between conditions, and/or.",
										Required:    true,
									},
									"conditions": schema.ListNestedAttribute{
										Description: "Conditions.",
										Required:    true,
										NestedObject: schema.NestedAttributeObject{
											Attributes: map[string]schema.Attribute{
												"alias": schema.StringAttribute{
													Description: "Detection object alias, i.e., the above targets[#].alias.",
													Required:    true,
												},
												"operator": schema.StringAttribute{
													Description: "Operator, =, >, <, etc.",
													Required:    true,
												},
												"operands": schema.ListAttribute{
													Description: "Operand array. (Operators like between, in need multiple operands).",
													Required:    true,
													ElementType: types.StringType,
												},
											},
										},
									},
									"status": schema.StringAttribute{
										Description: "When the condition is met, the status of the output event.",
										Required:    true,
									},
								},
							},
						},
					},
				},
				"channels": schema.ListAttribute{
					Description: "Channel UUID list.",
					Optional:    true,
					ElementType: types.StringType,
				},
				"at_accounts": schema.ListAttribute{
					Description: "Account UUID list to be @ under normal detection.",
					Optional:    true,
					ElementType: types.StringType,
				},
				"at_no_data_accounts": schema.ListAttribute{
					Description: "Account UUID list to be @ under no data circumstances.",
					Optional:    true,
					ElementType: types.StringType,
				},
			},
		},
		"open_permission_set": schema.BoolAttribute{
			Description: "Enable custom permission configuration, (default false: not enabled), after enabling, the operation permission of the rule is based on permissionSet.",
			Optional:    true,
		},
		"permission_set": schema.ListAttribute{
			Description: "Operation permission configuration, configurable (roles (except owners), member uuid, team uuid).",
			Optional:    true,
			ElementType: types.StringType,
		},
		"create_at": schema.Int64Attribute{
			Description: "The timestamp seconds of the resource created at.",
			Computed:    true,
			PlanModifiers: []planmodifier.Int64{
				int64planmodifier.UseStateForUnknown(),
			},
		},
		"update_at": schema.Int64Attribute{
			Description: "The timestamp seconds of the resource updated at.",
			Computed:    true,
		},
		"workspace_uuid": schema.StringAttribute{
			Description: "The uuid of the workspace.",
			Computed:    true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"monitor_uuid": schema.StringAttribute{
			Description: "The uuid of the monitor.",
			Computed:    true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"monitor_name": schema.StringAttribute{
			Description: "The name of the monitor.",
			Computed:    true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
	},
}
