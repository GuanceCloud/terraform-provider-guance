package default_region

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// defaultRegionResourceModel maps the resource schema data.
type defaultRegionResourceModel struct {
	UUID       types.String `tfsdk:"uuid"`
	Province   types.String `tfsdk:"province"`
	City       types.String `tfsdk:"city"`
	Country    types.String `tfsdk:"country"`
	Name       types.String `tfsdk:"name"`
	NameEn     types.String `tfsdk:"name_en"`
	ExtendInfo types.String `tfsdk:"extend_info"`
	Internal   types.Bool   `tfsdk:"internal"`
	Keycode    types.String `tfsdk:"keycode"`
	Isp        types.String `tfsdk:"isp"`
	Status     types.String `tfsdk:"status"`
	Region     types.String `tfsdk:"region"`
	Owner      types.String `tfsdk:"owner"`
	Heartbeat  types.String `tfsdk:"heartbeat"`
	Company    types.String `tfsdk:"company"`
	ExternalId types.String `tfsdk:"external_id"`
	ParentAk   types.String `tfsdk:"parent_ak"`
	CreateAt   types.String `tfsdk:"create_at"`
}

// defaultRegionDataSourceModel maps the resource schema data.
type defaultRegionDataSourceModel struct {
	Regions []defaultRegionResourceModel `tfsdk:"regions"`
}
