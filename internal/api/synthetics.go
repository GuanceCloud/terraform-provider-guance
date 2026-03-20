package api

import (
	"github.com/GuanceCloud/terraform-provider-guance/internal/consts"
)

// TagInfo represents a tag information

type TagInfo struct {
	ID          interface{} `json:"id,omitempty"`
	UUID        string      `json:"uuid,omitempty"`
	Name        string      `json:"name,omitempty"`
	Description string      `json:"description,omitempty"`
	Colour      string      `json:"colour,omitempty"`
}

// SyntheticsTest represents the synthetics test structure for API requests and responses
type SyntheticsTest struct {
	UUID          string      `json:"uuid,omitempty"`
	Type          string      `json:"type,omitempty"`
	Regions       []string    `json:"regions,omitempty"`
	Task          *TaskConfig `json:"task,omitempty"`
	Tags          []string    `json:"tags,omitempty"`
	TagInfo       []TagInfo   `json:"tagInfo,omitempty"`
	CreateAt      int64       `json:"createAt,omitempty"`
	UpdateAt      int64       `json:"updateAt,omitempty"`
	WorkspaceUUID string      `json:"workspaceUUID,omitempty"`
}

// TaskConfig represents the task configuration
type TaskConfig struct {
	URL                    string            `json:"url,omitempty"`
	Method                 string            `json:"method,omitempty"`
	Name                   string            `json:"name,omitempty"`
	Status                 string            `json:"status,omitempty"`
	Frequency              string            `json:"frequency,omitempty"`
	ScheduleType           string            `json:"schedule_type,omitempty"`
	Crontab                string            `json:"crontab,omitempty"`
	Desc                   string            `json:"desc,omitempty"`
	Steps                  []StepConfig      `json:"steps,omitempty"`
	AdvanceOptions         *AdvanceOptions   `json:"advance_options,omitempty"`
	AdvanceOptionsHeadless *AdvanceOptions   `json:"advance_options_headless,omitempty"`
	SuccessWhenLogic       string            `json:"success_when_logic,omitempty"`
	SuccessWhen            []SuccessWhenItem `json:"success_when,omitempty"`
	EnableTraceroute       bool              `json:"enable_traceroute,omitempty"`
	PacketCount            int               `json:"packet_count,omitempty"`
	Server                 string            `json:"server,omitempty"`
	Host                   string            `json:"host,omitempty"`
	Port                   string            `json:"port,omitempty"`
	Timeout                string            `json:"timeout,omitempty"`
	Message                string            `json:"message,omitempty"`
	PostMode               string            `json:"post_mode,omitempty"`
	PostScript             string            `json:"post_script,omitempty"`
	ExternalID             string            `json:"external_id,omitempty"`
	PostURL                string            `json:"post_url,omitempty"`
}

// AdvanceOptions represents the advanced options
type AdvanceOptions struct {
	RequestOptions *RequestOptions `json:"request_options,omitempty"`
	RequestBody    *RequestBody    `json:"request_body,omitempty"`
	Certificate    *Certificate    `json:"certificate,omitempty"`
	Proxy          *Proxy          `json:"proxy,omitempty"`
	Secret         *Secret         `json:"secret,omitempty"`
	Auth           *Auth           `json:"auth,omitempty"`
	RequestTimeout string          `json:"request_timeout,omitempty"`
	PostScript     string          `json:"post_script,omitempty"`
}

// RequestOptions represents the request options
type RequestOptions struct {
	FollowRedirect bool              `json:"follow_redirect,omitempty"`
	Headers        map[string]string `json:"headers,omitempty"`
	Cookies        string            `json:"cookies,omitempty"`
	Auth           *Auth             `json:"auth,omitempty"`
	Timeout        string            `json:"timeout,omitempty"`
	Metadata       map[string]string `json:"metadata,omitempty"`
	ProtoFiles     *ProtoFiles       `json:"proto_files,omitempty"`
	Reflection     *Reflection       `json:"reflection,omitempty"`
	HealthCheck    *HealthCheck      `json:"health_check,omitempty"`
}

// ProtoFiles represents proto files for gRPC testing
type ProtoFiles struct {
	ProtoFiles map[string]string `json:"protofiles,omitempty"`
	FullMethod string            `json:"full_method,omitempty"`
	Request    string            `json:"request,omitempty"`
}

// Reflection represents gRPC reflection configuration
type Reflection struct {
	FullMethod string `json:"full_method,omitempty"`
	Request    string `json:"request,omitempty"`
}

// HealthCheck represents gRPC health check configuration
type HealthCheck struct {
	Service string `json:"service,omitempty"`
}

// Auth represents the authentication configuration
type Auth struct {
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
}

// RequestBody represents the request body configuration
type RequestBody struct {
	BodyType string            `json:"body_type,omitempty"`
	Body     string            `json:"body,omitempty"`
	Form     map[string]string `json:"form,omitempty"`
}

// Certificate represents the certificate configuration
type Certificate struct {
	IgnoreServerCertificateError bool   `json:"ignore_server_certificate_error,omitempty"`
	PrivateKey                   string `json:"private_key,omitempty"`
	Certificate                  string `json:"certificate,omitempty"`
}

// Proxy represents the proxy configuration
type Proxy struct {
	URL     string            `json:"url,omitempty"`
	Headers map[string]string `json:"headers,omitempty"`
}

// Secret represents the secret configuration
type Secret struct {
	NotSave bool `json:"not_save,omitempty"`
}

// BodyCondition represents a body condition
type BodyCondition struct {
	Contains      string `json:"contains,omitempty"`
	NotContains   string `json:"not_contains,omitempty"`
	Is            string `json:"is,omitempty"`
	IsNot         string `json:"is_not,omitempty"`
	MatchRegex    string `json:"match_regex,omitempty"`
	NotMatchRegex string `json:"not_match_regex,omitempty"`
}

// HeaderCondition represents a header condition
type HeaderCondition struct {
	Contains      string `json:"contains,omitempty"`
	NotContains   string `json:"not_contains,omitempty"`
	Is            string `json:"is,omitempty"`
	IsNot         string `json:"is_not,omitempty"`
	MatchRegex    string `json:"match_regex,omitempty"`
	NotMatchRegex string `json:"not_match_regex,omitempty"`
}

// StatusCodeCondition represents a status code condition
type StatusCodeCondition struct {
	Is            string `json:"is,omitempty"`
	IsNot         string `json:"is_not,omitempty"`
	MatchRegex    string `json:"match_regex,omitempty"`
	NotMatchRegex string `json:"not_match_regex,omitempty"`
	Contains      string `json:"contains,omitempty"`
	NotContains   string `json:"not_contains,omitempty"`
}

// StepConfig represents a step in a multi-step test
type StepConfig struct {
	Type          string         `json:"type,omitempty"`
	Name          string         `json:"name,omitempty"`
	Task          any            `json:"task,omitempty"`
	AllowFailure  bool           `json:"allow_failure,omitempty"`
	ExtractedVars []ExtractedVar `json:"extracted_vars,omitempty"`
	Value         int            `json:"value,omitempty"`
	Retry         *RetryConfig   `json:"retry,omitempty"`
}

// ExtractedVar represents a variable extracted from a step
type ExtractedVar struct {
	Name   string `json:"name,omitempty"`
	Field  string `json:"field,omitempty"`
	Secure bool   `json:"secure"`
}

// RetryConfig represents the retry configuration for a step
type RetryConfig struct {
	Retry    int `json:"retry,omitempty"`
	Interval int `json:"interval,omitempty"`
}

// ResponseTimeCondition represents a response time condition
type ResponseTimeCondition struct {
	Target       string `json:"target,omitempty"`
	IsContainDNS bool   `json:"is_contain_dns,omitempty"`
	Func         string `json:"func,omitempty"`
	Op           string `json:"op,omitempty"`
}

// ResponseMessageCondition represents a response message condition
type ResponseMessageCondition struct {
	Contains      string `json:"contains,omitempty"`
	NotContains   string `json:"not_contains,omitempty"`
	Is            string `json:"is,omitempty"`
	IsNot         string `json:"is_not,omitempty"`
	MatchRegex    string `json:"match_regex,omitempty"`
	NotMatchRegex string `json:"not_match_regex,omitempty"`
}

// HopsCondition represents a hops condition
type HopsCondition struct {
	Op     string  `json:"op,omitempty"`
	Target float64 `json:"target,omitempty"`
}

// PacketLossCondition represents a packet loss condition
type PacketLossCondition struct {
	Op     string  `json:"op,omitempty"`
	Target float64 `json:"target,omitempty"`
}

// PacketsCondition represents a packets condition
type PacketsCondition struct {
	Op     string  `json:"op,omitempty"`
	Target float64 `json:"target,omitempty"`
}

// SuccessWhenItem represents a single success condition
type SuccessWhenItem struct {
	Body              []BodyCondition              `json:"body,omitempty"`
	Header            map[string][]HeaderCondition `json:"header,omitempty"`
	ResponseTime      any                          `json:"response_time,omitempty"`
	StatusCode        []StatusCodeCondition        `json:"status_code,omitempty"`
	ResponseMessage   []ResponseMessageCondition   `json:"response_message,omitempty"`
	Hops              []HopsCondition              `json:"hops,omitempty"`
	PacketLossPercent []PacketLossCondition        `json:"packet_loss_percent,omitempty"`
	Packets           []PacketsCondition           `json:"packets,omitempty"`
}

// DefaultRegion represents the default region structure for API responses
type DefaultRegion struct {
	UUID       string `json:"uuid,omitempty"`
	Province   string `json:"province,omitempty"`
	City       string `json:"city,omitempty"`
	Country    string `json:"country,omitempty"`
	Name       string `json:"name,omitempty"`
	NameEn     string `json:"name_en,omitempty"`
	ExtendInfo string `json:"extend_info,omitempty"`
	Internal   bool   `json:"internal,omitempty"`
	Keycode    string `json:"keycode,omitempty"`
	Isp        string `json:"isp,omitempty"`
	Status     string `json:"status,omitempty"`
	Region     string `json:"region,omitempty"`
	Owner      string `json:"owner,omitempty"`
	Heartbeat  int64  `json:"heartbeat,omitempty"`
	Company    string `json:"company,omitempty"`
	ExternalId string `json:"external_id,omitempty"`
	ParentAk   string `json:"parent_ak,omitempty"`
	CreateAt   int64  `json:"create_at,omitempty"`
}

// StatusUpdateRequest represents the request body for updating task status
type StatusUpdateRequest struct {
	TaskUUIDs []string `json:"taskUUIDs"`
	IsDisable bool     `json:"isDisable"`
}

func init() {
	apiURLs[consts.TypeNameSyntheticsTest] = map[string]string{
		ResourceCreate: "/dialing_task/add",
		ResourceDelete: "/dialing_task/delete",
		ResourceRead:   "/dialing_task/%s/get",
		ResourceUpdate: "/dialing_task/%s/modify",
		"multi_create": "/dialing_task/multi_task_add",
		"multi_update": "/dialing_task/%s/multi_task_modify",
		"set_status":   "/dialing_task/set_disable",
	}
}

// ReadDefaultRegion reads the default region list from the API
func (c *Client) ReadDefaultRegion() ([]*DefaultRegion, error) {
	var regions []*DefaultRegion

	if err := c.get("/dialing_task/default_region_list", &regions); err != nil {
		return nil, err
	}

	return regions, nil
}

// CustomRegion represents the custom region structure for API requests and responses
type CustomRegion struct {
	UUID       string `json:"uuid,omitempty"`
	Internal   bool   `json:"internal"`
	ISP        string `json:"isp,omitempty"`
	Country    string `json:"country,omitempty"`
	Province   string `json:"province,omitempty"`
	City       string `json:"city,omitempty"`
	Name       string `json:"name,omitempty"`
	NameEn     string `json:"name_en,omitempty"`
	Company    string `json:"company,omitempty"`
	Keycode    string `json:"keycode,omitempty"`
	CreateAt   int64  `json:"create_at,omitempty"`
	ExtendInfo string `json:"extend_info,omitempty"`
	ExternalID string `json:"external_id,omitempty"`
	Heartbeat  int64  `json:"heartbeat,omitempty"`
	Owner      string `json:"owner,omitempty"`
	ParentAK   string `json:"parent_ak,omitempty"`
	Region     string `json:"region,omitempty"`
	Status     string `json:"status,omitempty"`
}

// AKConfig represents the AK/SK configuration
type AKConfig struct {
	AK         string `json:"ak,omitempty"`
	SK         string `json:"sk,omitempty"`
	ExternalID string `json:"external_id,omitempty"`
	Owner      string `json:"owner,omitempty"`
	ParentAK   string `json:"parent_ak,omitempty"`
	Status     string `json:"status,omitempty"`
	UpdateAt   int64  `json:"update_at,omitempty"`
}

// CustomRegionResponse represents the response for custom region operations
type CustomRegionResponse struct {
	Code      int    `json:"code"`
	ErrorCode string `json:"errorCode"`
	Message   string `json:"message"`
	Success   bool   `json:"success"`
	TraceID   string `json:"traceId"`
	Content   struct {
		AK          *AKConfig         `json:"ak,omitempty"`
		Declaration map[string]string `json:"declaration,omitempty"`
		RegionInfo  *CustomRegion     `json:"regionInfo,omitempty"`
		Server      string            `json:"server,omitempty"`
	} `json:"content,omitempty"`
}
