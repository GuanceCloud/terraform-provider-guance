package synthetics_test

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// AtLeastOneFieldRequiredValidator validates that at least one of the specified fields is set
func AtLeastOneFieldRequiredValidator(fieldNames []string) validator.Object {
	return atLeastOneFieldRequiredValidator{fieldNames: fieldNames}
}

type atLeastOneFieldRequiredValidator struct {
	fieldNames []string
}

func (v atLeastOneFieldRequiredValidator) Description(ctx context.Context) string {
	return "at least one field must be set"
}

func (v atLeastOneFieldRequiredValidator) MarkdownDescription(ctx context.Context) string {
	return "at least one field must be set"
}

func (v atLeastOneFieldRequiredValidator) ValidateObject(ctx context.Context, req validator.ObjectRequest, resp *validator.ObjectResponse) {
	for _, fieldName := range v.fieldNames {
		if req.ConfigValue.Attributes()[fieldName].IsUnknown() || req.ConfigValue.Attributes()[fieldName].IsNull() {
			continue
		}
		return
	}

	resp.Diagnostics.AddError(
		"Missing required field",
		"At least one field must be set",
	)
}

// headerConditionObjectType defines the object type for header conditions
var headerConditionObjectType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"contains":        types.StringType,
		"not_contains":    types.StringType,
		"is":              types.StringType,
		"is_not":          types.StringType,
		"match_regex":     types.StringType,
		"not_match_regex": types.StringType,
	},
}

// stringMatchConditions returns a map of string match conditions for schema attributes
func stringMatchConditions() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"contains": schema.StringAttribute{
			Description: "The value should contain this string.",
			Optional:    true,
			Computed:    true,
			Default:     stringdefault.StaticString(""),
		},
		"not_contains": schema.StringAttribute{
			Description: "The value should not contain this string.",
			Optional:    true,
			Computed:    true,
			Default:     stringdefault.StaticString(""),
		},
		"is": schema.StringAttribute{
			Description: "The value should be exactly this string.",
			Optional:    true,
			Computed:    true,
			Default:     stringdefault.StaticString(""),
		},
		"is_not": schema.StringAttribute{
			Description: "The value should not be exactly this string.",
			Optional:    true,
			Computed:    true,
			Default:     stringdefault.StaticString(""),
		},
		"match_regex": schema.StringAttribute{
			Description: "The value should match this regex.",
			Optional:    true,
			Computed:    true,
			Default:     stringdefault.StaticString(""),
		},
		"not_match_regex": schema.StringAttribute{
			Description: "The value should not match this regex.",
			Optional:    true,
			Computed:    true,
			Default:     stringdefault.StaticString(""),
		},
	}
}

var resourceSchema = schema.Schema{
	Description: "Synthetics Test",
	Attributes: map[string]schema.Attribute{
		"uuid": schema.StringAttribute{
			Description: "The uuid of resource.",
			Computed:    true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"type": schema.StringAttribute{
			Description: "The type of synthetics test. Valid values: http, tcp, dns, browser, icmp, websocket, multi, grpc.",
			Required:    true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.RequiresReplace(),
				stringplanmodifier.UseStateForUnknown(),
			},
			Validators: []validator.String{
				stringvalidator.OneOf("http", "tcp", "dns", "browser", "icmp", "websocket", "multi", "grpc"),
			},
		},
		"regions": schema.ListAttribute{
			Description: "The regions where the test will be executed.",
			Required:    true,
			ElementType: types.StringType,
		},
		"task": schema.SingleNestedAttribute{
			Description: "The task configuration.",
			Required:    true,
			Attributes: map[string]schema.Attribute{
				"url": schema.StringAttribute{
					Description: "The URL to test. Only applicable for http tests.",
					Optional:    true,
				},
				"method": schema.StringAttribute{
					Description: "The HTTP method to use. Only applicable for http tests.",
					Optional:    true,
				},
				"name": schema.StringAttribute{
					Description: "The name of the task.",
					Required:    true,
					Validators: []validator.String{
						stringvalidator.LengthAtMost(256),
					},
				},
				"status": schema.StringAttribute{
					Description: "The status of the task. Valid values: ok, stop.",
					Optional:    true,
					Computed:    true,
					Default:     stringdefault.StaticString("ok"),
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.UseStateForUnknown(),
					},
					Validators: []validator.String{
						stringvalidator.OneOf("ok", "stop"),
					},
				},
				"frequency": schema.StringAttribute{
					Description: "The frequency of the test. Valid values: 1m, 5m, 15m, 30m, 1h, 6h, 12h, 24h.",
					Required:    true,
					Validators: []validator.String{
						stringvalidator.OneOf("1m", "5m", "15m", "30m", "1h", "6h", "12h", "24h"),
					},
				},
				"schedule_type": schema.StringAttribute{
					Description: "The schedule type of the test. Valid values: frequency, crontab.",
					Optional:    true,
				},
				"crontab": schema.StringAttribute{
					Description: "The crontab expression for the test when schedule_type is crontab.",
					Optional:    true,
					Computed:    true,
					Default:     stringdefault.StaticString(""),
				},
				"desc": schema.StringAttribute{
					Description: "The description of the task. Only applicable for multi tests.",
					Optional:    true,
				},
				"steps": schema.ListNestedAttribute{
					Description: "The steps for multi-step tests.",
					Optional:    true,
					NestedObject: schema.NestedAttributeObject{
						Attributes: map[string]schema.Attribute{
							"type": schema.StringAttribute{
								Description: "The type of the step. Valid values: http, wait.",
								Required:    true,
								Validators: []validator.String{
									stringvalidator.OneOf("http", "wait"),
								},
							},
							"task": schema.StringAttribute{
								Description: "The task configuration for HTTP steps as JSON string.",
								Optional:    true,
							},
							"allow_failure": schema.BoolAttribute{
								Description: "Whether to allow failure of this step.",
								Optional:    true,
							},
							"extracted_vars": schema.ListNestedAttribute{
								Description: "The variables extracted from the step.",
								Optional:    true,
								NestedObject: schema.NestedAttributeObject{
									Attributes: map[string]schema.Attribute{
										"name": schema.StringAttribute{
											Description: "The name of the extracted variable.",
											Required:    true,
										},
										"field": schema.StringAttribute{
											Description: "The field to extract the variable from.",
											Required:    true,
										},
										"secure": schema.BoolAttribute{
											Description: "Whether the variable is secure.",
											Optional:    true,
										},
									},
								},
							},
							"value": schema.Int64Attribute{
								Description: "The wait time in seconds for wait steps.",
								Optional:    true,
							},
							"retry": schema.SingleNestedAttribute{
								Description: "The retry configuration for the step.",
								Optional:    true,
								Attributes: map[string]schema.Attribute{
									"retry": schema.Int64Attribute{
										Description: "The number of retries.",
										Optional:    true,
									},
									"interval": schema.Int64Attribute{
										Description: "The retry interval in milliseconds.",
										Optional:    true,
									},
								},
							},
						},
					},
				},
				"advance_options": schema.SingleNestedAttribute{
					Description: "Advanced options for the task.",
					Optional:    true,
					Attributes: map[string]schema.Attribute{
						"request_options": schema.SingleNestedAttribute{
							Description: "Request options.",
							Optional:    true,
							Attributes: map[string]schema.Attribute{
								"follow_redirect": schema.BoolAttribute{
									Description: "Whether to follow redirects. Only applicable for http tests.",
									Computed:    true,
									Default:     booldefault.StaticBool(false),
									Optional:    true,
								},
								"headers": schema.MapAttribute{
									Description: "Request headers. Only applicable for http and websocket tests.",
									Optional:    true,
									ElementType: types.StringType,
								},
								"cookies": schema.StringAttribute{
									Description: "Cookies to send with the request. Only applicable for http tests.",
									Optional:    true,
									Computed:    true,
									Default:     stringdefault.StaticString(""),
								},
								"auth": schema.SingleNestedAttribute{
									Description: "Authentication credentials. Only applicable for http and websocket tests.",
									Optional:    true,
									Attributes: map[string]schema.Attribute{
										"username": schema.StringAttribute{
											Description: "Username for authentication.",
											Optional:    true,
										},
										"password": schema.StringAttribute{
											Description: "Password for authentication.",
											Optional:    true,
										},
									},
								},
								"timeout": schema.StringAttribute{
									Description: "The timeout for the request.",
									Optional:    true,
									Computed:    true,
									Default:     stringdefault.StaticString(""),
								},
								"metadata": schema.MapAttribute{
									Description: "gRPC metadata. Only applicable for grpc tests.",
									Optional:    true,
									ElementType: types.StringType,
								},
								"proto_files": schema.SingleNestedAttribute{
									Description: "Proto files for gRPC tests. Only applicable for grpc tests.",
									Optional:    true,
									Attributes: map[string]schema.Attribute{
										"protofiles": schema.MapAttribute{
											Description: "Proto files as map of file paths to content.",
											Optional:    true,
											ElementType: types.StringType,
										},
										"full_method": schema.StringAttribute{
											Description: "The full method name for gRPC tests.",
											Optional:    true,
										},
										"request": schema.StringAttribute{
											Description: "The request body for gRPC tests.",
											Optional:    true,
										},
									},
								},
								"reflection": schema.SingleNestedAttribute{
									Description: "Reflection configuration for gRPC tests. Only applicable for grpc tests.",
									Optional:    true,
									Attributes: map[string]schema.Attribute{
										"full_method": schema.StringAttribute{
											Description: "The full method name for gRPC tests.",
											Optional:    true,
										},
										"request": schema.StringAttribute{
											Description: "The request body for gRPC tests.",
											Optional:    true,
										},
									},
								},
								"health_check": schema.SingleNestedAttribute{
									Description: "Health check configuration for gRPC tests. Only applicable for grpc tests.",
									Optional:    true,
									Attributes: map[string]schema.Attribute{
										"service": schema.StringAttribute{
											Description: "The service name for health check.",
											Optional:    true,
										},
									},
								},
							},
						},
						"request_body": schema.SingleNestedAttribute{
							Description: "Request body.",
							Optional:    true,
							Attributes: map[string]schema.Attribute{
								"body_type": schema.StringAttribute{
									Description: "The type of the request body.",
									Optional:    true,
								},
								"body": schema.StringAttribute{
									Description: "The request body content.",
									Optional:    true,
								},
								"form": schema.MapAttribute{
									Description: "Form data for the request body.",
									Optional:    true,
									ElementType: types.StringType,
								},
							},
						},
						"certificate": schema.SingleNestedAttribute{
							Description: "Certificate configuration.",
							Optional:    true,
							Attributes: map[string]schema.Attribute{
								"ignore_server_certificate_error": schema.BoolAttribute{
									Description: "Whether to ignore server certificate errors.",
									Optional:    true,
									Computed:    true,
									Default:     booldefault.StaticBool(false),
								},
								"private_key": schema.StringAttribute{
									Description: "The private key for certificate authentication.",
									Optional:    true,
									Computed:    true,
									Default:     stringdefault.StaticString(""),
								},
								"certificate": schema.StringAttribute{
									Description: "The certificate for authentication.",
									Optional:    true,
									Computed:    true,
									Default:     stringdefault.StaticString(""),
								},
							},
						},
						"proxy": schema.SingleNestedAttribute{
							Description: "Proxy configuration.",
							Optional:    true,
							Attributes: map[string]schema.Attribute{
								"url": schema.StringAttribute{
									Description: "The proxy URL.",
									Optional:    true,
								},
								"headers": schema.MapAttribute{
									Description: "Proxy headers.",
									Optional:    true,
									ElementType: types.StringType,
								},
							},
						},
						"secret": schema.SingleNestedAttribute{
							Description: "Secret configuration.",
							Optional:    true,
							Attributes: map[string]schema.Attribute{
								"not_save": schema.BoolAttribute{
									Description: "Whether to not save the secret.",
									Optional:    true,
								},
							},
						},
						"request_timeout": schema.StringAttribute{
							Description: "The request timeout.",
							Optional:    true,
							Computed:    true,
							Default:     stringdefault.StaticString(""),
						},
						"post_script": schema.StringAttribute{
							Description: "The post script for the test. Only applicable for grpc tests.",
							Optional:    true,
							Computed:    true,
							Default:     stringdefault.StaticString(""),
						},
					},
				},
				"advance_options_headless": schema.SingleNestedAttribute{
					Description: "Advanced options for browser tasks.",
					Optional:    true,
					Attributes: map[string]schema.Attribute{
						"request_options": schema.SingleNestedAttribute{
							Description: "Request options.",
							Optional:    true,
							Attributes: map[string]schema.Attribute{
								"follow_redirect": schema.BoolAttribute{
									Description: "Whether to follow redirects.",
									Optional:    true,
								},
								"headers": schema.MapAttribute{
									Description: "Request headers.",
									Optional:    true,
									ElementType: types.StringType,
								},
								"cookies": schema.StringAttribute{
									Description: "Cookies to send with the request.",
									Optional:    true,
								},
								"auth": schema.SingleNestedAttribute{
									Description: "Authentication credentials.",
									Optional:    true,
									Attributes: map[string]schema.Attribute{
										"username": schema.StringAttribute{
											Description: "Username for authentication.",
											Optional:    true,
										},
										"password": schema.StringAttribute{
											Description: "Password for authentication.",
											Optional:    true,
										},
									},
								},
								"timeout": schema.StringAttribute{
									Description: "The timeout for the request.",
									Optional:    true,
								},
							},
						},
						"request_body": schema.SingleNestedAttribute{
							Description: "Request body.",
							Optional:    true,
							Attributes: map[string]schema.Attribute{
								"body_type": schema.StringAttribute{
									Description: "The type of the request body.",
									Optional:    true,
								},
								"body": schema.StringAttribute{
									Description: "The request body content.",
									Optional:    true,
								},
							},
						},
						"request_timeout": schema.StringAttribute{
							Description: "The request timeout.",
							Optional:    true,
						},
					},
				},
				"success_when_logic": schema.StringAttribute{
					Description: "The logic to use for success conditions. Valid values: and, or.",
					Optional:    true,
					Computed:    true,
					Default:     stringdefault.StaticString(""),
					Validators: []validator.String{
						stringvalidator.OneOf("and", "or"),
					},
				},
				"success_when": schema.ListNestedAttribute{
					Description: "The conditions that determine success.",
					Optional:    true,
					NestedObject: schema.NestedAttributeObject{
						Attributes: map[string]schema.Attribute{
							"body": schema.ListNestedAttribute{
								Description: "Body conditions.",
								Optional:    true,
								NestedObject: schema.NestedAttributeObject{
									Attributes: stringMatchConditions(),
								},
							},
							"header": schema.MapAttribute{
								Description: "Header conditions as JSON strings.",
								Optional:    true,
								ElementType: types.ListType{
									ElemType: types.StringType,
								},
							},
							"status_code": schema.ListNestedAttribute{
								Description: "Status code conditions.",
								Optional:    true,
								NestedObject: schema.NestedAttributeObject{
									Attributes: stringMatchConditions(),
								},
							},
							"response_time": schema.ListNestedAttribute{
								Description: "The response time conditions.",
								Optional:    true,
								NestedObject: schema.NestedAttributeObject{
									Attributes: map[string]schema.Attribute{
										"is_contain_dns": schema.BoolAttribute{
											Description: "Whether the response time should contain DNS. This field is only applicable for tcp and websocket test.",
											Optional:    true,
										},
										"target": schema.StringAttribute{
											Description: "The target value for comparison.",
											Optional:    true,
										},
										"func": schema.StringAttribute{
											Description: "The function to apply to the response time. Valid values: min, max, avg, std. This field is only applicable for icmp test.",
											Optional:    true,
										},
										"op": schema.StringAttribute{
											Description: "The comparison operator. Valid values: eq, lt, leq, gt, geq. This field is only applicable for icmp test.",
											Optional:    true,
										},
									},
								},
							},
							"response_message": schema.ListNestedAttribute{
								Description: "Response message conditions.",
								Optional:    true,
								NestedObject: schema.NestedAttributeObject{
									Attributes: stringMatchConditions(),
								},
							},
							"hops": schema.ListNestedAttribute{
								Description: "Network hops conditions.",
								Optional:    true,
								NestedObject: schema.NestedAttributeObject{
									Attributes: map[string]schema.Attribute{
										"op": schema.StringAttribute{
											Description: "The comparison operator. Valid values: eq, lt, leq, gt, geq.",
											Optional:    true,
										},
										"target": schema.Float64Attribute{
											Description: "The target value for comparison.",
											Optional:    true,
										},
									},
								},
							},
							"packet_loss_percent": schema.ListNestedAttribute{
								Description: "Packet loss percentage conditions.",
								Optional:    true,
								NestedObject: schema.NestedAttributeObject{
									Attributes: map[string]schema.Attribute{
										"op": schema.StringAttribute{
											Description: "The comparison operator. Valid values: eq, lt, leq, gt, geq.",
											Optional:    true,
										},
										"target": schema.Float64Attribute{
											Description: "The target value for comparison.",
											Optional:    true,
										},
									},
								},
							},
							"packets": schema.ListNestedAttribute{
								Description: "Packets conditions.",
								Optional:    true,
								NestedObject: schema.NestedAttributeObject{
									Attributes: map[string]schema.Attribute{
										"op": schema.StringAttribute{
											Description: "The comparison operator. Valid values: eq, lt, leq, gt, geq.",
											Optional:    true,
										},
										"target": schema.Float64Attribute{
											Description: "The target value for comparison.",
											Optional:    true,
										},
									},
								},
							},
						},
					},
				},
				"enable_traceroute": schema.BoolAttribute{
					Description: "Whether to enable traceroute. Only applicable for tcp and icmp tests.",
					Optional:    true,
					Computed:    true,
					Default:     booldefault.StaticBool(false),
				},
				"packet_count": schema.Int64Attribute{
					Description: "The number of packets to send for ICMP tests. Only applicable for icmp tests.",
					Optional:    true,
					Computed:    true,
					Default:     int64default.StaticInt64(0),
				},
				"server": schema.StringAttribute{
					Description: "The server to test for gRPC tests. Only applicable for grpc tests.",
					Optional:    true,
					Computed:    true,
					Default:     stringdefault.StaticString(""),
				},
				"host": schema.StringAttribute{
					Description: "The host to test for TCP/ICMP tests. Only applicable for tcp and icmp tests.",
					Optional:    true,
					Computed:    true,
					Default:     stringdefault.StaticString(""),
				},
				"port": schema.StringAttribute{
					Description: "The port to test for TCP tests. Only applicable for tcp tests.",
					Optional:    true,
					Computed:    true,
					Default:     stringdefault.StaticString(""),
				},
				"timeout": schema.StringAttribute{
					Description: "The timeout for TCP/ICMP tests. Only applicable for tcp and icmp tests.",
					Optional:    true,
					Computed:    true,
					Default:     stringdefault.StaticString(""),
				},
				"message": schema.StringAttribute{
					Description: "The message to send for WebSocket tests. Only applicable for websocket tests.",
					Optional:    true,
					Computed:    true,
					Default:     stringdefault.StaticString(""),
				},
				"post_mode": schema.StringAttribute{
					Description: "The post-test mode. Valid values: default, script. Only applicable for http tests.",
					Optional:    true,
					Computed:    true,
					Default:     stringdefault.StaticString(""),
					Validators: []validator.String{
						stringvalidator.OneOf("default", "script"),
					},
				},
				"post_script": schema.StringAttribute{
					Description: "The post-test script content. Only applicable for http and grpc tests.",
					Optional:    true,
					Computed:    true,
					Default:     stringdefault.StaticString(""),
				},
			},
		},
		"tags": schema.ListAttribute{
			Description: "The tags to associate with the test.",
			Optional:    true,
			ElementType: types.StringType,
			PlanModifiers: []planmodifier.List{
				listplanmodifier.UseStateForUnknown(),
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
