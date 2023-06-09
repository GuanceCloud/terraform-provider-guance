// Code generated by Guance Cloud Code Generation Pipeline. DO NOT EDIT.

package blacklist

import (
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/GuanceCloud/terraform-provider-guance/internal/consts"
)

// blackListResourceModel maps the resource schema data.
type blackListResourceModel struct {
	ID          types.String `tfsdk:"id"`
	CreatedAt   types.String `tfsdk:"created_at"`
	Source      *Source      `tfsdk:"source"`
	FilterRules []*Filter    `tfsdk:"filter_rules"`
}

// GetId returns the ID of the resource.
func (m *blackListResourceModel) GetId() string {
	return m.ID.ValueString()
}

// SetId sets the ID of the resource.
func (m *blackListResourceModel) SetId(id string) {
	m.ID = types.StringValue(id)
}

// GetResourceType returns the type of the resource.
func (m *blackListResourceModel) GetResourceType() string {
	return consts.TypeNameBlackList
}

// SetCreatedAt sets the creation time of the resource.
func (m *blackListResourceModel) SetCreatedAt(t string) {
	m.CreatedAt = types.StringValue(t)
}

// Filter maps the resource schema data.
type Filter struct {
	Name      types.String   `tfsdk:"name"`
	Operation types.String   `tfsdk:"operation"`
	Condition types.String   `tfsdk:"condition"`
	Values    []types.String `tfsdk:"values"`
}

// Source maps the resource schema data.
type Source struct {
	Type types.String `tfsdk:"type"`
	Name types.String `tfsdk:"name"`
}
