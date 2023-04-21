// Code generated by Guance Cloud Code Generation Pipeline. DO NOT EDIT.

package mute

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
)

var resourceSchema = schema.Schema{
	Description:         "Mute Rule",
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

		"mute_ranges": schema.ListNestedAttribute{
			Description: "Mute Ranges",

			Optional: true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: schemaMuteRange,
			},
		},

		"notify": schema.SingleNestedAttribute{
			Description: "Notify Options",

			Optional:   true,
			Attributes: schemaNotifyOptions,
		},

		"notify_targets": schema.ListNestedAttribute{
			Description: "Notify targets",

			Optional: true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: schemaNotifyTarget,
			},
		},

		"onetime": schema.SingleNestedAttribute{
			Description: "Onetime",

			Optional:   true,
			Attributes: schemaOnetimeOptions,
		},

		"repeat": schema.SingleNestedAttribute{
			Description: "Repeat",

			Optional:   true,
			Attributes: schemaRepeatOptions,
		},

		"mute_tags": schema.ListNestedAttribute{
			Description: "Tags",

			Optional: true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: schemaTag,
			},
		},
	},
}

// schemaMuteRange maps the resource schema data.
var schemaMuteRange = map[string]schema.Attribute{
	"type": schema.StringAttribute{
		Description: "Mute Range Type",
		Required:    true,
	},

	"monitor": schema.SingleNestedAttribute{
		Description: "Monitor configuration",

		Optional:   true,
		Attributes: schemaMuteRangeMonitor,
	},

	"alert_policy": schema.SingleNestedAttribute{
		Description: "Alert Policy configuration",

		Optional:   true,
		Attributes: schemaMuteRangeAlertPolicy,
	},
}

// schemaMuteRangeAlertPolicy maps the resource schema data.
var schemaMuteRangeAlertPolicy = map[string]schema.Attribute{
	"id": schema.StringAttribute{
		Description: "Alert Policy ID",
		Required:    true,
	},
}

// schemaMuteRangeMonitor maps the resource schema data.
var schemaMuteRangeMonitor = map[string]schema.Attribute{
	"id": schema.StringAttribute{
		Description: "Monitor ID",
		Required:    true,
	},
}

// schemaNotifyOptions maps the resource schema data.
var schemaNotifyOptions = map[string]schema.Attribute{
	"message": schema.StringAttribute{
		Description: "Notify Message",

		Optional: true,
	},

	"before_time": schema.StringAttribute{
		Description: "Notify Time",

		Optional: true,
	},
}

// schemaNotifyTarget maps the resource schema data.
var schemaNotifyTarget = map[string]schema.Attribute{
	"type": schema.StringAttribute{
		Description: "Notify Type",
		Required:    true,
	},

	"member_group": schema.SingleNestedAttribute{
		Description: "MemberGroup",

		Optional:   true,
		Attributes: schemaNotifyTargetMemberGroup,
	},

	"notification": schema.SingleNestedAttribute{
		Description: "Notification",

		Optional:   true,
		Attributes: schemaNotifyTargetNotification,
	},
}

// schemaNotifyTargetMemberGroup maps the resource schema data.
var schemaNotifyTargetMemberGroup = map[string]schema.Attribute{
	"id": schema.StringAttribute{
		Description: "MemberGroup ID",
		Required:    true,
	},
}

// schemaNotifyTargetNotification maps the resource schema data.
var schemaNotifyTargetNotification = map[string]schema.Attribute{
	"id": schema.StringAttribute{
		Description: "Notification ID",
		Required:    true,
	},
}

// schemaOnetimeOptions maps the resource schema data.
var schemaOnetimeOptions = map[string]schema.Attribute{
	"start": schema.StringAttribute{
		Description: "Start",

		Optional: true,
	},

	"end": schema.StringAttribute{
		Description: "End",

		Optional: true,
	},
}

// schemaRepeatCrontabSet maps the resource schema data.
var schemaRepeatCrontabSet = map[string]schema.Attribute{
	"min": schema.StringAttribute{
		Description: "Min",

		Optional: true,
	},

	"hour": schema.StringAttribute{
		Description: "Hour",

		Optional: true,
	},

	"day": schema.StringAttribute{
		Description: "Day",

		Optional: true,
	},

	"month": schema.StringAttribute{
		Description: "Month",

		Optional: true,
	},

	"week": schema.StringAttribute{
		Description: "Week",

		Optional: true,
	},
}

// schemaRepeatOptions maps the resource schema data.
var schemaRepeatOptions = map[string]schema.Attribute{
	"start": schema.StringAttribute{
		Description: "Start",

		Optional: true,
	},

	"end": schema.StringAttribute{
		Description: "End",

		Optional: true,
	},

	"crontab_duration": schema.StringAttribute{
		Description: "Crontab Duration",

		Optional: true,
	},

	"expire": schema.StringAttribute{
		Description: "Repeat Expire",

		Optional: true,
	},

	"crontab": schema.SingleNestedAttribute{
		Description: "Crontab configuration",

		Optional:   true,
		Attributes: schemaRepeatCrontabSet,
	},
}

// schemaTag maps the resource schema data.
var schemaTag = map[string]schema.Attribute{
	"key": schema.StringAttribute{
		Description: "Tag",
		Required:    true,
	},

	"value": schema.StringAttribute{
		Description: "Tag Value",
		Required:    true,
	},
}
