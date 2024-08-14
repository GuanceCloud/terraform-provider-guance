package pipeline

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type pipelineResourceModel struct {
	CreateAt  types.String    `tfsdk:"create_at"`
	Name      types.String    `tfsdk:"name"`
	Source    []types.String  `tfsdk:"source"`
	Content   types.String    `tfsdk:"content"`
	TestData  types.String    `tfsdk:"test_data"`
	IsForce   types.Bool      `tfsdk:"is_force"`
	Category  types.String    `tfsdk:"category"`
	AsDefault types.Int64     `tfsdk:"as_default"`
	Type      types.String    `tfsdk:"type"`
	UUID      types.String    `tfsdk:"uuid"`
	Extend    *pipelineExtend `tfsdk:"extend"`
}

type pipelineExtend struct {
	AppID       []types.String `tfsdk:"app_id"`
	Measurement []types.String `tfsdk:"measurement"`
}

var resourceSchema = schema.Schema{
	Description:         "Pipeline",
	MarkdownDescription: resourceDocument,
	Attributes: map[string]schema.Attribute{
		"uuid": schema.StringAttribute{
			Description: "The uuid of the pipeline.",
			Computed:    true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"name": schema.StringAttribute{
			Description: "The name of the pipeline.",
			Required:    true,
		},
		"type": schema.StringAttribute{
			MarkdownDescription: "The type of the pipeline.Valid value: `local`, `central`",
			Required:            true,
			Validators: []validator.String{
				stringvalidator.OneOf("local", "central"),
			},
		},
		"source": schema.ListAttribute{
			Description: "Data source list",
			Required:    true,
			ElementType: types.StringType,
		},
		"content": schema.StringAttribute{
			Description: "Pipeline file content",
			Required:    true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.RequiresReplace(),
			},
		},
		"test_data": schema.StringAttribute{
			Description: "Test data",
			Required:    true,
		},
		"is_force": schema.BoolAttribute{
			MarkdownDescription: "Is Force Overwrite. If the field `as_default` is true, `is_force` will be set to be true automatically.",
			Required:            true,
			PlanModifiers: []planmodifier.Bool{
				boolplanmodifier.RequiresReplace(),
			},
		},
		"category": schema.StringAttribute{
			Description: "Category",

			MarkdownDescription: `
		Category, value must be one of: *logging*, *object*, *custom_object*, *network*, *tracing*, *rum*, *security*, *keyevent*, *metric*, *profiling*, other value will be ignored.
		`,
			Required: true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.RequiresReplace(),
			},
			Validators: []validator.String{
				stringvalidator.OneOf("logging", "object", "custom_object", "network", "tracing", "rum", "security", "keyevent", "metric", "profiling"),
			},
		},
		"as_default": schema.Int64Attribute{
			Description: "Is Default Pipeline",
			Required:    true,
			PlanModifiers: []planmodifier.Int64{
				int64planmodifier.RequiresReplace(),
			},
		},
		"create_at": schema.StringAttribute{
			Description: "The creation time of the resource, in seconds as a timestamp.",
			Computed:    true,
		},
		"extend": schema.SingleNestedAttribute{
			Optional: true,
			Attributes: map[string]schema.Attribute{
				"app_id": schema.ListAttribute{
					Optional:    true,
					ElementType: types.StringType,
				},
				"measurement": schema.ListAttribute{
					Optional:    true,
					ElementType: types.StringType,
				},
			},
		},
	},
}
