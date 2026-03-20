package api

import (
	"github.com/GuanceCloud/terraform-provider-guance/internal/consts"
)

// Dashboard represents the dashboard structure for API requests
type Dashboard struct {
	Name                 string        `json:"name,omitempty"`
	Desc                 string        `json:"desc,omitempty"`
	Identifier           string        `json:"identifier,omitempty"`
	Extend               interface{}   `json:"extend,omitempty"`
	Mapping              []interface{} `json:"mapping,omitempty"`
	TagNames             []string      `json:"tagNames,omitempty"`
	TemplateInfo         interface{}   `json:"templateInfo,omitempty"`
	SpecifyDashboardUUID string        `json:"specifyDashboardUUID,omitempty"`
	IsPublic             int           `json:"isPublic"`
	PermissionSet        []string      `json:"permissionSet,omitempty"`
	ReadPermissionSet    []string      `json:"readPermissionSet,omitempty"`
}

// DashboardContent represents the dashboard structure for API responses
type DashboardContent struct {
	UUID              string        `json:"uuid,omitempty"`
	Name              string        `json:"name,omitempty"`
	Desc              string        `json:"desc,omitempty"`
	Identifier        string        `json:"identifier,omitempty"`
	Extend            interface{}   `json:"extend,omitempty"`
	Mapping           []interface{} `json:"mapping,omitempty"`
	CreateAt          float64       `json:"createAt,omitempty"`
	UpdateAt          float64       `json:"updateAt,omitempty"`
	WorkspaceUUID     string        `json:"workspaceUUID,omitempty"`
	IsPublic          int           `json:"isPublic"`
	PermissionSet     []string      `json:"permissionSet,omitempty"`
	ReadPermissionSet []string      `json:"readPermissionSet,omitempty"`
	TagInfo           interface{}   `json:"tagInfo,omitempty"`
}

func init() {
	apiURLs[consts.TypeNameDashboard] = map[string]string{
		ResourceCreate: "/dashboards/create",
		ResourceDelete: "/dashboards/%s/delete",
		ResourceRead:   "/dashboards/%s/get",
		ResourceUpdate: "/dashboards/%s/modify",
		"export":      "/dashboards/%s/export",
		"import":      "/dashboards/%s/import",
	}
}
