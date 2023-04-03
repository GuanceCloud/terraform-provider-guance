package functions

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

var dataSourceSchema = schema.Schema{
	Description:         "Function",
	MarkdownDescription: resourceDocument,
	Attributes: map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Description: "Identifier of the resource.",
			Computed:    true,
		},

		"max_results": schema.Int64Attribute{
			Description: "The max results count of the resource will be returned.",
			Optional:    true,
		},

		"type_name": schema.StringAttribute{
			Description: "The type name of the resource be queried",
			Optional:    true,
		},

		"items": schema.ListNestedAttribute{
			Description: "The list of the resource",
			Computed:    true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: map[string]schema.Attribute{

					"title": schema.StringAttribute{
						Description: "Title",

						Optional: true,
					},

					"description": schema.StringAttribute{
						Description: "Description",

						Optional: true,
					},

					"func_id": schema.StringAttribute{
						Description: "Function ID",

						Optional: true, Computed: true,
					},
				},
			},
		},
	},
}
