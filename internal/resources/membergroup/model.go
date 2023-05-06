// Code generated by Guance Cloud Code Generation Pipeline. DO NOT EDIT.

package membergroup

import (
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/GuanceCloud/terraform-provider-guance/internal/consts"
)

// memberGroupResourceModel maps the resource schema data.
type memberGroupResourceModel struct {
	ID        types.String   `tfsdk:"id"`
	CreatedAt types.String   `tfsdk:"created_at"`
	Name      types.String   `tfsdk:"name"`
	MemberIds []types.String `tfsdk:"member_ids"`
}

// GetId returns the ID of the resource.
func (m *memberGroupResourceModel) GetId() string {
	return m.ID.ValueString()
}

// SetId sets the ID of the resource.
func (m *memberGroupResourceModel) SetId(id string) {
	m.ID = types.StringValue(id)
}

// GetResourceType returns the type of the resource.
func (m *memberGroupResourceModel) GetResourceType() string {
	return consts.TypeNameMemberGroup
}

// SetCreatedAt sets the creation time of the resource.
func (m *memberGroupResourceModel) SetCreatedAt(t string) {
	m.CreatedAt = types.StringValue(t)
}