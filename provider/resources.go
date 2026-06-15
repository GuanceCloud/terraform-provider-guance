package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"

	"github.com/GuanceCloud/terraform-provider-guance/internal/resources/alert_policy"
	"github.com/GuanceCloud/terraform-provider-guance/internal/resources/alert_policy_notice_date"
	"github.com/GuanceCloud/terraform-provider-guance/internal/resources/blacklist"
	"github.com/GuanceCloud/terraform-provider-guance/internal/resources/dashboard"
	"github.com/GuanceCloud/terraform-provider-guance/internal/resources/membergroup"
	"github.com/GuanceCloud/terraform-provider-guance/internal/resources/monitor"
	"github.com/GuanceCloud/terraform-provider-guance/internal/resources/monitor_json"
	"github.com/GuanceCloud/terraform-provider-guance/internal/resources/mute"
	"github.com/GuanceCloud/terraform-provider-guance/internal/resources/notify_object"
	"github.com/GuanceCloud/terraform-provider-guance/internal/resources/pipeline"
	"github.com/GuanceCloud/terraform-provider-guance/internal/resources/role"
)

// Resources defines the resources implemented in the provider.
func (p *guanceProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		alert_policy.NewAlertPolicyResource,
		alert_policy_notice_date.NewAlertPolicyNoticeDateResource,
		blacklist.NewBlackListResource,
		// custom_region.NewCustomRegionResource,
		dashboard.NewDashboardResource,
		membergroup.NewMemberGroupResource,
		monitor.NewMonitorResource,
		monitor_json.NewMonitorJsonResource,
		mute.NewMuteResource,
		notify_object.NewNotifyObjectResource,
		pipeline.NewPipelineResource,
		role.NewRoleResource,
		// slo.NewSloResource,
		// synthetics_test.NewSyntheticsTestResource,
	}
}
