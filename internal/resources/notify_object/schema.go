package notify_object

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var resourceSchema = schema.Schema{
	Description:         "Notify Object",
	MarkdownDescription: resourceDocument,
	Attributes: map[string]schema.Attribute{
		"uuid": schema.StringAttribute{
			Description: "The uuid of resource.",
			Computed:    true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"type": schema.StringAttribute{
			Description: "The type of notify object.",
			Required:    true,
			Validators: []validator.String{
				stringvalidator.OneOf(
					"dingTalkRobot",
					"HTTPRequest",
					"wechatRobot",
					"mailGroup",
					"feishuRobot",
					"sms",
					"vms",
					"simpleHTTPRequest",
					"slackIncomingWebhook",
					"teamsWorkflowWebhook",
					"googleChatWebhook",
				),
			},
		},
		"name": schema.StringAttribute{
			Description: "The name of notify object.",
			Required:    true,
			Validators: []validator.String{
				stringvalidator.LengthAtMost(128),
			},
		},
		"opt_set": schema.StringAttribute{
			Description: "The option set of notify object in JSON format.",
			Required:    true,
		},
		"open_permission_set": schema.BoolAttribute{
			Description: "Whether to open permission set.",
			Optional:    true,
		},
		"permission_set": schema.ListAttribute{
			Description: "The permission set of resource.",
			Optional:    true,
			ElementType: types.StringType,
		},
		"create_at": schema.Int64Attribute{
			Description: "The timestamp seconds of the resource created at.",
			Computed:    true,
			PlanModifiers: []planmodifier.Int64{
				int64planmodifier.UseStateForUnknown(),
			},
		},
		"update_at": schema.Int64Attribute{
			Description: "The timestamp seconds of the resource updated at.",
			Computed:    true,
		},
		"workspace_uuid": schema.StringAttribute{
			Description: "The uuid of the workspace.",
			Computed:    true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
	},
}
