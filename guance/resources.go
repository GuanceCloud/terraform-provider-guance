package guance

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"

	"github.com/GuanceCloud/terraform-provider-guance/internal/resources/blacklist"
	"github.com/GuanceCloud/terraform-provider-guance/internal/resources/membergroup"
	"github.com/GuanceCloud/terraform-provider-guance/internal/resources/pipeline"
	"github.com/GuanceCloud/terraform-provider-guance/internal/resources/role"
)

// Resources defines the resources implemented in the provider.
func (p *guanceProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		blacklist.NewBlackListResource,
		membergroup.NewMemberGroupResource,
		pipeline.NewPipelineResource,
		role.NewRoleResource,
	}
}
