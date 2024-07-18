package membergroup

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// memberGroupResourceModel maps the resource schema data.
type memberGroupResourceModel struct {
	UUID         types.String   `tfsdk:"uuid"`
	Name         types.String   `tfsdk:"name"`
	AccountUUIDs []types.String `tfsdk:"account_uuids"`
}
