package resourceschemas

import (
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/GuanceCloud/terraform-provider-guance/internal/consts"
)

// resourceSchemaResourceModel maps the resource schema data.
type resourceSchemaResourceModel struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Title       *I18n        `tfsdk:"title"`
	Description *I18n        `tfsdk:"description"`
	Models      []*Model     `tfsdk:"models"`
}

// GetId returns the ID of the resource.
func (m *resourceSchemaResourceModel) GetId() string {
	return m.ID.ValueString()
}

// SetId sets the ID of the resource.
func (m *resourceSchemaResourceModel) SetId(id string) {
	m.ID = types.StringValue(id)
}

// GetResourceType returns the type of the resource.
func (m *resourceSchemaResourceModel) GetResourceType() string {
	return consts.TypeNameResourceSchema
}

// resourceSchemaDataSourceModel maps the resource schema data.
type resourceSchemaDataSourceModel struct {
	Items      []*resourceSchemaResourceModel `tfsdk:"items"`
	MaxResults types.Int64                    `tfsdk:"max_results"`
	ID         types.String                   `tfsdk:"id"`
}

// ElemSchema maps the resource schema data.
type ElemSchema struct {
	Type   types.String `tfsdk:"type"`
	Format types.String `tfsdk:"format"`
	Ref    types.String `tfsdk:"ref"`
	Enum   []*Enum      `tfsdk:"enum"`
}

// Enum maps the resource schema data.
type Enum struct {
	Name  types.String `tfsdk:"name"`
	Value types.String `tfsdk:"value"`
	Title *I18n        `tfsdk:"title"`
}

// I18n maps the resource schema data.
type I18n struct {
	Zh types.String `tfsdk:"zh"`
	En types.String `tfsdk:"en"`
}

// Model maps the resource schema data.
type Model struct {
	Name        types.String `tfsdk:"name"`
	Title       *I18n        `tfsdk:"title"`
	Description *I18n        `tfsdk:"description"`
	Properties  []*Property  `tfsdk:"properties"`
}

// PropMeta maps the resource schema data.
type PropMeta struct {
	Dynamic   types.Bool `tfsdk:"dynamic"`
	Immutable types.Bool `tfsdk:"immutable"`
}

// PropSchema maps the resource schema data.
type PropSchema struct {
	Type     types.String `tfsdk:"type"`
	Format   types.String `tfsdk:"format"`
	Required types.Bool   `tfsdk:"required"`
	Elem     *ElemSchema  `tfsdk:"elem"`
	Enum     []*Enum      `tfsdk:"enum"`
	Model    types.String `tfsdk:"model"`
	Ref      types.String `tfsdk:"ref"`
}

// Property maps the resource schema data.
type Property struct {
	Name        types.String `tfsdk:"name"`
	Title       *I18n        `tfsdk:"title"`
	Description *I18n        `tfsdk:"description"`
	Schema      *PropSchema  `tfsdk:"schema"`
	Meta        *PropMeta    `tfsdk:"meta"`
}
