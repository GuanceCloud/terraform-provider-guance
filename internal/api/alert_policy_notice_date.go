package api

import (
	"fmt"
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

const alertPolicyNoticeDateListPageSize = 100

func (c *Client) ListAlertPolicyNoticeDates(search string, content *AlertPolicyNoticeDateListContent) error {
	content.Data = nil
	for pageIndex := 1; ; pageIndex++ {
		query := alertPolicyNoticeDateListQuery(search, pageIndex)
		var page AlertPolicyNoticeDateListContent
		if err := c.get("/notice/date/list?"+query.Encode(), &page); err != nil {
			return err
		}
		content.Data = append(content.Data, page.Data...)
		if len(page.Data) < alertPolicyNoticeDateListPageSize {
			return nil
		}
	}
}

func alertPolicyNoticeDateListQuery(search string, pageIndex int) url.Values {
	query := url.Values{}
	query.Set("pageIndex", fmt.Sprintf("%d", pageIndex))
	query.Set("pageSize", fmt.Sprintf("%d", alertPolicyNoticeDateListPageSize))
	if search != "" {
		query.Set("search", search)
	}
	return query
}

func init() {
	apiURLs[consts.TypeNameAlertPolicyNoticeDate] = map[string]string{
		ResourceCreate: "/notice/date/add",
		ResourceRead:   "/notice/date/%s/get",
		ResourceUpdate: "/notice/date/%s/modify",
		ResourceDelete: "/notice/date/delete",
	}
}
