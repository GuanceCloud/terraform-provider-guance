package dashboard

import (
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/GuanceCloud/terraform-provider-guance/internal/consts"
)

// dashboardResourceModel maps the resource schema data.
type dashboardResourceModel struct {
	ID       types.String   `tfsdk:"id"`
	Name     types.String   `tfsdk:"name"`
	Extend   types.String   `tfsdk:"extend"`
	Mapping  []*Mapping     `tfsdk:"mapping"`
	Tags     []types.String `tfsdk:"tags"`
	Template *Template      `tfsdk:"template"`
}

// GetId returns the ID of the resource.
func (m *dashboardResourceModel) GetId() string {
	return m.ID.ValueString()
}

// SetId sets the ID of the resource.
func (m *dashboardResourceModel) SetId(id string) {
	m.ID = types.StringValue(id)
}

// GetResourceType returns the type of the resource.
func (m *dashboardResourceModel) GetResourceType() string {
	return consts.TypeNameDashboard
}

// Chart maps the resource schema data.
type Chart struct {
	Name    types.String  `tfsdk:"name"`
	Type    types.String  `tfsdk:"type"`
	Group   types.String  `tfsdk:"group"`
	Pos     *ChartPos     `tfsdk:"pos"`
	Queries []*ChartQuery `tfsdk:"queries"`
}

// ChartPos maps the resource schema data.
type ChartPos struct {
	I types.String `tfsdk:"i"`
	X types.Int64  `tfsdk:"x"`
	Y types.Int64  `tfsdk:"y"`
	W types.Int64  `tfsdk:"w"`
	H types.Int64  `tfsdk:"h"`
}

// ChartQuery maps the resource schema data.
type ChartQuery struct {
	Checked    types.Bool   `tfsdk:"checked"`
	Datasource types.String `tfsdk:"datasource"`
	Qtype      types.String `tfsdk:"qtype"`
	Query      *Query       `tfsdk:"query"`
	Unit       types.String `tfsdk:"unit"`
}

// IconSet maps the resource schema data.
type IconSet struct {
	Sm types.String `tfsdk:"sm"`
	Md types.String `tfsdk:"md"`
}

// Main maps the resource schema data.
type Main struct {
	Type   types.String   `tfsdk:"type"`
	Vars   []*Var         `tfsdk:"vars"`
	Groups []types.String `tfsdk:"groups"`
	Charts []*Chart       `tfsdk:"charts"`
}

// Mapping maps the resource schema data.
type Mapping struct {
	Class      types.String `tfsdk:"class"`
	Field      types.String `tfsdk:"field"`
	Mapping    types.String `tfsdk:"mapping"`
	Datasource types.String `tfsdk:"datasource"`
}

// Query maps the resource schema data.
type Query struct {
	Density     types.String   `tfsdk:"density"`
	Filter      []*QueryFilter `tfsdk:"filter"`
	GroupBy     types.String   `tfsdk:"group_by"`
	GroupByTime types.String   `tfsdk:"group_by_time"`
	Q           types.String   `tfsdk:"q"`
}

// QueryFilter maps the resource schema data.
type QueryFilter struct {
	Logic types.String `tfsdk:"logic"`
	Name  types.String `tfsdk:"name"`
	Op    types.String `tfsdk:"op"`
	Value types.String `tfsdk:"value"`
}

// Template maps the resource schema data.
type Template struct {
	Title     types.String       `tfsdk:"title"`
	Summary   types.String       `tfsdk:"summary"`
	Dashboard *TemplateDashboard `tfsdk:"dashboard"`
	IconSet   *IconSet           `tfsdk:"icon_set"`
	Main      *Main              `tfsdk:"main"`
}

// TemplateDashboard maps the resource schema data.
type TemplateDashboard struct {
	Extend  types.String `tfsdk:"extend"`
	Mapping []*Mapping   `tfsdk:"mapping"`
}

// Var maps the resource schema data.
type Var struct {
	Name  types.String `tfsdk:"name"`
	Type  types.String `tfsdk:"type"`
	Value types.String `tfsdk:"value"`
	Label types.String `tfsdk:"label"`
}
