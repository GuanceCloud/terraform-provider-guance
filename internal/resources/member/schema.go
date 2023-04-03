package member

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
)

var resourceSchema = schema.Schema{
	Description:         "Workspace Member",
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

		"workspace_id": schema.StringAttribute{
			Description: "Workspace ID",
			Required:    true,
		},

		"username": schema.StringAttribute{
			Description: "Username",

			Optional: true,
		},

		"name": schema.StringAttribute{
			Description: "Name",

			Optional: true,
		},

		"email": schema.StringAttribute{
			Description: "Email",

			Optional: true,
		},

		"mobile": schema.StringAttribute{
			Description: "Mobile",

			Optional: true,
		},

		"role": schema.StringAttribute{
			Description: "Role",

			Optional: true,
		},
	},
}
