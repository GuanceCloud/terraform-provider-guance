package blacklist

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
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
		"name": schema.StringAttribute{
			Description: "The name of resource.",
			Required:    true,
			Validators: []validator.String{
				stringvalidator.LengthAtMost(50),
			},
		},
		"desc": schema.StringAttribute{
			Description: "The description of resource.",
			Optional:    true,
			Validators: []validator.String{
				stringvalidator.LengthAtMost(256),
			},
		},

		"create_at": schema.StringAttribute{
			Description: "The timestamp seconds of the resource created at.",
			Computed:    true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"update_at": schema.StringAttribute{
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

		"source": schema.StringAttribute{
			Description: "The source of the resource.",
			Optional:    true,
		},
		"sources": schema.ListAttribute{
			Description: "The source list of the resource.",
			Optional:    true,
			ElementType: types.StringType,
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
	},

	"operation": schema.StringAttribute{
		Description: "The operation of the filter.",
		MarkdownDescription: `
		Operation, value must be one of: *in*, *not_in*, *match*, *not_match*, other value will be ignored.
		`,
		Required: true,
		Validators: []validator.String{
			stringvalidator.OneOf("in", "not_in", "match", "not_match"),
		},
	},

	"condition": schema.StringAttribute{
		Description: "The condition of the filter.",
		Required:    true,
	},

	"values": schema.ListAttribute{
		Description: "The values of the filter.",
		Optional:    true,
		ElementType: types.StringType,
	},
}
