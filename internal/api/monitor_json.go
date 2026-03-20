package api

import (
	"github.com/GuanceCloud/terraform-provider-guance/internal/consts"
)

// MonitorJson represents the monitor json structure for API requests
type MonitorJson struct {
	Checkers []interface{} `json:"checkers"`
	Type     string        `json:"type,omitempty"`
}

// MonitorJsonExportRequest represents the request body for exporting monitors
type MonitorJsonExportRequest struct {
	Checkers []string `json:"checkers"`
}

// MonitorJsonImportRequest represents the request body for importing monitors
type MonitorJsonImportRequest struct {
	Checkers             []interface{} `json:"checkers"`
	Type                 string        `json:"type,omitempty"`
	SkipRepeatNameCheck  bool          `json:"skipRepeatNameCheck"`
	SkipRepeatNameCreate bool          `json:"skipRepeatNameCreate"`
	DeleteRepeatName     bool          `json:"deleteRepeatName"`
}

// MonitorJsonReplaceRequest represents the request body for replacing a monitor
type MonitorJsonReplaceRequest struct {
	Checker interface{} `json:"checker"`
}

// MonitorJsonContent represents the monitor json structure for API responses
type MonitorJsonContent struct {
	Checkers []interface{} `json:"checkers"`
}

// MonitorJsonImportContent represents the import response content
type MonitorJsonImportContent struct {
	Declaration interface{}      `json:"declaration"`
	Rule        []MonitorContent `json:"rule"`
}

func init() {
	apiURLs[consts.TypeNameMonitorJson] = map[string]string{
		"export":  "/checker/export",
		"import":  "/checker/import",
		"replace": "/checker/%s/replace",
	}
}
