// Code generated by Guance Cloud Code Generation Pipeline. DO NOT EDIT.

package membergroup

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var resourceSchema = schema.Schema{
	Description:         "Member Group",
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
			Description: "Name",

			Required: true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.RequiresReplace(),
			},
		},

		"member_ids": schema.ListAttribute{
			Description: "Member id list",

			Optional:    true,
			ElementType: types.StringType,
			PlanModifiers: []planmodifier.List{
				listplanmodifier.RequiresReplace(),
			},
		},
	},
}
