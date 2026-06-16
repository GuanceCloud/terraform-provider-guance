package api

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/GuanceCloud/terraform-provider-guance/internal/consts"
)

type Mute struct {
	Name             string              `json:"name,omitempty"`
	Description      string              `json:"description,omitempty"`
	Type             string              `json:"type,omitempty"`
	MuteRanges       []MuteRange         `json:"muteRanges"`
	Tags             map[string][]string `json:"tags,omitempty"`
	FilterString     string              `json:"filterString,omitempty"`
	NotifyTargets    []MuteNotifyTarget  `json:"notifyTargets,omitempty"`
	NotifyMessage    string              `json:"notifyMessage,omitempty"`
	NotifyTimeStr    string              `json:"notifyTimeStr,omitempty"`
	StartTime        string              `json:"startTime,omitempty"`
	EndTime          string              `json:"endTime,omitempty"`
	RepeatTimeSet    int                 `json:"repeatTimeSet"`
	RepeatCrontabSet *RepeatCrontabSet   `json:"repeatCrontabSet,omitempty"`
	CrontabDuration  int                 `json:"crontabDuration,omitempty"`
	RepeatExpireTime string              `json:"repeatExpireTime,omitempty"`
	Timezone         string              `json:"timezone,omitempty"`
	Declaration      map[string]string   `json:"declaration,omitempty"`
}

type MuteContent struct {
	UUID             string              `json:"uuid,omitempty"`
	Name             string              `json:"name,omitempty"`
	Description      string              `json:"description,omitempty"`
	Type             string              `json:"type,omitempty"`
	MuteRanges       []MuteRange         `json:"muteRanges,omitempty"`
	Tags             map[string][]string `json:"tags,omitempty"`
	FilterString     string              `json:"filterString,omitempty"`
	NotifyTargets    []MuteNotifyTarget  `json:"notifyTargets,omitempty"`
	NotifyMessage    string              `json:"notifyMessage,omitempty"`
	NotifyTimeStr    string              `json:"notifyTimeStr,omitempty"`
	StartTime        string              `json:"startTime,omitempty"`
	EndTime          string              `json:"endTime,omitempty"`
	RepeatTimeSet    int                 `json:"repeatTimeSet,omitempty"`
	RepeatCrontabSet *RepeatCrontabSet   `json:"repeatCrontabSet,omitempty"`
	CrontabDuration  int                 `json:"crontabDuration,omitempty"`
	RepeatExpireTime string              `json:"repeatExpireTime,omitempty"`
	Timezone         string              `json:"timezone,omitempty"`
	Declaration      map[string]any      `json:"declaration,omitempty"`
	Status           int                 `json:"status,omitempty"`
	CreateAt         float64             `json:"createAt,omitempty"`
	UpdateAt         float64             `json:"updateAt,omitempty"`
	WorkspaceUUID    string              `json:"workspaceUUID,omitempty"`
	fieldPresent     map[string]bool     `json:"-"`
}

func (m *MuteContent) UnmarshalJSON(data []byte) error {
	type muteContentAlias MuteContent
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	var content muteContentAlias
	if err := json.Unmarshal(data, &content); err != nil {
		return err
	}
	*m = MuteContent(content)
	m.fieldPresent = make(map[string]bool, len(raw))
	for field := range raw {
		m.fieldPresent[field] = true
	}
	return nil
}

func (m *MuteContent) FieldPresent(field string) bool {
	return m != nil && m.fieldPresent[field]
}

type MuteRange struct {
	Name            string `json:"name,omitempty"`
	Type            string `json:"type,omitempty"`
	CheckerUUID     string `json:"checkerUUID,omitempty"`
	MonitorUUID     string `json:"monitorUUID,omitempty"`
	SLOUUID         string `json:"sloUUID,omitempty"`
	AlertPolicyUUID string `json:"alertPolicyUUID,omitempty"`
	TagUUID         string `json:"tagUUID,omitempty"`
}

type MuteNotifyTarget struct {
	To   []string `json:"to,omitempty"`
	Type string   `json:"type,omitempty"`
}

type RepeatCrontabSet struct {
	Min   string `json:"min,omitempty"`
	Hour  string `json:"hour,omitempty"`
	Day   string `json:"day,omitempty"`
	Month string `json:"month,omitempty"`
	Week  string `json:"week,omitempty"`
}

type MuteListContent struct {
	Data []json.RawMessage `json:"data,omitempty"`
}

type MuteListOptions struct {
	Search     string
	WorkStatus string
	IsEnable   string
	Type       string
	Creator    string
	Updator    string
}

func (c *Client) ListMutes(search string, content *MuteListContent) error {
	return c.ListMutesWithOptions(MuteListOptions{Search: search}, content)
}

func (c *Client) ListMutesWithOptions(options MuteListOptions, content *MuteListContent) error {
	content.Data = nil
	for pageIndex := 1; ; pageIndex++ {
		query := muteListQuery(options, pageIndex)
		var page MuteListContent
		if err := c.get("/monitor/mute/list?"+query.Encode(), &page); err != nil {
			return err
		}
		content.Data = append(content.Data, page.Data...)
		if len(page.Data) < muteLookupPageSize {
			return nil
		}
	}
}

func muteListQuery(options MuteListOptions, pageIndex int) url.Values {
	query := url.Values{}
	query.Set("pageIndex", fmt.Sprintf("%d", pageIndex))
	query.Set("pageSize", fmt.Sprintf("%d", muteLookupPageSize))
	if options.Search != "" {
		query.Set("search", options.Search)
	}
	if options.WorkStatus != "" {
		query.Set("workStatus", options.WorkStatus)
	}
	if options.IsEnable != "" {
		query.Set("isEnable", options.IsEnable)
	}
	if options.Type != "" {
		query.Set("type", options.Type)
	}
	if options.Creator != "" {
		query.Set("creator", options.Creator)
	}
	if options.Updator != "" {
		query.Set("updator", options.Updator)
	}
	return query
}

const muteLookupPageSize = 100

func (c *Client) GetMute(uuid string, content *MuteContent) error {
	for pageIndex := 1; ; pageIndex++ {
		path := fmt.Sprintf("/monitor/mute/list?pageIndex=%d&pageSize=%d", pageIndex, muteLookupPageSize)

		var list MuteListContent
		if err := c.get(path, &list); err != nil {
			return err
		}
		for _, raw := range list.Data {
			var item struct {
				UUID string `json:"uuid,omitempty"`
			}
			if err := json.Unmarshal(raw, &item); err != nil {
				continue
			}
			if item.UUID == uuid {
				if err := json.Unmarshal(raw, content); err != nil {
					return err
				}
				return nil
			}
		}
		if len(list.Data) < muteLookupPageSize {
			break
		}
	}

	return Error404
}

func init() {
	apiURLs[consts.TypeNameMute] = map[string]string{
		ResourceCreate: "/monitor/mute/create",
		ResourceUpdate: "/monitor/mute/%s/modify",
		ResourceDelete: "/monitor/mute/%s/delete",
		"enable":       "/monitor/mute/%s/enable",
		"disable":      "/monitor/mute/%s/disable",
	}
}

func (c *Client) SetMuteEnabled(uuid string, enabled bool) error {
	op := "disable"
	if enabled {
		op = "enable"
	}
	path, err := getResourcePath(consts.TypeNameMute, op)
	if err != nil {
		return fmt.Errorf("api path for mute %s not found: %w", op, err)
	}
	return c.post(fmt.Sprintf(path, uuid), nil, nil)
}
