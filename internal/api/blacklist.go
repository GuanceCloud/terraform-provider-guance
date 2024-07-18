package api

import (
	"github.com/GuanceCloud/terraform-provider-guance/internal/consts"
)

type Filter struct {
	Name      string   `json:"name"`
	Value     []string `json:"value"`
	Condition string   `json:"condition"`
	Operation string   `json:"operation"`
}

type Blacklist struct {
	ID         int64    `json:"id,omitempty"`
	UUID       string   `json:"uuid,omitempty"`
	CreateAt   int64    `json:"createAt,omitempty"`
	Conditions string   `json:"conditions,omitempty"`
	Source     string   `json:"source"`
	Type       string   `json:"type"`
	Filters    []Filter `json:"filters"`
}

func init() {
	apiURLs[consts.TypeNameBlackList] = map[string]string{
		ResourceCreate: "/blacklist/add",
		ResourceDelete: "/blacklist/delete?blacklistUUIDs=%s",
		ResourceRead:   "/blacklist/%s/get",
		ResourceUpdate: "/blacklist/%s/modify",
	}
}
