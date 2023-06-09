// Code generated by Guance Cloud Code Generation Pipeline. DO NOT EDIT.

package notification

import (
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/GuanceCloud/terraform-provider-guance/internal/consts"
)

// notificationResourceModel maps the resource schema data.
type notificationResourceModel struct {
	ID            types.String   `tfsdk:"id"`
	CreatedAt     types.String   `tfsdk:"created_at"`
	Name          types.String   `tfsdk:"name"`
	Type          types.String   `tfsdk:"type"`
	DingTalkRobot *DingTalkRobot `tfsdk:"ding_talk_robot"`
	HttpRequest   *HTTPRequest   `tfsdk:"http_request"`
	WechatRobot   *WeChatRobot   `tfsdk:"wechat_robot"`
	MailGroup     *MailGroup     `tfsdk:"mail_group"`
	FeishuRobot   *FeishuRobot   `tfsdk:"feishu_robot"`
	Sms           *SMS           `tfsdk:"sms"`
}

// GetId returns the ID of the resource.
func (m *notificationResourceModel) GetId() string {
	return m.ID.ValueString()
}

// SetId sets the ID of the resource.
func (m *notificationResourceModel) SetId(id string) {
	m.ID = types.StringValue(id)
}

// GetResourceType returns the type of the resource.
func (m *notificationResourceModel) GetResourceType() string {
	return consts.TypeNameNotification
}

// SetCreatedAt sets the creation time of the resource.
func (m *notificationResourceModel) SetCreatedAt(t string) {
	m.CreatedAt = types.StringValue(t)
}

// DingTalkRobot maps the resource schema data.
type DingTalkRobot struct {
	Webhook types.String `tfsdk:"webhook"`
	Secret  types.String `tfsdk:"secret"`
}

// FeishuRobot maps the resource schema data.
type FeishuRobot struct {
	Webhook types.String `tfsdk:"webhook"`
	Secret  types.String `tfsdk:"secret"`
}

// HTTPRequest maps the resource schema data.
type HTTPRequest struct {
	Url types.String `tfsdk:"url"`
}

// MailGroup maps the resource schema data.
type MailGroup struct {
	To []types.String `tfsdk:"to"`
}

// SMS maps the resource schema data.
type SMS struct {
	To []types.String `tfsdk:"to"`
}

// WeChatRobot maps the resource schema data.
type WeChatRobot struct {
	Webhook types.String `tfsdk:"webhook"`
}
