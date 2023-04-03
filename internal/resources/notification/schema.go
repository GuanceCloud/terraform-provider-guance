package notification

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var resourceSchema = schema.Schema{
	Description:         "Notification",
	MarkdownDescription: resourceDocument,
	Attributes: map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Description: "Numeric identifier of the order.",
			Computed:    true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},

		"created_at": schema.StringAttribute{
			Description: "Timestamp of the last Terraform update of the order.",
			Computed:    true,
		},

		"name": schema.StringAttribute{
			Description: "Notification object name",
			Required:    true,
		},

		"type": schema.StringAttribute{
			Description: "Trigger rule type",
			Required:    true,
		},

		"ding_talk_robot": schema.SingleNestedAttribute{
			Description: "DingTalk Robot",

			Optional:   true,
			Attributes: schemaDingTalkRobot,
		},

		"http_request": schema.SingleNestedAttribute{
			Description: "HTTP Request",

			Optional:   true,
			Attributes: schemaHTTPRequest,
		},

		"wechat_robot": schema.SingleNestedAttribute{
			Description: "WeChat Robot",

			Optional:   true,
			Attributes: schemaWeChatRobot,
		},

		"mail_group": schema.SingleNestedAttribute{
			Description: "Mail Group",

			Optional:   true,
			Attributes: schemaMailGroup,
		},

		"feishu_robot": schema.SingleNestedAttribute{
			Description: "Feishu Robot",

			Optional:   true,
			Attributes: schemaFeishuRobot,
		},

		"sms": schema.SingleNestedAttribute{
			Description: "SMS",

			Optional:   true,
			Attributes: schemaSMS,
		},
	},
}

// schemaDingTalkRobot maps the resource schema data.
var schemaDingTalkRobot = map[string]schema.Attribute{
	"webhook": schema.StringAttribute{
		Description: "DingTalk Robot Call Address",
		Required:    true,
	},

	"secret": schema.StringAttribute{
		Description: "DingTalk Robot Call Secret",
		Required:    true,
	},
}

// schemaFeishuRobot maps the resource schema data.
var schemaFeishuRobot = map[string]schema.Attribute{
	"webhook": schema.StringAttribute{
		Description: "Feishu Robot Call Address",
		Required:    true,
	},

	"secret": schema.StringAttribute{
		Description: "Feishu Robot Call Secret",
		Required:    true,
	},
}

// schemaHTTPRequest maps the resource schema data.
var schemaHTTPRequest = map[string]schema.Attribute{
	"url": schema.StringAttribute{
		Description: "HTTP Call Address",
		Required:    true,
	},
}

// schemaMailGroup maps the resource schema data.
var schemaMailGroup = map[string]schema.Attribute{
	"to": schema.ListAttribute{
		Description: "Member Account List",
		Required:    true,
		ElementType: types.StringType,
	},
}

// schemaSMS maps the resource schema data.
var schemaSMS = map[string]schema.Attribute{
	"to": schema.ListAttribute{
		Description: "Phone Number List",
		Required:    true,
		ElementType: types.StringType,
	},
}

// schemaWeChatRobot maps the resource schema data.
var schemaWeChatRobot = map[string]schema.Attribute{
	"webhook": schema.StringAttribute{
		Description: "Robot Call Address",
		Required:    true,
	},
}