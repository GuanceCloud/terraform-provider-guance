package default_region

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

var dataSourceSchema = schema.Schema{
	Description: "Default Region",
	Attributes: map[string]schema.Attribute{
		"regions": schema.ListNestedAttribute{
			Description: "The list of the default regions.",
			Computed:    true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: map[string]schema.Attribute{
					"uuid": schema.StringAttribute{
						Description: "The uuid of the region.",
						Computed:    true,
					},
					"province": schema.StringAttribute{
						Description: "The province of the region.",
						Computed:    true,
					},
					"city": schema.StringAttribute{
						Description: "The city of the region.",
						Computed:    true,
					},
					"country": schema.StringAttribute{
						Description: "The country of the region.",
						Computed:    true,
					},
					"name": schema.StringAttribute{
						Description: "The name of the region.",
						Computed:    true,
					},
					"name_en": schema.StringAttribute{
						Description: "The English name of the region.",
						Computed:    true,
					},
					"extend_info": schema.StringAttribute{
						Description: "The extend info of the region.",
						Computed:    true,
					},
					"internal": schema.BoolAttribute{
						Description: "Whether the region is internal.",
						Computed:    true,
					},
					"keycode": schema.StringAttribute{
						Description: "The keycode of the region.",
						Computed:    true,
					},
					"isp": schema.StringAttribute{
						Description: "The ISP of the region.",
						Computed:    true,
					},
					"status": schema.StringAttribute{
						Description: "The status of the region.",
						Computed:    true,
					},
					"region": schema.StringAttribute{
						Description: "The region code.",
						Computed:    true,
					},
					"owner": schema.StringAttribute{
						Description: "The owner of the region.",
						Computed:    true,
					},
					"heartbeat": schema.StringAttribute{
						Description: "The heartbeat timestamp of the region.",
						Computed:    true,
					},
					"company": schema.StringAttribute{
						Description: "The company of the region.",
						Computed:    true,
					},
					"external_id": schema.StringAttribute{
						Description: "The external ID of the region.",
						Computed:    true,
					},
					"parent_ak": schema.StringAttribute{
						Description: "The parent AK of the region.",
						Computed:    true,
					},
					"create_at": schema.StringAttribute{
						Description: "The create timestamp of the region.",
						Computed:    true,
					},
				},
			},
		},
	},
}
