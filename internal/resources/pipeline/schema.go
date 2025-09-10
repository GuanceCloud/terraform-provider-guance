package pipeline

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type pipelineResourceModel struct {
	CreateAt          types.String    `tfsdk:"create_at"`
	UpdateAt          types.String    `tfsdk:"update_at"`
	Name              types.String    `tfsdk:"name"`
	Source            []types.String  `tfsdk:"source"`
	Content           types.String    `tfsdk:"content"`
	TestData          types.String    `tfsdk:"test_data"`
	IsForce           types.Bool      `tfsdk:"is_force"`
	IsDisabled        types.Bool      `tfsdk:"is_disabled"`
	Category          types.String    `tfsdk:"category"`
	AsDefault         types.Int64     `tfsdk:"as_default"`
	EnableByLogBackup types.Int64     `tfsdk:"enable_by_log_backup"`
	Status            types.Int64     `tfsdk:"status"`
	Type              types.String    `tfsdk:"type"`
	DataType          types.String    `tfsdk:"data_type"`
	UUID              types.String    `tfsdk:"uuid"`
	WorkspaceUUID     types.String    `tfsdk:"workspace_uuid"`
	Extend            *pipelineExtend `tfsdk:"extend"`
}

type pipelineExtend struct {
	AppID        []types.String `tfsdk:"app_id"`
	Measurement  []types.String `tfsdk:"measurement"`
	LoggingIndex types.String   `tfsdk:"logging_index"`
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
		"workspace_uuid": schema.StringAttribute{
			Description: "The uuid of the workspace.",
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
		"data_type": schema.StringAttribute{
			MarkdownDescription: "The type of the data. Valid value: `line_protocol`, `json`, `message`.",
			Optional:            true,
			Validators: []validator.String{
				stringvalidator.OneOf("line_protocol", "json", "message"),
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
		},
		"test_data": schema.StringAttribute{
			Description: "Test data",
			Required:    true,
		},
		"is_force": schema.BoolAttribute{
			MarkdownDescription: "Is Force Overwrite. If the field `as_default` is true, `is_force` will be set to be true automatically.",
			Optional:            true,
		},
		"is_disabled": schema.BoolAttribute{
			Description: "Is Disabled",
			Optional:    true,
		},
		"category": schema.StringAttribute{
			Description: "Category",

			MarkdownDescription: `
		Category, value must be one of: *logging*, *object*, *custom_object*, *network*, *tracing*, *rum*, *security*, *keyevent*, *metric*, *profiling*, *dialtesting*, *billing*, *keyevent*, other value will be ignored.
		`,
			Required: true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.RequiresReplace(),
			},
			Validators: []validator.String{
				stringvalidator.OneOf("logging", "object", "custom_object", "network", "tracing", "rum", "security", "keyevent", "metric", "profiling", "dialtesting", "billing", "keyevent"),
			},
		},
		"as_default": schema.Int64Attribute{
			Description: "Is Default Pipeline",
			Optional:    true,
		},
		"enable_by_log_backup": schema.Int64Attribute{
			Description: "Enable pipeline processing (1=enable, 0=disable) for forwarded data.",
			Optional:    true,
		},
		"status": schema.Int64Attribute{
			Description: "The status of the pipeline. Valid value: `0`, `2`.",
			Computed:    true,
		},
		"create_at": schema.StringAttribute{
			Description: "The creation time of the resource, in seconds as a timestamp.",
			Computed:    true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"update_at": schema.StringAttribute{
			Description: "The update time of the resource, in seconds as a timestamp.",
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
				"logging_index": schema.StringAttribute{
					Optional: true,
				},
			},
		},
	},
}
