package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type RoleItem struct {
	Name string `json:"name"`
	UUID string `json:"uuid"`
}

type Member struct {
	ID       int64      `json:"id"`
	UUID     string     `json:"uuid,omitempty"`
	CreateAt int64      `json:"createAt,omitempty"`
	Email    string     `json:"email"`
	Name     string     `json:"name"`
	Roles    []RoleItem `json:"roles"`
}

type MemberData struct {
	Data []*Member `json:"data"`
}

type MemberResponse struct {
	Response
	Content MemberData `json:"content"`
}

func (c *Client) ReadMember(search string) ([]*Member, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/v1/workspace/members/list?search=%s", c.EndPoint, search), nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	res := MemberResponse{}
	err = json.Unmarshal(body, &res)
	if err != nil {
		return nil, err
	}

	if res.Code == 404 {
		return nil, nil
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return res.Content.Data, nil
}
