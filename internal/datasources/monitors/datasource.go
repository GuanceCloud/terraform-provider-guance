package monitors

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	dsschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/GuanceCloud/terraform-provider-guance/internal/api"
	"github.com/GuanceCloud/terraform-provider-guance/internal/helpers/tfconvert"
)

var (
	_ datasource.DataSource              = &monitorsDataSource{}
	_ datasource.DataSourceWithConfigure = &monitorsDataSource{}
)

func NewMonitorsDataSource() datasource.DataSource {
	return &monitorsDataSource{}
}

type monitorsDataSource struct {
	client *api.Client
}

func (d *monitorsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_monitors"
}

func (d *monitorsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = dsschema.Schema{
		Description:         "List monitors by search and optional filters.",
		MarkdownDescription: "The `guance_monitors` data source lists Guance monitors/checkers by search and optional filters.",
		Attributes: map[string]dsschema.Attribute{
			"search": dsschema.StringAttribute{
				Description: "Monitor search keyword.",
				Optional:    true,
			},
			"type": dsschema.StringAttribute{
				Description: "Checker type filter, such as simpleCheck. Leave empty to list all monitor/checker types.",
				Optional:    true,
			},
			"status": dsschema.StringAttribute{
				Description: "Monitor status filter, such as 0 or 2.",
				Optional:    true,
			},
			"tags_uuid": dsschema.StringAttribute{
				Description: "Comma-separated monitor tag UUID filter.",
				Optional:    true,
			},
			"alert_policy_uuid": dsschema.StringAttribute{
				Description: "Alert policy UUID filter.",
				Optional:    true,
			},
			"dashboard_uuid": dsschema.StringAttribute{
				Description: "Dashboard UUID filter.",
				Optional:    true,
			},
			"checker_uuid": dsschema.StringAttribute{
				Description: "Comma-separated checker UUID filter.",
				Optional:    true,
			},
			"monitors": dsschema.ListNestedAttribute{
				Description: "Matched monitors.",
				Computed:    true,
				NestedObject: dsschema.NestedAttributeObject{
					Attributes: map[string]dsschema.Attribute{
						"uuid": dsschema.StringAttribute{
							Description: "Monitor/checker UUID.",
							Computed:    true,
						},
						"name": dsschema.StringAttribute{
							Description: "Monitor/checker name.",
							Computed:    true,
						},
						"type": dsschema.StringAttribute{
							Description: "Monitor type, such as trigger or smartMonitor.",
							Computed:    true,
						},
						"status": dsschema.Int64Attribute{
							Description: "Monitor status.",
							Computed:    true,
						},
						"alert_policy_uuids": dsschema.ListAttribute{
							Description: "Attached alert policy UUIDs.",
							Computed:    true,
							ElementType: types.StringType,
						},
						"dashboard_uuid": dsschema.StringAttribute{
							Description: "Associated dashboard UUID.",
							Computed:    true,
						},
						"tags": dsschema.ListAttribute{
							Description: "Monitor tag names.",
							Computed:    true,
							ElementType: types.StringType,
						},
						"workspace_uuid": dsschema.StringAttribute{
							Description: "Workspace UUID.",
							Computed:    true,
						},
						"monitor_uuid": dsschema.StringAttribute{
							Description: "Monitor group UUID.",
							Computed:    true,
						},
						"monitor_name": dsschema.StringAttribute{
							Description: "Monitor group name.",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func (d *monitorsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	d.client = req.ProviderData.(*api.Client)
}

func (d *monitorsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state monitorsDataSourceModel
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	list := &api.MonitorListContent{}
	options := api.MonitorListOptions{
		Search:          state.Search.ValueString(),
		Type:            state.Type.ValueString(),
		Status:          state.Status.ValueString(),
		TagsUUID:        state.TagsUUID.ValueString(),
		AlertPolicyUUID: state.AlertPolicyUUID.ValueString(),
		DashboardUUID:   state.DashboardUUID.ValueString(),
		CheckerUUID:     state.CheckerUUID.ValueString(),
	}
	if err := d.client.ListMonitorsWithOptions(options, list); err != nil {
		resp.Diagnostics.AddError("Error listing monitors", err.Error())
		return
	}

	state.Monitors = make([]monitorListModel, 0, len(list.Data))
	for _, item := range list.Data {
		state.Monitors = append(state.Monitors, monitorFromContent(item))
	}
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

func monitorFromContent(content api.MonitorContent) monitorListModel {
	return monitorListModel{
		UUID:             types.StringValue(content.UUID),
		Name:             tfconvert.StringValueOrNull(monitorName(content)),
		Type:             tfconvert.StringValueOrNull(content.Type),
		Status:           types.Int64Value(int64(content.Status)),
		AlertPolicyUUIDs: tfconvert.StringsToTypes(content.AlertPolicyUUIDs),
		DashboardUUID:    tfconvert.StringValueOrNull(content.DashboardUUID),
		Tags:             tfconvert.StringsToTypes(content.Tags),
		WorkspaceUUID:    types.StringValue(content.WorkspaceUUID),
		MonitorUUID:      tfconvert.StringValueOrNull(content.MonitorUUID),
		MonitorName:      tfconvert.StringValueOrNull(content.MonitorName),
	}
}

func monitorName(content api.MonitorContent) string {
	if content.MonitorName != "" {
		return content.MonitorName
	}
	if script, ok := content.JsonScript.(map[string]any); ok {
		if title, ok := script["title"].(string); ok && title != "" {
			return title
		}
		if name, ok := script["name"].(string); ok && name != "" {
			return name
		}
	}
	return ""
}
