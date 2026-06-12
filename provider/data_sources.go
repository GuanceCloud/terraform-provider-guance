package provider

import (
	"context"

	"github.com/GuanceCloud/terraform-provider-guance/internal/datasources/alert_policy"
	"github.com/GuanceCloud/terraform-provider-guance/internal/datasources/alert_policy_notice_date"
	"github.com/GuanceCloud/terraform-provider-guance/internal/datasources/members"
	"github.com/GuanceCloud/terraform-provider-guance/internal/datasources/monitor"
	"github.com/GuanceCloud/terraform-provider-guance/internal/datasources/monitors"
	"github.com/GuanceCloud/terraform-provider-guance/internal/datasources/mute"
	"github.com/GuanceCloud/terraform-provider-guance/internal/datasources/notify_object"
	"github.com/GuanceCloud/terraform-provider-guance/internal/datasources/permissions"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

// DataSources defines the data sources implemented in the provider.
func (p *guanceProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		members.NewMembersDataSource,
		permissions.NewPermissionsDataSource,
		monitor.NewMonitorDataSource,
		monitors.NewMonitorsDataSource,
		notify_object.NewNotifyObjectDataSource,
		alert_policy.NewAlertPolicyDataSource,
		alert_policy_notice_date.NewAlertPolicyNoticeDateDataSource,
		mute.NewMuteDataSource,
		// default_region.NewDefaultRegionDataSource,
	}
}
