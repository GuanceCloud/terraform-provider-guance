package monitor

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var resourceSchema = schema.Schema{
	Description:         "Monitor",
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

		"dashboard_id": schema.StringAttribute{
			Description: "<no value>",

			Optional: true,
		},

		"script": schema.SingleNestedAttribute{
			Description: "<no value>",

			Optional:   true,
			Attributes: schemaMonitorScript,
		},
	},
}

// schemaAPMCheck maps the resource schema data.
var schemaAPMCheck = map[string]schema.Attribute{
	"name": schema.StringAttribute{
		Description: "check item name",

		Optional: true,
	},

	"title": schema.StringAttribute{
		Description: "event title",

		Optional: true,
	},

	"message": schema.StringAttribute{
		Description: "event message",

		Optional: true,
	},

	"every": schema.StringAttribute{
		Description: "check frequency",

		Optional: true,
	},

	"interval": schema.Int64Attribute{
		Description: "Query interval",

		Optional: true,
	},

	"recover_need_period_count": schema.Int64Attribute{
		Description: "recover need period count",

		Optional: true,
	},

	"no_data_interval": schema.Int64Attribute{
		Description: "no data interval",

		Optional: true,
	},

	"targets": schema.ListNestedAttribute{
		Description: "targets for checking",

		Optional: true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: schemaTarget,
		},
	},

	"checker": schema.SingleNestedAttribute{
		Description: "condition configuration for checking",

		Optional:   true,
		Attributes: schemaChecker,
	},
}

// schemaCheckFunc maps the resource schema data.
var schemaCheckFunc = map[string]schema.Attribute{
	"func_id": schema.StringAttribute{
		Description: "function ID",
		Required:    true,
	},

	"kwargs": schema.StringAttribute{
		Description: "parameters",

		Optional: true,
	},
}

// schemaChecker maps the resource schema data.
var schemaChecker = map[string]schema.Attribute{
	"rules": schema.ListNestedAttribute{
		Description: "rules",

		Optional: true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: schemaRule,
		},
	},
}

// schemaCloudDialCheck maps the resource schema data.
var schemaCloudDialCheck = map[string]schema.Attribute{
	"name": schema.StringAttribute{
		Description: "check item name",

		Optional: true,
	},

	"title": schema.StringAttribute{
		Description: "event title",

		Optional: true,
	},

	"message": schema.StringAttribute{
		Description: "event message",

		Optional: true,
	},

	"every": schema.StringAttribute{
		Description: "check frequency",

		Optional: true,
	},

	"interval": schema.Int64Attribute{
		Description: "Query interval",

		Optional: true,
	},

	"recover_need_period_count": schema.Int64Attribute{
		Description: "recover need period count",

		Optional: true,
	},

	"no_data_interval": schema.Int64Attribute{
		Description: "no data interval",

		Optional: true,
	},

	"targets": schema.ListNestedAttribute{
		Description: "targets for checking",

		Optional: true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: schemaTarget,
		},
	},

	"checker": schema.SingleNestedAttribute{
		Description: "condition configuration for checking",

		Optional:   true,
		Attributes: schemaChecker,
	},
}

// schemaCondition maps the resource schema data.
var schemaCondition = map[string]schema.Attribute{
	"alias": schema.StringAttribute{
		Description: "alias of target",

		Optional: true,
	},

	"operator": schema.StringAttribute{
		Description: "operator",

		Optional: true,
	},

	"operands": schema.ListAttribute{
		Description: "operands",

		Optional:    true,
		ElementType: types.StringType,
	},
}

// schemaLoggingCheck maps the resource schema data.
var schemaLoggingCheck = map[string]schema.Attribute{
	"name": schema.StringAttribute{
		Description: "check item name",

		Optional: true,
	},

	"title": schema.StringAttribute{
		Description: "event title",

		Optional: true,
	},

	"message": schema.StringAttribute{
		Description: "event message",

		Optional: true,
	},

	"every": schema.StringAttribute{
		Description: "check frequency",

		Optional: true,
	},

	"interval": schema.Int64Attribute{
		Description: "Query interval",

		Optional: true,
	},

	"recover_need_period_count": schema.Int64Attribute{
		Description: "recover need period count",

		Optional: true,
	},

	"no_data_interval": schema.Int64Attribute{
		Description: "no data interval",

		Optional: true,
	},

	"targets": schema.ListNestedAttribute{
		Description: "targets for checking",

		Optional: true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: schemaTarget,
		},
	},

	"checker": schema.SingleNestedAttribute{
		Description: "condition configuration for checking",

		Optional:   true,
		Attributes: schemaChecker,
	},
}

// schemaMonitorScript maps the resource schema data.
var schemaMonitorScript = map[string]schema.Attribute{
	"type": schema.StringAttribute{
		Description: "check method type",

		Optional: true,
	},

	"simple_check": schema.SingleNestedAttribute{
		Description: "simple check",

		Optional:   true,
		Attributes: schemaSimpleCheck,
	},

	"senior_check": schema.SingleNestedAttribute{
		Description: "senior check",

		Optional:   true,
		Attributes: schemaSeniorCheck,
	},

	"logging_check": schema.SingleNestedAttribute{
		Description: "logging check",

		Optional:   true,
		Attributes: schemaLoggingCheck,
	},

	"mutations_check": schema.SingleNestedAttribute{
		Description: "mutations check",

		Optional:   true,
		Attributes: schemaMutationsCheck,
	},

	"water_level_check": schema.SingleNestedAttribute{
		Description: "water level check",

		Optional:   true,
		Attributes: schemaWaterLevelCheck,
	},

	"range_check": schema.SingleNestedAttribute{
		Description: "range check",

		Optional:   true,
		Attributes: schemaRangeCheck,
	},

	"security_check": schema.SingleNestedAttribute{
		Description: "security check",

		Optional:   true,
		Attributes: schemaSecurityCheck,
	},

	"apm_check": schema.SingleNestedAttribute{
		Description: "APM check",

		Optional:   true,
		Attributes: schemaAPMCheck,
	},

	"rum_check": schema.SingleNestedAttribute{
		Description: "RUM check",

		Optional:   true,
		Attributes: schemaRUMCheck,
	},

	"process_check": schema.SingleNestedAttribute{
		Description: "process check",

		Optional:   true,
		Attributes: schemaProcessCheck,
	},

	"cloud_dial_check": schema.SingleNestedAttribute{
		Description: "cloud dial check",

		Optional:   true,
		Attributes: schemaCloudDialCheck,
	},
}

// schemaMutationsCheck maps the resource schema data.
var schemaMutationsCheck = map[string]schema.Attribute{
	"name": schema.StringAttribute{
		Description: "check item name",

		Optional: true,
	},

	"title": schema.StringAttribute{
		Description: "event title",

		Optional: true,
	},

	"message": schema.StringAttribute{
		Description: "event message",

		Optional: true,
	},

	"every": schema.StringAttribute{
		Description: "check frequency",

		Optional: true,
	},

	"interval": schema.Int64Attribute{
		Description: "Query interval",

		Optional: true,
	},

	"recover_need_period_count": schema.Int64Attribute{
		Description: "recover need period count",

		Optional: true,
	},

	"no_data_interval": schema.Int64Attribute{
		Description: "no data interval",

		Optional: true,
	},

	"targets": schema.ListNestedAttribute{
		Description: "targets for checking",

		Optional: true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: schemaTarget,
		},
	},

	"checker": schema.SingleNestedAttribute{
		Description: "condition configuration for checking",

		Optional:   true,
		Attributes: schemaChecker,
	},
}

// schemaProcessCheck maps the resource schema data.
var schemaProcessCheck = map[string]schema.Attribute{
	"name": schema.StringAttribute{
		Description: "check item name",

		Optional: true,
	},

	"title": schema.StringAttribute{
		Description: "event title",

		Optional: true,
	},

	"message": schema.StringAttribute{
		Description: "event message",

		Optional: true,
	},

	"every": schema.StringAttribute{
		Description: "check frequency",

		Optional: true,
	},

	"interval": schema.Int64Attribute{
		Description: "Query interval",

		Optional: true,
	},

	"recover_need_period_count": schema.Int64Attribute{
		Description: "recover need period count",

		Optional: true,
	},

	"no_data_interval": schema.Int64Attribute{
		Description: "no data interval",

		Optional: true,
	},

	"targets": schema.ListNestedAttribute{
		Description: "targets for checking",

		Optional: true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: schemaTarget,
		},
	},

	"checker": schema.SingleNestedAttribute{
		Description: "condition configuration for checking",

		Optional:   true,
		Attributes: schemaChecker,
	},
}

// schemaRUMCheck maps the resource schema data.
var schemaRUMCheck = map[string]schema.Attribute{
	"name": schema.StringAttribute{
		Description: "check item name",

		Optional: true,
	},

	"title": schema.StringAttribute{
		Description: "event title",

		Optional: true,
	},

	"message": schema.StringAttribute{
		Description: "event message",

		Optional: true,
	},

	"every": schema.StringAttribute{
		Description: "check frequency",

		Optional: true,
	},

	"interval": schema.Int64Attribute{
		Description: "Query interval",

		Optional: true,
	},

	"recover_need_period_count": schema.Int64Attribute{
		Description: "recover need period count",

		Optional: true,
	},

	"no_data_interval": schema.Int64Attribute{
		Description: "no data interval",

		Optional: true,
	},

	"targets": schema.ListNestedAttribute{
		Description: "targets for checking",

		Optional: true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: schemaTarget,
		},
	},

	"checker": schema.SingleNestedAttribute{
		Description: "condition configuration for checking",

		Optional:   true,
		Attributes: schemaChecker,
	},
}

// schemaRangeCheck maps the resource schema data.
var schemaRangeCheck = map[string]schema.Attribute{
	"name": schema.StringAttribute{
		Description: "check item name",

		Optional: true,
	},

	"title": schema.StringAttribute{
		Description: "event title",

		Optional: true,
	},

	"message": schema.StringAttribute{
		Description: "event message",

		Optional: true,
	},

	"every": schema.StringAttribute{
		Description: "check frequency",

		Optional: true,
	},

	"interval": schema.Int64Attribute{
		Description: "Query interval",

		Optional: true,
	},

	"recover_need_period_count": schema.Int64Attribute{
		Description: "recover need period count",

		Optional: true,
	},

	"no_data_interval": schema.Int64Attribute{
		Description: "no data interval",

		Optional: true,
	},

	"targets": schema.ListNestedAttribute{
		Description: "targets for checking",

		Optional: true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: schemaTarget,
		},
	},

	"checker": schema.SingleNestedAttribute{
		Description: "condition configuration for checking",

		Optional:   true,
		Attributes: schemaChecker,
	},
}

// schemaRule maps the resource schema data.
var schemaRule = map[string]schema.Attribute{
	"conditions": schema.ListNestedAttribute{
		Description: "conditions",

		Optional: true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: schemaCondition,
		},
	},

	"condition_logic": schema.StringAttribute{
		Description: "condition logic",

		Optional: true,
	},

	"status": schema.StringAttribute{
		Description: "Fit the condition, the output event's status. The value is the same as the event's status",

		Optional: true,
	},

	"direction": schema.StringAttribute{
		Description: "[interval/water level/mutation parameter] check direction",

		Optional: true,
	},

	"period_num": schema.Int64Attribute{
		Description: "[interval/water level/mutation parameter] only check the latest data point number",

		Optional: true,
	},

	"check_percent": schema.Int64Attribute{
		Description: "[interval parameter] abnormal percentage threshold",

		Optional: true,
	},

	"check_count": schema.Int64Attribute{
		Description: "[water level/mutation parameter] continuous abnormal point number",

		Optional: true,
	},

	"strength": schema.Int64Attribute{
		Description: "[water level/mutation parameter] strength",

		Optional: true,
	},
}

// schemaSecurityCheck maps the resource schema data.
var schemaSecurityCheck = map[string]schema.Attribute{
	"name": schema.StringAttribute{
		Description: "check item name",

		Optional: true,
	},

	"title": schema.StringAttribute{
		Description: "event title",

		Optional: true,
	},

	"message": schema.StringAttribute{
		Description: "event message",

		Optional: true,
	},

	"every": schema.StringAttribute{
		Description: "check frequency",

		Optional: true,
	},

	"interval": schema.Int64Attribute{
		Description: "Query interval",

		Optional: true,
	},

	"recover_need_period_count": schema.Int64Attribute{
		Description: "recover need period count",

		Optional: true,
	},

	"no_data_interval": schema.Int64Attribute{
		Description: "no data interval",

		Optional: true,
	},

	"targets": schema.ListNestedAttribute{
		Description: "targets for checking",

		Optional: true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: schemaTarget,
		},
	},

	"checker": schema.SingleNestedAttribute{
		Description: "condition configuration for checking",

		Optional:   true,
		Attributes: schemaChecker,
	},
}

// schemaSeniorCheck maps the resource schema data.
var schemaSeniorCheck = map[string]schema.Attribute{
	"name": schema.StringAttribute{
		Description: "rule name",
		Required:    true,
	},

	"title": schema.StringAttribute{
		Description: "event title",
		Required:    true,
	},

	"message": schema.StringAttribute{
		Description: "event message",
		Required:    true,
	},

	"type": schema.StringAttribute{
		Description: "rule type",
		Required:    true,
	},

	"every": schema.StringAttribute{
		Description: "check frequency",
		Required:    true,
	},

	"check_funcs": schema.ListNestedAttribute{
		Description: "check functions",
		Required:    true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: schemaCheckFunc,
		},
	},
}

// schemaSimpleCheck maps the resource schema data.
var schemaSimpleCheck = map[string]schema.Attribute{
	"name": schema.StringAttribute{
		Description: "check item name",

		Optional: true,
	},

	"title": schema.StringAttribute{
		Description: "event title",

		Optional: true,
	},

	"message": schema.StringAttribute{
		Description: "event message",

		Optional: true,
	},

	"every": schema.StringAttribute{
		Description: "check frequency",

		Optional: true,
	},

	"interval": schema.Int64Attribute{
		Description: "Query interval",

		Optional: true,
	},

	"recover_need_period_count": schema.Int64Attribute{
		Description: "recover need period count",

		Optional: true,
	},

	"no_data_interval": schema.Int64Attribute{
		Description: "no data interval",

		Optional: true,
	},

	"targets": schema.ListNestedAttribute{
		Description: "targets for checking",

		Optional: true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: schemaTarget,
		},
	},

	"checker": schema.SingleNestedAttribute{
		Description: "condition configuration for checking",

		Optional:   true,
		Attributes: schemaChecker,
	},
}

// schemaTarget maps the resource schema data.
var schemaTarget = map[string]schema.Attribute{
	"alias": schema.StringAttribute{
		Description: "alias",

		Optional: true,
	},

	"dql": schema.StringAttribute{
		Description: "dql",

		Optional: true,
	},
}

// schemaWaterLevelCheck maps the resource schema data.
var schemaWaterLevelCheck = map[string]schema.Attribute{
	"name": schema.StringAttribute{
		Description: "check item name",

		Optional: true,
	},

	"title": schema.StringAttribute{
		Description: "event title",

		Optional: true,
	},

	"message": schema.StringAttribute{
		Description: "event message",

		Optional: true,
	},

	"every": schema.StringAttribute{
		Description: "check frequency",

		Optional: true,
	},

	"interval": schema.Int64Attribute{
		Description: "Query interval",

		Optional: true,
	},

	"recover_need_period_count": schema.Int64Attribute{
		Description: "recover need period count",

		Optional: true,
	},

	"no_data_interval": schema.Int64Attribute{
		Description: "no data interval",

		Optional: true,
	},

	"targets": schema.ListNestedAttribute{
		Description: "targets for checking",

		Optional: true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: schemaTarget,
		},
	},

	"checker": schema.SingleNestedAttribute{
		Description: "condition configuration for checking",

		Optional:   true,
		Attributes: schemaChecker,
	},
}