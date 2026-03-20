package api

import (
	"fmt"

	"github.com/GuanceCloud/terraform-provider-guance/internal/consts"
)

// AlertPolicy represents the alert policy structure for API requests
type AlertPolicy struct {
	Name               string           `json:"name,omitempty"`
	Desc               string           `json:"desc,omitempty"`
	OpenPermissionSet  bool             `json:"openPermissionSet,omitempty"`
	PermissionSet      []string         `json:"permissionSet,omitempty"`
	CheckerUUIDs       []string         `json:"checkerUUIDs,omitempty"`
	SecurityRuleUUIDs  []string         `json:"securityRuleUUIDs,omitempty"`
	RuleTimezone       string           `json:"ruleTimezone,omitempty"`
	AlertOpt           *AlertOpt        `json:"alertOpt,omitempty"`
}

// AlertPolicyContent represents the alert policy structure for API responses
type AlertPolicyContent struct {
	UUID               string           `json:"uuid,omitempty"`
	Name               string           `json:"name,omitempty"`
	Desc               string           `json:"desc,omitempty"`
	OpenPermissionSet  bool             `json:"openPermissionSet,omitempty"`
	PermissionSet      []string         `json:"permissionSet,omitempty"`
	CheckerUUIDs       []string         `json:"checkerUUIDs,omitempty"`
	SecurityRuleUUIDs  []string         `json:"securityRuleUUIDs,omitempty"`
	RuleTimezone       string           `json:"ruleTimezone,omitempty"`
	AlertOpt           *AlertOpt        `json:"alertOpt,omitempty"`
	CreateAt           float64          `json:"createAt,omitempty"`
	UpdateAt           float64          `json:"updateAt,omitempty"`
	WorkspaceUUID      string           `json:"workspaceUUID,omitempty"`
}

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

// SilentTimeoutByStatus represents the silentTimeoutByStatus structure
type SilentTimeoutByStatus struct {
	Status        string `json:"status,omitempty"`
	SilentTimeout int    `json:"silentTimeout,omitempty"`
}

// AlertTarget represents the alertTarget structure
type AlertTarget struct {
	Name             string      `json:"name,omitempty"`
	Targets          []Target    `json:"targets,omitempty"`
	Crontab          string      `json:"crontab,omitempty"`
	CrontabDuration  int         `json:"crontabDuration,omitempty"`
	CustomDateUUIDs  []string    `json:"customDateUUIDs,omitempty"`
	CustomStartTime  string      `json:"customStartTime,omitempty"`
	CustomDuration   int         `json:"customDuration,omitempty"`
	AlertInfo        []AlertInfo `json:"alertInfo,omitempty"`
}

// Target represents the targets structure
type Target struct {
	To              []string          `json:"to,omitempty"`
	Status          string            `json:"status,omitempty"`
	DfSource        string            `json:"df_source,omitempty"`
	UpgradeTargets  []UpgradeTarget   `json:"upgradeTargets,omitempty"`
	Tags            map[string][]string `json:"tags,omitempty"`
	FilterString    string            `json:"filterString,omitempty"`
}

// UpgradeTarget represents the upgradeTargets structure
type UpgradeTarget struct {
	To       []string `json:"to,omitempty"`
	Duration int      `json:"duration,omitempty"`
	ToWay    []string `json:"toWay,omitempty"`
}

// AlertInfo represents the alertInfo structure
type AlertInfo struct {
	Name         string    `json:"name,omitempty"`
	Targets      []Target  `json:"targets,omitempty"`
	FilterString string    `json:"filterString,omitempty"`
	MemberInfo   []string  `json:"memberInfo,omitempty"`
}

// AlertPolicyDeleteRequest represents the request body for deleting alert policies
type AlertPolicyDeleteRequest struct {
	AlertPolicyUUIDs []string `json:"alertPolicyUUIDs"`
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
