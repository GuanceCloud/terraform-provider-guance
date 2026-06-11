package api

import (
	"net/url"

	"github.com/GuanceCloud/terraform-provider-guance/internal/consts"
)

type AlertPolicyNoticeDate struct {
	Name        string   `json:"name,omitempty"`
	NoticeDates []string `json:"noticeDates,omitempty"`
}

type AlertPolicyNoticeDateContent struct {
	UUID          string   `json:"uuid,omitempty"`
	Name          string   `json:"name,omitempty"`
	Dates         []string `json:"dates,omitempty"`
	CreateAt      float64  `json:"createAt,omitempty"`
	UpdateAt      float64  `json:"updateAt,omitempty"`
	WorkspaceUUID string   `json:"workspaceUUID,omitempty"`
}

type AlertPolicyNoticeDateListContent struct {
	Data []AlertPolicyNoticeDateContent `json:"data,omitempty"`
}

func (c *Client) ListAlertPolicyNoticeDates(search string, content *AlertPolicyNoticeDateListContent) error {
	query := url.Values{}
	query.Set("pageIndex", "1")
	query.Set("pageSize", "100")
	if search != "" {
		query.Set("search", search)
	}
	return c.get("/notice/date/list?"+query.Encode(), content)
}

func init() {
	apiURLs[consts.TypeNameAlertPolicyNoticeDate] = map[string]string{
		ResourceCreate: "/notice/date/add",
		ResourceRead:   "/notice/date/%s/get",
		ResourceUpdate: "/notice/date/%s/modify",
		ResourceDelete: "/notice/date/delete",
	}
}
