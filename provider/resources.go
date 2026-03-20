package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"

	"github.com/GuanceCloud/terraform-provider-guance/internal/resources/blacklist"
	"github.com/GuanceCloud/terraform-provider-guance/internal/resources/dashboard"
	"github.com/GuanceCloud/terraform-provider-guance/internal/resources/membergroup"
	"github.com/GuanceCloud/terraform-provider-guance/internal/resources/monitor_json"
	"github.com/GuanceCloud/terraform-provider-guance/internal/resources/pipeline"
	"github.com/GuanceCloud/terraform-provider-guance/internal/resources/role"
)

// Resources defines the resources implemented in the provider.
func (p *guanceProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		// alert_policy.NewAlertPolicyResource,
		blacklist.NewBlackListResource,
		// custom_region.NewCustomRegionResource,
		dashboard.NewDashboardResource,
		membergroup.NewMemberGroupResource,
		monitor_json.NewMonitorJsonResource,
		// monitor.NewMonitorResource,
		// notify_object.NewNotifyObjectResource,
		pipeline.NewPipelineResource,
		role.NewRoleResource,
		// slo.NewSloResource,
		// synthetics_test.NewSyntheticsTestResource,
	}
}
