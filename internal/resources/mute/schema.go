package mute

import (
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var datetimePattern = regexp.MustCompile(`^$|^\d{4}/\d{2}/\d{2} \d{2}:\d{2}:\d{2}$|^0$`)

var resourceSchema = schema.Schema{
	Description:         "Mute rule.",
	MarkdownDescription: resourceDocument,
	Attributes: map[string]schema.Attribute{
		"uuid": schema.StringAttribute{
			Description: "The UUID of the mute rule.",
			Computed:    true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"name": schema.StringAttribute{
			Description: "The name of the mute rule.",
			Required:    true,
			Validators: []validator.String{
				stringvalidator.LengthAtMost(128),
			},
		},
		"description": schema.StringAttribute{
			Description: "The description of the mute rule.",
			Optional:    true,
			Validators: []validator.String{
				stringvalidator.LengthAtMost(256),
			},
		},
		"type": schema.StringAttribute{
			Description: "Mute rule type.",
			Required:    true,
			Validators: []validator.String{
				stringvalidator.OneOf("checker", "alertPolicy", "tag", "custom"),
			},
		},
		"mute_ranges": schema.ListNestedAttribute{
			Description: "Mute ranges. Empty means all ranges for the selected type.",
			Required:    true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: map[string]schema.Attribute{
					"name": schema.StringAttribute{
						Description: "The display name of the muted resource.",
						Optional:    true,
					},
					"type": schema.StringAttribute{
						Description: "The muted resource type returned by the API.",
						Optional:    true,
					},
					"checker_uuid": schema.StringAttribute{
						Description: "Monitor/checker UUID.",
						Optional:    true,
					},
					"monitor_uuid": schema.StringAttribute{
						Description: "Monitor UUID.",
						Optional:    true,
					},
					"slo_uuid": schema.StringAttribute{
						Description: "SLO UUID.",
						Optional:    true,
					},
					"alert_policy_uuid": schema.StringAttribute{
						Description: "Alert policy UUID.",
						Optional:    true,
					},
					"tag_uuid": schema.StringAttribute{
						Description: "Monitor tag UUID.",
						Optional:    true,
					},
				},
			},
		},
		"tags": schema.MapAttribute{
			Description: "Event attribute filters. Prefix a key with '-' for negative matching.",
			Optional:    true,
			ElementType: types.ListType{ElemType: types.StringType},
		},
		"filter_string": schema.StringAttribute{
			Description: "Event attribute filter expression. Takes precedence over tags.",
			Optional:    true,
			Validators: []validator.String{
				stringvalidator.LengthAtMost(2048),
			},
		},
		"notify_targets": schema.ListNestedAttribute{
			Description: "Notification targets for mute notifications.",
			Optional:    true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: map[string]schema.Attribute{
					"to": schema.ListAttribute{
						Description: "Notification target UUIDs.",
						Required:    true,
						ElementType: types.StringType,
					},
					"type": schema.StringAttribute{
						Description: "Notification target type, such as mail or notifyObject.",
						Required:    true,
					},
				},
			},
		},
		"notify_message": schema.StringAttribute{
			Description: "Notification message.",
			Optional:    true,
			Validators: []validator.String{
				stringvalidator.LengthAtMost(3000),
			},
		},
		"notify_time_str": schema.StringAttribute{
			Description: "Notification time in YYYY/MM/DD HH:mm:ss, empty for no scheduled notification.",
			Optional:    true,
			Validators: []validator.String{
				stringvalidator.RegexMatches(datetimePattern, "must use YYYY/MM/DD HH:mm:ss, empty string, or 0"),
			},
		},
		"start_time": schema.StringAttribute{
			Description: "One-time mute start time in YYYY/MM/DD HH:mm:ss.",
			Optional:    true,
			Validators: []validator.String{
				stringvalidator.RegexMatches(datetimePattern, "must use YYYY/MM/DD HH:mm:ss"),
			},
		},
		"end_time": schema.StringAttribute{
			Description: "One-time mute end time in YYYY/MM/DD HH:mm:ss.",
			Optional:    true,
			Validators: []validator.String{
				stringvalidator.RegexMatches(datetimePattern, "must use YYYY/MM/DD HH:mm:ss or empty string"),
			},
		},
		"repeat_time_set": schema.Int64Attribute{
			Description: "Whether to repeat the mute rule. 0 for one-time, 1 for repeated.",
			Optional:    true,
			Computed:    true,
			Default:     int64default.StaticInt64(0),
			Validators: []validator.Int64{
				int64validator.OneOf(0, 1),
			},
		},
		"repeat_crontab_set": schema.SingleNestedAttribute{
			Description: "Repeated mute crontab settings.",
			Optional:    true,
			Attributes: map[string]schema.Attribute{
				"min": schema.StringAttribute{
					Description: "Minute field.",
					Required:    true,
				},
				"hour": schema.StringAttribute{
					Description: "Hour field.",
					Required:    true,
				},
				"day": schema.StringAttribute{
					Description: "Day field.",
					Required:    true,
				},
				"month": schema.StringAttribute{
					Description: "Month field.",
					Required:    true,
				},
				"week": schema.StringAttribute{
					Description: "Week field.",
					Required:    true,
				},
			},
		},
		"crontab_duration": schema.Int64Attribute{
			Description: "Repeated mute duration in seconds.",
			Optional:    true,
			Validators: []validator.Int64{
				int64validator.AtLeast(0),
			},
		},
		"repeat_expire_time": schema.StringAttribute{
			Description: "Repeated mute expiration time in YYYY/MM/DD HH:mm:ss, or 0 for never expires.",
			Optional:    true,
			Validators: []validator.String{
				stringvalidator.RegexMatches(datetimePattern, "must use YYYY/MM/DD HH:mm:ss, empty string, or 0"),
			},
		},
		"timezone": schema.StringAttribute{
			Description: "Mute rule timezone.",
			Optional:    true,
			Computed:    true,
			Default:     stringdefault.StaticString("Asia/Shanghai"),
		},
		"declaration": schema.MapAttribute{
			Description: "Custom declaration information.",
			Optional:    true,
			ElementType: types.StringType,
		},
		"enabled": schema.BoolAttribute{
			Description: "Whether the mute rule is enabled. This manages API status 0/2 through the mute enable/disable endpoints.",
			Optional:    true,
			Computed:    true,
			Default:     booldefault.StaticBool(true),
		},
		"status": schema.Int64Attribute{
			Description: "Mute rule status returned by the API.",
			Computed:    true,
		},
		"create_at": schema.Int64Attribute{
			Description: "The timestamp seconds of the resource created at.",
			Computed:    true,
			PlanModifiers: []planmodifier.Int64{
				int64planmodifier.UseStateForUnknown(),
			},
		},
		"update_at": schema.Int64Attribute{
			Description: "The timestamp seconds of the resource updated at.",
			Computed:    true,
		},
		"workspace_uuid": schema.StringAttribute{
			Description: "The uuid of the workspace.",
			Computed:    true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
	},
}
