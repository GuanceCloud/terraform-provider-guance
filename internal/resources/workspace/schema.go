package workspace

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
)

var resourceSchema = schema.Schema{
	Description:         "Workspace",
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
			Description: "Workspace name",
			Required:    true,
		},

		"description": schema.StringAttribute{
			Description: "Workspace description information",

			Optional: true,
		},
	},
}
