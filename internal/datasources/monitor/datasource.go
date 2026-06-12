package monitor

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
)

var (
	_ datasource.DataSource                     = &monitorDataSource{}
	_ datasource.DataSourceWithConfigure        = &monitorDataSource{}
	_ datasource.DataSourceWithConfigValidators = &monitorDataSource{}
)

func NewMonitorDataSource() datasource.DataSource {
	return &monitorDataSource{}
}

type monitorDataSource struct {
	client *api.Client
}

func (d *monitorDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_monitor"
}

func (d *monitorDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = dsschema.Schema{
		Description:         "Lookup a monitor by UUID or exact name.",
		MarkdownDescription: "The `guance_monitor` data source reads an existing Guance monitor/checker by `uuid` or exact `name`.",
		Attributes: map[string]dsschema.Attribute{
			"uuid": dsschema.StringAttribute{
				Description: "The UUID of the monitor/checker.",
				Optional:    true,
				Computed:    true,
			},
			"name": dsschema.StringAttribute{
				Description: "The exact monitor/checker name.",
				Optional:    true,
				Computed:    true,
			},
			"type": dsschema.StringAttribute{
				Description: "Monitor type.",
				Optional:    true,
				Computed:    true,
			},
			"status": dsschema.Int64Attribute{
				Description: "Monitor status, 0 enabled and 2 disabled.",
				Computed:    true,
			},
			"extend": dsschema.StringAttribute{
				Description: "Monitor extend payload as canonical JSON.",
				Computed:    true,
			},
			"alert_policy_uuids": dsschema.ListAttribute{
				Description: "Alert policy UUIDs attached to this monitor.",
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
			"secret": dsschema.StringAttribute{
				Description: "Monitor secret.",
				Computed:    true,
			},
			"json_script": dsschema.StringAttribute{
				Description: "Monitor jsonScript payload as canonical JSON.",
				Computed:    true,
			},
			"open_permission_set": dsschema.BoolAttribute{
				Description: "Whether custom operation permissions are enabled.",
				Computed:    true,
			},
			"permission_set": dsschema.ListAttribute{
				Description: "Operation permission configuration.",
				Computed:    true,
				ElementType: types.StringType,
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
			"monitor_uuid": dsschema.StringAttribute{
				Description: "The monitor group UUID.",
				Computed:    true,
			},
			"monitor_name": dsschema.StringAttribute{
				Description: "The monitor group name.",
				Computed:    true,
			},
		},
	}
}

func (d *monitorDataSource) ConfigValidators(_ context.Context) []datasource.ConfigValidator {
	return []datasource.ConfigValidator{
		datasourcevalidator.ExactlyOneOf(path.MatchRoot("uuid"), path.MatchRoot("name")),
	}
}

func (d *monitorDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	d.client = req.ProviderData.(*api.Client)
}

func (d *monitorDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state monitorDataSourceModel
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	content := &api.MonitorContent{}
	if !state.UUID.IsNull() && state.UUID.ValueString() != "" {
		if err := d.client.Read(consts.TypeNameMonitor, state.UUID.ValueString(), content); err != nil {
			resp.Diagnostics.AddError("Error reading monitor", err.Error())
			return
		}
	} else {
		list := &api.MonitorListContent{}
		options := api.MonitorListOptions{
			Search: state.Name.ValueString(),
			Type:   state.Type.ValueString(),
		}
		if err := d.client.ListMonitorsWithOptions(options, list); err != nil {
			resp.Diagnostics.AddError("Error listing monitors", err.Error())
			return
		}
		matched := make([]api.MonitorContent, 0, 1)
		for _, item := range list.Data {
			if monitorName(item) == state.Name.ValueString() {
				matched = append(matched, item)
			}
		}
		if len(matched) != 1 {
			resp.Diagnostics.AddError("Unable to find unique monitor", fmt.Sprintf("found %d monitors with name %q", len(matched), state.Name.ValueString()))
			return
		}
		*content = matched[0]
	}

	if err := stateFromMonitorContent(&state, content); err != nil {
		resp.Diagnostics.AddError("Error reading monitor JSON", err.Error())
		return
	}
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

func stateFromMonitorContent(state *monitorDataSourceModel, content *api.MonitorContent) error {
	state.UUID = types.StringValue(content.UUID)
	state.Name = tfconvert.StringValueOrNull(monitorName(*content))
	state.Type = tfconvert.StringValueOrNull(content.Type)
	state.Status = types.Int64Value(int64(content.Status))
	state.AlertPolicyUUIDs = tfconvert.StringsToTypes(content.AlertPolicyUUIDs)
	state.DashboardUUID = tfconvert.StringValueOrNull(content.DashboardUUID)
	state.Tags = tfconvert.StringsToTypes(content.Tags)
	state.Secret = tfconvert.StringValueOrNull(content.Secret)
	state.OpenPermissionSet = types.BoolValue(content.OpenPermissionSet)
	state.PermissionSet = tfconvert.StringsToTypes(content.PermissionSet)
	state.CreateAt = types.Int64Value(int64(content.CreateAt))
	state.UpdateAt = types.Int64Value(int64(content.UpdateAt))
	state.WorkspaceUUID = types.StringValue(content.WorkspaceUUID)
	state.MonitorUUID = tfconvert.StringValueOrNull(content.MonitorUUID)
	state.MonitorName = tfconvert.StringValueOrNull(content.MonitorName)

	extend, err := canonicalJSONFromValue(content.Extend)
	if err != nil {
		return err
	}
	state.Extend = extend
	jsonScript, err := canonicalJSONFromValue(content.JsonScript)
	if err != nil {
		return err
	}
	state.JsonScript = jsonScript
	return nil
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

func canonicalJSONFromValue(value any) (types.String, error) {
	if value == nil {
		return types.StringNull(), nil
	}
	body, err := tfconvert.CanonicalJSONFromValue(value)
	if err != nil {
		return types.StringNull(), err
	}
	return types.StringValue(body), nil
}
