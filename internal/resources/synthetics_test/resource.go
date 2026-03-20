package synthetics_test

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/GuanceCloud/cliutils/dialtesting"
	"github.com/GuanceCloud/terraform-provider-guance/internal/api"
	"github.com/GuanceCloud/terraform-provider-guance/internal/consts"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &syntheticsTestResource{}
	_ resource.ResourceWithConfigure   = &syntheticsTestResource{}
	_ resource.ResourceWithImportState = &syntheticsTestResource{}
)

// NewSyntheticsTestResource creates a new synthetics test resource
func NewSyntheticsTestResource() resource.Resource {
	return &syntheticsTestResource{}
}

// syntheticsTestResource is the resource implementation.
type syntheticsTestResource struct {
	client *api.Client
}

// Schema defines the schema for the data source.
func (r *syntheticsTestResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = resourceSchema
}

// Configure adds the provider configured client to the data source.
func (r *syntheticsTestResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.client = req.ProviderData.(*api.Client)
}

// Metadata returns the data source type name.
func (r *syntheticsTestResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_synthetics_test"
}

func getTask(taskType string) dialtesting.ITask {
	switch taskType {
	case "http":
		return &dialtesting.HTTPTask{}
	default:
		return nil
	}
}

// Create creates the resource and sets the initial Terraform state.
func (r *syntheticsTestResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan syntheticsTestResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	item := r.getSyntheticsTestFromPlan(ctx, &plan)
	content := &api.SyntheticsTest{}
	var err error

	if err != nil {
		resp.Diagnostics.AddError(
			"Error marshaling item",
			fmt.Sprintf("Failed to marshal item to JSON for synthetics test: %s", err),
		)
		return
	}
	// Check if it's a multi-step test
	if plan.Type.ValueString() == "multi" {
		err = r.client.MultiCreate(consts.TypeNameSyntheticsTest, item, content)
	} else {
		err = r.client.Create(consts.TypeNameSyntheticsTest, item, content)
	}

	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating synthetics test",
			"Could not create synthetics test, unexpected error: "+err.Error(),
		)
		return
	}

	// Update plan with response data - only update computed fields
	plan.UUID = types.StringValue(content.UUID)
	plan.WorkspaceUUID = types.StringValue(content.WorkspaceUUID)
	plan.CreateAt = types.Int64Value(content.CreateAt)
	plan.UpdateAt = types.Int64Value(content.CreateAt) // fix updateat -1

	status := "ok"
	if plan.Task != nil && plan.Task.Status.ValueString() != "" {
		status = plan.Task.Status.ValueString()
	}
	// Convert status to isDisable bool
	isDisable := status == "stop"
	// Call the SetStatus API to set the status
	err = r.client.SetStatus(consts.TypeNameSyntheticsTest, []string{plan.UUID.ValueString()}, isDisable)
	if err != nil {
		resp.Diagnostics.AddError("Error setting status", fmt.Sprintf("Failed to set status: %s", err))
		return
	}

	// Update tags from API response
	// Only update tags if API returns tags
	if len(content.TagInfo) > 0 || len(content.Tags) > 0 {
		tagNames := getTags(content.TagInfo)
		if !isSameTags(plan.Tags, tagNames) {
			plan.Tags = []types.String{}
			for _, tag := range tagNames {
				plan.Tags = append(plan.Tags, types.StringValue(tag))
			}
		}
	}

	// Set state
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *syntheticsTestResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state syntheticsTestResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get client from provider
	if r.client == nil {
		resp.Diagnostics.AddError("Client not configured", "Client not configured for synthetics test resource")
		return
	}

	content := &api.SyntheticsTest{}
	err := r.client.Read(consts.TypeNameSyntheticsTest, state.UUID.ValueString(), content)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading synthetics test",
			"Could not read synthetics test, unexpected error: "+err.Error(),
		)
		return
	}

	testType := content.Type

	if state.Task == nil {
		state.Task = &taskConfig{}
	}
	// update common fields
	r.updateTaskState(ctx, state.Task, content.Task)

	switch testType {
	case TaskTypeHTTP:
		r.updateHTTPTaskState(ctx, &state, content)
	case TaskTypeTCP:
		r.updateTCPTaskState(ctx, &state, content)
	case TaskTypeWebSocket:
		r.updateWebSocketTaskState(ctx, &state, content)
	case TaskTypeICMP:
		r.updateICMPTaskState(ctx, &state, content)
	case TaskTypeGRPC:
		r.updateGRPCTaskState(ctx, &state, content)
	case "multi":
		// For multi-step tests, update basic task state
		r.updateMultiTaskState(ctx, &state, content)
	default:
		resp.Diagnostics.AddError(
			"Error reading synthetics test",
			"Could not read synthetics test, unexpected error: "+fmt.Sprintf("unsupported test type: %s", testType),
		)
		return
	}

	state.Type = types.StringValue(testType)
	state.Regions = []types.String{}
	for _, region := range content.Regions {
		state.Regions = append(state.Regions, types.StringValue(region))
	}

	tagNames := getTags(content.TagInfo)
	if !isSameTags(state.Tags, tagNames) {
		state.Tags = []types.String{}
		for _, tag := range tagNames {
			state.Tags = append(state.Tags, types.StringValue(tag))
		}
	}

	// Update state with all fields from API response
	state.UUID = types.StringValue(content.UUID)
	state.WorkspaceUUID = types.StringValue(content.WorkspaceUUID)
	state.CreateAt = types.Int64Value(content.CreateAt)
	state.UpdateAt = types.Int64Value(content.UpdateAt)

	// Set state
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *syntheticsTestResource) updateTaskState(ctx context.Context, task *taskConfig, apiTask *api.TaskConfig) {
	if task == nil || apiTask == nil {
		return
	}

	task.Status = types.StringValue(apiTask.Status)
	task.Name = types.StringValue(apiTask.Name)
	task.Frequency = types.StringValue(apiTask.Frequency)
	task.Crontab = types.StringValue(apiTask.Crontab)
	task.PostMode = types.StringValue(apiTask.PostMode)
	task.Timeout = types.StringValue(apiTask.Timeout)
	task.PostScript = types.StringValue(apiTask.PostScript)
	task.SuccessWhenLogic = types.StringValue(apiTask.SuccessWhenLogic)
}

func (r *syntheticsTestResource) updateHTTPTaskState(ctx context.Context, state *syntheticsTestResourceModel, content *api.SyntheticsTest) {
	if content.Task != nil {

		state.Task.URL = types.StringValue(content.Task.URL)
		state.Task.Method = types.StringValue(content.Task.Method)

		// update advance options
		if content.Task.AdvanceOptions != nil {
			if state.Task.AdvanceOptions == nil {
				state.Task.AdvanceOptions = &advanceOptions{}
			}

			// update proxy
			if content.Task.AdvanceOptions.Proxy != nil {
				headers := map[string]types.String{}
				for k, v := range content.Task.AdvanceOptions.Proxy.Headers {
					headers[k] = types.StringValue(v)
				}
				state.Task.AdvanceOptions.Proxy = &proxy{
					URL:     types.StringValue(content.Task.AdvanceOptions.Proxy.URL),
					Headers: headers,
				}
			} else {
				state.Task.AdvanceOptions.Proxy = nil
			}

			// update certificate
			if content.Task.AdvanceOptions.Certificate != nil {
				state.Task.AdvanceOptions.Certificate = &certificate{
					IgnoreServerCertificateError: types.BoolValue(content.Task.AdvanceOptions.Certificate.IgnoreServerCertificateError),
					PrivateKey:                   types.StringValue(content.Task.AdvanceOptions.Certificate.PrivateKey),
					Certificate:                  types.StringValue(content.Task.AdvanceOptions.Certificate.Certificate),
				}
			} else {
				state.Task.AdvanceOptions.Certificate = nil
			}

			// update request body
			if content.Task.AdvanceOptions.RequestBody != nil {
				state.Task.AdvanceOptions.RequestBody = &requestBody{
					BodyType: types.StringValue(content.Task.AdvanceOptions.RequestBody.BodyType),
					Body:     types.StringValue(content.Task.AdvanceOptions.RequestBody.Body),
				}

				if content.Task.AdvanceOptions.RequestBody.Form != nil {
					state.Task.AdvanceOptions.RequestBody.Form = map[string]types.String{}
					for k, v := range content.Task.AdvanceOptions.RequestBody.Form {
						state.Task.AdvanceOptions.RequestBody.Form[k] = types.StringValue(v)
					}
				}
			} else {
				state.Task.AdvanceOptions.RequestBody = nil
			}

			// update request options
			if content.Task.AdvanceOptions.RequestOptions != nil {
				state.Task.AdvanceOptions.RequestOptions = &requestOptions{
					FollowRedirect: types.BoolValue(content.Task.AdvanceOptions.RequestOptions.FollowRedirect),
					Cookies:        types.StringValue(content.Task.AdvanceOptions.RequestOptions.Cookies),
					Timeout:        types.StringValue(content.Task.AdvanceOptions.RequestOptions.Timeout),
				}
				if content.Task.AdvanceOptions.RequestOptions.Headers != nil {
					state.Task.AdvanceOptions.RequestOptions.Headers = map[string]types.String{}
					for k, v := range content.Task.AdvanceOptions.RequestOptions.Headers {
						state.Task.AdvanceOptions.RequestOptions.Headers[k] = types.StringValue(v)
					}
				}
				if content.Task.AdvanceOptions.RequestOptions.Auth != nil {
					state.Task.AdvanceOptions.RequestOptions.Auth = &auth{
						Username: types.StringValue(content.Task.AdvanceOptions.RequestOptions.Auth.Username),
						Password: types.StringValue(content.Task.AdvanceOptions.RequestOptions.Auth.Password),
					}
				}
				// metadata
				if content.Task.AdvanceOptions.RequestOptions.Metadata != nil {
					state.Task.AdvanceOptions.RequestOptions.Metadata = map[string]types.String{}
					for k, v := range content.Task.AdvanceOptions.RequestOptions.Metadata {
						state.Task.AdvanceOptions.RequestOptions.Metadata[k] = types.StringValue(v)
					}
				}
				// proto files
				if content.Task.AdvanceOptions.RequestOptions.ProtoFiles != nil {
					state.Task.AdvanceOptions.RequestOptions.ProtoFiles = &protoFiles{
						ProtoFiles: map[string]types.String{},
						FullMethod: types.StringValue(content.Task.AdvanceOptions.RequestOptions.ProtoFiles.FullMethod),
						Request:    types.StringValue(content.Task.AdvanceOptions.RequestOptions.ProtoFiles.Request),
					}
					if content.Task.AdvanceOptions.RequestOptions.ProtoFiles.ProtoFiles != nil {
						for k, v := range content.Task.AdvanceOptions.RequestOptions.ProtoFiles.ProtoFiles {
							state.Task.AdvanceOptions.RequestOptions.ProtoFiles.ProtoFiles[k] = types.StringValue(v)
						}
					}
				}
				// reflection
				if content.Task.AdvanceOptions.RequestOptions.Reflection != nil {
					state.Task.AdvanceOptions.RequestOptions.Reflection = &reflection{
						FullMethod: types.StringValue(content.Task.AdvanceOptions.RequestOptions.Reflection.FullMethod),
						Request:    types.StringValue(content.Task.AdvanceOptions.RequestOptions.Reflection.Request),
					}
				}
				// health check
				if content.Task.AdvanceOptions.RequestOptions.HealthCheck != nil {
					state.Task.AdvanceOptions.RequestOptions.HealthCheck = &healthCheck{
						Service: types.StringValue(content.Task.AdvanceOptions.RequestOptions.HealthCheck.Service),
					}
				}
			} else {
				state.Task.AdvanceOptions.RequestOptions = nil
			}

			// update request timeout
			if content.Task.AdvanceOptions.RequestTimeout != "" {
				state.Task.AdvanceOptions.RequestTimeout = types.StringValue(content.Task.AdvanceOptions.RequestTimeout)
			} else {
				state.Task.AdvanceOptions.RequestTimeout = types.StringValue("")
			}
		} else {
			state.Task.AdvanceOptions = nil
		}

		// update success when
		if len(content.Task.SuccessWhen) > 0 {
			state.Task.SuccessWhen = []successWhenItem{}
			for _, swItem := range content.Task.SuccessWhen {
				successWhenItem := successWhenItem{}

				// update body conditions
				if len(swItem.Body) > 0 {
					successWhenItem.Body = []bodyCondition{}
					for _, bodyItem := range swItem.Body {
						bodyCondition := bodyCondition{
							Contains:      types.StringValue(bodyItem.Contains),
							NotContains:   types.StringValue(bodyItem.NotContains),
							Is:            types.StringValue(bodyItem.Is),
							IsNot:         types.StringValue(bodyItem.IsNot),
							MatchRegex:    types.StringValue(bodyItem.MatchRegex),
							NotMatchRegex: types.StringValue(bodyItem.NotMatchRegex),
						}
						successWhenItem.Body = append(successWhenItem.Body, bodyCondition)
					}
				}

				// update status code conditions
				if len(swItem.StatusCode) > 0 {
					successWhenItem.StatusCode = []statusCodeCondition{}
					for _, statusCodeItem := range swItem.StatusCode {
						statusCodeCondition := statusCodeCondition{
							Contains:      types.StringValue(statusCodeItem.Contains),
							NotContains:   types.StringValue(statusCodeItem.NotContains),
							Is:            types.StringValue(statusCodeItem.Is),
							IsNot:         types.StringValue(statusCodeItem.IsNot),
							MatchRegex:    types.StringValue(statusCodeItem.MatchRegex),
							NotMatchRegex: types.StringValue(statusCodeItem.NotMatchRegex),
						}
						successWhenItem.StatusCode = append(successWhenItem.StatusCode, statusCodeCondition)
					}
				}

				// update response time conditions
				if swItem.ResponseTime != nil {
					if responseTimeStr, ok := swItem.ResponseTime.(string); ok && responseTimeStr != "" {
						successWhenItem.ResponseTime = []responseTimeCondition{
							{
								Target: types.StringValue(responseTimeStr),
							},
						}
					}
				}

				// update header conditions
				if len(swItem.Header) > 0 {
					successWhenItem.Header = make(map[string][]types.String)
					for k, headerItems := range swItem.Header {
						stringValues := []types.String{}
						for _, headerItem := range headerItems {
							// Marshal the headerItem directly into a JSON string
							headerJSON, err := json.Marshal(headerItem)
							if err == nil {
								stringValues = append(stringValues, types.StringValue(string(headerJSON)))
							}
						}
						if len(stringValues) > 0 {
							successWhenItem.Header[k] = stringValues
						}
					}
				}

				state.Task.SuccessWhen = append(state.Task.SuccessWhen, successWhenItem)
			}
		} else {
			state.Task.SuccessWhen = []successWhenItem{}
		}

	} else {
		state.Task = nil
	}
}

func (r *syntheticsTestResource) updateTCPTaskState(ctx context.Context, state *syntheticsTestResourceModel, content *api.SyntheticsTest) {
	if content.Task != nil {
		state.Task.Host = types.StringValue(content.Task.Host)
		state.Task.Port = types.StringValue(content.Task.Port)
		state.Task.Timeout = types.StringValue(content.Task.Timeout)
		state.Task.EnableTraceroute = types.BoolValue(content.Task.EnableTraceroute)

		// update advance options
		if content.Task.AdvanceOptions != nil {
			if state.Task.AdvanceOptions == nil {
				state.Task.AdvanceOptions = &advanceOptions{}
			}

			// update request options
			if content.Task.AdvanceOptions.RequestOptions != nil {
				state.Task.AdvanceOptions.RequestOptions = &requestOptions{
					Timeout: types.StringValue(content.Task.AdvanceOptions.RequestOptions.Timeout),
				}
			} else {
				state.Task.AdvanceOptions.RequestOptions = nil
			}

			// update request timeout
			if content.Task.AdvanceOptions.RequestTimeout != "" {
				state.Task.AdvanceOptions.RequestTimeout = types.StringValue(content.Task.AdvanceOptions.RequestTimeout)
			} else {
				state.Task.AdvanceOptions.RequestTimeout = types.StringValue("")
			}
		} else {
			state.Task.AdvanceOptions = nil
		}

		// update success when
		if len(content.Task.SuccessWhen) > 0 {
			state.Task.SuccessWhen = []successWhenItem{}
			for _, swItem := range content.Task.SuccessWhen {
				successWhenItem := successWhenItem{}

				// update response time conditions
				if swItem.ResponseTime != nil {
					successWhenItem.ResponseTime = []responseTimeCondition{}

					// Handle []interface{} type from JSON unmarshaling
					if responseTimeSlice, ok := swItem.ResponseTime.([]interface{}); ok && len(responseTimeSlice) > 0 {
						for _, rtItem := range responseTimeSlice {
							if rtMap, ok := rtItem.(map[string]interface{}); ok {
								rtCondition := responseTimeCondition{}

								// Extract target
								if target, ok := rtMap["target"].(string); ok {
									rtCondition.Target = types.StringValue(target)
								}

								// Extract is_contain_dns
								if isContainDNS, ok := rtMap["is_contain_dns"].(bool); ok {
									rtCondition.IsContainDNS = types.BoolValue(isContainDNS)
								}

								// Extract func
								if fn, ok := rtMap["func"].(string); ok && fn != "" {
									rtCondition.Func = types.StringValue(fn)
								}

								// Extract op
								if op, ok := rtMap["op"].(string); ok && op != "" {
									rtCondition.Op = types.StringValue(op)
								}

								successWhenItem.ResponseTime = append(successWhenItem.ResponseTime, rtCondition)
							}
						}
					}
				}

				// update response message conditions
				if len(swItem.ResponseMessage) > 0 {
					successWhenItem.ResponseMessage = []responseMessageCondition{}
					for _, responseMessageItem := range swItem.ResponseMessage {
						responseMessageCondition := responseMessageCondition{
							Contains:      types.StringValue(responseMessageItem.Contains),
							NotContains:   types.StringValue(responseMessageItem.NotContains),
							Is:            types.StringValue(responseMessageItem.Is),
							IsNot:         types.StringValue(responseMessageItem.IsNot),
							MatchRegex:    types.StringValue(responseMessageItem.MatchRegex),
							NotMatchRegex: types.StringValue(responseMessageItem.NotMatchRegex),
						}
						successWhenItem.ResponseMessage = append(successWhenItem.ResponseMessage, responseMessageCondition)
					}
				}

				// update header conditions
				if len(swItem.Header) > 0 {
					successWhenItem.Header = make(map[string][]types.String)
					for k, headerItems := range swItem.Header {
						stringValues := []types.String{}
						for _, headerItem := range headerItems {
							// Marshal the headerItem directly into a JSON string
							headerJSON, err := json.Marshal(headerItem)
							if err == nil {
								stringValues = append(stringValues, types.StringValue(string(headerJSON)))
							}
						}
						if len(stringValues) > 0 {
							successWhenItem.Header[k] = stringValues
						}
					}
				}

				// update hops conditions
				if len(swItem.Hops) > 0 {
					successWhenItem.Hops = []hopsCondition{}
					for _, hopsItem := range swItem.Hops {
						hopsCondition := hopsCondition{
							Op:     types.StringValue(hopsItem.Op),
							Target: types.Float64Value(hopsItem.Target),
						}
						successWhenItem.Hops = append(successWhenItem.Hops, hopsCondition)
					}
				}

				state.Task.SuccessWhen = append(state.Task.SuccessWhen, successWhenItem)
			}
		} else {
			state.Task.SuccessWhen = []successWhenItem{}
		}

	} else {
		state.Task = nil
	}
}

func (r *syntheticsTestResource) updateICMPTaskState(ctx context.Context, state *syntheticsTestResourceModel, content *api.SyntheticsTest) {
	if content.Task != nil {
		state.Task.Host = types.StringValue(content.Task.Host)
		state.Task.Timeout = types.StringValue(content.Task.Timeout)
		state.Task.EnableTraceroute = types.BoolValue(content.Task.EnableTraceroute)
		state.Task.PacketCount = types.Int64Value(int64(content.Task.PacketCount))

		// update advance options
		if content.Task.AdvanceOptions != nil {
			if state.Task.AdvanceOptions == nil {
				state.Task.AdvanceOptions = &advanceOptions{}
			}

			// update request timeout
			if content.Task.AdvanceOptions.RequestTimeout != "" {
				state.Task.AdvanceOptions.RequestTimeout = types.StringValue(content.Task.AdvanceOptions.RequestTimeout)
			} else {
				state.Task.AdvanceOptions.RequestTimeout = types.StringValue("")
			}
		} else {
			state.Task.AdvanceOptions = nil
		}

		// update success when
		if len(content.Task.SuccessWhen) > 0 {
			state.Task.SuccessWhen = []successWhenItem{}
			for _, swItem := range content.Task.SuccessWhen {
				successWhenItem := successWhenItem{}

				// update response time conditions
				if swItem.ResponseTime != nil {
					successWhenItem.ResponseTime = []responseTimeCondition{}

					// Handle []interface{} type from JSON unmarshaling
					if responseTimeSlice, ok := swItem.ResponseTime.([]interface{}); ok && len(responseTimeSlice) > 0 {
						for _, rtItem := range responseTimeSlice {
							if rtMap, ok := rtItem.(map[string]interface{}); ok {
								rtCondition := responseTimeCondition{}

								// Extract target
								if target, ok := rtMap["target"].(string); ok {
									rtCondition.Target = types.StringValue(target)
								}

								// Extract is_contain_dns
								if isContainDNS, ok := rtMap["is_contain_dns"].(bool); ok {
									rtCondition.IsContainDNS = types.BoolValue(isContainDNS)
								}

								// Extract func
								if fn, ok := rtMap["func"].(string); ok && fn != "" {
									rtCondition.Func = types.StringValue(fn)
								}

								// Extract op
								if op, ok := rtMap["op"].(string); ok && op != "" {
									rtCondition.Op = types.StringValue(op)
								}

								successWhenItem.ResponseTime = append(successWhenItem.ResponseTime, rtCondition)
							}
						}
					}
				}

				// update hops conditions
				if len(swItem.Hops) > 0 {
					successWhenItem.Hops = []hopsCondition{}
					for _, hopsItem := range swItem.Hops {
						hopsCondition := hopsCondition{
							Op:     types.StringValue(hopsItem.Op),
							Target: types.Float64Value(hopsItem.Target),
						}
						successWhenItem.Hops = append(successWhenItem.Hops, hopsCondition)
					}
				}

				// update packet loss percent conditions
				if len(swItem.PacketLossPercent) > 0 {
					successWhenItem.PacketLossPercent = []packetLossCondition{}
					for _, packetLossItem := range swItem.PacketLossPercent {
						packetLossCondition := packetLossCondition{
							Op:     types.StringValue(packetLossItem.Op),
							Target: types.Float64Value(packetLossItem.Target),
						}
						successWhenItem.PacketLossPercent = append(successWhenItem.PacketLossPercent, packetLossCondition)
					}
				}

				// update packets conditions
				if len(swItem.Packets) > 0 {
					successWhenItem.Packets = []packetsCondition{}
					for _, packetsItem := range swItem.Packets {
						packetsCondition := packetsCondition{
							Op:     types.StringValue(packetsItem.Op),
							Target: types.Float64Value(packetsItem.Target),
						}
						successWhenItem.Packets = append(successWhenItem.Packets, packetsCondition)
					}
				}

				state.Task.SuccessWhen = append(state.Task.SuccessWhen, successWhenItem)
			}
		} else {
			state.Task.SuccessWhen = []successWhenItem{}
		}

	} else {
		state.Task = nil
	}
}

func (r *syntheticsTestResource) updateWebSocketTaskState(ctx context.Context, state *syntheticsTestResourceModel, content *api.SyntheticsTest) {
	if content.Task != nil {

		state.Task.URL = types.StringValue(content.Task.URL)
		state.Task.Timeout = types.StringValue(content.Task.Timeout)
		state.Task.EnableTraceroute = types.BoolValue(content.Task.EnableTraceroute)
		if content.Task.Message != "" {
			state.Task.Message = types.StringValue(content.Task.Message)
		} else {
			state.Task.Message = types.StringValue("")
		}

		// update advance options
		if content.Task.AdvanceOptions != nil {
			if state.Task.AdvanceOptions == nil {
				state.Task.AdvanceOptions = &advanceOptions{}
			}

			// update request options
			if content.Task.AdvanceOptions.RequestOptions != nil {
				state.Task.AdvanceOptions.RequestOptions = &requestOptions{
					Timeout:        types.StringValue(content.Task.AdvanceOptions.RequestOptions.Timeout),
					Cookies:        types.StringValue(""),
					FollowRedirect: types.BoolValue(false),
					Headers:        make(map[string]types.String),
				}
				for k, v := range content.Task.AdvanceOptions.RequestOptions.Headers {
					state.Task.AdvanceOptions.RequestOptions.Headers[k] = types.StringValue(v)
				}
			} else {
				state.Task.AdvanceOptions.RequestOptions = nil
			}

			// update request timeout
			if content.Task.AdvanceOptions.RequestTimeout != "" {
				state.Task.AdvanceOptions.RequestTimeout = types.StringValue(content.Task.AdvanceOptions.RequestTimeout)
			} else {
				state.Task.AdvanceOptions.RequestTimeout = types.StringValue("")
			}
		} else {
			state.Task.AdvanceOptions = nil
		}

		// update success when
		if len(content.Task.SuccessWhen) > 0 {
			state.Task.SuccessWhen = []successWhenItem{}
			for _, swItem := range content.Task.SuccessWhen {
				successWhenItem := successWhenItem{}

				// update response time conditions
				if swItem.ResponseTime != nil {
					successWhenItem.ResponseTime = []responseTimeCondition{}

					// Handle []interface{} type from JSON unmarshaling
					if responseTimeSlice, ok := swItem.ResponseTime.([]interface{}); ok && len(responseTimeSlice) > 0 {
						for _, rtItem := range responseTimeSlice {
							if rtMap, ok := rtItem.(map[string]interface{}); ok {
								rtCondition := responseTimeCondition{}

								// Extract target
								if target, ok := rtMap["target"].(string); ok {
									rtCondition.Target = types.StringValue(target)
								}

								// Extract is_contain_dns
								if isContainDNS, ok := rtMap["is_contain_dns"].(bool); ok {
									rtCondition.IsContainDNS = types.BoolValue(isContainDNS)
								}

								// Extract func
								if fn, ok := rtMap["func"].(string); ok && fn != "" {
									rtCondition.Func = types.StringValue(fn)
								}

								// Extract op
								if op, ok := rtMap["op"].(string); ok && op != "" {
									rtCondition.Op = types.StringValue(op)
								}

								successWhenItem.ResponseTime = append(successWhenItem.ResponseTime, rtCondition)
							}
						}
					}
				}

				// update response message conditions
				if len(swItem.ResponseMessage) > 0 {
					successWhenItem.ResponseMessage = []responseMessageCondition{}
					for _, responseMessageItem := range swItem.ResponseMessage {
						responseMessageCondition := responseMessageCondition{
							Contains:      types.StringValue(responseMessageItem.Contains),
							NotContains:   types.StringValue(responseMessageItem.NotContains),
							Is:            types.StringValue(responseMessageItem.Is),
							IsNot:         types.StringValue(responseMessageItem.IsNot),
							MatchRegex:    types.StringValue(responseMessageItem.MatchRegex),
							NotMatchRegex: types.StringValue(responseMessageItem.NotMatchRegex),
						}
						successWhenItem.ResponseMessage = append(successWhenItem.ResponseMessage, responseMessageCondition)
					}
				}

				// update header conditions
				if len(swItem.Header) > 0 {
					successWhenItem.Header = make(map[string][]types.String)
					for k, headerItems := range swItem.Header {
						stringValues := []types.String{}
						for _, headerItem := range headerItems {
							// Marshal the headerItem directly into a JSON string
							headerJSON, err := json.Marshal(headerItem)
							if err == nil {
								stringValues = append(stringValues, types.StringValue(string(headerJSON)))
							}
						}
						if len(stringValues) > 0 {
							successWhenItem.Header[k] = stringValues
						}
					}
				}

				// update hops conditions
				if len(swItem.Hops) > 0 {
					successWhenItem.Hops = []hopsCondition{}
					for _, hopsItem := range swItem.Hops {
						hopsCondition := hopsCondition{
							Op:     types.StringValue(hopsItem.Op),
							Target: types.Float64Value(hopsItem.Target),
						}
						successWhenItem.Hops = append(successWhenItem.Hops, hopsCondition)
					}
				}

				state.Task.SuccessWhen = append(state.Task.SuccessWhen, successWhenItem)
			}
		} else {
			state.Task.SuccessWhen = []successWhenItem{}
		}

	} else {
		state.Task = nil
	}
}

// updateMultiTaskState updates the state for multi-step tests
func (r *syntheticsTestResource) updateMultiTaskState(ctx context.Context, state *syntheticsTestResourceModel, content *api.SyntheticsTest) {
	if content.Task != nil {
		// Update steps
		if len(content.Task.Steps) > 0 {
			stepObjects := []stepConfig{}
			for _, step := range content.Task.Steps {
				// Create a stepConfig object
				stepObj := stepConfig{
					Type: types.StringValue(step.Type),
				}

				// Add task only for http steps
				if step.Type == "http" {
					stepObj.Task = types.StringValue("") // Initialize task field
					stepObj.AllowFailure = types.BoolValue(step.AllowFailure)
					if step.Task != nil {
						// 处理不同类型的 Task 值
						switch v := step.Task.(type) {
						case string:
							// 尝试解析字符串为 JSON 对象，然后重新序列化为标准 JSON 字符串
							var taskObj interface{}
							if err := json.Unmarshal([]byte(v), &taskObj); err == nil {
								// 如果解析成功，重新序列化为标准 JSON 字符串
								taskBytes, err := json.Marshal(taskObj)
								if err == nil {
									stepObj.Task = types.StringValue(string(taskBytes))
								} else {
									// 解析失败，使用原始字符串
									stepObj.Task = types.StringValue(v)
								}
							} else {
								// 解析失败，使用原始字符串
								stepObj.Task = types.StringValue(v)
							}
						case map[string]interface{}:
							// 如果是 map，直接序列化为 JSON 字符串
							taskBytes, err := json.Marshal(v)
							if err == nil {
								stepObj.Task = types.StringValue(string(taskBytes))
							}
						default:
							// 其他类型，尝试转换为字符串
							stepObj.Task = types.StringValue(fmt.Sprintf("%v", v))
						}
					}
					// Add extracted_vars if present
					if step.ExtractedVars != nil {
						extractedVarsObjects := []extractedVar{}
						for _, extractedVarItem := range step.ExtractedVars {
							extractedVarObj := extractedVar{
								Name:   types.StringValue(extractedVarItem.Name),
								Field:  types.StringValue(extractedVarItem.Field),
								Secure: types.BoolValue(extractedVarItem.Secure),
							}
							extractedVarsObjects = append(extractedVarsObjects, extractedVarObj)
						}
						stepObj.ExtractedVars = extractedVarsObjects
					}

				}

				if step.Type != "http" {
					stepObj.Value = types.Int64Value(int64(step.Value))
				}

				// Add retry if present
				if step.Retry != nil {
					stepObj.Retry = &retryConfig{
						Retry:    types.Int64Value(int64(step.Retry.Retry)),
						Interval: types.Int64Value(int64(step.Retry.Interval)),
					}
				}

				stepObjects = append(stepObjects, stepObj)
			}
			state.Task.Steps = stepObjects
		}
	} else {
		state.Task = nil
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *syntheticsTestResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan syntheticsTestResourceModel
	var err error
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get client from provider
	if r.client == nil {
		resp.Diagnostics.AddError("Client not configured", "Client not configured for synthetics test resource")
		return
	}

	// Convert plan to API request
	item := r.getSyntheticsTestFromPlan(ctx, &plan)

	// omit fields when update
	item.Type = ""

	// Update synthetics test
	content := &api.SyntheticsTest{}

	var updateErr error
	if plan.Type.ValueString() == "multi" {
		updateErr = r.client.MultiUpdate(consts.TypeNameSyntheticsTest, plan.UUID.ValueString(), item, content)
	} else {
		updateErr = r.client.Update(consts.TypeNameSyntheticsTest, plan.UUID.ValueString(), item, content)
	}

	if updateErr != nil {
		resp.Diagnostics.AddError(
			"Error updating synthetics test",
			"Could not update synthetics test, unexpected error: "+updateErr.Error(),
		)
		return
	}

	// Retrieve values from state
	var state syntheticsTestResourceModel
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Update plan with response data
	plan.UpdateAt = types.Int64Value(content.UpdateAt)

	// Check if status needs to be updated
	// Compare plan and state to determine if status should be changed
	planStatus := "ok"
	if plan.Task != nil && !plan.Task.Status.IsNull() && !plan.Task.Status.IsUnknown() {
		planStatus = plan.Task.Status.ValueString()
	}

	stateStatus := "ok"
	if state.Task != nil && !state.Task.Status.IsNull() && !state.Task.Status.IsUnknown() {
		stateStatus = state.Task.Status.ValueString()
	}

	// Only update status if it has changed
	if planStatus != stateStatus {
		isDisable := planStatus == "stop"
		err = r.client.SetStatus(consts.TypeNameSyntheticsTest, []string{plan.UUID.ValueString()}, isDisable)
		if err != nil {
			resp.Diagnostics.AddError("Error setting status", fmt.Sprintf("Failed to set status: %s", err))
			return
		}
	}

	// Set state
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *syntheticsTestResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state syntheticsTestResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get client from provider
	if r.client == nil {
		resp.Diagnostics.AddError("Client not configured", "Client not configured for synthetics test resource")
		return
	}

	// Delete synthetics test using POST method
	body := map[string][]string{
		"taskUUIDs": {state.UUID.ValueString()},
	}
	err := r.client.DeleteByPost(consts.TypeNameSyntheticsTest, "", body, nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting synthetics test",
			"Could not delete synthetics test, unexpected error: "+err.Error(),
		)
		return
	}
}

// ImportState imports a resource into the Terraform state.
func (r *syntheticsTestResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Set the UUID from the import ID
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("uuid"), req.ID)...)
}

// getSyntheticsTestFromPlan converts the plan to an API request
func (r *syntheticsTestResource) getSyntheticsTestFromPlan(ctx context.Context, plan *syntheticsTestResourceModel) *api.SyntheticsTest {
	testType := plan.Type.ValueString()
	test := &api.SyntheticsTest{}

	// Only set type for non-multi tests
	if testType != "multi" {
		test.Type = testType
	}

	// Add regions
	if len(plan.Regions) > 0 {
		regions := make([]string, len(plan.Regions))
		for i, region := range plan.Regions {
			regions[i] = region.ValueString()
		}
		test.Regions = regions
	}

	// Add tags
	if len(plan.Tags) > 0 {
		tags := make([]string, len(plan.Tags))
		for i, tag := range plan.Tags {
			tags[i] = tag.ValueString()
		}
		test.Tags = tags
	}

	if plan.Task == nil {
		return test
	}

	switch testType {
	case TaskTypeHTTP:
		test.Task = r.getHTTPTaskConfigFromPlanTask(ctx, plan.Task)
	case TaskTypeTCP:
		test.Task = r.getTCPTaskConfigFromPlanTask(ctx, plan.Task)
	case TaskTypeWebSocket:
		test.Task = r.getWebSocketTaskConfigFromPlanTask(ctx, plan.Task)
	case TaskTypeICMP:
		test.Task = r.getICMPTaskConfigFromPlanTask(ctx, plan.Task)
	case TaskTypeGRPC:
		test.Task = r.getGRPCTaskConfigFromPlanTask(ctx, plan.Task)
	case "multi":
		// For multi-step tests, create a basic task config
		test.Task = r.getMultiTaskConfigFromPlanTask(ctx, plan.Task)
	default:
		err := fmt.Errorf("unsupported test type: %s", testType)
		tflog.Error(ctx, err.Error())
		return nil
	}

	return test
}

func getBasicTask(ctx context.Context, task *taskConfig) *api.TaskConfig {
	if task == nil {
		return nil
	}
	apiTask := &api.TaskConfig{
		Name:             task.Name.ValueString(),
		Frequency:        task.Frequency.ValueString(),
		ScheduleType:     task.ScheduleType.ValueString(),
		Crontab:          task.Crontab.ValueString(),
		SuccessWhenLogic: task.SuccessWhenLogic.ValueString(),
	}

	return apiTask
}

func (r *syntheticsTestResource) getMultiTaskConfigFromPlanTask(ctx context.Context, task *taskConfig) *api.TaskConfig {
	if task == nil {
		return nil
	}
	apiTask := getBasicTask(ctx, task)
	if apiTask == nil {
		return nil
	}
	apiTask.Desc = task.Desc.ValueString()

	// Add steps for multi-step tests
	if len(task.Steps) > 0 {
		apiTask.Steps = []api.StepConfig{}
		for _, step := range task.Steps {
			stepConfig := api.StepConfig{}

			// Extract type
			if !step.Type.IsNull() {
				stepConfig.Type = step.Type.ValueString()
			}

			// Extract task
			if stepConfig.Type == "http" {
				if !step.Task.IsNull() {
					// Use the JSON string directly as raw message
					stepConfig.Task = json.RawMessage(step.Task.ValueString())
				}
				// Extract allow_failure
				allowFailureValue := false // Default to false if not provided
				if !step.AllowFailure.IsNull() {
					allowFailureValue = step.AllowFailure.ValueBool()
				}
				stepConfig.AllowFailure = allowFailureValue
			}

			// Extract value
			if stepConfig.Type == "wait" && !step.Value.IsNull() {
				stepConfig.Value = int(step.Value.ValueInt64())
			}

			// Extract retry
			if step.Retry != nil {
				stepConfig.Retry = &api.RetryConfig{}
				if !step.Retry.Retry.IsNull() {
					stepConfig.Retry.Retry = int(step.Retry.Retry.ValueInt64())
				}
				if !step.Retry.Interval.IsNull() {
					stepConfig.Retry.Interval = int(step.Retry.Interval.ValueInt64())
				}
			}

			// Extract extracted_vars
			if step.ExtractedVars != nil {
				stepConfig.ExtractedVars = []api.ExtractedVar{}
				for _, extractedVar := range step.ExtractedVars {
					extractedVarConfig := api.ExtractedVar{}

					// Extract name
					if !extractedVar.Name.IsNull() {
						extractedVarConfig.Name = extractedVar.Name.ValueString()
					} else {
						// Name is required, set a default if not provided
						extractedVarConfig.Name = ""
					}

					// Extract field
					if !extractedVar.Field.IsNull() {
						extractedVarConfig.Field = extractedVar.Field.ValueString()
					} else {
						// Field is required, set a default if not provided
						extractedVarConfig.Field = ""
					}

					// Extract secure
					secureValue := false // Default to false if not provided
					if !extractedVar.Secure.IsNull() {
						secureValue = extractedVar.Secure.ValueBool()
					}
					extractedVarConfig.Secure = secureValue

					stepConfig.ExtractedVars = append(stepConfig.ExtractedVars, extractedVarConfig)
				}
			}

			apiTask.Steps = append(apiTask.Steps, stepConfig)
		}
	}

	return apiTask
}

func (r *syntheticsTestResource) getHTTPTaskConfigFromPlanTask(ctx context.Context, task *taskConfig) *api.TaskConfig {
	if task == nil {
		return nil
	}
	apiTask := getBasicTask(ctx, task)
	if apiTask == nil {
		return nil
	}
	apiTask.URL = task.URL.ValueString()
	apiTask.Method = task.Method.ValueString()
	apiTask.PostMode = task.PostMode.ValueString()
	apiTask.PostScript = task.PostScript.ValueString()

	// Add advanced options
	if task.AdvanceOptions != nil {
		apiTask.AdvanceOptions = &api.AdvanceOptions{}

		// Add request timeout
		if !task.AdvanceOptions.RequestTimeout.IsNull() {
			apiTask.AdvanceOptions.RequestTimeout = task.AdvanceOptions.RequestTimeout.ValueString()
		}

		// Add request options
		if task.AdvanceOptions.RequestOptions != nil {
			requestOptionsConfig := &api.RequestOptions{}

			if !task.AdvanceOptions.RequestOptions.FollowRedirect.IsNull() {
				requestOptionsConfig.FollowRedirect = task.AdvanceOptions.RequestOptions.FollowRedirect.ValueBool()
			}

			if !task.AdvanceOptions.RequestOptions.Cookies.IsNull() {
				requestOptionsConfig.Cookies = task.AdvanceOptions.RequestOptions.Cookies.ValueString()
			}

			if !task.AdvanceOptions.RequestOptions.Timeout.IsNull() {
				requestOptionsConfig.Timeout = task.AdvanceOptions.RequestOptions.Timeout.ValueString()
			}

			// Add headers
			if len(task.AdvanceOptions.RequestOptions.Headers) > 0 {
				headersMap := make(map[string]string)
				for k, v := range task.AdvanceOptions.RequestOptions.Headers {
					if !v.IsNull() {
						headersMap[k] = v.ValueString()
					}
				}
				requestOptionsConfig.Headers = headersMap
			}

			// Add metadata
			if len(task.AdvanceOptions.RequestOptions.Metadata) > 0 {
				metadataMap := make(map[string]string)
				for k, v := range task.AdvanceOptions.RequestOptions.Metadata {
					if !v.IsNull() {
						metadataMap[k] = v.ValueString()
					}
				}
				requestOptionsConfig.Metadata = metadataMap
			}

			// Add auth
			if task.AdvanceOptions.RequestOptions.Auth != nil {
				authConfig := &api.Auth{}

				if !task.AdvanceOptions.RequestOptions.Auth.Username.IsNull() {
					authConfig.Username = task.AdvanceOptions.RequestOptions.Auth.Username.ValueString()
				}

				if !task.AdvanceOptions.RequestOptions.Auth.Password.IsNull() {
					authConfig.Password = task.AdvanceOptions.RequestOptions.Auth.Password.ValueString()
				}

				requestOptionsConfig.Auth = authConfig
			}

			apiTask.AdvanceOptions.RequestOptions = requestOptionsConfig
		}

		// Add request body
		if task.AdvanceOptions.RequestBody != nil {
			requestBodyConfig := &api.RequestBody{}

			if !task.AdvanceOptions.RequestBody.BodyType.IsNull() {
				requestBodyConfig.BodyType = task.AdvanceOptions.RequestBody.BodyType.ValueString()
			}

			if !task.AdvanceOptions.RequestBody.Body.IsNull() {
				requestBodyConfig.Body = task.AdvanceOptions.RequestBody.Body.ValueString()
			}

			// Add form data
			if len(task.AdvanceOptions.RequestBody.Form) > 0 {
				formMap := make(map[string]string)
				for k, v := range task.AdvanceOptions.RequestBody.Form {
					if !v.IsNull() {
						formMap[k] = v.ValueString()
					}
				}
				requestBodyConfig.Form = formMap
			}

			apiTask.AdvanceOptions.RequestBody = requestBodyConfig
		}

		// Add certificate
		if task.AdvanceOptions.Certificate != nil {
			certificateConfig := &api.Certificate{}

			if !task.AdvanceOptions.Certificate.IgnoreServerCertificateError.IsNull() {
				certificateConfig.IgnoreServerCertificateError = task.AdvanceOptions.Certificate.IgnoreServerCertificateError.ValueBool()
			}

			if !task.AdvanceOptions.Certificate.PrivateKey.IsNull() {
				certificateConfig.PrivateKey = task.AdvanceOptions.Certificate.PrivateKey.ValueString()
			}

			if !task.AdvanceOptions.Certificate.Certificate.IsNull() {
				certificateConfig.Certificate = task.AdvanceOptions.Certificate.Certificate.ValueString()
			}

			apiTask.AdvanceOptions.Certificate = certificateConfig
		}

		// Add proxy
		if task.AdvanceOptions.Proxy != nil {
			proxyConfig := &api.Proxy{}

			if !task.AdvanceOptions.Proxy.URL.IsNull() {
				proxyConfig.URL = task.AdvanceOptions.Proxy.URL.ValueString()
			}

			// Add headers
			if len(task.AdvanceOptions.Proxy.Headers) > 0 {
				headersMap := make(map[string]string)
				for k, v := range task.AdvanceOptions.Proxy.Headers {
					if !v.IsNull() {
						headersMap[k] = v.ValueString()
					}
				}
				proxyConfig.Headers = headersMap
			}

			apiTask.AdvanceOptions.Proxy = proxyConfig
		}

		// Add secret
		if task.AdvanceOptions.Secret != nil {
			secretConfig := &api.Secret{}

			if !task.AdvanceOptions.Secret.NotSave.IsNull() {
				secretConfig.NotSave = task.AdvanceOptions.Secret.NotSave.ValueBool()
			}

			apiTask.AdvanceOptions.Secret = secretConfig
		}
	}

	// Add success when conditions
	if len(task.SuccessWhen) > 0 {
		apiTask.SuccessWhen = []api.SuccessWhenItem{}
		for _, swItem := range task.SuccessWhen {
			swConfig := api.SuccessWhenItem{}

			// Add response time conditions
			if len(swItem.ResponseTime) > 0 {
				swConfig.ResponseTime = swItem.ResponseTime[0].Target.ValueString()
			}

			// Add body conditions
			if len(swItem.Body) > 0 {
				swConfig.Body = []api.BodyCondition{}
				for _, bodyItem := range swItem.Body {
					bodyCondition := api.BodyCondition{}

					if !bodyItem.Contains.IsNull() {
						bodyCondition.Contains = bodyItem.Contains.ValueString()
					}

					if !bodyItem.NotContains.IsNull() {
						bodyCondition.NotContains = bodyItem.NotContains.ValueString()
					}

					if !bodyItem.Is.IsNull() {
						bodyCondition.Is = bodyItem.Is.ValueString()
					}

					if !bodyItem.IsNot.IsNull() {
						bodyCondition.IsNot = bodyItem.IsNot.ValueString()
					}

					if !bodyItem.MatchRegex.IsNull() {
						bodyCondition.MatchRegex = bodyItem.MatchRegex.ValueString()
					}

					if !bodyItem.NotMatchRegex.IsNull() {
						bodyCondition.NotMatchRegex = bodyItem.NotMatchRegex.ValueString()
					}

					swConfig.Body = append(swConfig.Body, bodyCondition)
				}
			}

			// Add status code conditions
			if len(swItem.StatusCode) > 0 {
				swConfig.StatusCode = []api.StatusCodeCondition{}
				for _, statusCodeItem := range swItem.StatusCode {
					statusCodeCondition := api.StatusCodeCondition{}

					if !statusCodeItem.Is.IsNull() {
						statusCodeCondition.Is = statusCodeItem.Is.ValueString()
					}

					if !statusCodeItem.IsNot.IsNull() {
						statusCodeCondition.IsNot = statusCodeItem.IsNot.ValueString()
					}

					if !statusCodeItem.MatchRegex.IsNull() {
						statusCodeCondition.MatchRegex = statusCodeItem.MatchRegex.ValueString()
					}

					if !statusCodeItem.NotMatchRegex.IsNull() {
						statusCodeCondition.NotMatchRegex = statusCodeItem.NotMatchRegex.ValueString()
					}

					if !statusCodeItem.Contains.IsNull() {
						statusCodeCondition.Contains = statusCodeItem.Contains.ValueString()
					}

					if !statusCodeItem.NotContains.IsNull() {
						statusCodeCondition.NotContains = statusCodeItem.NotContains.ValueString()
					}

					swConfig.StatusCode = append(swConfig.StatusCode, statusCodeCondition)
				}
			}

			// Add header conditions
			if len(swItem.Header) > 0 {
				swConfig.Header = make(map[string][]api.HeaderCondition)
				for k, headerItems := range swItem.Header {
					if len(headerItems) > 0 {
						headerConditions := []api.HeaderCondition{}
						for _, headerItem := range headerItems {
							if !headerItem.IsNull() {
								headerJSON := headerItem.ValueString()
								var hc api.HeaderCondition
								if err := json.Unmarshal([]byte(headerJSON), &hc); err == nil {
									headerConditions = append(headerConditions, hc)
								}
							}
						}
						if len(headerConditions) > 0 {
							swConfig.Header[k] = headerConditions
						}
					}
				}
			}

			// Add response message conditions
			if len(swItem.ResponseMessage) > 0 {
				swConfig.ResponseMessage = []api.ResponseMessageCondition{}
				for _, responseMessageItem := range swItem.ResponseMessage {
					responseMessageCondition := api.ResponseMessageCondition{}

					if !responseMessageItem.Contains.IsNull() {
						responseMessageCondition.Contains = responseMessageItem.Contains.ValueString()
					}

					if !responseMessageItem.NotContains.IsNull() {
						responseMessageCondition.NotContains = responseMessageItem.NotContains.ValueString()
					}

					if !responseMessageItem.Is.IsNull() {
						responseMessageCondition.Is = responseMessageItem.Is.ValueString()
					}

					if !responseMessageItem.IsNot.IsNull() {
						responseMessageCondition.IsNot = responseMessageItem.IsNot.ValueString()
					}

					if !responseMessageItem.MatchRegex.IsNull() {
						responseMessageCondition.MatchRegex = responseMessageItem.MatchRegex.ValueString()
					}

					if !responseMessageItem.NotMatchRegex.IsNull() {
						responseMessageCondition.NotMatchRegex = responseMessageItem.NotMatchRegex.ValueString()
					}

					swConfig.ResponseMessage = append(swConfig.ResponseMessage, responseMessageCondition)
				}
			}

			apiTask.SuccessWhen = append(apiTask.SuccessWhen, swConfig)
		}
	}

	return apiTask
}

func (r *syntheticsTestResource) getTCPTaskConfigFromPlanTask(ctx context.Context, task *taskConfig) *api.TaskConfig {
	if task == nil {
		return nil
	}
	apiTask := getBasicTask(ctx, task)
	if apiTask == nil {
		return nil
	}
	apiTask.Host = task.Host.ValueString()
	apiTask.Port = task.Port.ValueString()
	apiTask.Timeout = task.Timeout.ValueString()
	if !task.EnableTraceroute.IsNull() {
		apiTask.EnableTraceroute = task.EnableTraceroute.ValueBool()
	}

	// Add advanced options
	if task.AdvanceOptions != nil {
		apiTask.AdvanceOptions = &api.AdvanceOptions{}

		// Add request timeout
		if !task.AdvanceOptions.RequestTimeout.IsNull() {
			apiTask.AdvanceOptions.RequestTimeout = task.AdvanceOptions.RequestTimeout.ValueString()
		}

		// Add request options
		if task.AdvanceOptions.RequestOptions != nil {
			requestOptionsConfig := &api.RequestOptions{}

			if !task.AdvanceOptions.RequestOptions.Timeout.IsNull() {
				requestOptionsConfig.Timeout = task.AdvanceOptions.RequestOptions.Timeout.ValueString()
			}

			// Add headers
			if len(task.AdvanceOptions.RequestOptions.Headers) > 0 {
				headersMap := make(map[string]string)
				for k, v := range task.AdvanceOptions.RequestOptions.Headers {
					if !v.IsNull() {
						headersMap[k] = v.ValueString()
					}
				}
				requestOptionsConfig.Headers = headersMap
			}

			apiTask.AdvanceOptions.RequestOptions = requestOptionsConfig
		}
	}

	// Add success when conditions
	if len(task.SuccessWhen) > 0 {
		apiTask.SuccessWhen = []api.SuccessWhenItem{}
		for _, swItem := range task.SuccessWhen {
			swConfig := api.SuccessWhenItem{}

			// Add response time conditions
			if len(swItem.ResponseTime) > 0 {
				responseTimeConditions := []api.ResponseTimeCondition{}
				for _, rtItem := range swItem.ResponseTime {
					rtCondition := api.ResponseTimeCondition{
						Target:       rtItem.Target.ValueString(),
						IsContainDNS: rtItem.IsContainDNS.ValueBool(),
					}

					// Only set Func if it's not null
					if !rtItem.Func.IsNull() {
						rtCondition.Func = rtItem.Func.ValueString()
					}

					// Only set Op if it's not null
					if !rtItem.Op.IsNull() {
						rtCondition.Op = rtItem.Op.ValueString()
					}

					responseTimeConditions = append(responseTimeConditions, rtCondition)
				}
				swConfig.ResponseTime = responseTimeConditions
			}

			// Add response message conditions
			if len(swItem.ResponseMessage) > 0 {
				swConfig.ResponseMessage = []api.ResponseMessageCondition{}
				for _, responseMessageItem := range swItem.ResponseMessage {
					responseMessageCondition := api.ResponseMessageCondition{}

					if !responseMessageItem.Contains.IsNull() {
						responseMessageCondition.Contains = responseMessageItem.Contains.ValueString()
					}

					if !responseMessageItem.NotContains.IsNull() {
						responseMessageCondition.NotContains = responseMessageItem.NotContains.ValueString()
					}

					if !responseMessageItem.Is.IsNull() {
						responseMessageCondition.Is = responseMessageItem.Is.ValueString()
					}

					if !responseMessageItem.IsNot.IsNull() {
						responseMessageCondition.IsNot = responseMessageItem.IsNot.ValueString()
					}

					if !responseMessageItem.MatchRegex.IsNull() {
						responseMessageCondition.MatchRegex = responseMessageItem.MatchRegex.ValueString()
					}

					if !responseMessageItem.NotMatchRegex.IsNull() {
						responseMessageCondition.NotMatchRegex = responseMessageItem.NotMatchRegex.ValueString()
					}

					swConfig.ResponseMessage = append(swConfig.ResponseMessage, responseMessageCondition)
				}
			}

			// Add header conditions
			if len(swItem.Header) > 0 {
				swConfig.Header = make(map[string][]api.HeaderCondition)
				for k, headerItems := range swItem.Header {
					if len(headerItems) > 0 {
						headerConditions := []api.HeaderCondition{}
						for _, headerItem := range headerItems {
							if !headerItem.IsNull() {
								headerJSON := headerItem.ValueString()
								var hc api.HeaderCondition
								if err := json.Unmarshal([]byte(headerJSON), &hc); err == nil {
									headerConditions = append(headerConditions, hc)
								}
							}
						}
						if len(headerConditions) > 0 {
							swConfig.Header[k] = headerConditions
						}
					}
				}
			}

			// Add hops conditions
			if len(swItem.Hops) > 0 {
				swConfig.Hops = []api.HopsCondition{}
				for _, hopsItem := range swItem.Hops {
					hopsCondition := api.HopsCondition{
						Op:     hopsItem.Op.ValueString(),
						Target: hopsItem.Target.ValueFloat64(),
					}
					swConfig.Hops = append(swConfig.Hops, hopsCondition)
				}
			}

			apiTask.SuccessWhen = append(apiTask.SuccessWhen, swConfig)
		}
	}

	return apiTask
}

func (r *syntheticsTestResource) getICMPTaskConfigFromPlanTask(ctx context.Context, task *taskConfig) *api.TaskConfig {
	if task == nil {
		return nil
	}
	apiTask := getBasicTask(ctx, task)
	if apiTask == nil {
		return nil
	}
	apiTask.Host = task.Host.ValueString()
	apiTask.Timeout = task.Timeout.ValueString()
	if !task.EnableTraceroute.IsNull() {
		apiTask.EnableTraceroute = task.EnableTraceroute.ValueBool()
	}
	if !task.PacketCount.IsNull() {
		apiTask.PacketCount = int(task.PacketCount.ValueInt64())
	}

	// Add advanced options
	if task.AdvanceOptions != nil {
		apiTask.AdvanceOptions = &api.AdvanceOptions{}

		// Add request timeout
		if !task.AdvanceOptions.RequestTimeout.IsNull() {
			apiTask.AdvanceOptions.RequestTimeout = task.AdvanceOptions.RequestTimeout.ValueString()
		}
	}

	// Add success when conditions
	if len(task.SuccessWhen) > 0 {
		apiTask.SuccessWhen = []api.SuccessWhenItem{}
		for _, swItem := range task.SuccessWhen {
			swConfig := api.SuccessWhenItem{}

			// Add response time conditions
			if len(swItem.ResponseTime) > 0 {
				responseTimeConditions := []api.ResponseTimeCondition{}
				for _, rtItem := range swItem.ResponseTime {
					rtCondition := api.ResponseTimeCondition{
						Target:       rtItem.Target.ValueString(),
						IsContainDNS: rtItem.IsContainDNS.ValueBool(),
					}

					// Only set Func if it's not null
					if !rtItem.Func.IsNull() {
						rtCondition.Func = rtItem.Func.ValueString()
					}

					// Only set Op if it's not null
					if !rtItem.Op.IsNull() {
						rtCondition.Op = rtItem.Op.ValueString()
					}

					responseTimeConditions = append(responseTimeConditions, rtCondition)
				}
				swConfig.ResponseTime = responseTimeConditions
			}

			// Add hops conditions
			if len(swItem.Hops) > 0 {
				swConfig.Hops = []api.HopsCondition{}
				for _, hopsItem := range swItem.Hops {
					hopsCondition := api.HopsCondition{
						Op:     hopsItem.Op.ValueString(),
						Target: hopsItem.Target.ValueFloat64(),
					}
					swConfig.Hops = append(swConfig.Hops, hopsCondition)
				}
			}

			// Add packet loss percent conditions
			if len(swItem.PacketLossPercent) > 0 {
				swConfig.PacketLossPercent = []api.PacketLossCondition{}
				for _, packetLossItem := range swItem.PacketLossPercent {
					packetLossCondition := api.PacketLossCondition{
						Op:     packetLossItem.Op.ValueString(),
						Target: packetLossItem.Target.ValueFloat64(),
					}
					swConfig.PacketLossPercent = append(swConfig.PacketLossPercent, packetLossCondition)
				}
			}

			// Add packets conditions
			if len(swItem.Packets) > 0 {
				swConfig.Packets = []api.PacketsCondition{}
				for _, packetsItem := range swItem.Packets {
					packetsCondition := api.PacketsCondition{
						Op:     packetsItem.Op.ValueString(),
						Target: packetsItem.Target.ValueFloat64(),
					}
					swConfig.Packets = append(swConfig.Packets, packetsCondition)
				}
			}

			apiTask.SuccessWhen = append(apiTask.SuccessWhen, swConfig)
		}
	}

	return apiTask
}

func (r *syntheticsTestResource) getWebSocketTaskConfigFromPlanTask(ctx context.Context, task *taskConfig) *api.TaskConfig {
	if task == nil {
		return nil
	}
	apiTask := getBasicTask(ctx, task)
	if apiTask == nil {
		return nil
	}
	apiTask.URL = task.URL.ValueString()
	apiTask.Timeout = task.Timeout.ValueString()
	if !task.EnableTraceroute.IsNull() {
		apiTask.EnableTraceroute = task.EnableTraceroute.ValueBool()
	}
	if !task.Message.IsNull() {
		apiTask.Message = task.Message.ValueString()
	}

	// Add advanced options
	if task.AdvanceOptions != nil {
		apiTask.AdvanceOptions = &api.AdvanceOptions{}

		// Add request timeout
		if !task.AdvanceOptions.RequestTimeout.IsNull() {
			apiTask.AdvanceOptions.RequestTimeout = task.AdvanceOptions.RequestTimeout.ValueString()
		}

		// Add request options
		if task.AdvanceOptions.RequestOptions != nil {
			requestOptionsConfig := &api.RequestOptions{}

			if !task.AdvanceOptions.RequestOptions.Timeout.IsNull() {
				requestOptionsConfig.Timeout = task.AdvanceOptions.RequestOptions.Timeout.ValueString()
			}

			// Add follow redirect
			if !task.AdvanceOptions.RequestOptions.FollowRedirect.IsNull() {
				requestOptionsConfig.FollowRedirect = task.AdvanceOptions.RequestOptions.FollowRedirect.ValueBool()
			}

			// Add cookies
			if !task.AdvanceOptions.RequestOptions.Cookies.IsNull() {
				requestOptionsConfig.Cookies = task.AdvanceOptions.RequestOptions.Cookies.ValueString()
			}

			// Add headers
			if len(task.AdvanceOptions.RequestOptions.Headers) > 0 {
				headersMap := make(map[string]string)
				for k, v := range task.AdvanceOptions.RequestOptions.Headers {
					if !v.IsNull() {
						headersMap[k] = v.ValueString()
					}
				}
				requestOptionsConfig.Headers = headersMap
			}

			apiTask.AdvanceOptions.RequestOptions = requestOptionsConfig
		}
	}

	// Add success when conditions
	if len(task.SuccessWhen) > 0 {
		apiTask.SuccessWhen = []api.SuccessWhenItem{}
		for _, swItem := range task.SuccessWhen {
			swConfig := api.SuccessWhenItem{}

			// Add response time conditions
			if len(swItem.ResponseTime) > 0 {
				responseTimeConditions := []api.ResponseTimeCondition{}
				for _, rtItem := range swItem.ResponseTime {
					rtCondition := api.ResponseTimeCondition{
						Target:       rtItem.Target.ValueString(),
						IsContainDNS: rtItem.IsContainDNS.ValueBool(),
					}

					// Only set Func if it's not null
					if !rtItem.Func.IsNull() {
						rtCondition.Func = rtItem.Func.ValueString()
					}

					// Only set Op if it's not null
					if !rtItem.Op.IsNull() {
						rtCondition.Op = rtItem.Op.ValueString()
					}

					responseTimeConditions = append(responseTimeConditions, rtCondition)
				}
				swConfig.ResponseTime = responseTimeConditions
			}

			// Add response message conditions
			if len(swItem.ResponseMessage) > 0 {
				swConfig.ResponseMessage = []api.ResponseMessageCondition{}
				for _, responseMessageItem := range swItem.ResponseMessage {
					responseMessageCondition := api.ResponseMessageCondition{}

					if !responseMessageItem.Contains.IsNull() {
						responseMessageCondition.Contains = responseMessageItem.Contains.ValueString()
					}

					if !responseMessageItem.NotContains.IsNull() {
						responseMessageCondition.NotContains = responseMessageItem.NotContains.ValueString()
					}

					if !responseMessageItem.Is.IsNull() {
						responseMessageCondition.Is = responseMessageItem.Is.ValueString()
					}

					if !responseMessageItem.IsNot.IsNull() {
						responseMessageCondition.IsNot = responseMessageItem.IsNot.ValueString()
					}

					if !responseMessageItem.MatchRegex.IsNull() {
						responseMessageCondition.MatchRegex = responseMessageItem.MatchRegex.ValueString()
					}

					if !responseMessageItem.NotMatchRegex.IsNull() {
						responseMessageCondition.NotMatchRegex = responseMessageItem.NotMatchRegex.ValueString()
					}

					swConfig.ResponseMessage = append(swConfig.ResponseMessage, responseMessageCondition)
				}
			}

			// Add header conditions
			if len(swItem.Header) > 0 {
				swConfig.Header = make(map[string][]api.HeaderCondition)
				for k, headerItems := range swItem.Header {
					if len(headerItems) > 0 {
						headerConditions := []api.HeaderCondition{}
						for _, headerItem := range headerItems {
							if !headerItem.IsNull() {
								headerJSON := headerItem.ValueString()
								var hc api.HeaderCondition
								if err := json.Unmarshal([]byte(headerJSON), &hc); err == nil {
									headerConditions = append(headerConditions, hc)
								}
							}
						}
						if len(headerConditions) > 0 {
							swConfig.Header[k] = headerConditions
						}
					}
				}
			}

			// Add hops conditions
			if len(swItem.Hops) > 0 {
				swConfig.Hops = []api.HopsCondition{}
				for _, hopsItem := range swItem.Hops {
					hopsCondition := api.HopsCondition{
						Op:     hopsItem.Op.ValueString(),
						Target: hopsItem.Target.ValueFloat64(),
					}
					swConfig.Hops = append(swConfig.Hops, hopsCondition)
				}
			}

			apiTask.SuccessWhen = append(apiTask.SuccessWhen, swConfig)
		}
	}

	return apiTask
}

func (r *syntheticsTestResource) getGRPCTaskConfigFromPlanTask(ctx context.Context, task *taskConfig) *api.TaskConfig {
	if task == nil {
		return nil
	}
	apiTask := getBasicTask(ctx, task)
	if apiTask == nil {
		return nil
	}
	apiTask.Server = task.Server.ValueString()

	// Add advanced options
	if task.AdvanceOptions != nil {
		apiTask.AdvanceOptions = &api.AdvanceOptions{}

		// Add request timeout
		if !task.AdvanceOptions.RequestTimeout.IsNull() {
			apiTask.AdvanceOptions.RequestTimeout = task.AdvanceOptions.RequestTimeout.ValueString()
		}

		// Add request options
		if task.AdvanceOptions.RequestOptions != nil {
			requestOptionsConfig := &api.RequestOptions{}

			// Add proto files
			if task.AdvanceOptions.RequestOptions.ProtoFiles != nil {
				protoFilesConfig := &api.ProtoFiles{}

				if !task.AdvanceOptions.RequestOptions.ProtoFiles.FullMethod.IsNull() {
					protoFilesConfig.FullMethod = task.AdvanceOptions.RequestOptions.ProtoFiles.FullMethod.ValueString()
				}

				if !task.AdvanceOptions.RequestOptions.ProtoFiles.Request.IsNull() {
					protoFilesConfig.Request = task.AdvanceOptions.RequestOptions.ProtoFiles.Request.ValueString()
				}

				if len(task.AdvanceOptions.RequestOptions.ProtoFiles.ProtoFiles) > 0 {
					protoFilesMap := make(map[string]string)
					for k, v := range task.AdvanceOptions.RequestOptions.ProtoFiles.ProtoFiles {
						if !v.IsNull() {
							protoFilesMap[k] = v.ValueString()
						}
					}
					protoFilesConfig.ProtoFiles = protoFilesMap
				}

				requestOptionsConfig.ProtoFiles = protoFilesConfig
			}

			// Add reflection
			if task.AdvanceOptions.RequestOptions.Reflection != nil {
				reflectionConfig := &api.Reflection{}

				if !task.AdvanceOptions.RequestOptions.Reflection.FullMethod.IsNull() {
					reflectionConfig.FullMethod = task.AdvanceOptions.RequestOptions.Reflection.FullMethod.ValueString()
				}

				if !task.AdvanceOptions.RequestOptions.Reflection.Request.IsNull() {
					reflectionConfig.Request = task.AdvanceOptions.RequestOptions.Reflection.Request.ValueString()
				}

				requestOptionsConfig.Reflection = reflectionConfig
			}

			// Add health check
			if task.AdvanceOptions.RequestOptions.HealthCheck != nil {
				healthCheckConfig := &api.HealthCheck{}

				if !task.AdvanceOptions.RequestOptions.HealthCheck.Service.IsNull() {
					healthCheckConfig.Service = task.AdvanceOptions.RequestOptions.HealthCheck.Service.ValueString()
				}

				requestOptionsConfig.HealthCheck = healthCheckConfig
			}

			apiTask.AdvanceOptions.RequestOptions = requestOptionsConfig
		}

		// Add certificate
		if task.AdvanceOptions.Certificate != nil {
			certificateConfig := &api.Certificate{}

			if !task.AdvanceOptions.Certificate.IgnoreServerCertificateError.IsNull() {
				certificateConfig.IgnoreServerCertificateError = task.AdvanceOptions.Certificate.IgnoreServerCertificateError.ValueBool()
			}

			apiTask.AdvanceOptions.Certificate = certificateConfig
		}

		// Add post script
		if !task.AdvanceOptions.PostScript.IsNull() {
			apiTask.AdvanceOptions.PostScript = task.AdvanceOptions.PostScript.ValueString()
		}
	}

	// Add success when conditions
	if len(task.SuccessWhen) > 0 {
		apiTask.SuccessWhen = []api.SuccessWhenItem{}
		for _, swItem := range task.SuccessWhen {
			swConfig := api.SuccessWhenItem{}

			// Add body conditions
			if len(swItem.Body) > 0 {
				swConfig.Body = []api.BodyCondition{}
				for _, bodyItem := range swItem.Body {
					bodyCondition := api.BodyCondition{}

					if !bodyItem.Contains.IsNull() {
						bodyCondition.Contains = bodyItem.Contains.ValueString()
					}

					if !bodyItem.NotContains.IsNull() {
						bodyCondition.NotContains = bodyItem.NotContains.ValueString()
					}

					if !bodyItem.Is.IsNull() {
						bodyCondition.Is = bodyItem.Is.ValueString()
					}

					if !bodyItem.IsNot.IsNull() {
						bodyCondition.IsNot = bodyItem.IsNot.ValueString()
					}

					if !bodyItem.MatchRegex.IsNull() {
						bodyCondition.MatchRegex = bodyItem.MatchRegex.ValueString()
					}

					if !bodyItem.NotMatchRegex.IsNull() {
						bodyCondition.NotMatchRegex = bodyItem.NotMatchRegex.ValueString()
					}

					swConfig.Body = append(swConfig.Body, bodyCondition)
				}
			}

			// Add response time conditions (like HTTP)
			if len(swItem.ResponseTime) > 0 {
				responseTimeStr := swItem.ResponseTime[0].Target.ValueString()
				if responseTimeStr != "" {
					swConfig.ResponseTime = responseTimeStr
				}
			}

			apiTask.SuccessWhen = append(apiTask.SuccessWhen, swConfig)
		}
	}

	return apiTask
}

func (r *syntheticsTestResource) updateGRPCTaskState(ctx context.Context, state *syntheticsTestResourceModel, content *api.SyntheticsTest) {
	if content.Task != nil {
		state.Task.Server = types.StringValue(content.Task.Server)

		// update advance options
		if content.Task.AdvanceOptions != nil {
			if state.Task.AdvanceOptions == nil {
				state.Task.AdvanceOptions = &advanceOptions{}
			}

			// update request timeout
			if content.Task.AdvanceOptions.RequestTimeout != "" {
				state.Task.AdvanceOptions.RequestTimeout = types.StringValue(content.Task.AdvanceOptions.RequestTimeout)
			} else {
				state.Task.AdvanceOptions.RequestTimeout = types.StringValue("")
			}

			// update post script
			if content.Task.AdvanceOptions.PostScript != "" {
				state.Task.AdvanceOptions.PostScript = types.StringValue(content.Task.AdvanceOptions.PostScript)
			} else {
				state.Task.AdvanceOptions.PostScript = types.StringValue("")
			}

			// update certificate
			if content.Task.AdvanceOptions.Certificate != nil {
				state.Task.AdvanceOptions.Certificate = &certificate{
					IgnoreServerCertificateError: types.BoolValue(content.Task.AdvanceOptions.Certificate.IgnoreServerCertificateError),
				}
			} else {
				state.Task.AdvanceOptions.Certificate = nil
			}

			// update request options
			if content.Task.AdvanceOptions.RequestOptions != nil {
				if state.Task.AdvanceOptions.RequestOptions == nil {
					state.Task.AdvanceOptions.RequestOptions = &requestOptions{}
				}

				// update timeout
				if content.Task.AdvanceOptions.RequestOptions.Timeout != "" {
					state.Task.AdvanceOptions.RequestOptions.Timeout = types.StringValue(content.Task.AdvanceOptions.RequestOptions.Timeout)
				} else {
					state.Task.AdvanceOptions.RequestOptions.Timeout = types.StringValue("")
				}

				// update metadata
				if content.Task.AdvanceOptions.RequestOptions.Metadata != nil {
					state.Task.AdvanceOptions.RequestOptions.Metadata = map[string]types.String{}
					for k, v := range content.Task.AdvanceOptions.RequestOptions.Metadata {
						state.Task.AdvanceOptions.RequestOptions.Metadata[k] = types.StringValue(v)
					}
				}

				// update proto files
				if content.Task.AdvanceOptions.RequestOptions.ProtoFiles != nil {
					if state.Task.AdvanceOptions.RequestOptions.ProtoFiles == nil {
						state.Task.AdvanceOptions.RequestOptions.ProtoFiles = &protoFiles{}
					}
					state.Task.AdvanceOptions.RequestOptions.ProtoFiles.FullMethod = types.StringValue(content.Task.AdvanceOptions.RequestOptions.ProtoFiles.FullMethod)
					state.Task.AdvanceOptions.RequestOptions.ProtoFiles.Request = types.StringValue(content.Task.AdvanceOptions.RequestOptions.ProtoFiles.Request)

					if content.Task.AdvanceOptions.RequestOptions.ProtoFiles.ProtoFiles != nil {
						state.Task.AdvanceOptions.RequestOptions.ProtoFiles.ProtoFiles = map[string]types.String{}
						for k, v := range content.Task.AdvanceOptions.RequestOptions.ProtoFiles.ProtoFiles {
							state.Task.AdvanceOptions.RequestOptions.ProtoFiles.ProtoFiles[k] = types.StringValue(v)
						}
					}
				} else {
					state.Task.AdvanceOptions.RequestOptions.ProtoFiles = nil
				}

				// update reflection
				if content.Task.AdvanceOptions.RequestOptions.Reflection != nil {
					if state.Task.AdvanceOptions.RequestOptions.Reflection == nil {
						state.Task.AdvanceOptions.RequestOptions.Reflection = &reflection{}
					}
					state.Task.AdvanceOptions.RequestOptions.Reflection.FullMethod = types.StringValue(content.Task.AdvanceOptions.RequestOptions.Reflection.FullMethod)
					state.Task.AdvanceOptions.RequestOptions.Reflection.Request = types.StringValue(content.Task.AdvanceOptions.RequestOptions.Reflection.Request)
				} else {
					state.Task.AdvanceOptions.RequestOptions.Reflection = nil
				}

				// update health check
				if content.Task.AdvanceOptions.RequestOptions.HealthCheck != nil {
					if state.Task.AdvanceOptions.RequestOptions.HealthCheck == nil {
						state.Task.AdvanceOptions.RequestOptions.HealthCheck = &healthCheck{}
					}
					state.Task.AdvanceOptions.RequestOptions.HealthCheck.Service = types.StringValue(content.Task.AdvanceOptions.RequestOptions.HealthCheck.Service)
				} else {
					state.Task.AdvanceOptions.RequestOptions.HealthCheck = nil
				}
			} else {
				state.Task.AdvanceOptions.RequestOptions = nil
			}
		} else {
			state.Task.AdvanceOptions = nil
		}

		// update success when
		if len(content.Task.SuccessWhen) > 0 {
			state.Task.SuccessWhen = []successWhenItem{}
			for _, swItem := range content.Task.SuccessWhen {
				successWhenItem := successWhenItem{}

				// update body conditions
				if len(swItem.Body) > 0 {
					successWhenItem.Body = []bodyCondition{}
					for _, bodyItem := range swItem.Body {
						bodyCondition := bodyCondition{}

						if bodyItem.Contains != "" {
							bodyCondition.Contains = types.StringValue(bodyItem.Contains)
						}

						if bodyItem.NotContains != "" {
							bodyCondition.NotContains = types.StringValue(bodyItem.NotContains)
						}

						if bodyItem.Is != "" {
							bodyCondition.Is = types.StringValue(bodyItem.Is)
						}

						if bodyItem.IsNot != "" {
							bodyCondition.IsNot = types.StringValue(bodyItem.IsNot)
						}

						if bodyItem.MatchRegex != "" {
							bodyCondition.MatchRegex = types.StringValue(bodyItem.MatchRegex)
						}

						if bodyItem.NotMatchRegex != "" {
							bodyCondition.NotMatchRegex = types.StringValue(bodyItem.NotMatchRegex)
						}

						successWhenItem.Body = append(successWhenItem.Body, bodyCondition)
					}
				}

				// update response time conditions (like HTTP)
				if swItem.ResponseTime != nil {
					if responseTimeStr, ok := swItem.ResponseTime.(string); ok && responseTimeStr != "" {
						successWhenItem.ResponseTime = []responseTimeCondition{}
						rtCondition := responseTimeCondition{
							Target: types.StringValue(responseTimeStr),
						}
						successWhenItem.ResponseTime = append(successWhenItem.ResponseTime, rtCondition)
					}
				}

				state.Task.SuccessWhen = append(state.Task.SuccessWhen, successWhenItem)
			}
		} else {
			state.Task.SuccessWhen = []successWhenItem{}
		}

	} else {
		state.Task = nil
	}
}

// stepConfigAttrTypes returns the attribute types for step configuration
func stepConfigAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"type":           types.StringType,
		"name":           types.StringType,
		"task":           types.StringType,
		"allow_failure":  types.BoolType,
		"extracted_vars": types.ListType{ElemType: types.ObjectType{AttrTypes: extractedVarAttrTypes()}},
		"value":          types.Int64Type,
	}
}

// extractedVarAttrTypes returns the attribute types for extracted variables
func extractedVarAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"name":   types.StringType,
		"field":  types.StringType,
		"secure": types.BoolType,
	}
}

// responseTimeConditionAttrTypes returns the attribute types for response time conditions
func responseTimeConditionAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"is_contain_dns": types.BoolType,
		"target":         types.StringType,
		"func":           types.StringType,
		"op":             types.StringType,
	}
}

// successWhenItemAttrTypes returns the attribute types for success when conditions
func successWhenItemAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"response_time":       types.ListType{ElemType: types.ObjectType{AttrTypes: responseTimeConditionAttrTypes()}},
		"body":                types.ListType{ElemType: types.ObjectType{AttrTypes: bodyConditionAttrTypes()}},
		"header":              types.MapType{ElemType: types.ObjectType{AttrTypes: headerConditionAttrTypes()}},
		"status_code":         types.ListType{ElemType: types.ObjectType{AttrTypes: statusCodeConditionAttrTypes()}},
		"response_message":    types.ListType{ElemType: types.ObjectType{AttrTypes: responseMessageConditionAttrTypes()}},
		"hops":                types.ListType{ElemType: types.ObjectType{AttrTypes: hopsConditionAttrTypes()}},
		"packet_loss_percent": types.ListType{ElemType: types.ObjectType{AttrTypes: packetLossConditionAttrTypes()}},
		"packets":             types.ListType{ElemType: types.ObjectType{AttrTypes: packetsConditionAttrTypes()}},
	}
}

// bodyConditionAttrTypes returns the attribute types for body conditions
func bodyConditionAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"contains":        types.StringType,
		"not_contains":    types.StringType,
		"is":              types.StringType,
		"is_not":          types.StringType,
		"match_regex":     types.StringType,
		"not_match_regex": types.StringType,
	}
}

// headerConditionAttrTypes returns the attribute types for header conditions
func headerConditionAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"contains":        types.StringType,
		"not_contains":    types.StringType,
		"is":              types.StringType,
		"is_not":          types.StringType,
		"match_regex":     types.StringType,
		"not_match_regex": types.StringType,
	}
}

// statusCodeConditionAttrTypes returns the attribute types for status code conditions
func statusCodeConditionAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"is":              types.StringType,
		"is_not":          types.StringType,
		"match_regex":     types.StringType,
		"not_match_regex": types.StringType,
		"contains":        types.StringType,
		"not_contains":    types.StringType,
	}
}

// responseMessageConditionAttrTypes returns the attribute types for response message conditions
func responseMessageConditionAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"contains":        types.StringType,
		"not_contains":    types.StringType,
		"is":              types.StringType,
		"is_not":          types.StringType,
		"match_regex":     types.StringType,
		"not_match_regex": types.StringType,
	}
}

// hopsConditionAttrTypes returns the attribute types for hops conditions
func hopsConditionAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"op":     types.StringType,
		"target": types.Float64Type,
	}
}

// packetLossConditionAttrTypes returns the attribute types for packet loss conditions
func packetLossConditionAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"op":     types.StringType,
		"target": types.Float64Type,
	}
}

// packetsConditionAttrTypes returns the attribute types for packets conditions
func packetsConditionAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"op":     types.StringType,
		"target": types.Float64Type,
	}
}

// authAttrTypes returns the attribute types for authentication
func authAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"username": types.StringType,
		"password": types.StringType,
	}
}

// requestBodyAttrTypes returns the attribute types for request body
func requestBodyAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"body_type": types.StringType,
		"body":      types.StringType,
		"form":      types.MapType{ElemType: types.StringType},
	}
}

// certificateAttrTypes returns the attribute types for certificate
func certificateAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"ignore_server_certificate_error": types.BoolType,
		"private_key":                     types.StringType,
		"certificate":                     types.StringType,
	}
}

// proxyAttrTypes returns the attribute types for proxy
func proxyAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"url":     types.StringType,
		"headers": types.MapType{ElemType: types.StringType},
	}
}

// secretAttrTypes returns the attribute types for secret
func secretAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"not_save": types.BoolType,
	}
}

// protoFilesAttrTypes returns the attribute types for proto files
func protoFilesAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"protofiles":  types.MapType{ElemType: types.StringType},
		"full_method": types.StringType,
		"request":     types.StringType,
	}
}

// reflectionAttrTypes returns the attribute types for reflection
func reflectionAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"full_method": types.StringType,
		"request":     types.StringType,
	}
}

// healthCheckAttrTypes returns the attribute types for health check
func healthCheckAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"service": types.StringType,
	}
}

// requestOptionsAttrTypes returns the attribute types for request options
func requestOptionsAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"follow_redirect": types.BoolType,
		"headers":         types.MapType{ElemType: types.StringType},
		"cookies":         types.StringType,
		"auth":            types.ObjectType{AttrTypes: authAttrTypes()},
		"timeout":         types.StringType,
		"metadata":        types.MapType{ElemType: types.StringType},
		"proto_files":     types.ObjectType{AttrTypes: protoFilesAttrTypes()},
		"reflection":      types.ObjectType{AttrTypes: reflectionAttrTypes()},
		"health_check":    types.ObjectType{AttrTypes: healthCheckAttrTypes()},
	}
}

// advanceOptionsAttrTypes returns the attribute types for advance options
func advanceOptionsAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"request_options": types.ObjectType{AttrTypes: requestOptionsAttrTypes()},
		"request_body":    types.ObjectType{AttrTypes: requestBodyAttrTypes()},
		"certificate":     types.ObjectType{AttrTypes: certificateAttrTypes()},
		"proxy":           types.ObjectType{AttrTypes: proxyAttrTypes()},
		"secret":          types.ObjectType{AttrTypes: secretAttrTypes()},
		"request_timeout": types.StringType,
	}
}

// requestOptionsHeadlessAttrTypes returns the attribute types for headless request options
func requestOptionsHeadlessAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"follow_redirect": types.BoolType,
		"headers":         types.MapType{ElemType: types.StringType},
		"cookies":         types.StringType,
		"auth":            types.ObjectType{AttrTypes: authAttrTypes()},
		"timeout":         types.StringType,
	}
}

// requestBodyHeadlessAttrTypes returns the attribute types for headless request body
func requestBodyHeadlessAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"body_type": types.StringType,
		"body":      types.StringType,
	}
}

// advanceOptionsHeadlessAttrTypes returns the attribute types for headless advance options
func advanceOptionsHeadlessAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"request_options": types.ObjectType{AttrTypes: requestOptionsHeadlessAttrTypes()},
		"request_body":    types.ObjectType{AttrTypes: requestBodyHeadlessAttrTypes()},
		"request_timeout": types.StringType,
	}
}

func getTags(tags []api.TagInfo) []string {
	var tagNames []string
	for _, tag := range tags {
		tagNames = append(tagNames, tag.Name)
	}

	return tagNames
}

func isSameTags(stateTags []types.String, tagNames []string) bool {
	if len(stateTags) != len(tagNames) {
		return false
	}

	tags := []string{}
	for _, tag := range stateTags {
		tags = append(tags, tag.ValueString())
	}

	sort.Strings(tags)
	sort.Strings(tagNames)

	for i := range tags {
		if tags[i] != tagNames[i] {
			return false
		}
	}

	return true
}
