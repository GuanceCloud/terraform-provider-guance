package members

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

var dataSourceSchema = schema.Schema{
	Description:         "Workspace Member",
	MarkdownDescription: resourceDocument,
	Attributes: map[string]schema.Attribute{
		"search": schema.StringAttribute{
			Description: "Search the member by email or name.",
			Optional:    true,
		},
		"members": schema.ListNestedAttribute{
			Description: "The list of the members.",
			Computed:    true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: map[string]schema.Attribute{
					"uuid": schema.StringAttribute{
						Description: "The uuid of the member.",
						Computed:    true,
					},

					"create_at": schema.StringAttribute{
						Description: "The unix timestamp in seconds of the member creation.",
						Computed:    true,
					},

					"email": schema.StringAttribute{
						Description: "Email",

						Optional: true,
					},
					"name": schema.StringAttribute{
						Description: "User name",

						Optional: true,
					},
					"roles": roleAttribute,
				},
			},
		},
	},
}

var roleAttribute = schema.ListNestedAttribute{
	Description: "Role",

	MarkdownDescription: `
	Role, value must be one of: *owner*, *wsAdmin*, *general*, *readOnly*, other value will be ignored.
	`,
	Optional: true,
	NestedObject: schema.NestedAttributeObject{
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Optional:    true,
				Description: "The name of the role.",
			},
			"uuid": schema.StringAttribute{
				Optional:    true,
				Description: "The UUID of the role.",
			},
		},
	},
}
