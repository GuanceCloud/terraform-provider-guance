package api

import (
	"github.com/GuanceCloud/terraform-provider-guance/internal/consts"
)

// SLO represents the SLO structure for API requests
type SLO struct {
	Name              string   `json:"name,omitempty"`
	Interval          string   `json:"interval,omitempty"`
	Goal              float64  `json:"goal,omitempty"`
	MinGoal           float64  `json:"minGoal,omitempty"`
	SliUUIDs          []string `json:"sliUUIDs,omitempty"`
	Describe          string   `json:"describe,omitempty"`
	AlertPolicyUUIDs  []string `json:"alertPolicyUUIDs,omitempty"`
	Tags              []string `json:"tags,omitempty"`
}

// SLOContent represents the SLO structure for API responses
type SLOContent struct {
	UUID              string      `json:"uuid,omitempty"`
	Name              string      `json:"name,omitempty"`
	Status            int         `json:"status,omitempty"`
	CreateAt          float64     `json:"createAt,omitempty"`
	UpdateAt          float64     `json:"updateAt,omitempty"`
	WorkspaceUUID     string      `json:"workspaceUUID,omitempty"`
	Config            SLOConfig   `json:"config,omitempty"`
	AlertPolicyInfos  []interface{} `json:"alertPolicyInfos,omitempty"`
}

// SLOConfig represents the SLO config structure
type SLOConfig struct {
	CheckRange   int             `json:"checkRange,omitempty"`
	Describe     string          `json:"describe,omitempty"`
	Goal         float64         `json:"goal,omitempty"`
	Interval     string          `json:"interval,omitempty"`
	MinGoal      float64         `json:"minGoal,omitempty"`
	SliInfos     []SliInfo       `json:"sli_infos,omitempty"`
}

// SliInfo represents the SLI info structure
type SliInfo struct {
	ID     string `json:"id,omitempty"`
	Name   string `json:"name,omitempty"`
	Status int    `json:"status,omitempty"`
}

func init() {
	apiURLs[consts.TypeNameSLO] = map[string]string{
		ResourceCreate: "/slo/add",
		ResourceDelete: "/slo/%s/delete",
		ResourceRead:   "/slo/%s/get",
		ResourceUpdate: "/slo/%s/modify",
	}
}
