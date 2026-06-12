package mute

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/datasourcevalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	dsschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/GuanceCloud/terraform-provider-guance/internal/api"
	resourcemute "github.com/GuanceCloud/terraform-provider-guance/internal/resources/mute"
)

var (
	_ datasource.DataSource                     = &muteDataSource{}
	_ datasource.DataSourceWithConfigure        = &muteDataSource{}
	_ datasource.DataSourceWithConfigValidators = &muteDataSource{}
)

func NewMuteDataSource() datasource.DataSource {
	return &muteDataSource{}
}

type muteDataSource struct {
	client *api.Client
}

type muteDataSourceModel struct {
	UUID             types.String                   `tfsdk:"uuid"`
	Name             types.String                   `tfsdk:"name"`
	Description      types.String                   `tfsdk:"description"`
	Type             types.String                   `tfsdk:"type"`
	WorkStatus       types.String                   `tfsdk:"work_status"`
	IsEnable         types.String                   `tfsdk:"is_enable"`
	Creator          types.String                   `tfsdk:"creator"`
	Updator          types.String                   `tfsdk:"updator"`
	MuteRanges       []resourcemute.MuteRange       `tfsdk:"mute_ranges"`
	Tags             map[string][]string            `tfsdk:"tags"`
	FilterString     types.String                   `tfsdk:"filter_string"`
	NotifyTargets    []resourcemute.NotifyTarget    `tfsdk:"notify_targets"`
	NotifyMessage    types.String                   `tfsdk:"notify_message"`
	NotifyTimeStr    types.String                   `tfsdk:"notify_time_str"`
	StartTime        types.String                   `tfsdk:"start_time"`
	EndTime          types.String                   `tfsdk:"end_time"`
	RepeatTimeSet    types.Int64                    `tfsdk:"repeat_time_set"`
	RepeatCrontabSet *resourcemute.RepeatCrontabSet `tfsdk:"repeat_crontab_set"`
	CrontabDuration  types.Int64                    `tfsdk:"crontab_duration"`
	RepeatExpireTime types.String                   `tfsdk:"repeat_expire_time"`
	Timezone         types.String                   `tfsdk:"timezone"`
	Declaration      map[string]string              `tfsdk:"declaration"`
	Enabled          types.Bool                     `tfsdk:"enabled"`
	Status           types.Int64                    `tfsdk:"status"`
	CreateAt         types.Int64                    `tfsdk:"create_at"`
	UpdateAt         types.Int64                    `tfsdk:"update_at"`
	WorkspaceUUID    types.String                   `tfsdk:"workspace_uuid"`
}

func (d *muteDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_mute"
}

func (d *muteDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = dsschema.Schema{
		Description:         "Lookup a mute rule by UUID or exact name.",
		MarkdownDescription: "The `guance_mute` data source reads an existing Guance mute rule by `uuid` or exact `name`.",
		Attributes: map[string]dsschema.Attribute{
			"uuid": dsschema.StringAttribute{
				Description: "The UUID of the mute rule.",
				Optional:    true,
				Computed:    true,
			},
			"name": dsschema.StringAttribute{
				Description: "The exact name of the mute rule.",
				Optional:    true,
				Computed:    true,
			},
			"description": dsschema.StringAttribute{
				Description: "The description of the mute rule.",
				Computed:    true,
			},
			"type": dsschema.StringAttribute{
				Description: "Mute rule type. When configured with name lookup, filters the mute list by type.",
				Optional:    true,
				Computed:    true,
			},
			"work_status": dsschema.StringAttribute{
				Description: "Filter mute rules by work status when looking up by name.",
				Optional:    true,
			},
			"is_enable": dsschema.StringAttribute{
				Description: "Filter mute rules by enabled flag when looking up by name.",
				Optional:    true,
			},
			"creator": dsschema.StringAttribute{
				Description: "Filter mute rules by creator when looking up by name.",
				Optional:    true,
			},
			"updator": dsschema.StringAttribute{
				Description: "Filter mute rules by updator when looking up by name.",
				Optional:    true,
			},
			"mute_ranges": dsschema.ListNestedAttribute{
				Description: "Mute ranges.",
				Computed:    true,
				NestedObject: dsschema.NestedAttributeObject{
					Attributes: map[string]dsschema.Attribute{
						"name": dsschema.StringAttribute{
							Description: "The display name of the muted resource.",
							Computed:    true,
						},
						"type": dsschema.StringAttribute{
							Description: "The muted resource type returned by the API.",
							Computed:    true,
						},
						"checker_uuid": dsschema.StringAttribute{
							Description: "Monitor/checker UUID.",
							Computed:    true,
						},
						"monitor_uuid": dsschema.StringAttribute{
							Description: "Monitor UUID.",
							Computed:    true,
						},
						"slo_uuid": dsschema.StringAttribute{
							Description: "SLO UUID.",
							Computed:    true,
						},
						"alert_policy_uuid": dsschema.StringAttribute{
							Description: "Alert policy UUID.",
							Computed:    true,
						},
						"tag_uuid": dsschema.StringAttribute{
							Description: "Monitor tag UUID.",
							Computed:    true,
						},
					},
				},
			},
			"tags": dsschema.MapAttribute{
				Description: "Event attribute filters.",
				Computed:    true,
				ElementType: types.ListType{ElemType: types.StringType},
			},
			"filter_string": dsschema.StringAttribute{
				Description: "Event attribute filter expression.",
				Computed:    true,
			},
			"notify_targets": dsschema.ListNestedAttribute{
				Description: "Notification targets for mute notifications.",
				Computed:    true,
				NestedObject: dsschema.NestedAttributeObject{
					Attributes: map[string]dsschema.Attribute{
						"to": dsschema.ListAttribute{
							Description: "Notification target UUIDs.",
							Computed:    true,
							ElementType: types.StringType,
						},
						"type": dsschema.StringAttribute{
							Description: "Notification target type.",
							Computed:    true,
						},
					},
				},
			},
			"notify_message": dsschema.StringAttribute{
				Description: "Notification message.",
				Computed:    true,
			},
			"notify_time_str": dsschema.StringAttribute{
				Description: "Notification time.",
				Computed:    true,
			},
			"start_time": dsschema.StringAttribute{
				Description: "One-time mute start time.",
				Computed:    true,
			},
			"end_time": dsschema.StringAttribute{
				Description: "One-time mute end time.",
				Computed:    true,
			},
			"repeat_time_set": dsschema.Int64Attribute{
				Description: "Whether the mute rule repeats.",
				Computed:    true,
			},
			"repeat_crontab_set": dsschema.SingleNestedAttribute{
				Description: "Repeated mute crontab settings.",
				Computed:    true,
				Attributes: map[string]dsschema.Attribute{
					"min": dsschema.StringAttribute{
						Description: "Minute field.",
						Computed:    true,
					},
					"hour": dsschema.StringAttribute{
						Description: "Hour field.",
						Computed:    true,
					},
					"day": dsschema.StringAttribute{
						Description: "Day field.",
						Computed:    true,
					},
					"month": dsschema.StringAttribute{
						Description: "Month field.",
						Computed:    true,
					},
					"week": dsschema.StringAttribute{
						Description: "Week field.",
						Computed:    true,
					},
				},
			},
			"crontab_duration": dsschema.Int64Attribute{
				Description: "Repeated mute duration in seconds.",
				Computed:    true,
			},
			"repeat_expire_time": dsschema.StringAttribute{
				Description: "Repeated mute expiration time.",
				Computed:    true,
			},
			"timezone": dsschema.StringAttribute{
				Description: "Mute rule timezone.",
				Computed:    true,
			},
			"declaration": dsschema.MapAttribute{
				Description: "Custom declaration information.",
				Computed:    true,
				ElementType: types.StringType,
			},
			"enabled": dsschema.BoolAttribute{
				Description: "Whether the mute rule is enabled. API status 0 maps to true, and status 2 maps to false.",
				Computed:    true,
			},
			"status": dsschema.Int64Attribute{
				Description: "Mute rule status returned by the API.",
				Computed:    true,
			},
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

func (d *muteDataSource) ConfigValidators(_ context.Context) []datasource.ConfigValidator {
	return []datasource.ConfigValidator{
		datasourcevalidator.ExactlyOneOf(path.MatchRoot("uuid"), path.MatchRoot("name")),
	}
}

func (d *muteDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	d.client = req.ProviderData.(*api.Client)
}

func (d *muteDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state muteDataSourceModel
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	content := &api.MuteContent{}
	if !state.UUID.IsNull() && state.UUID.ValueString() != "" {
		if err := d.client.GetMute(state.UUID.ValueString(), content); err != nil {
			resp.Diagnostics.AddError("Error reading mute rule", err.Error())
			return
		}
	} else {
		list := &api.MuteListContent{}
		options := api.MuteListOptions{
			Search:     state.Name.ValueString(),
			WorkStatus: state.WorkStatus.ValueString(),
			IsEnable:   state.IsEnable.ValueString(),
			Type:       state.Type.ValueString(),
			Creator:    state.Creator.ValueString(),
			Updator:    state.Updator.ValueString(),
		}
		if err := d.client.ListMutesWithOptions(options, list); err != nil {
			resp.Diagnostics.AddError("Error listing mute rules", err.Error())
			return
		}
		matched := make([]api.MuteContent, 0, 1)
		for _, raw := range list.Data {
			var item struct {
				Name string `json:"name,omitempty"`
			}
			if err := json.Unmarshal(raw, &item); err != nil {
				continue
			}
			if item.Name == state.Name.ValueString() {
				var content api.MuteContent
				if err := json.Unmarshal(raw, &content); err != nil {
					resp.Diagnostics.AddError("Error decoding mute rule", err.Error())
					return
				}
				matched = append(matched, content)
			}
		}
		if len(matched) != 1 {
			resp.Diagnostics.AddError("Unable to find unique mute rule", fmt.Sprintf("found %d mute rules with name %q", len(matched), state.Name.ValueString()))
			return
		}
		*content = matched[0]
	}

	state.UUID = types.StringValue(content.UUID)
	state.Name = types.StringValue(content.Name)
	state.Description = resourcemute.StringPointerValue(content.Description)
	state.Type = types.StringValue(content.Type)
	state.MuteRanges = resourcemute.MuteRangesFromContent(content.MuteRanges, nil, true)
	state.Tags = content.Tags
	state.FilterString = resourcemute.StringPointerValue(content.FilterString)
	state.NotifyTargets = resourcemute.NotifyTargetsFromContent(content.NotifyTargets, nil, true)
	state.NotifyMessage = resourcemute.StringPointerValue(content.NotifyMessage)
	state.NotifyTimeStr = resourcemute.StringPointerValue(content.NotifyTimeStr)
	state.StartTime = resourcemute.StringPointerValue(content.StartTime)
	state.EndTime = resourcemute.StringPointerValue(content.EndTime)
	state.RepeatTimeSet = types.Int64Value(int64(resourcemute.RepeatTimeSetFromContent(content)))
	if content.RepeatCrontabSet != nil {
		state.RepeatCrontabSet = resourcemute.RepeatCrontabSetFromContent(content.RepeatCrontabSet)
	}
	state.CrontabDuration = types.Int64Value(int64(content.CrontabDuration))
	state.RepeatExpireTime = resourcemute.RepeatExpireTimeValueOrExisting(content.RepeatExpireTime, types.StringNull())
	state.Timezone = resourcemute.StringPointerValue(content.Timezone)
	state.Declaration = resourcemute.DeclarationFromContent(content.Declaration)
	state.Status = types.Int64Value(int64(content.Status))
	state.Enabled = resourcemute.MuteEnabledFromStatus(content.Status, state.Enabled)
	state.CreateAt = types.Int64Value(int64(content.CreateAt))
	state.UpdateAt = types.Int64Value(int64(content.UpdateAt))
	state.WorkspaceUUID = types.StringValue(content.WorkspaceUUID)
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}
