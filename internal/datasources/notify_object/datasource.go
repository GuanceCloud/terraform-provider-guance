package notify_object

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/datasourcevalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	dsschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/GuanceCloud/terraform-provider-guance/internal/api"
	"github.com/GuanceCloud/terraform-provider-guance/internal/helpers/tfconvert"
)

var (
	_ datasource.DataSource                     = &notifyObjectDataSource{}
	_ datasource.DataSourceWithConfigure        = &notifyObjectDataSource{}
	_ datasource.DataSourceWithConfigValidators = &notifyObjectDataSource{}
)

func NewNotifyObjectDataSource() datasource.DataSource {
	return &notifyObjectDataSource{}
}

type notifyObjectDataSource struct {
	client *api.Client
}

type notifyObjectDataSourceModel struct {
	UUID              types.String   `tfsdk:"uuid"`
	Type              types.String   `tfsdk:"type"`
	Name              types.String   `tfsdk:"name"`
	OptSet            types.String   `tfsdk:"opt_set"`
	OpenPermissionSet types.Bool     `tfsdk:"open_permission_set"`
	PermissionSet     []types.String `tfsdk:"permission_set"`
	CreateAt          types.Int64    `tfsdk:"create_at"`
	UpdateAt          types.Int64    `tfsdk:"update_at"`
	WorkspaceUUID     types.String   `tfsdk:"workspace_uuid"`
}

func (d *notifyObjectDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_notify_object"
}

func (d *notifyObjectDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = dsschema.Schema{
		Description:         "Lookup a notify object by UUID or exact name.",
		MarkdownDescription: "The `guance_notify_object` data source reads an existing Guance alert notification object by `uuid` or exact `name`.",
		Attributes: map[string]dsschema.Attribute{
			"uuid": dsschema.StringAttribute{
				Description: "The UUID of the notify object.",
				Optional:    true,
				Computed:    true,
			},
			"name": dsschema.StringAttribute{
				Description: "The exact name of the notify object.",
				Optional:    true,
				Computed:    true,
			},
			"type": dsschema.StringAttribute{
				Description: "The type of notify object.",
				Computed:    true,
			},
			"opt_set": dsschema.StringAttribute{
				Description: "The option set of notify object in JSON format.",
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

func (d *notifyObjectDataSource) ConfigValidators(_ context.Context) []datasource.ConfigValidator {
	return []datasource.ConfigValidator{
		datasourcevalidator.ExactlyOneOf(path.MatchRoot("uuid"), path.MatchRoot("name")),
	}
}

func (d *notifyObjectDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	d.client = req.ProviderData.(*api.Client)
}

func (d *notifyObjectDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state notifyObjectDataSourceModel
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	content := &api.NotifyObjectContent{}
	if !state.UUID.IsNull() && state.UUID.ValueString() != "" {
		if err := d.client.GetNotifyObject(state.UUID.ValueString(), content); err != nil {
			resp.Diagnostics.AddError("Error reading notify object", err.Error())
			return
		}
	} else {
		list := &api.NotifyObjectListContent{}
		if err := d.client.ListNotifyObjects(state.Name.ValueString(), list); err != nil {
			resp.Diagnostics.AddError("Error listing notify objects", err.Error())
			return
		}
		matched := make([]api.NotifyObjectContent, 0, 1)
		for _, item := range list.Data {
			if item.Name == state.Name.ValueString() {
				matched = append(matched, item)
			}
		}
		if len(matched) != 1 {
			resp.Diagnostics.AddError("Unable to find unique notify object", fmt.Sprintf("found %d notify objects with name %q", len(matched), state.Name.ValueString()))
			return
		}
		*content = matched[0]
	}

	if err := stateFromNotifyObjectContent(&state, content); err != nil {
		resp.Diagnostics.AddError("Error reading notify object opt_set", err.Error())
		return
	}
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

func stateFromNotifyObjectContent(state *notifyObjectDataSourceModel, content *api.NotifyObjectContent) error {
	state.UUID = types.StringValue(content.UUID)
	state.Type = types.StringValue(content.Type)
	state.Name = types.StringValue(content.Name)
	if content.OptSet != nil {
		optSet, err := tfconvert.CanonicalJSONFromValue(content.OptSet)
		if err != nil {
			return err
		}
		state.OptSet = types.StringValue(optSet)
	}
	state.OpenPermissionSet = types.BoolValue(content.OpenPermissionSet)
	state.PermissionSet = tfconvert.StringsToTypes(content.PermissionSet)
	state.CreateAt = types.Int64Value(int64(content.CreateAt))
	state.UpdateAt = types.Int64Value(int64(content.UpdateAt))
	state.WorkspaceUUID = types.StringValue(content.WorkspaceUUID)
	return nil
}
