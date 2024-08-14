package members

import (
	"context"
	_ "embed"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/GuanceCloud/terraform-provider-guance/internal/api"
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
	client *api.Client
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

	d.client = req.ProviderData.(*api.Client)
}

// Read refreshes the Terraform state with the latest data.
func (d *membersDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state memberDataSourceModel

	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	results, err := d.client.ReadMember(state.Search.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to List",
			err.Error(),
		)
		return
	}

	state.Members = make([]memberResourceModel, 0)

	for _, member := range results {
		m := memberResourceModel{
			CreateAt: types.StringValue(fmt.Sprintf("%d", member.CreateAt)),
			Email:    types.StringValue(member.Email),
			Name:     types.StringValue(member.Name),
			UUID:     types.StringValue(member.UUID),
		}

		for _, r := range member.Roles {
			m.Roles = append(m.Roles, roleModel{
				Name: types.StringValue(r.Name),
				UUID: types.StringValue(r.UUID),
			})
		}
		tflog.Info(ctx, "member", map[string]interface{}{"member": member})
		state.Members = append(state.Members, m)
	}

	// Set state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
