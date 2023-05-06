// Code generated by Guance Cloud Code Generation Pipeline. DO NOT EDIT.

package members

import (
	"context"
	_ "embed"
	"fmt"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/GuanceCloud/terraform-provider-guance/internal/consts"
	"github.com/GuanceCloud/terraform-provider-guance/internal/helpers/tfcodec"
	"github.com/GuanceCloud/terraform-provider-guance/internal/sdk"
)

//go:embed README.md
var resourceDocument string

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &membersDataSource{}
	_ datasource.DataSourceWithConfigure = &membersDataSource{}
)

// NewMembersDataSource is a helper function to simplify the provider implementation.
func NewMembersDataSource() datasource.DataSource {
	return &membersDataSource{}
}

// membersDataSource is the data source implementation.
type membersDataSource struct {
	client *sdk.Client[sdk.Resource]
}

// Metadata returns the data source type name.
func (d *membersDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_members"
}

// Schema defines the schema for the data source.
func (d *membersDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = dataSourceSchema
}

// Configure adds the provider configured client to the data source.
func (d *membersDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*sdk.Client[sdk.Resource])
}

// Read refreshes the Terraform state with the latest data.
func (d *membersDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state memberDataSourceModel
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	results, err := d.client.List(ctx, &sdk.ListOptions{
		MaxResults: state.MaxResults.ValueInt64(),
		TypeName:   consts.TypeNameMember,
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to List",
			err.Error(),
		)
		return
	}

	var mErr error
	var items []*memberResourceModel
	for _, rd := range results.ResourceDescriptions {
		if !sdk.FilterAllSuccess(rd.Properties, state.Filters...) {
			continue
		}

		item := &memberResourceModel{}
		if err := tfcodec.DecodeJSON([]byte(rd.Properties), item); err != nil {
			mErr = multierror.Append(mErr, fmt.Errorf("unable to decode properties: %w", err))
			continue
		}
		item.SetId(rd.Identifier)
		item.SetCreatedAt(rd.CreatedAt)

		items = append(items, item)
	}
	if mErr != nil {
		resp.Diagnostics.AddError(
			"Unable to List resources",
			mErr.Error(),
		)
		return
	}
	state.Items = items
	state.ID = types.StringValue("placeholder")

	// Set state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}