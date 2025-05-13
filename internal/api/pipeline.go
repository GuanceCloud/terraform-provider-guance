package api

import (
	"github.com/GuanceCloud/terraform-provider-guance/internal/consts"
)

type PipelineExtend struct {
	AppID        []string `json:"appID"`
	Measurement  []string `json:"measurement"`
	LoggingIndex string   `json:"loggingIndex"`
}
type Pipeline struct {
	UUID              string          `json:"uuid,omitempty"`
	WorkspaceUUID     string          `json:"workspaceUUID,omitempty"`
	Category          string          `json:"category"`
	Name              string          `json:"name"`
	Source            []string        `json:"source"`
	Content           string          `json:"content"`
	TestData          string          `json:"testData"`
	Type              string          `json:"type"`
	DataType          string          `json:"dataType"`
	IsForce           bool            `json:"isForce"`
	AsDefault         int64           `json:"asDefault"`
	EnableByLogBackup int64           `json:"enableByLogBackup"`
	Status            int64           `json:"status,omitempty"`
	Extend            *PipelineExtend `json:"extend,omitempty"`
	ID                int64           `json:"id,omitempty"`
	CreateAt          int64           `json:"createAt,omitempty"`
	UpdateAt          float64         `json:"updateAt,omitempty"`
}

type PipelineResponse struct {
	Response
	Content *Pipeline `json:"content"`
}

func (c *Client) DisablePipeline(resourceUUIDs []string, isDisable bool, resp any) error {
	path := "/pipeline/batch_set_disable"

	item := struct {
		UUIDs     []string `json:"plUUIDs"`
		IsDisable bool     `json:"isDisable"`
	}{
		UUIDs:     resourceUUIDs,
		IsDisable: isDisable,
	}

	return c.post(path, item, resp)
}

func init() {
	apiURLs[consts.TypeNamePipeline] = map[string]string{
		ResourceCreate: "/pipeline/add",
		ResourceDelete: "/pipeline/delete?pipelineUUIDs=%s",
		ResourceRead:   "/pipeline/%s/get",
		ResourceUpdate: "/pipeline/%s/modify",
	}
}
