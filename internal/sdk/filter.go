package sdk

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/tidwall/gjson"
)

type Filter struct {
	Name   types.String   `tfsdk:"name"`
	Values []types.String `tfsdk:"values"`
}

func (f *Filter) IsOK(state string) bool {
	for _, value := range f.Values {
		if gjson.Get(state, f.Name.ValueString()).String() == value.ValueString() {
			return true
		}
	}
	return false
}

func FilterAllSuccess(state string, filters ...*Filter) bool {
	for _, filter := range filters {
		if !filter.IsOK(state) {
			return false
		}
	}
	return true
}
