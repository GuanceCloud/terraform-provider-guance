package functions

import (
	"context"
	_ "embed"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/GuanceCloud/terraform-provider-guance/internal/consts"
	"github.com/GuanceCloud/terraform-provider-guance/internal/sdk"
)

//go:embed README.md
var resourceDocument string

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &functionsDataSource{}
	_ datasource.DataSourceWithConfigure = &functionsDataSource{}
)

// NewFunctionsDataSource is a helper function to simplify the provider implementation.
func NewFunctionsDataSource() datasource.DataSource {
	return &functionsDataSource{}
}

// functionsDataSource is the data source implementation.
type functionsDataSource struct {
	client *sdk.Client[*functionResourceModel]
}

// Metadata returns the data source type name.
func (d *functionsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_functions"
}

// Schema defines the schema for the data source.
func (d *functionsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = dataSourceSchema
}

// Configure adds the provider configured client to the data source.
func (d *functionsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*sdk.Client[*functionResourceModel])
}

// Read refreshes the Terraform state with the latest data.
func (d *functionsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state functionDataSourceModel
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	results, err := d.client.List(ctx, &sdk.ListOptions{
		MaxResults: state.MaxResults.ValueInt64(),
		TypeName:   consts.TypeNameFunction,
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to List",
			err.Error(),
		)
		return
	}
	state.Items = results
	state.ID = types.StringValue("placeholder")

	// Set state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
