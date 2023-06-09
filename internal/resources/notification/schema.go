// Code generated by Guance Cloud Code Generation Pipeline. DO NOT EDIT.

package notification

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var resourceSchema = schema.Schema{
	Description:         "Notification",
	MarkdownDescription: resourceDocument,
	Attributes: map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Description: "The Guance Resource Name (GRN) of cloud resource.",
			Computed:    true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},

		"created_at": schema.StringAttribute{
			Description: "The RFC3339/ISO8601 time string of resource created at.",
			Computed:    true,
		},

		"name": schema.StringAttribute{
			Description: "Notification object name",

			Required: true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.RequiresReplace(),
			},
		},

		"type": schema.StringAttribute{
			Description: "Trigger rule type",

			MarkdownDescription: `
		Trigger rule type, value must be one of: *ding_talk_robot*, *http_request*, *wechat_robot*, *mail_group*, *feishu_robot*, *sms*, other value will be ignored.
		`,
			Required: true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.RequiresReplace(),
			},
			Validators: []validator.String{
				stringvalidator.OneOf("ding_talk_robot", "http_request", "wechat_robot", "mail_group", "feishu_robot", "sms"),
			},
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

		Required: true,
		PlanModifiers: []planmodifier.String{
			stringplanmodifier.RequiresReplace(),
		},
	},

	"secret": schema.StringAttribute{
		Description: "DingTalk Robot Call Secret",

		Required: true,
		PlanModifiers: []planmodifier.String{
			stringplanmodifier.RequiresReplace(),
		},
	},
}

// schemaFeishuRobot maps the resource schema data.
var schemaFeishuRobot = map[string]schema.Attribute{
	"webhook": schema.StringAttribute{
		Description: "Feishu Robot Call Address",

		Required: true,
		PlanModifiers: []planmodifier.String{
			stringplanmodifier.RequiresReplace(),
		},
	},

	"secret": schema.StringAttribute{
		Description: "Feishu Robot Call Secret",

		Required: true,
		PlanModifiers: []planmodifier.String{
			stringplanmodifier.RequiresReplace(),
		},
	},
}

// schemaHTTPRequest maps the resource schema data.
var schemaHTTPRequest = map[string]schema.Attribute{
	"url": schema.StringAttribute{
		Description: "HTTP Call Address",

		Required: true,
		PlanModifiers: []planmodifier.String{
			stringplanmodifier.RequiresReplace(),
		},
	},
}

// schemaMailGroup maps the resource schema data.
var schemaMailGroup = map[string]schema.Attribute{
	"to": schema.ListAttribute{
		Description: "Member Account List",

		Required:    true,
		ElementType: types.StringType,
		PlanModifiers: []planmodifier.List{
			listplanmodifier.RequiresReplace(),
		},
	},
}

// schemaSMS maps the resource schema data.
var schemaSMS = map[string]schema.Attribute{
	"to": schema.ListAttribute{
		Description: "Phone Number List",

		Required:    true,
		ElementType: types.StringType,
		PlanModifiers: []planmodifier.List{
			listplanmodifier.RequiresReplace(),
		},
	},
}

// schemaWeChatRobot maps the resource schema data.
var schemaWeChatRobot = map[string]schema.Attribute{
	"webhook": schema.StringAttribute{
		Description: "Robot Call Address",

		Required: true,
		PlanModifiers: []planmodifier.String{
			stringplanmodifier.RequiresReplace(),
		},
	},
}
