package functions

import (
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/GuanceCloud/terraform-provider-guance/internal/consts"
)

// functionResourceModel maps the resource schema data.
type functionResourceModel struct {
	ID          types.String `tfsdk:"id"`
	Title       types.String `tfsdk:"title"`
	Description types.String `tfsdk:"description"`
	FuncId      types.String `tfsdk:"func_id"`
}

// GetId returns the ID of the resource.
func (m *functionResourceModel) GetId() string {
	return m.ID.ValueString()
}

// SetId sets the ID of the resource.
func (m *functionResourceModel) SetId(id string) {
	m.ID = types.StringValue(id)
}

// GetResourceType returns the type of the resource.
func (m *functionResourceModel) GetResourceType() string {
	return consts.TypeNameFunction
}

// functionDataSourceModel maps the resource schema data.
type functionDataSourceModel struct {
	Items      []*functionResourceModel `tfsdk:"items"`
	MaxResults types.Int64              `tfsdk:"max_results"`
	ID         types.String             `tfsdk:"id"`
}
