package member

import (
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/GuanceCloud/terraform-provider-guance/internal/consts"
)

// memberResourceModel maps the resource schema data.
type memberResourceModel struct {
	ID          types.String `tfsdk:"id"`
	WorkspaceId types.String `tfsdk:"workspace_id"`
	Username    types.String `tfsdk:"username"`
	Name        types.String `tfsdk:"name"`
	Email       types.String `tfsdk:"email"`
	Mobile      types.String `tfsdk:"mobile"`
	Role        types.String `tfsdk:"role"`
}

// GetId returns the ID of the resource.
func (m *memberResourceModel) GetId() string {
	return m.ID.ValueString()
}

// SetId sets the ID of the resource.
func (m *memberResourceModel) SetId(id string) {
	m.ID = types.StringValue(id)
}

// GetResourceType returns the type of the resource.
func (m *memberResourceModel) GetResourceType() string {
	return consts.TypeNameMember
}
