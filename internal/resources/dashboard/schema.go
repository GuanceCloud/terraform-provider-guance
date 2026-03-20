package dashboard

import (
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var resourceSchema = schema.Schema{
	Description:         "Dashboard",
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
			Description: "The name of resource.",
			Required:    true,
			Validators: []validator.String{
				stringvalidator.LengthAtMost(128),
			},
		},
		"desc": schema.StringAttribute{
			Description: "The description of resource.",
			Optional:    true,
			Validators: []validator.String{
				stringvalidator.LengthAtMost(2048),
			},
		},
		"identifier": schema.StringAttribute{
			Description: "The identifier of resource.",
			Optional:    true,
			Validators: []validator.String{
				stringvalidator.LengthAtMost(128),
			},
		},

		"tag_names": schema.ListAttribute{
			Description: "The tag names of resource.",
			Optional:    true,
			ElementType: types.StringType,
		},
		"template_info": schema.StringAttribute{
			Description: "The template info of resource.",
			Optional:    true,
		},

		"template_info_export": schema.StringAttribute{
			Description: "The exported template info of resource.",
			Optional:    true,
			Computed:    true,
			Sensitive:   true,
		},
		"specify_dashboard_uuid": schema.StringAttribute{
			Description: "The specify dashboard uuid of resource.",
			Optional:    true,
			Validators: []validator.String{
				stringvalidator.RegexMatches(
					regexp.MustCompile(`^dsbd_custom_[a-z0-9]{32}$`),
					"Must start with 'dsbd_custom_' followed by 32 lowercase alphanumeric characters",
				),
			},
		},
		"is_public": schema.Int64Attribute{
			Description: "Whether the resource is public.",
			Optional:    true,
		},
		"permission_set": schema.ListAttribute{
			Description: "The permission set of resource.",
			Optional:    true,
			ElementType: types.StringType,
		},
		"read_permission_set": schema.ListAttribute{
			Description: "The read permission set of resource.",
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
