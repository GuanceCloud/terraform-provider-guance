// Code generated by Guance Cloud Code Generation Pipeline. DO NOT EDIT.

package monitor

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
)

var resourceSchema = schema.Schema{
	Description:         "Monitor",
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

		"manifest": schema.StringAttribute{
			Description: "Monitor Configuration",

			Required: true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.RequiresReplace(),
			},
		},

		"alert_policy": schema.SingleNestedAttribute{
			Description: "Alert Policy Configuration",

			Required:   true,
			Attributes: schemaAlertPolicy,
		},

		"dashboard": schema.SingleNestedAttribute{
			Description: "Dashboard Configuration",

			Optional:   true,
			Attributes: schemaDashboard,
		},
	},
}

// schemaAlertPolicy maps the resource schema data.
var schemaAlertPolicy = map[string]schema.Attribute{
	"id": schema.StringAttribute{
		Description: "Alert Policy ID",

		Optional: true,
		PlanModifiers: []planmodifier.String{
			stringplanmodifier.RequiresReplace(),
		},
	},
}

// schemaDashboard maps the resource schema data.
var schemaDashboard = map[string]schema.Attribute{
	"id": schema.StringAttribute{
		Description: "Dashboard ID",

		Optional: true,
		PlanModifiers: []planmodifier.String{
			stringplanmodifier.RequiresReplace(),
		},
	},
}
