package slo

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/float64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var resourceSchema = schema.Schema{
	Description:         "SLO (Service Level Objective)",
	MarkdownDescription: resourceDocument,
	Attributes: map[string]schema.Attribute{
		"uuid": schema.StringAttribute{
			Description: "The uuid of resource.",
			Computed:    true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"name": schema.StringAttribute{
			Description: "SLO name",
			Required:    true,
			Validators: []validator.String{
				stringvalidator.LengthAtMost(256),
			},
		},
		"interval": schema.StringAttribute{
			Description: "Detection frequency",
			Required:    true,
			Validators: []validator.String{
				stringvalidator.OneOf("5m", "10m"),
			},
		},
		"goal": schema.Float64Attribute{
			Description: "SLO expected goal, range: 0-100",
			Required:    true,
			Validators: []validator.Float64{
				float64validator.Between(0, 100),
			},
		},
		"min_goal": schema.Float64Attribute{
			Description: "SLO minimum goal, range: 0-100, must be less than goal",
			Required:    true,
			Validators: []validator.Float64{
				float64validator.Between(0, 100),
			},
		},
		"sli_uuids": schema.ListAttribute{
			Description: "SLI UUID list",
			Required:    true,
			ElementType: types.StringType,
		},
		"describe": schema.StringAttribute{
			Description: "SLO description",
			Optional:    true,
			Validators: []validator.String{
				stringvalidator.LengthAtMost(3000),
			},
		},
		"alert_policy_uuids": schema.ListAttribute{
			Description: "Alert policy UUIDs",
			Optional:    true,
			ElementType: types.StringType,
		},
		"tags": schema.ListAttribute{
			Description: "Tag names for filtering",
			Optional:    true,
			ElementType: types.StringType,
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
