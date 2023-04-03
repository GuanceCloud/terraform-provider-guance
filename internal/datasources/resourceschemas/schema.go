package resourceschemas

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

var dataSourceSchema = schema.Schema{
	Description:         "Resource Schema",
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

					"name": schema.StringAttribute{
						Description: "Resource name",
						Required:    true,
					},

					"title": schema.SingleNestedAttribute{
						Description: "Resource title",

						Optional:   true,
						Attributes: schemaI18n,
					},

					"description": schema.SingleNestedAttribute{
						Description: "Resource description",

						Optional:   true,
						Attributes: schemaI18n,
					},

					"models": schema.ListNestedAttribute{
						Description: "Resource dependends on model",

						Optional: true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: schemaModel,
						},
					},
				},
			},
		},
	},
}

// schemaElemSchema maps the resource schema data.
var schemaElemSchema = map[string]schema.Attribute{
	"type": schema.StringAttribute{
		Description: "Element type",
		Required:    true,
	},

	"format": schema.StringAttribute{
		Description: "Element format",

		Optional: true,
	},

	"ref": schema.StringAttribute{
		Description: "Element reference model",

		Optional: true,
	},

	"enum": schema.ListNestedAttribute{
		Description: "Element enum",

		Optional: true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: schemaEnum,
		},
	},
}

// schemaEnum maps the resource schema data.
var schemaEnum = map[string]schema.Attribute{
	"name": schema.StringAttribute{
		Description: "Enum name",
		Required:    true,
	},

	"value": schema.StringAttribute{
		Description: "Enum value",
		Required:    true,
	},

	"title": schema.SingleNestedAttribute{
		Description: "Enum title",

		Optional:   true,
		Attributes: schemaI18n,
	},
}

// schemaI18n maps the resource schema data.
var schemaI18n = map[string]schema.Attribute{
	"zh": schema.StringAttribute{
		Description: "Chinese",

		Optional: true,
	},

	"en": schema.StringAttribute{
		Description: "English",

		Optional: true,
	},
}

// schemaModel maps the resource schema data.
var schemaModel = map[string]schema.Attribute{
	"name": schema.StringAttribute{
		Description: "Model name",
		Required:    true,
	},

	"title": schema.SingleNestedAttribute{
		Description: "Model title",

		Optional:   true,
		Attributes: schemaI18n,
	},

	"description": schema.SingleNestedAttribute{
		Description: "Model description",

		Optional:   true,
		Attributes: schemaI18n,
	},

	"properties": schema.ListNestedAttribute{
		Description: "Model properties",

		Optional: true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: schemaProperty,
		},
	},
}

// schemaPropMeta maps the resource schema data.
var schemaPropMeta = map[string]schema.Attribute{
	"dynamic": schema.BoolAttribute{
		Description: "Property is dynamic",

		Optional: true,
	},

	"immutable": schema.BoolAttribute{
		Description: "Property is immutable",

		Optional: true,
	},
}

// schemaPropSchema maps the resource schema data.
var schemaPropSchema = map[string]schema.Attribute{
	"type": schema.StringAttribute{
		Description: "Property type",
		Required:    true,
	},

	"format": schema.StringAttribute{
		Description: "Property format",

		Optional: true,
	},

	"required": schema.BoolAttribute{
		Description: "Property is required",

		Optional: true,
	},

	"elem": schema.SingleNestedAttribute{
		Description: "Property element schema",

		Optional:   true,
		Attributes: schemaElemSchema,
	},

	"enum": schema.ListNestedAttribute{
		Description: "Property enum",

		Optional: true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: schemaEnum,
		},
	},

	"model": schema.StringAttribute{
		Description: "Property reference model",

		Optional: true,
	},

	"ref": schema.StringAttribute{
		Description: "Property reference resource",

		Optional: true,
	},
}

// schemaProperty maps the resource schema data.
var schemaProperty = map[string]schema.Attribute{
	"name": schema.StringAttribute{
		Description: "Property name",
		Required:    true,
	},

	"title": schema.SingleNestedAttribute{
		Description: "Property title",

		Optional:   true,
		Attributes: schemaI18n,
	},

	"description": schema.SingleNestedAttribute{
		Description: "Property description",

		Optional:   true,
		Attributes: schemaI18n,
	},

	"schema": schema.SingleNestedAttribute{
		Description: "Property schema",

		Optional:   true,
		Attributes: schemaPropSchema,
	},

	"meta": schema.SingleNestedAttribute{
		Description: "Property meta",

		Optional:   true,
		Attributes: schemaPropMeta,
	},
}
