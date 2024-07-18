package role

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var resourceSchema = schema.Schema{
	Description:         "Role is a set of permissions that can be assigned to a user or a group of users.",
	MarkdownDescription: resourceDocument,
	Attributes: map[string]schema.Attribute{
		"uuid": schema.StringAttribute{
			Description: "The UUID of the role.",
			Computed:    true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"name": schema.StringAttribute{
			Description: "The name of the role.",
			Required:    true,
		},

		"desc": schema.StringAttribute{
			Description: "The description of the role.",
			Computed:    true,
			Optional:    true,
		},

		"keys": schema.ListAttribute{
			Description: "The permission keys.",
			Required:    true,
			ElementType: types.StringType,
		},
	},
}
