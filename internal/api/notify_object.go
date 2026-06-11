package api

import (
	"net/url"

	"github.com/GuanceCloud/terraform-provider-guance/internal/consts"
)

// NotifyObject represents the notify object structure for API requests
type NotifyObject struct {
	Type              string      `json:"type,omitempty"`
	Name              string      `json:"name,omitempty"`
	OptSet            interface{} `json:"optSet,omitempty"`
	OpenPermissionSet bool        `json:"openPermissionSet,omitempty"`
	PermissionSet     []string    `json:"permissionSet,omitempty"`
}

// NotifyObjectContent represents the notify object structure for API responses
type NotifyObjectContent struct {
	UUID              string      `json:"uuid,omitempty"`
	Type              string      `json:"type,omitempty"`
	Name              string      `json:"name,omitempty"`
	OptSet            interface{} `json:"optSet,omitempty"`
	OpenPermissionSet bool        `json:"openPermissionSet,omitempty"`
	PermissionSet     []string    `json:"permissionSet,omitempty"`
	CreateAt          float64     `json:"createAt,omitempty"`
	UpdateAt          float64     `json:"updateAt,omitempty"`
	WorkspaceUUID     string      `json:"workspaceUUID,omitempty"`
}

type NotifyObjectListContent struct {
	Data []NotifyObjectContent `json:"data,omitempty"`
}

func (c *Client) ListNotifyObjects(search string, content *NotifyObjectListContent) error {
	query := url.Values{}
	query.Set("pageIndex", "1")
	query.Set("pageSize", "100")
	if search != "" {
		query.Set("search", search)
	}
	return c.get("/notify_object/list?"+query.Encode(), content)
}

// UpdateNotifyObject updates a notify object
func (c *Client) UpdateNotifyObject(body any, content any) error {
	return c.post("/notify_object/modify", body, content)
}

// DeleteNotifyObject deletes a notify object by UUID
func (c *Client) DeleteNotifyObject(uuid string) error {
	body := map[string]string{
		"notifyObjectUUID": uuid,
	}
	err := c.post("/notify_object/delete", body, nil)
	if err == Error404 {
		return nil
	}
	return err
}

// GetNotifyObject gets a notify object by UUID
func (c *Client) GetNotifyObject(uuid string, content any) error {
	path := "/notify_object/get?notifyObjectUUID=" + url.QueryEscape(uuid)
	return c.get(path, content)
}

func init() {
	apiURLs[consts.TypeNameNotifyObject] = map[string]string{
		ResourceCreate: "/notify_object/create",
		ResourceDelete: "/notify_object/delete",
		ResourceRead:   "/notify_object/get",
		ResourceUpdate: "/notify_object/modify",
	}
}
