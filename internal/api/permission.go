package api

import (
	"fmt"
)

type Permission struct {
	Desc                string        `json:"desc"`
	Disabled            uint8         `json:"disabled"`
	IsSupportCustomRole uint8         `json:"isSupportCustomRole"`
	IsSupportGeneral    uint8         `json:"isSupportGeneral"`
	IsSupportOwner      uint8         `json:"isSupportOwner"`
	IsSupportReadOnly   uint8         `json:"isSupportReadOnly"`
	IsSupportWsAdmin    uint8         `json:"isSupportWsAdmin"`
	Key                 string        `json:"key"`
	Name                string        `json:"name"`
	Subs                []*Permission `json:"subs"`
}

func (c *Client) ReadPermission(isSupportCustomRole bool) ([]*Permission, error) {
	content := []*Permission{}
	err := c.get(fmt.Sprintf("/permission/list?isSupportCustomRole=%t", isSupportCustomRole), &content)
	return content, err
}
