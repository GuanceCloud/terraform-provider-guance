package default_region

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/GuanceCloud/terraform-provider-guance/internal/api"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &defaultRegionDataSource{}
	_ datasource.DataSourceWithConfigure = &defaultRegionDataSource{}
)

// NewDefaultRegionDataSource is a helper function to simplify the provider implementation.
func NewDefaultRegionDataSource() datasource.DataSource {
	return &defaultRegionDataSource{}
}

// defaultRegionDataSource is the data source implementation.
type defaultRegionDataSource struct {
	client *api.Client
}

// Metadata returns the data source type name.
func (d *defaultRegionDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_default_region"
}

// Schema defines the schema for the data source.
func (d *defaultRegionDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = dataSourceSchema
}

// Configure adds the provider configured client to the data source.
func (d *defaultRegionDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*api.Client)
}

// Read refreshes the Terraform state with the latest data.
func (d *defaultRegionDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state defaultRegionDataSourceModel

	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	results, err := d.client.ReadDefaultRegion()
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to List Default Regions",
			err.Error(),
		)
		return
	}

	state.Regions = make([]defaultRegionResourceModel, 0)

	for _, region := range results {
		r := defaultRegionResourceModel{
			UUID:       types.StringValue(region.UUID),
			Province:   types.StringValue(region.Province),
			City:       types.StringValue(region.City),
			Country:    types.StringValue(region.Country),
			Name:       types.StringValue(region.Name),
			NameEn:     types.StringValue(region.NameEn),
			ExtendInfo: types.StringValue(region.ExtendInfo),
			Internal:   types.BoolValue(region.Internal),
			Keycode:    types.StringValue(region.Keycode),
			Isp:        types.StringValue(region.Isp),
			Status:     types.StringValue(region.Status),
			Region:     types.StringValue(region.Region),
			Owner:      types.StringValue(region.Owner),
			Heartbeat:  types.StringValue(fmt.Sprintf("%d", region.Heartbeat)),
			Company:    types.StringValue(region.Company),
			ExternalId: types.StringValue(region.ExternalId),
			ParentAk:   types.StringValue(region.ParentAk),
			CreateAt:   types.StringValue(fmt.Sprintf("%d", region.CreateAt)),
		}

		state.Regions = append(state.Regions, r)
	}

	// Set state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
