package dashboard

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var resourceSchema = schema.Schema{
	Description:         "Dashboard",
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
			Description: "Dashboard Name",
			Required:    true,
		},

		"extend": schema.StringAttribute{
			Description: "Dashboard Extend",

			Optional: true,
		},

		"mapping": schema.ListNestedAttribute{
			Description: "Dashboard Mapping",

			Optional: true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: schemaMapping,
			},
		},

		"tags": schema.ListAttribute{
			Description: "Dashboard Tag Names",

			Optional:    true,
			ElementType: types.StringType,
		},

		"template": schema.SingleNestedAttribute{
			Description: "Dashboard Template Info",

			Optional:   true,
			Attributes: schemaTemplate,
		},
	},
}

// schemaChart maps the resource schema data.
var schemaChart = map[string]schema.Attribute{
	"name": schema.StringAttribute{
		Description: "Chart Name",
		Required:    true,
	},

	"type": schema.StringAttribute{
		Description: "Chart Type",
		Required:    true,
	},

	"group": schema.StringAttribute{
		Description: "Chart Group",

		Optional: true,
	},

	"pos": schema.SingleNestedAttribute{
		Description: "Chart Position Info",

		Optional:   true,
		Attributes: schemaChartPos,
	},

	"queries": schema.ListNestedAttribute{
		Description: "Chart Query Info",

		Optional: true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: schemaChartQuery,
		},
	},
}

// schemaChartPos maps the resource schema data.
var schemaChartPos = map[string]schema.Attribute{
	"i": schema.StringAttribute{
		Description: "TODO: What is i?",
		Required:    true,
	},

	"x": schema.Int64Attribute{
		Description: "Chart X",
		Required:    true,
	},

	"y": schema.Int64Attribute{
		Description: "Chart Y",
		Required:    true,
	},

	"w": schema.Int64Attribute{
		Description: "Chart Width",
		Required:    true,
	},

	"h": schema.Int64Attribute{
		Description: "Chart Height",
		Required:    true,
	},
}

// schemaChartQuery maps the resource schema data.
var schemaChartQuery = map[string]schema.Attribute{
	"checked": schema.BoolAttribute{
		Description: "Checked",
		Required:    true,
	},

	"datasource": schema.StringAttribute{
		Description: "Datasource",
		Required:    true,
	},

	"qtype": schema.StringAttribute{
		Description: "Query Type",
		Required:    true,
	},

	"query": schema.SingleNestedAttribute{
		Description: "Query",
		Required:    true,
		Attributes:  schemaQuery,
	},

	"unit": schema.StringAttribute{
		Description: "Unit",

		Optional: true,
	},
}

// schemaIconSet maps the resource schema data.
var schemaIconSet = map[string]schema.Attribute{
	"sm": schema.StringAttribute{
		Description: "Small Icon",

		Optional: true,
	},

	"md": schema.StringAttribute{
		Description: "Middle Icon",

		Optional: true,
	},
}

// schemaMain maps the resource schema data.
var schemaMain = map[string]schema.Attribute{
	"type": schema.StringAttribute{
		Description: "Dashboard Type",
		Required:    true,
	},

	"vars": schema.ListNestedAttribute{
		Description: "Dashboard Vars",

		Optional: true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: schemaVar,
		},
	},

	"groups": schema.ListAttribute{
		Description: "Dashboard Groups",

		Optional:    true,
		ElementType: types.StringType,
	},

	"charts": schema.ListNestedAttribute{
		Description: "Dashboard Charts",

		Optional: true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: schemaChart,
		},
	},
}

// schemaMapping maps the resource schema data.
var schemaMapping = map[string]schema.Attribute{
	"class": schema.StringAttribute{
		Description: "Class Name",
		Required:    true,
	},

	"field": schema.StringAttribute{
		Description: "Field Name",
		Required:    true,
	},

	"mapping": schema.StringAttribute{
		Description: "Mapping Field Name",
		Required:    true,
	},

	"datasource": schema.StringAttribute{
		Description: "Data Source",
		Required:    true,
	},
}

// schemaQuery maps the resource schema data.
var schemaQuery = map[string]schema.Attribute{
	"density": schema.StringAttribute{
		Description: "Density",
		Required:    true,
	},

	"filter": schema.ListNestedAttribute{
		Description: "Filter",
		Required:    true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: schemaQueryFilter,
		},
	},

	"group_by": schema.StringAttribute{
		Description: "Group By",
		Required:    true,
	},

	"group_by_time": schema.StringAttribute{
		Description: "Group By Time",
		Required:    true,
	},

	"q": schema.StringAttribute{
		Description: "Query",
		Required:    true,
	},
}

// schemaQueryFilter maps the resource schema data.
var schemaQueryFilter = map[string]schema.Attribute{
	"logic": schema.StringAttribute{
		Description: "Logic",
		Required:    true,
	},

	"name": schema.StringAttribute{
		Description: "Field Name",
		Required:    true,
	},

	"op": schema.StringAttribute{
		Description: "Operator",
		Required:    true,
	},

	"value": schema.StringAttribute{
		Description: "Value",
		Required:    true,
	},
}

// schemaTemplate maps the resource schema data.
var schemaTemplate = map[string]schema.Attribute{
	"title": schema.StringAttribute{
		Description: "Dashboard Title",
		Required:    true,
	},

	"summary": schema.StringAttribute{
		Description: "Dashboard Summary",

		Optional: true,
	},

	"dashboard": schema.SingleNestedAttribute{
		Description: "Dashboard Info",

		Optional:   true,
		Attributes: schemaTemplateDashboard,
	},

	"icon_set": schema.SingleNestedAttribute{
		Description: "Dashboard Icon Set",

		Optional:   true,
		Attributes: schemaIconSet,
	},

	"main": schema.SingleNestedAttribute{
		Description: "Dashboard Main",

		Optional:   true,
		Attributes: schemaMain,
	},
}

// schemaTemplateDashboard maps the resource schema data.
var schemaTemplateDashboard = map[string]schema.Attribute{
	"extend": schema.StringAttribute{
		Description: "Dashboard Extend",

		Optional: true,
	},

	"mapping": schema.ListNestedAttribute{
		Description: "Dashboard Mapping",

		Optional: true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: schemaMapping,
		},
	},
}

// schemaVar maps the resource schema data.
var schemaVar = map[string]schema.Attribute{
	"name": schema.StringAttribute{
		Description: "Var Name",
		Required:    true,
	},

	"type": schema.StringAttribute{
		Description: "Var Type",
		Required:    true,
	},

	"value": schema.StringAttribute{
		Description: "Var Value",
		Required:    true,
	},

	"label": schema.StringAttribute{
		Description: "Var Label",

		Optional: true,
	},
}
