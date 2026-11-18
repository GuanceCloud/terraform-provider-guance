package api

import (
	"encoding/json"
	"fmt"

	"github.com/GuanceCloud/terraform-provider-guance/internal/consts"
)

const (
	pipelineExtendAppID        = "appID"
	pipelineExtendMeasurement  = "measurement"
	pipelineExtendLoggingIndex = "loggingIndex"
)

// PipelineExtend is the "extend" object of a pipeline.
//
// Only appID, measurement and loggingIndex are managed by Terraform. Any other
// key returned by the API (for example keys written by the console) is kept in
// Extra, so that it can be sent back unchanged on modify. Unset managed keys
// (nil) are left out of the JSON instead of being sent as null, and a
// LoggingIndex pointing to "" is sent as "loggingIndex": "".
type PipelineExtend struct {
	AppID        []string
	Measurement  []string
	LoggingIndex *string
	Extra        map[string]json.RawMessage
}

// MarshalJSON writes Extra plus the managed keys that are set.
func (e PipelineExtend) MarshalJSON() ([]byte, error) {
	out := make(map[string]any, len(e.Extra)+3)
	for k, v := range e.Extra {
		out[k] = v
	}
	if e.AppID != nil {
		out[pipelineExtendAppID] = e.AppID
	}
	if e.Measurement != nil {
		out[pipelineExtendMeasurement] = e.Measurement
	}
	if e.LoggingIndex != nil {
		out[pipelineExtendLoggingIndex] = *e.LoggingIndex
	}
	return json.Marshal(out)
}

// UnmarshalJSON reads the managed keys and keeps every other key in Extra.
// A managed key that is absent or null stays nil.
func (e *PipelineExtend) UnmarshalJSON(data []byte) error {
	raw := map[string]json.RawMessage{}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	*e = PipelineExtend{}
	for k, v := range raw {
		var err error
		switch k {
		case pipelineExtendAppID:
			err = json.Unmarshal(v, &e.AppID)
		case pipelineExtendMeasurement:
			err = json.Unmarshal(v, &e.Measurement)
		case pipelineExtendLoggingIndex:
			err = json.Unmarshal(v, &e.LoggingIndex)
		default:
			if e.Extra == nil {
				e.Extra = map[string]json.RawMessage{}
			}
			e.Extra[k] = v
		}
		if err != nil {
			return fmt.Errorf("pipeline extend %q: %w", k, err)
		}
	}
	return nil
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
