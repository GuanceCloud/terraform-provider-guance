package alert_policy_notice_date

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
	_ datasource.DataSource                     = &alertPolicyNoticeDateDataSource{}
	_ datasource.DataSourceWithConfigure        = &alertPolicyNoticeDateDataSource{}
	_ datasource.DataSourceWithConfigValidators = &alertPolicyNoticeDateDataSource{}
)

func NewAlertPolicyNoticeDateDataSource() datasource.DataSource {
	return &alertPolicyNoticeDateDataSource{}
}

type alertPolicyNoticeDateDataSource struct {
	client *api.Client
}

type alertPolicyNoticeDateDataSourceModel struct {
	UUID          types.String   `tfsdk:"uuid"`
	Name          types.String   `tfsdk:"name"`
	NoticeDates   []types.String `tfsdk:"notice_dates"`
	CreateAt      types.Int64    `tfsdk:"create_at"`
	UpdateAt      types.Int64    `tfsdk:"update_at"`
	WorkspaceUUID types.String   `tfsdk:"workspace_uuid"`
}

func (d *alertPolicyNoticeDateDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_alert_policy_notice_date"
}

func (d *alertPolicyNoticeDateDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = dsschema.Schema{
		Description:         "Lookup an alert policy custom notice date by UUID or exact name.",
		MarkdownDescription: "The `guance_alert_policy_notice_date` data source reads an existing custom notice date by `uuid` or exact `name`.",
		Attributes: map[string]dsschema.Attribute{
			"uuid": dsschema.StringAttribute{
				Description: "The UUID of the notice date.",
				Optional:    true,
				Computed:    true,
			},
			"name": dsschema.StringAttribute{
				Description: "The exact name of the notice date.",
				Optional:    true,
				Computed:    true,
			},
			"notice_dates": dsschema.ListAttribute{
				Description: "Custom notice dates.",
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
		},
	}
}

func (d *alertPolicyNoticeDateDataSource) ConfigValidators(_ context.Context) []datasource.ConfigValidator {
	return []datasource.ConfigValidator{
		datasourcevalidator.ExactlyOneOf(path.MatchRoot("uuid"), path.MatchRoot("name")),
	}
}

func (d *alertPolicyNoticeDateDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	d.client = req.ProviderData.(*api.Client)
}

func (d *alertPolicyNoticeDateDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state alertPolicyNoticeDateDataSourceModel
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	content := &api.AlertPolicyNoticeDateContent{}
	if !state.UUID.IsNull() && state.UUID.ValueString() != "" {
		if err := d.client.Read(consts.TypeNameAlertPolicyNoticeDate, state.UUID.ValueString(), content); err != nil {
			resp.Diagnostics.AddError("Error reading alert policy notice date", err.Error())
			return
		}
	} else {
		list := &api.AlertPolicyNoticeDateListContent{}
		if err := d.client.ListAlertPolicyNoticeDates(state.Name.ValueString(), list); err != nil {
			resp.Diagnostics.AddError("Error listing alert policy notice dates", err.Error())
			return
		}
		matched := make([]api.AlertPolicyNoticeDateContent, 0, 1)
		for _, item := range list.Data {
			if item.Name == state.Name.ValueString() {
				matched = append(matched, item)
			}
		}
		if len(matched) != 1 {
			resp.Diagnostics.AddError("Unable to find unique alert policy notice date", fmt.Sprintf("found %d notice dates with name %q", len(matched), state.Name.ValueString()))
			return
		}
		*content = matched[0]
	}

	state.UUID = types.StringValue(content.UUID)
	state.Name = types.StringValue(content.Name)
	state.NoticeDates = tfconvert.StringsToTypes(content.Dates)
	state.CreateAt = types.Int64Value(int64(content.CreateAt))
	state.UpdateAt = types.Int64Value(int64(content.UpdateAt))
	state.WorkspaceUUID = types.StringValue(content.WorkspaceUUID)
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}
