package alertpolicy

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var resourceSchema = schema.Schema{
	Description:         "Alert Policy",
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

		"type": schema.StringAttribute{
			Description: "Trigger rule type, default is custom",

			Optional: true,
		},

		"name": schema.StringAttribute{
			Description: "Alert Policy Name",
			Required:    true,
		},

		"silent_timeout": schema.Int64Attribute{
			Description: "Silent timeout timestamp",

			Optional: true,
		},

		"alert_target": schema.ListNestedAttribute{
			Description: "Alert Action",

			Optional: true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: schemaAlertTarget,
			},
		},
	},
}

// schemaAlertTarget maps the resource schema data.
var schemaAlertTarget = map[string]schema.Attribute{
	"type": schema.StringAttribute{
		Description: "Alert type",
		Required:    true,
	},

	"status": schema.ListAttribute{
		Description: "The status value of the event to be sent",
		Required:    true,
		ElementType: types.StringType,
	},

	"min_interval": schema.Int64Attribute{
		Description: "The minimum alert interval, in seconds. 0 / null means always send an alert",

		Optional: true,
	},

	"allow_week_days": schema.ListAttribute{
		Description: "Allowed to send alerts on weekdays",

		Optional:    true,
		ElementType: types.Int64Type,
	},

	"notification": schema.SingleNestedAttribute{
		Description: "Notification",

		Optional:   true,
		Attributes: schemaNotificationTarget,
	},
}

// schemaNotificationTarget maps the resource schema data.
var schemaNotificationTarget = map[string]schema.Attribute{
	"to": schema.ListAttribute{
		Description: "Notification",
		Required:    true,
		ElementType: types.StringType,
	},
}
