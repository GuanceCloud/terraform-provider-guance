package alert_policy

import "github.com/GuanceCloud/terraform-provider-guance/internal/api"

type AlertOptModel = alertOptModel

func AlertOptFromContentForDataSource(content *api.AlertOptContent) *AlertOptModel {
	return alertOptFromContentForDataSource(content)
}
