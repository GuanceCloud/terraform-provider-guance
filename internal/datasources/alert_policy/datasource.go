package alert_policy

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/datasourcevalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	dsschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/GuanceCloud/terraform-provider-guance/internal/api"
	"github.com/GuanceCloud/terraform-provider-guance/internal/consts"
	"github.com/GuanceCloud/terraform-provider-guance/internal/helpers/tfconvert"
	resourcealert "github.com/GuanceCloud/terraform-provider-guance/internal/resources/alert_policy"
)

var (
	_ datasource.DataSource                     = &alertPolicyDataSource{}
	_ datasource.DataSourceWithConfigure        = &alertPolicyDataSource{}
	_ datasource.DataSourceWithConfigValidators = &alertPolicyDataSource{}
)

func NewAlertPolicyDataSource() datasource.DataSource {
	return &alertPolicyDataSource{}
}

type alertPolicyDataSource struct {
	client *api.Client
}

type alertPolicyDataSourceModel struct {
	UUID              types.String                 `tfsdk:"uuid"`
	Name              types.String                 `tfsdk:"name"`
	NotifyObjectUUIDs []types.String               `tfsdk:"notify_object_uuids"`
	Desc              types.String                 `tfsdk:"desc"`
	OpenPermissionSet types.Bool                   `tfsdk:"open_permission_set"`
	PermissionSet     []types.String               `tfsdk:"permission_set"`
	CheckerUUIDs      []types.String               `tfsdk:"checker_uuids"`
	SecurityRuleUUIDs []types.String               `tfsdk:"security_rule_uuids"`
	RuleTimezone      types.String                 `tfsdk:"rule_timezone"`
	AlertOpt          *resourcealert.AlertOptModel `tfsdk:"alert_opt"`
	CreateAt          types.Int64                  `tfsdk:"create_at"`
	UpdateAt          types.Int64                  `tfsdk:"update_at"`
	WorkspaceUUID     types.String                 `tfsdk:"workspace_uuid"`
}

func (d *alertPolicyDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_alert_policy"
}

func (d *alertPolicyDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = dsschema.Schema{
		Description:         "Lookup an alert policy by UUID or exact name.",
		MarkdownDescription: "The `guance_alert_policy` data source reads an existing Guance alert policy by `uuid` or exact `name`.",
		Attributes: map[string]dsschema.Attribute{
			"uuid": dsschema.StringAttribute{
				Description: "The UUID of the alert policy.",
				Optional:    true,
				Computed:    true,
			},
			"name": dsschema.StringAttribute{
				Description: "The exact name of the alert policy.",
				Optional:    true,
				Computed:    true,
			},
			"notify_object_uuids": dsschema.ListAttribute{
				Description: "Filter alert policies by notification object UUIDs when looking up by name.",
				Optional:    true,
				ElementType: types.StringType,
			},
			"desc": dsschema.StringAttribute{
				Description: "The description of the alert policy.",
				Computed:    true,
			},
			"open_permission_set": dsschema.BoolAttribute{
				Description: "Whether custom permission configuration is enabled.",
				Computed:    true,
			},
			"permission_set": dsschema.ListAttribute{
				Description: "Operation permission configuration.",
				Computed:    true,
				ElementType: types.StringType,
			},
			"checker_uuids": dsschema.ListAttribute{
				Description: "Monitor/smart monitor/smart inspection/slo UUIDs.",
				Computed:    true,
				ElementType: types.StringType,
			},
			"security_rule_uuids": dsschema.ListAttribute{
				Description: "Security monitoring UUIDs.",
				Computed:    true,
				ElementType: types.StringType,
			},
			"rule_timezone": dsschema.StringAttribute{
				Description: "The timezone of the alert policy.",
				Computed:    true,
			},
			"alert_opt": alertOptDataSourceAttribute(),
			"create_at": dsschema.Int64Attribute{
				Description: "The timestamp seconds of the resource created at.",
				Computed:    true,
			},
			"update_at": dsschema.Int64Attribute{
				Description: "The timestamp seconds of the resource updated at.",
				Computed:    true,
			},
			"workspace_uuid": dsschema.StringAttribute{
				Description: "The UUID of the workspace.",
				Computed:    true,
			},
		},
	}
}

func (d *alertPolicyDataSource) ConfigValidators(_ context.Context) []datasource.ConfigValidator {
	return []datasource.ConfigValidator{
		datasourcevalidator.ExactlyOneOf(path.MatchRoot("uuid"), path.MatchRoot("name")),
	}
}

func (d *alertPolicyDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	d.client = req.ProviderData.(*api.Client)
}

func (d *alertPolicyDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state alertPolicyDataSourceModel
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	content := &api.AlertPolicyContent{}
	if !state.UUID.IsNull() && state.UUID.ValueString() != "" {
		if err := d.client.Read(consts.TypeNameAlertPolicy, state.UUID.ValueString(), content); err != nil {
			resp.Diagnostics.AddError("Error reading alert policy", err.Error())
			return
		}
	} else {
		list := &api.AlertPolicyListContent{}
		options := api.AlertPolicyListOptions{
			Search:            state.Name.ValueString(),
			NotifyObjectUUIDs: tfconvert.TypesToStrings(state.NotifyObjectUUIDs),
		}
		if err := d.client.ListAlertPoliciesWithOptions(options, list); err != nil {
			resp.Diagnostics.AddError("Error listing alert policies", err.Error())
			return
		}
		matched := make([]api.AlertPolicyContent, 0, 1)
		for _, item := range list.Data {
			if item.Name == state.Name.ValueString() {
				matched = append(matched, item)
			}
		}
		if len(matched) != 1 {
			resp.Diagnostics.AddError("Unable to find unique alert policy", fmt.Sprintf("found %d alert policies with name %q", len(matched), state.Name.ValueString()))
			return
		}
		*content = matched[0]
	}

	state.UUID = types.StringValue(content.UUID)
	state.Name = types.StringValue(content.Name)
	state.Desc = tfconvert.StringValueOrNull(content.Desc)
	state.OpenPermissionSet = types.BoolValue(content.OpenPermissionSet)
	state.PermissionSet = tfconvert.StringsToTypes(content.PermissionSet)
	state.CheckerUUIDs = tfconvert.StringsToTypes(content.CheckerUUIDs)
	state.SecurityRuleUUIDs = tfconvert.StringsToTypes(content.SecurityRuleUUIDs)
	state.RuleTimezone = tfconvert.StringValueOrNull(content.RuleTimezone)
	state.AlertOpt = resourcealert.AlertOptFromContentForDataSource(content.AlertOpt)
	state.CreateAt = types.Int64Value(int64(content.CreateAt))
	state.UpdateAt = types.Int64Value(int64(content.UpdateAt))
	state.WorkspaceUUID = types.StringValue(content.WorkspaceUUID)
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

func alertOptDataSourceAttribute() dsschema.SingleNestedAttribute {
	return dsschema.SingleNestedAttribute{
		Description: "Alert settings.",
		Computed:    true,
		Attributes: map[string]dsschema.Attribute{
			"agg_type": dsschema.StringAttribute{
				Description: "Alert aggregation type.",
				Computed:    true,
			},
			"ignore_ok": dsschema.BoolAttribute{
				Description: "Advanced configuration, normal level only generates events, does not send notifications.",
				Computed:    true,
			},
			"alert_type": dsschema.StringAttribute{
				Description: "Alert policy notification type, level(status)/member(member).",
				Computed:    true,
			},
			"silent_timeout": dsschema.Int64Attribute{
				Description: "Minimum alert interval in seconds.",
				Computed:    true,
			},
			"silent_timeout_by_status_enable": dsschema.BoolAttribute{
				Description: "Whether level-specific repeated alert configuration is enabled.",
				Computed:    true,
			},
			"silent_timeout_by_status": dsschema.ListNestedAttribute{
				Description: "Level-specific minimum alert interval configuration.",
				Computed:    true,
				NestedObject: dsschema.NestedAttributeObject{
					Attributes: map[string]dsschema.Attribute{
						"status": dsschema.StringAttribute{
							Description: "Status value.",
							Computed:    true,
						},
						"silent_timeout": dsschema.Int64Attribute{
							Description: "Minimum alert interval, in seconds.",
							Computed:    true,
						},
					},
				},
			},
			"alert_target": alertTargetDataSourceAttribute(),
			"agg_interval": dsschema.Int64Attribute{
				Description: "Alert aggregation interval, in seconds.",
				Computed:    true,
			},
			"agg_fields": dsschema.ListAttribute{
				Description: "Aggregation field list.",
				Computed:    true,
				ElementType: types.StringType,
			},
			"agg_labels": dsschema.ListAttribute{
				Description: "Label value list when aggregating by label.",
				Computed:    true,
				ElementType: types.StringType,
			},
			"agg_cluster_fields": dsschema.ListAttribute{
				Description: "Field list for smart aggregation.",
				Computed:    true,
				ElementType: types.StringType,
			},
			"agg_send_first": dsschema.BoolAttribute{
				Description: "Whether the first alert is sent directly when aggregating.",
				Computed:    true,
			},
		},
	}
}

func alertTargetDataSourceAttribute() dsschema.ListNestedAttribute {
	return dsschema.ListNestedAttribute{
		Description: "Trigger action configuration.",
		Computed:    true,
		NestedObject: dsschema.NestedAttributeObject{
			Attributes: map[string]dsschema.Attribute{
				"name": dsschema.StringAttribute{
					Description: "Configuration name.",
					Computed:    true,
				},
				"targets": targetDataSourceAttribute(),
				"crontab": dsschema.StringAttribute{
					Description: "Repeated time period start crontab.",
					Computed:    true,
				},
				"crontab_duration": dsschema.Int64Attribute{
					Description: "Repeated time duration in seconds.",
					Computed:    true,
				},
				"custom_date_uuids": dsschema.ListAttribute{
					Description: "Custom notification date UUID list.",
					Computed:    true,
					ElementType: types.StringType,
				},
				"custom_start_time": dsschema.StringAttribute{
					Description: "Daily custom start time.",
					Computed:    true,
				},
				"custom_duration": dsschema.Int64Attribute{
					Description: "Custom time duration in seconds.",
					Computed:    true,
				},
				"alert_info": dsschema.ListNestedAttribute{
					Description: "Member-type alert notification configuration.",
					Computed:    true,
					NestedObject: dsschema.NestedAttributeObject{
						Attributes: map[string]dsschema.Attribute{
							"name": dsschema.StringAttribute{
								Description: "Configuration name.",
								Computed:    true,
							},
							"targets": targetDataSourceAttribute(),
							"filter_string": dsschema.StringAttribute{
								Description: "Filter condition original string.",
								Computed:    true,
							},
							"member_info": dsschema.ListAttribute{
								Description: "Team or member UUID list.",
								Computed:    true,
								ElementType: types.StringType,
							},
						},
					},
				},
			},
		},
	}
}

func targetDataSourceAttribute() dsschema.ListNestedAttribute {
	return dsschema.ListNestedAttribute{
		Description: "Notification target configuration.",
		Computed:    true,
		NestedObject: dsschema.NestedAttributeObject{
			Attributes: map[string]dsschema.Attribute{
				"to": dsschema.ListAttribute{
					Description: "Notification objects/members/teams.",
					Computed:    true,
					ElementType: types.StringType,
				},
				"status": dsschema.StringAttribute{
					Description: "Alert status values.",
					Computed:    true,
				},
				"df_source": dsschema.StringAttribute{
					Description: "Data source, such as security for security monitoring statuses.",
					Computed:    true,
				},
				"upgrade_targets": dsschema.ListNestedAttribute{
					Description: "Upgrade notification target configuration.",
					Computed:    true,
					NestedObject: dsschema.NestedAttributeObject{
						Attributes: map[string]dsschema.Attribute{
							"to": dsschema.ListAttribute{
								Description: "Notification objects/members/teams.",
								Computed:    true,
								ElementType: types.StringType,
							},
							"duration": dsschema.Int64Attribute{
								Description: "Upgrade duration in seconds.",
								Computed:    true,
							},
							"to_way": dsschema.ListAttribute{
								Description: "Notification ways.",
								Computed:    true,
								ElementType: types.StringType,
							},
						},
					},
				},
				"tags": dsschema.MapAttribute{
					Description: "Filter conditions.",
					Computed:    true,
					ElementType: types.ListType{ElemType: types.StringType},
				},
				"filter_string": dsschema.StringAttribute{
					Description: "Filter condition original string.",
					Computed:    true,
				},
			},
		},
	}
}
