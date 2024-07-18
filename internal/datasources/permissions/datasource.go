package permissions

import (
	"context"
	_ "embed"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/GuanceCloud/terraform-provider-guance/internal/api"
)

//go:embed README.md
var resourceDocument string

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &permissionsDataSource{}
	_ datasource.DataSourceWithConfigure = &permissionsDataSource{}
)

// NewPermissionsDataSource is a helper function to simplify the provider implementation.
func NewPermissionsDataSource() datasource.DataSource {
	return &permissionsDataSource{}
}

// permissionsDataSource is the data source implementation.
type permissionsDataSource struct {
	client *api.Client
}

// Metadata returns the data source type name.
func (d *permissionsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_permissions"
}

// Schema defines the schema for the data source.
func (d *permissionsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = dataSourceSchema
}

// Configure adds the provider configured client to the data source.
func (d *permissionsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*api.Client)
}

// Read refreshes the Terraform state with the latest data.
func (d *permissionsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state permissionDataSourceModel

	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	results, err := d.client.ReadPermission(state.IsSupportCustomRole.ValueBool())
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to List",
			err.Error(),
		)
		return
	}

	state.Permissions = make([]permissionResourceModel, 0)

	for _, p := range results {
		m := permissionResourceModel{
			Desc:                types.StringValue(p.Desc),
			Disabled:            types.Int64Value(int64(p.Disabled)),
			IsSupportCustomRole: types.Int64Value(int64(p.IsSupportCustomRole)),
			IsSupportGeneral:    types.Int64Value(int64(p.IsSupportGeneral)),
			IsSupportOwner:      types.Int64Value(int64(p.IsSupportOwner)),
			IsSupportReadOnly:   types.Int64Value(int64(p.IsSupportReadOnly)),
			IsSupportWsAdmin:    types.Int64Value(int64(p.IsSupportWsAdmin)),
			Key:                 types.StringValue(p.Key),
			Name:                types.StringValue(p.Name),
		}
		m.Subs = make([]*permissionSubResourceModel, 0)

		for _, s := range p.Subs {
			sub := permissionSubResourceModel{
				Desc:                types.StringValue(s.Desc),
				Disabled:            types.Int64Value(int64(s.Disabled)),
				IsSupportCustomRole: types.Int64Value(int64(s.IsSupportCustomRole)),
				IsSupportGeneral:    types.Int64Value(int64(s.IsSupportGeneral)),
				IsSupportOwner:      types.Int64Value(int64(s.IsSupportOwner)),
				IsSupportReadOnly:   types.Int64Value(int64(s.IsSupportReadOnly)),
				IsSupportWsAdmin:    types.Int64Value(int64(s.IsSupportWsAdmin)),
				Key:                 types.StringValue(s.Key),
				Name:                types.StringValue(s.Name),
			}
			m.Subs = append(m.Subs, &sub)
		}

		state.Permissions = append(state.Permissions, m)
	}

	// Set state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
