package mute

import (
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/GuanceCloud/terraform-provider-guance/internal/consts"
)

// muteResourceModel maps the resource schema data.
type muteResourceModel struct {
	ID         types.String   `tfsdk:"id"`
	MuteRanges []types.String `tfsdk:"mute_ranges"`
	Tags       *Tag           `tfsdk:"tags"`
	Start      types.Int64    `tfsdk:"start"`
	End        types.Int64    `tfsdk:"end"`
	Notify     *NotifyOptions `tfsdk:"notify"`
	Repeat     *RepeatOptions `tfsdk:"repeat"`
}

// GetId returns the ID of the resource.
func (m *muteResourceModel) GetId() string {
	return m.ID.ValueString()
}

// SetId sets the ID of the resource.
func (m *muteResourceModel) SetId(id string) {
	m.ID = types.StringValue(id)
}

// GetResourceType returns the type of the resource.
func (m *muteResourceModel) GetResourceType() string {
	return consts.TypeNameMute
}

// NotifyOptions maps the resource schema data.
type NotifyOptions struct {
	Targets []*NotifyTarget `tfsdk:"targets"`
	Message types.String    `tfsdk:"message"`
	Time    types.Int64     `tfsdk:"time"`
}

// NotifyTarget maps the resource schema data.
type NotifyTarget struct {
	Type types.String `tfsdk:"type"`
	To   types.String `tfsdk:"to"`
}

// RepeatCrontabSet maps the resource schema data.
type RepeatCrontabSet struct {
	Min   types.String `tfsdk:"min"`
	Hour  types.String `tfsdk:"hour"`
	Day   types.String `tfsdk:"day"`
	Month types.String `tfsdk:"month"`
	Week  types.String `tfsdk:"week"`
}

// RepeatOptions maps the resource schema data.
type RepeatOptions struct {
	Time            types.Int64       `tfsdk:"time"`
	Crontab         *RepeatCrontabSet `tfsdk:"crontab"`
	CrontabDuration types.Int64       `tfsdk:"crontab_duration"`
	Expire          types.Int64       `tfsdk:"expire"`
}

// Tag maps the resource schema data.
type Tag struct {
	Key   types.String `tfsdk:"key"`
	Value types.String `tfsdk:"value"`
}
