package custom_region

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// customRegionResourceModel maps the resource schema data
type customRegionResourceModel struct {
	UUID        types.String `tfsdk:"uuid"`
	Internal    types.Bool   `tfsdk:"internal"`
	ISP         types.String `tfsdk:"isp"`
	Country     types.String `tfsdk:"country"`
	Province    types.String `tfsdk:"province"`
	City        types.String `tfsdk:"city"`
	Name        types.String `tfsdk:"name"`
	NameEn      types.String `tfsdk:"name_en"`
	Company     types.String `tfsdk:"company"`
	Keycode     types.String `tfsdk:"keycode"`
	CreateAt    types.Int64  `tfsdk:"create_at"`
	ExtendInfo  types.String `tfsdk:"extend_info"`
	ExternalID  types.String `tfsdk:"external_id"`
	Heartbeat   types.Int64  `tfsdk:"heartbeat"`
	Owner       types.String `tfsdk:"owner"`
	ParentAK    types.String `tfsdk:"parent_ak"`
	Region      types.String `tfsdk:"region"`
	Status      types.String `tfsdk:"status"`
	AK          types.Object `tfsdk:"ak"`
	Server      types.String `tfsdk:"server"`
	Declaration types.Map    `tfsdk:"declaration"`
}
