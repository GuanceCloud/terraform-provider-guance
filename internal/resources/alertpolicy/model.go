package alertpolicy

import (
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/GuanceCloud/terraform-provider-guance/internal/consts"
)

// alertPolicyResourceModel maps the resource schema data.
type alertPolicyResourceModel struct {
	ID            types.String   `tfsdk:"id"`
	Type          types.String   `tfsdk:"type"`
	Name          types.String   `tfsdk:"name"`
	SilentTimeout types.Int64    `tfsdk:"silent_timeout"`
	AlertTarget   []*AlertTarget `tfsdk:"alert_target"`
}

// GetId returns the ID of the resource.
func (m *alertPolicyResourceModel) GetId() string {
	return m.ID.ValueString()
}

// SetId sets the ID of the resource.
func (m *alertPolicyResourceModel) SetId(id string) {
	m.ID = types.StringValue(id)
}

// GetResourceType returns the type of the resource.
func (m *alertPolicyResourceModel) GetResourceType() string {
	return consts.TypeNameAlertPolicy
}

// AlertTarget maps the resource schema data.
type AlertTarget struct {
	Type          types.String        `tfsdk:"type"`
	Status        []types.String      `tfsdk:"status"`
	MinInterval   types.Int64         `tfsdk:"min_interval"`
	AllowWeekDays []types.Int64       `tfsdk:"allow_week_days"`
	Notification  *NotificationTarget `tfsdk:"notification"`
}

// NotificationTarget maps the resource schema data.
type NotificationTarget struct {
	To []types.String `tfsdk:"to"`
}
