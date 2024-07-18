package blacklist

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
	Description:         "BlackList",
	MarkdownDescription: resourceDocument,
	Attributes: map[string]schema.Attribute{
		"uuid": schema.StringAttribute{
			Description: "The uuid of resource.",
			Computed:    true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},

		"created_at": schema.StringAttribute{
			Description: "The timestamp seconds of the resource created at.",
			Computed:    true,
		},

		"source": schema.StringAttribute{
			Description: "The source of the resource.",
			Required:    true,
		},
		"type": schema.StringAttribute{
			Description: "The type of the resource.",
			Required:    true,
			Validators: []validator.String{
				stringvalidator.OneOf("object", "custom_object", "logging", "keyevent", "tracing", "rum", "network", "security", "profiling", "metric"),
			},
		},
		"filters": schema.ListNestedAttribute{
			Description: "The filters of the resource.",
			Optional:    true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: schemaFilter,
			},
		},
	},
}

// schemaFilter maps the resource schema data.
var schemaFilter = map[string]schema.Attribute{
	"name": schema.StringAttribute{
		Description: "The name of the filter.",

		Required: true,
		PlanModifiers: []planmodifier.String{
			stringplanmodifier.RequiresReplace(),
		},
	},

	"operation": schema.StringAttribute{
		Description: "The operation of the filter.",
		MarkdownDescription: `
		Operation, value must be one of: *in*, *not in*, *match*, *not match*, other value will be ignored.
		`,
		Required: true,
		PlanModifiers: []planmodifier.String{
			stringplanmodifier.RequiresReplace(),
		},
		Validators: []validator.String{
			stringvalidator.OneOf("in", "not in", "match", "not match"),
		},
	},

	"condition": schema.StringAttribute{
		Description: "The condition of the filter.",
		Required:    true,
		PlanModifiers: []planmodifier.String{
			stringplanmodifier.RequiresReplace(),
		},
	},

	"values": schema.ListAttribute{
		Description: "The values of the filter.",
		Optional:    true,
		ElementType: types.StringType,
		PlanModifiers: []planmodifier.List{
			listplanmodifier.RequiresReplace(),
		},
	},
}
