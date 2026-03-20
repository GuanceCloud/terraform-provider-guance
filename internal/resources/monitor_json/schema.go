package monitor_json

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

var resourceSchema = schema.Schema{
	Description:         "Monitor Json",
	MarkdownDescription: resourceDocument,
	Attributes: map[string]schema.Attribute{
		"uuid": schema.StringAttribute{
			Description: "The uuid of resource.",
			Computed:    true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"checker_json": schema.StringAttribute{
			Description: "The checker json configuration for import.",
			Optional:    true,
		},
		"checker_json_export": schema.StringAttribute{
			Description: "The exported checker json configuration of resource.",
			Optional:    true,
			Computed:    true,
			Sensitive:   true,
		},
		"type": schema.StringAttribute{
			Description: "Monitor type, default trigger, trigger: normal monitor, smartMonitor: smart monitor.",
			Optional:    true,
			Validators: []validator.String{
				stringvalidator.OneOf("trigger", "smartMonitor"),
			},
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.RequiresReplace(),
			},
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
