package api

import (
	"fmt"

	"github.com/GuanceCloud/terraform-provider-guance/internal/consts"
)

type Role struct {
	UUID        string   `json:"uuid,omitempty"`
	Name        string   `json:"name"`
	Desc        string   `json:"desc"`
	Keys        []string `json:"keys,omitempty"`
	Permissions []string `json:"permissions,omitempty"`
}

func (c *Client) DeleteRole(key string) error {
	path, err := getResourcePath(consts.TypeNameRole, ResourceDelete)
	if err != nil {
		return fmt.Errorf("api path for delete not found: %w", err)
	}
	path = fmt.Sprintf(path, key)
	err = c.post(path, nil, nil)

	// just ignore 404 error
	if err == Error404 {
		return nil
	}
	return err
}

func init() {
	apiURLs[consts.TypeNameRole] = map[string]string{
		ResourceCreate: "/role/add",
		ResourceDelete: "/role/%s/delete",
		ResourceRead:   "/role/%s/get",
		ResourceUpdate: "/role/%s/modify",
	}
}
