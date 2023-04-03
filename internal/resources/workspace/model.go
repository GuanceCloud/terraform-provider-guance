package workspace

import (
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/GuanceCloud/terraform-provider-guance/internal/consts"
)

// workspaceResourceModel maps the resource schema data.
type workspaceResourceModel struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
}

// GetId returns the ID of the resource.
func (m *workspaceResourceModel) GetId() string {
	return m.ID.ValueString()
}

// SetId sets the ID of the resource.
func (m *workspaceResourceModel) SetId(id string) {
	m.ID = types.StringValue(id)
}

// GetResourceType returns the type of the resource.
func (m *workspaceResourceModel) GetResourceType() string {
	return consts.TypeNameWorkspace
}
