package intelligentinspection

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var resourceSchema = schema.Schema{
	Description:         "Intelligent Inspection",
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

		"monitor_id": schema.StringAttribute{
			Description: "Monitor ID",

			Optional: true,
		},

		"ref_key": schema.StringAttribute{
			Description: "Ref Key",

			Optional: true,
		},

		"ref_func": schema.SingleNestedAttribute{
			Description: "Ref Func Info",

			Optional:   true,
			Attributes: schemaRefFunc,
		},
	},
}

// schemaRefFunc maps the resource schema data.
var schemaRefFunc = map[string]schema.Attribute{
	"func_id": schema.StringAttribute{
		Description: "Func ID",

		Optional: true,
	},

	"title": schema.StringAttribute{
		Description: "Title",

		Optional: true,
	},

	"description": schema.StringAttribute{
		Description: "Description",

		Optional: true,
	},

	"definition": schema.StringAttribute{
		Description: "Definition",

		Optional: true,
	},

	"category": schema.StringAttribute{
		Description: "Category",

		Optional: true,
	},

	"args": schema.ListAttribute{
		Description: "Args",

		Optional:    true,
		ElementType: types.StringType,
	},

	"kwargs": schema.StringAttribute{
		Description: "Kwargs",

		Optional: true,
	},

	"disabled": schema.BoolAttribute{
		Description: "Is Disabled",

		Optional: true,
	},
}
