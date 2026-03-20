package alert_policy

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var resourceSchema = schema.Schema{
	Description:         "Alert Policy",
	MarkdownDescription: resourceDocument,
	Attributes: map[string]schema.Attribute{
		"uuid": schema.StringAttribute{
			Description: "The uuid of resource.",
			Computed:    true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"name": schema.StringAttribute{
			Description: "The name of the alert policy.",
			Required:    true,
			Validators: []validator.String{
				stringvalidator.LengthAtMost(128),
			},
		},
		"desc": schema.StringAttribute{
			Description: "The description of the alert policy.",
			Optional:    true,
			Validators: []validator.String{
				stringvalidator.LengthAtMost(256),
			},
		},
		"open_permission_set": schema.BoolAttribute{
			Description: "Whether to open custom permission configuration.",
			Optional:    true,
		},
		"permission_set": schema.ListAttribute{
			Description: "Operation permission configuration, can configure roles, member uuids, team uuids.",
			Optional:    true,
			ElementType: types.StringType,
		},
		"checker_uuids": schema.ListAttribute{
			Description: "Monitor/smart monitor/smart inspection/slo uuid.",
			Optional:    true,
			ElementType: types.StringType,
		},
		"security_rule_uuids": schema.ListAttribute{
			Description: "Security monitoring (cspm, siem) uuids.",
			Optional:    true,
			ElementType: types.StringType,
		},
		"rule_timezone": schema.StringAttribute{
			Description: "The timezone of the alert policy.",
			Required:    true,
			Validators: []validator.String{
				stringvalidator.LengthAtMost(64),
			},
		},
		"alert_opt": schema.SingleNestedAttribute{
			Description: "Alert settings.",
			Optional:    true,
			Attributes: map[string]schema.Attribute{
				"agg_type": schema.StringAttribute{
					Description: "Alert aggregation type.",
					Optional:    true,
					Validators: []validator.String{
						stringvalidator.OneOf("byFields", "byCluster", "byAI", "byCustom"),
					},
				},
				"ignore_ok": schema.BoolAttribute{
					Description: "Advanced configuration, normal level only generates events, does not send notifications.",
					Optional:    true,
				},
				"alert_type": schema.StringAttribute{
					Description: "Alert policy notification type, level(status)/member(member).",
					Optional:    true,
					Validators: []validator.String{
						stringvalidator.OneOf("status", "member"),
					},
				},
				"silent_timeout": schema.Int64Attribute{
					Description: "Minimum alert interval, how long the same alert will not be sent repeatedly (i.e., alert silence duration), in seconds/s, 0/null means only send once alert (i.e., interval is infinitely long).",
					Optional:    true,
				},
				"silent_timeout_by_status_enable": schema.BoolAttribute{
					Description: "Whether to enable level-specific repeated alert configuration, default false, use silentTimeout configuration.",
					Optional:    true,
				},
				"silent_timeout_by_status": schema.ListNestedAttribute{
					Description: "Level-specific minimum alert interval configuration.",
					Optional:    true,
					NestedObject: schema.NestedAttributeObject{
						Attributes: map[string]schema.Attribute{
							"status": schema.StringAttribute{
								Description: "Status value for non-security monitoring types: fatal, critical, error, warning, nodata, info. For security monitoring types: security_critical, security_high, security_medium, security_low, security_info.",
								Required:    true,
							},
							"silent_timeout": schema.Int64Attribute{
								Description: "Minimum alert interval, in seconds.",
								Required:    true,
							},
						},
					},
				},
				"alert_target": schema.ListNestedAttribute{
					Description: "Trigger action, note the trigger time parameter processing.",
					Optional:    true,
					NestedObject: schema.NestedAttributeObject{
						Attributes: map[string]schema.Attribute{
							"name": schema.StringAttribute{
								Description: "Configuration name.",
								Optional:    true,
							},
							"targets": schema.ListNestedAttribute{
								Description: "Notification target configuration.",
								Optional:    true,
								NestedObject: schema.NestedAttributeObject{
									Attributes: map[string]schema.Attribute{
										"to": schema.ListAttribute{
											Description: "Notification objects/members/teams.",
											Required:    true,
											ElementType: types.StringType,
										},
										"status": schema.StringAttribute{
											Description: "Status values for which alerts need to be sent (multiple statuses can be separated by commas, All means all).",
											Required:    true,
										},
										"df_source": schema.StringAttribute{
											Description: "When status needs to be security monitoring status, this must be specified as security, default not passing indicates non-security monitoring status.",
											Optional:    true,
										},
										"upgrade_targets": schema.ListNestedAttribute{
											Description: "Upgrade notification for each alert configuration status.",
											Optional:    true,
											NestedObject: schema.NestedAttributeObject{
												Attributes: map[string]schema.Attribute{
													"to": schema.ListAttribute{
														Description: "Notification objects/members/teams.",
														Required:    true,
														ElementType: types.StringType,
													},
													"duration": schema.Int64Attribute{
														Description: "Duration, continuous generation of events at this level status triggers upgrade notification.",
														Optional:    true,
													},
													"to_way": schema.ListAttribute{
														Description: "When alertType is member type, only notification objects and fixed fields email, sms can be selected.",
														Optional:    true,
														ElementType: types.StringType,
													},
												},
											},
										},
										"tags": schema.MapAttribute{
											Description: "Filter conditions.",
											Optional:    true,
											ElementType: types.ListType{
												ElemType: types.StringType,
											},
										},
										"filter_string": schema.StringAttribute{
											Description: "Filter condition original string, can replace tags, filterString use priority is greater than tags.",
											Optional:    true,
										},
									},
								},
							},
							"crontab": schema.StringAttribute{
								Description: "When selecting repeated time period, start Crontab (Crontab syntax).",
								Optional:    true,
							},
							"crontab_duration": schema.Int64Attribute{
								Description: "Select repeated time, from Crontab start, duration (seconds).",
								Optional:    true,
							},
							"custom_date_uuids": schema.ListAttribute{
								Description: "When selecting custom time, custom notification date UUID list.",
								Optional:    true,
								ElementType: types.StringType,
							},
							"custom_start_time": schema.StringAttribute{
								Description: "When selecting custom time, daily start time, format: HH:mm:ss.",
								Optional:    true,
							},
							"custom_duration": schema.Int64Attribute{
								Description: "When selecting custom time period, from customStartTime custom start time, duration (seconds).",
								Optional:    true,
							},
							"alert_info": schema.ListNestedAttribute{
								Description: "When the alert policy is of member type, the notification related information configuration.",
								Optional:    true,
								NestedObject: schema.NestedAttributeObject{
									Attributes: map[string]schema.Attribute{
										"name": schema.StringAttribute{
											Description: "Configuration name.",
											Optional:    true,
										},
										"targets": schema.ListNestedAttribute{
											Description: "Notification target configuration.",
											Required:    true,
											NestedObject: schema.NestedAttributeObject{
												Attributes: map[string]schema.Attribute{
													"to": schema.ListAttribute{
														Description: "Notification objects/members/teams.",
														Required:    true,
														ElementType: types.StringType,
													},
													"status": schema.StringAttribute{
														Description: "Status values for which alerts need to be sent (multiple statuses can be separated by commas, All means all).",
														Required:    true,
													},
													"df_source": schema.StringAttribute{
														Description: "When status needs to be security monitoring status, this must be specified as security, default not passing indicates non-security monitoring status.",
														Optional:    true,
													},
													"upgrade_targets": schema.ListNestedAttribute{
														Description: "Upgrade notification for each alert configuration status.",
														Optional:    true,
														NestedObject: schema.NestedAttributeObject{
															Attributes: map[string]schema.Attribute{
																"to": schema.ListAttribute{
																	Description: "Notification objects/members/teams.",
																	Required:    true,
																	ElementType: types.StringType,
																},
																"duration": schema.Int64Attribute{
																	Description: "Duration, continuous generation of events at this level status triggers upgrade notification.",
																	Optional:    true,
																},
																"to_way": schema.ListAttribute{
																	Description: "When alertType is member type, only notification objects and fixed fields email, sms can be selected.",
																	Optional:    true,
																	ElementType: types.StringType,
																},
															},
															},
														},
														"tags": schema.MapAttribute{
															Description: "Filter conditions.",
															Optional:    true,
															ElementType: types.ListType{
																ElemType: types.StringType,
															},
														},
														"filter_string": schema.StringAttribute{
															Description: "Filter condition original string, can replace tags, filterString use priority is greater than tags.",
															Optional:    true,
														},
													},
												},
										},
										"filter_string": schema.StringAttribute{
											Description: "When alertType is member, use this field, filter condition original string.",
											Optional:    true,
										},
										"member_info": schema.ListAttribute{
											Description: "When alertType is member, use this field (team UUID member UUID).",
											Optional:    true,
											ElementType: types.StringType,
										},
									},
								},
							},
						},
					},
				},
				"agg_interval": schema.Int64Attribute{
					Description: "Alert aggregation interval, in seconds, 0 means no aggregation.",
					Optional:    true,
					Validators: []validator.Int64{
						int64validator.Between(0, 1800),
					},
				},
				"agg_fields": schema.ListAttribute{
					Description: "Aggregation field list, keep empty list [] means \"aggregation rule: all\", df_monitor_checker_id: monitor/smart inspection/SLO, df_dimension_tags: detection dimension, df_label: label, CLUSTER: smart aggregation, or when aggType=byCustom means custom aggregation field.",
					Optional:    true,
					ElementType: types.StringType,
				},
				"agg_labels": schema.ListAttribute{
					Description: "Label value list when aggregating by label, need to specify df_label in aggFields to take effect.",
					Optional:    true,
					ElementType: types.StringType,
				},
				"agg_cluster_fields": schema.ListAttribute{
					Description: "Field list for smart aggregation, need to specify CLUSTER in aggFields to take effect, optional values \"df_title\": title, \"df_message\": content.",
					Optional:    true,
					ElementType: types.StringType,
				},
				"agg_send_first": schema.BoolAttribute{
					Description: "When aggregating, whether the first alert is sent directly (added in 2025-09-03 iteration).",
					Optional:    true,
				},
			},
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
	},
}
