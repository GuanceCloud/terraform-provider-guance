package synthetics_test

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const (
	TaskTypeHTTP      = "http"
	TaskTypeTCP       = "tcp"
	TaskTypeWebSocket = "websocket"
	TaskTypeICMP      = "icmp"
	TaskTypeGRPC      = "grpc"
)

// syntheticsTestResourceModel maps the resource schema data
type syntheticsTestResourceModel struct {
	UUID          types.String   `tfsdk:"uuid"`
	Type          types.String   `tfsdk:"type"`
	Regions       []types.String `tfsdk:"regions"`
	Task          *taskConfig    `tfsdk:"task"`
	Tags          []types.String `tfsdk:"tags"`
	CreateAt      types.Int64    `tfsdk:"create_at"`
	UpdateAt      types.Int64    `tfsdk:"update_at"`
	WorkspaceUUID types.String   `tfsdk:"workspace_uuid"`
}

// taskConfig represents the task configuration
type taskConfig struct {
	URL                    types.String      `tfsdk:"url"`
	Method                 types.String      `tfsdk:"method"`
	Name                   types.String      `tfsdk:"name"`
	Status                 types.String      `tfsdk:"status"`
	Frequency              types.String      `tfsdk:"frequency"`
	ScheduleType           types.String      `tfsdk:"schedule_type"`
	Crontab                types.String      `tfsdk:"crontab"`
	Desc                   types.String      `tfsdk:"desc"`
	Steps                  []stepConfig      `tfsdk:"steps"`
	AdvanceOptions         *advanceOptions   `tfsdk:"advance_options"`
	AdvanceOptionsHeadless types.Object      `tfsdk:"advance_options_headless"`
	SuccessWhenLogic       types.String      `tfsdk:"success_when_logic"`
	SuccessWhen            []successWhenItem `tfsdk:"success_when"`
	EnableTraceroute       types.Bool        `tfsdk:"enable_traceroute"`
	PacketCount            types.Int64       `tfsdk:"packet_count"`
	Server                 types.String      `tfsdk:"server"`
	Host                   types.String      `tfsdk:"host"`
	Port                   types.String      `tfsdk:"port"`
	Timeout                types.String      `tfsdk:"timeout"`
	Message                types.String      `tfsdk:"message"`
	PostMode               types.String      `tfsdk:"post_mode"`
	PostScript             types.String      `tfsdk:"post_script"`
}

// stepConfig represents a step in a multi-step test
type stepConfig struct {
	Type          types.String   `tfsdk:"type"`
	Task          types.String   `tfsdk:"task"`
	AllowFailure  types.Bool     `tfsdk:"allow_failure"`
	ExtractedVars []extractedVar `tfsdk:"extracted_vars"`
	Value         types.Int64    `tfsdk:"value"`
	Retry         *retryConfig   `tfsdk:"retry"`
}

// retryConfig represents the retry configuration for a step
type retryConfig struct {
	Retry    types.Int64 `tfsdk:"retry"`
	Interval types.Int64 `tfsdk:"interval"`
}

// extractedVar represents a variable extracted from a step
type extractedVar struct {
	Name   types.String `tfsdk:"name"`
	Field  types.String `tfsdk:"field"`
	Secure types.Bool   `tfsdk:"secure"`
}

// advanceOptions represents the advanced options
type advanceOptions struct {
	RequestOptions *requestOptions `tfsdk:"request_options"`
	RequestBody    *requestBody    `tfsdk:"request_body"`
	Certificate    *certificate    `tfsdk:"certificate"`
	Proxy          *proxy          `tfsdk:"proxy"`
	Secret         *secret         `tfsdk:"secret"`
	RequestTimeout types.String    `tfsdk:"request_timeout"`
	PostScript     types.String    `tfsdk:"post_script"`
}

// requestOptions represents the request options
type requestOptions struct {
	FollowRedirect types.Bool              `tfsdk:"follow_redirect"`
	Headers        map[string]types.String `tfsdk:"headers"`
	Cookies        types.String            `tfsdk:"cookies"`
	Auth           *auth                   `tfsdk:"auth"`
	Timeout        types.String            `tfsdk:"timeout"`
	Metadata       map[string]types.String `tfsdk:"metadata"`
	ProtoFiles     *protoFiles             `tfsdk:"proto_files"`
	Reflection     *reflection             `tfsdk:"reflection"`
	HealthCheck    *healthCheck            `tfsdk:"health_check"`
}

// auth represents the authentication configuration
type auth struct {
	Username types.String `tfsdk:"username"`
	Password types.String `tfsdk:"password"`
}

// requestBody represents the request body configuration
type requestBody struct {
	BodyType types.String            `tfsdk:"body_type"`
	Body     types.String            `tfsdk:"body"`
	Form     map[string]types.String `tfsdk:"form"`
}

// protoFiles represents proto files for gRPC testing
type protoFiles struct {
	ProtoFiles map[string]types.String `tfsdk:"protofiles"`
	FullMethod types.String            `tfsdk:"full_method"`
	Request    types.String            `tfsdk:"request"`
}

// reflection represents gRPC reflection configuration
type reflection struct {
	FullMethod types.String `tfsdk:"full_method"`
	Request    types.String `tfsdk:"request"`
}

// healthCheck represents gRPC health check configuration
type healthCheck struct {
	Service types.String `tfsdk:"service"`
}

// certificate represents the certificate configuration
type certificate struct {
	IgnoreServerCertificateError types.Bool   `tfsdk:"ignore_server_certificate_error"`
	PrivateKey                   types.String `tfsdk:"private_key"`
	Certificate                  types.String `tfsdk:"certificate"`
}

// proxy represents the proxy configuration
type proxy struct {
	URL     types.String            `tfsdk:"url"`
	Headers map[string]types.String `tfsdk:"headers"`
}

// secret represents the secret configuration
type secret struct {
	NotSave types.Bool `tfsdk:"not_save"`
}

// responseTimeCondition represents a response time condition
type responseTimeCondition struct {
	IsContainDNS types.Bool   `tfsdk:"is_contain_dns"`
	Target       types.String `tfsdk:"target"`
	Func         types.String `tfsdk:"func"`
	Op           types.String `tfsdk:"op"`
}

// successWhenItem represents a single success condition
type successWhenItem struct {
	Body              []bodyCondition            `tfsdk:"body"`
	Header            map[string][]types.String  `tfsdk:"header"`
	ResponseTime      []responseTimeCondition    `tfsdk:"response_time"`
	StatusCode        []statusCodeCondition      `tfsdk:"status_code"`
	ResponseMessage   []responseMessageCondition `tfsdk:"response_message"`
	Hops              []hopsCondition            `tfsdk:"hops"`
	PacketLossPercent []packetLossCondition      `tfsdk:"packet_loss_percent"`
	Packets           []packetsCondition         `tfsdk:"packets"`
}

// responseMessageCondition represents a response message condition
type responseMessageCondition struct {
	Contains      types.String `tfsdk:"contains"`
	NotContains   types.String `tfsdk:"not_contains"`
	Is            types.String `tfsdk:"is"`
	IsNot         types.String `tfsdk:"is_not"`
	MatchRegex    types.String `tfsdk:"match_regex"`
	NotMatchRegex types.String `tfsdk:"not_match_regex"`
}

// hopsCondition represents a hops condition
type hopsCondition struct {
	Op     types.String  `tfsdk:"op"`
	Target types.Float64 `tfsdk:"target"`
}

// packetLossCondition represents a packet loss condition
type packetLossCondition struct {
	Op     types.String  `tfsdk:"op"`
	Target types.Float64 `tfsdk:"target"`
}

// packetsCondition represents a packets condition
type packetsCondition struct {
	Op     types.String  `tfsdk:"op"`
	Target types.Float64 `tfsdk:"target"`
}

// bodyCondition represents a body condition
type bodyCondition struct {
	Contains      types.String `tfsdk:"contains"`
	NotContains   types.String `tfsdk:"not_contains"`
	Is            types.String `tfsdk:"is"`
	IsNot         types.String `tfsdk:"is_not"`
	MatchRegex    types.String `tfsdk:"match_regex"`
	NotMatchRegex types.String `tfsdk:"not_match_regex"`
}

// headerCondition represents a header condition
type headerCondition struct {
	Contains      types.String `tfsdk:"contains"`
	NotContains   types.String `tfsdk:"not_contains"`
	Is            types.String `tfsdk:"is"`
	IsNot         types.String `tfsdk:"is_not"`
	MatchRegex    types.String `tfsdk:"match_regex"`
	NotMatchRegex types.String `tfsdk:"not_match_regex"`
}

// statusCodeCondition represents a status code condition
type statusCodeCondition struct {
	Is            types.String `tfsdk:"is"`
	IsNot         types.String `tfsdk:"is_not"`
	MatchRegex    types.String `tfsdk:"match_regex"`
	NotMatchRegex types.String `tfsdk:"not_match_regex"`
	Contains      types.String `tfsdk:"contains"`
	NotContains   types.String `tfsdk:"not_contains"`
}
