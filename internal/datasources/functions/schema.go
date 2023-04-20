// Code generated by Guance Cloud Code Generation Pipeline. DO NOT EDIT.

package functions

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
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

		"filter": schema.ListNestedAttribute{
			Description: "The list of the resource",
			Optional:    true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: map[string]schema.Attribute{
					"name": schema.StringAttribute{
						Description: "The filter path, represent as json path.",
						Required:    true,
					},

					"values": schema.ListAttribute{
						ElementType: types.StringType,
						Required:    true,
						Description: "The filter values",
					},
				},
			},
		},

		"items": schema.ListNestedAttribute{
			Description: "The list of the resource",
			Computed:    true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: map[string]schema.Attribute{
					"id": schema.StringAttribute{
						Description: "Numeric identifier of the order.",
						Computed:    true,
					},

					"created_at": schema.StringAttribute{
						Description: "Timestamp of the last Terraform update of the order.",
						Computed:    true,
					},

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
