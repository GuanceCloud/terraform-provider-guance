package permissions

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

var dataSourceSchema = schema.Schema{
	Description:         "Role permissions",
	MarkdownDescription: resourceDocument,
	Attributes: map[string]schema.Attribute{
		"is_support_custom_role": schema.BoolAttribute{
			Description: "Filter the permission list that supports custom role.",
			Optional:    true,
		},
		"permissions": schema.ListNestedAttribute{
			Description: "The list of the permissions.",
			Computed:    true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: permissionAttribute,
			},
		},
	},
}

var permissionAttribute = map[string]schema.Attribute{
	"desc": schema.StringAttribute{
		Description: "The description of the permission.",
		Computed:    true,
	},
	"disabled": schema.Int64Attribute{
		Description: "The disabled status of the permission.",
		Computed:    true,
	},

	"is_support_custom_role": schema.Int64Attribute{
		Description: "Whether support custom role.",
		Computed:    true,
	},
	"is_support_general": schema.Int64Attribute{
		Description: "Whether support general.",
		Computed:    true,
	},
	"is_support_owner": schema.Int64Attribute{
		Description: "Whether support owner.",
		Computed:    true,
	},
	"is_support_read_only": schema.Int64Attribute{
		Description: "Whether support readonly.",
		Computed:    true,
	},
	"is_support_ws_admin": schema.Int64Attribute{
		Description: "Whether support WsAdmin.",
		Computed:    true,
	},
	"key": schema.StringAttribute{
		Description: "The key of the permission.",
		Computed:    true,
	},
	"name": schema.StringAttribute{
		Description: "The name of the permission.",
		Computed:    true,
	},
	"subs": schema.ListNestedAttribute{
		Description: "The list of the sub permissions.",
		Computed:    true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: map[string]schema.Attribute{
				"desc": schema.StringAttribute{
					Description: "The description of the permission.",
					Computed:    true,
				},
				"disabled": schema.Int64Attribute{
					Description: "The disabled status of the permission.",
					Computed:    true,
				},
				"is_support_custom_role": schema.Int64Attribute{
					Description: "Whether support custom role.",
					Computed:    true,
				},
				"is_support_general": schema.Int64Attribute{
					Description: "Whether support general.",
					Computed:    true,
				},
				"is_support_owner": schema.Int64Attribute{
					Description: "Whether support owner.",
					Computed:    true,
				},
				"is_support_read_only": schema.Int64Attribute{
					Description: "Whether support readonly.",
					Computed:    true,
				},
				"is_support_ws_admin": schema.Int64Attribute{
					Description: "Whether support WsAdmin.",
					Computed:    true,
				},
				"key": schema.StringAttribute{
					Description: "The key of the permission.",
					Computed:    true,
				},
				"name": schema.StringAttribute{
					Description: "The name of the permission.",
					Computed:    true,
				},
			},
		},
	},
}
