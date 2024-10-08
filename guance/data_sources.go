package guance

import (
	"context"

	"github.com/GuanceCloud/terraform-provider-guance/internal/datasources/members"
	"github.com/GuanceCloud/terraform-provider-guance/internal/datasources/permissions"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

// DataSources defines the data sources implemented in the provider.
func (p *guanceProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		members.NewMembersDataSource,
		permissions.NewPermissionsDataSource,
	}
}
