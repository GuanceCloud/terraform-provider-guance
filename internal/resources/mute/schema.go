package mute

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
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

		"mute_ranges": schema.ListAttribute{
			Description: "Mute Ranges",

			Optional:    true,
			ElementType: types.StringType,
		},

		"tags": schema.SingleNestedAttribute{
			Description: "Tags",

			Optional:   true,
			Attributes: schemaTag,
		},

		"start": schema.Int64Attribute{
			Description: "Start",

			Optional: true,
		},

		"end": schema.Int64Attribute{
			Description: "End",

			Optional: true,
		},

		"notify": schema.SingleNestedAttribute{
			Description: "Notify Options",

			Optional:   true,
			Attributes: schemaNotifyOptions,
		},

		"repeat": schema.SingleNestedAttribute{
			Description: "Repeat",

			Optional:   true,
			Attributes: schemaRepeatOptions,
		},
	},
}

// schemaNotifyOptions maps the resource schema data.
var schemaNotifyOptions = map[string]schema.Attribute{
	"targets": schema.ListNestedAttribute{
		Description: "Notify Targets",

		Optional: true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: schemaNotifyTarget,
		},
	},

	"message": schema.StringAttribute{
		Description: "Notify Message",

		Optional: true,
	},

	"time": schema.Int64Attribute{
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

	"to": schema.StringAttribute{
		Description: "Notify Target",
		Required:    true,
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
	"time": schema.Int64Attribute{
		Description: "Repeat Time Set",

		Optional: true,
	},

	"crontab": schema.SingleNestedAttribute{
		Description: "Repeat Crontab Set",

		Optional:   true,
		Attributes: schemaRepeatCrontabSet,
	},

	"crontab_duration": schema.Int64Attribute{
		Description: "Crontab Duration",

		Optional: true,
	},

	"expire": schema.Int64Attribute{
		Description: "Repeat Expire",

		Optional: true,
	},
}

// schemaTag maps the resource schema data.
var schemaTag = map[string]schema.Attribute{
	"key": schema.StringAttribute{
		Description: "<no value>",
		Required:    true,
	},

	"value": schema.StringAttribute{
		Description: "<no value>",
		Required:    true,
	},
}
