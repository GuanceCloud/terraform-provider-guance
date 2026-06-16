package api

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/GuanceCloud/terraform-provider-guance/internal/consts"
)

// AlertPolicy represents the alert policy structure for API requests
type AlertPolicy struct {
	Name              string    `json:"name,omitempty"`
	Desc              string    `json:"desc,omitempty"`
	OpenPermissionSet bool      `json:"openPermissionSet,omitempty"`
	PermissionSet     []string  `json:"permissionSet,omitempty"`
	CheckerUUIDs      []string  `json:"checkerUUIDs,omitempty"`
	SecurityRuleUUIDs []string  `json:"securityRuleUUIDs,omitempty"`
	RuleTimezone      string    `json:"ruleTimezone,omitempty"`
	AlertOpt          *AlertOpt `json:"alertOpt,omitempty"`
}

// AlertPolicyContent represents the alert policy structure for API responses
type AlertPolicyContent struct {
	UUID              string           `json:"uuid,omitempty"`
	Name              string           `json:"name,omitempty"`
	Desc              string           `json:"desc,omitempty"`
	OpenPermissionSet bool             `json:"openPermissionSet,omitempty"`
	PermissionSet     []string         `json:"permissionSet,omitempty"`
	CheckerUUIDs      []string         `json:"checkerUUIDs,omitempty"`
	SecurityRuleUUIDs []string         `json:"securityRuleUUIDs,omitempty"`
	RuleTimezone      string           `json:"ruleTimezone,omitempty"`
	AlertOpt          *AlertOptContent `json:"alertOpt,omitempty"`
	CreateAt          float64          `json:"createAt,omitempty"`
	UpdateAt          float64          `json:"updateAt,omitempty"`
	WorkspaceUUID     string           `json:"workspaceUUID,omitempty"`
}

type AlertPolicyListContent struct {
	Data []AlertPolicyContent `json:"data,omitempty"`
}

type AlertPolicyListOptions struct {
	Search            string
	NotifyObjectUUIDs []string
}

const alertPolicyListPageSize = 100

// AlertOpt represents the alertOpt structure
type AlertOpt struct {
	AggType                     string                  `json:"aggType,omitempty"`
	IgnoreOK                    bool                    `json:"ignoreOK,omitempty"`
	AlertType                   string                  `json:"alertType,omitempty"`
	SilentTimeout               int                     `json:"silentTimeout,omitempty"`
	SilentTimeoutByStatusEnable bool                    `json:"silentTimeoutByStatusEnable,omitempty"`
	SilentTimeoutByStatus       []SilentTimeoutByStatus `json:"silentTimeoutByStatus,omitempty"`
	AlertTarget                 []AlertTarget           `json:"alertTarget,omitempty"`
	AggInterval                 int                     `json:"aggInterval,omitempty"`
	AggFields                   []string                `json:"aggFields,omitempty"`
	AggLabels                   []string                `json:"aggLabels,omitempty"`
	AggClusterFields            []string                `json:"aggClusterFields,omitempty"`
	AggSendFirst                bool                    `json:"aggSendFirst,omitempty"`
}

// AlertOptContent represents the alertOpt structure for API responses.
type AlertOptContent struct {
	AggType                     string                  `json:"aggType,omitempty"`
	IgnoreOK                    *bool                   `json:"ignoreOK,omitempty"`
	AlertType                   string                  `json:"alertType,omitempty"`
	SilentTimeout               *int                    `json:"silentTimeout,omitempty"`
	SilentTimeoutByStatusEnable *bool                   `json:"silentTimeoutByStatusEnable,omitempty"`
	SilentTimeoutByStatus       []SilentTimeoutByStatus `json:"silentTimeoutByStatus,omitempty"`
	AlertTarget                 []AlertTargetContent    `json:"alertTarget,omitempty"`
	AggInterval                 *int                    `json:"aggInterval,omitempty"`
	AggFields                   []string                `json:"aggFields,omitempty"`
	AggLabels                   []string                `json:"aggLabels,omitempty"`
	AggClusterFields            []string                `json:"aggClusterFields,omitempty"`
	AggSendFirst                *bool                   `json:"aggSendFirst,omitempty"`
}

// SilentTimeoutByStatus represents the silentTimeoutByStatus structure
type SilentTimeoutByStatus struct {
	Status        string `json:"status,omitempty"`
	SilentTimeout int    `json:"silentTimeout,omitempty"`
}

// AlertTarget represents the alertTarget structure
type AlertTarget struct {
	Name            string      `json:"name,omitempty"`
	Targets         []Target    `json:"targets,omitempty"`
	Crontab         string      `json:"crontab,omitempty"`
	CrontabDuration int         `json:"crontabDuration,omitempty"`
	CustomDateUUIDs []string    `json:"customDateUUIDs,omitempty"`
	CustomStartTime string      `json:"customStartTime,omitempty"`
	CustomDuration  int         `json:"customDuration,omitempty"`
	AlertInfo       []AlertInfo `json:"alertInfo,omitempty"`
}

// AlertTargetContent represents the alertTarget structure for API responses.
type AlertTargetContent struct {
	Name            string             `json:"name,omitempty"`
	Targets         []TargetContent    `json:"targets,omitempty"`
	Crontab         string             `json:"crontab,omitempty"`
	CrontabDuration *int               `json:"crontabDuration,omitempty"`
	CustomDateUUIDs []string           `json:"customDateUUIDs,omitempty"`
	CustomStartTime string             `json:"customStartTime,omitempty"`
	CustomDuration  *int               `json:"customDuration,omitempty"`
	AlertInfo       []AlertInfoContent `json:"alertInfo,omitempty"`
}

// Target represents the targets structure
type Target struct {
	To             []string            `json:"to,omitempty"`
	Status         string              `json:"status,omitempty"`
	DfSource       string              `json:"df_source,omitempty"`
	UpgradeTargets []UpgradeTarget     `json:"upgradeTargets,omitempty"`
	Tags           map[string][]string `json:"tags,omitempty"`
	FilterString   string              `json:"filterString,omitempty"`
}

// TargetContent represents the targets structure for API responses.
type TargetContent struct {
	To             []string               `json:"to,omitempty"`
	Status         string                 `json:"status,omitempty"`
	DfSource       string                 `json:"df_source,omitempty"`
	UpgradeTargets []UpgradeTargetContent `json:"upgradeTargets,omitempty"`
	Tags           map[string][]string    `json:"tags,omitempty"`
	FilterString   string                 `json:"filterString,omitempty"`
}

// UpgradeTarget represents the upgradeTargets structure
type UpgradeTarget struct {
	To       []string `json:"to,omitempty"`
	Duration int      `json:"duration,omitempty"`
	ToWay    []string `json:"toWay,omitempty"`
}

// UpgradeTargetContent represents the upgradeTargets structure for API responses.
type UpgradeTargetContent struct {
	To       []string `json:"to,omitempty"`
	Duration *int     `json:"duration,omitempty"`
	ToWay    []string `json:"toWay,omitempty"`
}

// AlertInfo represents the alertInfo structure
type AlertInfo struct {
	Name         string   `json:"name,omitempty"`
	Targets      []Target `json:"targets,omitempty"`
	FilterString string   `json:"filterString,omitempty"`
	MemberInfo   []string `json:"memberInfo,omitempty"`
}

// AlertInfoContent represents the alertInfo structure for API responses.
type AlertInfoContent struct {
	Name         string          `json:"name,omitempty"`
	Targets      []TargetContent `json:"targets,omitempty"`
	FilterString string          `json:"filterString,omitempty"`
	MemberInfo   []string        `json:"memberInfo,omitempty"`
}

// AlertPolicyDeleteRequest represents the request body for deleting alert policies
type AlertPolicyDeleteRequest struct {
	AlertPolicyUUIDs []string `json:"alertPolicyUUIDs"`
}

func (c *Client) ListAlertPolicies(search string, content *AlertPolicyListContent) error {
	return c.ListAlertPoliciesWithOptions(AlertPolicyListOptions{Search: search}, content)
}

func (c *Client) ListAlertPoliciesWithOptions(options AlertPolicyListOptions, content *AlertPolicyListContent) error {
	content.Data = nil
	for pageIndex := 1; ; pageIndex++ {
		query := alertPolicyListQuery(options, pageIndex)
		var page AlertPolicyListContent
		if err := c.get("/alert_policy/list?"+query.Encode(), &page); err != nil {
			return err
		}
		content.Data = append(content.Data, page.Data...)
		if len(page.Data) < alertPolicyListPageSize {
			return nil
		}
	}
}

func alertPolicyListQuery(options AlertPolicyListOptions, pageIndex int) url.Values {
	query := url.Values{}
	query.Set("pageIndex", fmt.Sprintf("%d", pageIndex))
	query.Set("pageSize", fmt.Sprintf("%d", alertPolicyListPageSize))
	if options.Search != "" {
		query.Set("search", options.Search)
	}
	notifyObjectUUIDs := compactStrings(options.NotifyObjectUUIDs)
	if len(notifyObjectUUIDs) > 0 {
		query.Set("notifyObjectUUIDs", strings.Join(notifyObjectUUIDs, ","))
	}
	return query
}

func compactStrings(values []string) []string {
	result := make([]string, 0, len(values))
	for _, value := range values {
		if value != "" {
			result = append(result, value)
		}
	}
	return result
}

func (c *Client) DeleteAlertPolicy(key string) error {
	path, err := getResourcePath(consts.TypeNameAlertPolicy, ResourceDelete)
	if err != nil {
		return fmt.Errorf("api path for delete not found: %w", err)
	}
	body := AlertPolicyDeleteRequest{
		AlertPolicyUUIDs: []string{key},
	}

	err = c.post(path, body, nil)

	// just ignore 404 error
	if err == Error404 {
		return nil
	}
	return err
}

func init() {
	apiURLs[consts.TypeNameAlertPolicy] = map[string]string{
		ResourceCreate: "/alert_policy/add_v2",
		ResourceRead:   "/alert_policy/%s/get",
		ResourceUpdate: "/alert_policy/%s/modify_v2",
		ResourceDelete: "/alert_policy/delete",
	}
}
