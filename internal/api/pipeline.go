package api

import (
	"github.com/GuanceCloud/terraform-provider-guance/internal/consts"
)

type PipelineExtend struct {
	AppID       []string `json:"appID"`
	Measurement []string `json:"measurement"`
}
type Pipeline struct {
	UUID      string          `json:"uuid,omitempty"`
	Category  string          `json:"category"`
	Name      string          `json:"name"`
	Source    []string        `json:"source"`
	Content   string          `json:"content"`
	TestData  string          `json:"testData"`
	Type      string          `json:"type"`
	IsForce   bool            `json:"isForce"`
	AsDefault int64           `json:"asDefault"`
	Extend    *PipelineExtend `json:"extend,omitempty"`
	ID        int64           `json:"id,omitempty"`
	CreateAt  int64           `json:"createAt,omitempty"`
}

type PipelineResponse struct {
	Response
	Content *Pipeline `json:"content"`
}

func init() {
	apiURLs[consts.TypeNamePipeline] = map[string]string{
		ResourceCreate: "/pipeline/add",
		ResourceDelete: "/pipeline/delete?pipelineUUIDs=%s",
		ResourceRead:   "/pipeline/%s/get",
		ResourceUpdate: "/pipeline/%s/modify",
	}
}
