package intelligentinspection

import (
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/GuanceCloud/terraform-provider-guance/internal/consts"
)

// intelligentInspectionResourceModel maps the resource schema data.
type intelligentInspectionResourceModel struct {
	ID        types.String `tfsdk:"id"`
	MonitorId types.String `tfsdk:"monitor_id"`
	RefKey    types.String `tfsdk:"ref_key"`
	RefFunc   *RefFunc     `tfsdk:"ref_func"`
}

// GetId returns the ID of the resource.
func (m *intelligentInspectionResourceModel) GetId() string {
	return m.ID.ValueString()
}

// SetId sets the ID of the resource.
func (m *intelligentInspectionResourceModel) SetId(id string) {
	m.ID = types.StringValue(id)
}

// GetResourceType returns the type of the resource.
func (m *intelligentInspectionResourceModel) GetResourceType() string {
	return consts.TypeNameIntelligentInspection
}

// RefFunc maps the resource schema data.
type RefFunc struct {
	FuncId      types.String   `tfsdk:"func_id"`
	Title       types.String   `tfsdk:"title"`
	Description types.String   `tfsdk:"description"`
	Definition  types.String   `tfsdk:"definition"`
	Category    types.String   `tfsdk:"category"`
	Args        []types.String `tfsdk:"args"`
	Kwargs      types.String   `tfsdk:"kwargs"`
	Disabled    types.Bool     `tfsdk:"disabled"`
}
