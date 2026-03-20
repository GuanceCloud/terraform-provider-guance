package api

import (
	"fmt"

	"github.com/GuanceCloud/terraform-provider-guance/internal/consts"
)

// Monitor represents the monitor structure for API requests
type Monitor struct {
	Type              string      `json:"type,omitempty"`
	Status            int         `json:"status,omitempty"`
	Extend            interface{} `json:"extend,omitempty"`
	AlertPolicyUUIDs  []string    `json:"alertPolicyUUIDs,omitempty"`
	DashboardUUID     string      `json:"dashboardUUID,omitempty"`
	Tags              []string    `json:"tags,omitempty"`
	Secret            string      `json:"secret,omitempty"`
	JsonScript        interface{} `json:"jsonScript,omitempty"`
	OpenPermissionSet bool        `json:"openPermissionSet,omitempty"`
	PermissionSet     []string    `json:"permissionSet,omitempty"`
}

// MonitorContent represents the monitor structure for API responses
type MonitorContent struct {
	UUID              string      `json:"uuid,omitempty"`
	Type              string      `json:"type,omitempty"`
	Status            int         `json:"status,omitempty"`
	Extend            interface{} `json:"extend,omitempty"`
	AlertPolicyUUIDs  []string    `json:"alertPolicyUUIDs,omitempty"`
	DashboardUUID     string      `json:"dashboardUUID,omitempty"`
	Tags              []string    `json:"tags,omitempty"`
	Secret            string      `json:"secret,omitempty"`
	JsonScript        interface{} `json:"jsonScript,omitempty"`
	OpenPermissionSet bool        `json:"openPermissionSet,omitempty"`
	PermissionSet     []string    `json:"permissionSet,omitempty"`
	CreateAt          float64     `json:"createAt,omitempty"`
	UpdateAt          float64     `json:"updateAt,omitempty"`
	WorkspaceUUID     string      `json:"workspaceUUID,omitempty"`
	MonitorUUID       string      `json:"monitorUUID,omitempty"`
	MonitorName       string      `json:"monitorName,omitempty"`
	Creator           string      `json:"creator,omitempty"`
	Updator           string      `json:"updator,omitempty"`
	CreatedWay        string      `json:"createdWay,omitempty"`
	RefKey            string      `json:"refKey,omitempty"`
	CrontabInfo       interface{} `json:"crontabInfo,omitempty"`
	TagInfo           interface{} `json:"tagInfo,omitempty"`
	AlertPolicyInfos  interface{} `json:"alertPolicyInfos,omitempty"`
	CreatorInfo       interface{} `json:"creatorInfo,omitempty"`
	UpdatorInfo       interface{} `json:"updatorInfo,omitempty"`
}

func (c *Client) DeleteMonitor(key string) error {
	path, err := getResourcePath(consts.TypeNameMonitor, ResourceDelete)
	if err != nil {
		return fmt.Errorf("api path for delete not found: %w", err)
	}
	body := MonitorDeleteRequest{
		RuleUUIDs: []string{key},
	}

	err = c.post(path, body, nil)

	// just ignore 404 error
	if err == Error404 {
		return nil
	}
	return err
}

// MonitorDeleteRequest represents the request body for deleting monitors
type MonitorDeleteRequest struct {
	RuleUUIDs []string `json:"ruleUUIDs"`
}

func init() {
	apiURLs[consts.TypeNameMonitor] = map[string]string{
		ResourceCreate: "/checker/add",
		ResourceRead:   "/checker/%s/get",
		ResourceUpdate: "/checker/%s/modify",
		ResourceDelete: "/checker/delete",
	}
}
