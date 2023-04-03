package guance

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"

	"github.com/GuanceCloud/terraform-provider-guance/internal/resources/alertpolicy"
	"github.com/GuanceCloud/terraform-provider-guance/internal/resources/dashboard"
	"github.com/GuanceCloud/terraform-provider-guance/internal/resources/intelligentinspection"
	"github.com/GuanceCloud/terraform-provider-guance/internal/resources/member"
	"github.com/GuanceCloud/terraform-provider-guance/internal/resources/monitor"
	"github.com/GuanceCloud/terraform-provider-guance/internal/resources/mute"
	"github.com/GuanceCloud/terraform-provider-guance/internal/resources/notification"
	"github.com/GuanceCloud/terraform-provider-guance/internal/resources/workspace"
)

// Resources defines the resources implemented in the provider.
func (p *guanceProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		alertpolicy.NewAlertPolicyResource,
		dashboard.NewDashboardResource,
		intelligentinspection.NewIntelligentInspectionResource,
		member.NewMemberResource,
		monitor.NewMonitorResource,
		mute.NewMuteResource,
		notification.NewNotificationResource,
		workspace.NewWorkspaceResource,
	}
}
