package monitor

import (
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/GuanceCloud/terraform-provider-guance/internal/consts"
)

// monitorResourceModel maps the resource schema data.
type monitorResourceModel struct {
	ID          types.String   `tfsdk:"id"`
	DashboardId types.String   `tfsdk:"dashboard_id"`
	Script      *MonitorScript `tfsdk:"script"`
}

// GetId returns the ID of the resource.
func (m *monitorResourceModel) GetId() string {
	return m.ID.ValueString()
}

// SetId sets the ID of the resource.
func (m *monitorResourceModel) SetId(id string) {
	m.ID = types.StringValue(id)
}

// GetResourceType returns the type of the resource.
func (m *monitorResourceModel) GetResourceType() string {
	return consts.TypeNameMonitor
}

// APMCheck maps the resource schema data.
type APMCheck struct {
	Name                   types.String `tfsdk:"name"`
	Title                  types.String `tfsdk:"title"`
	Message                types.String `tfsdk:"message"`
	Every                  types.String `tfsdk:"every"`
	Interval               types.Int64  `tfsdk:"interval"`
	RecoverNeedPeriodCount types.Int64  `tfsdk:"recover_need_period_count"`
	NoDataInterval         types.Int64  `tfsdk:"no_data_interval"`
	Targets                []*Target    `tfsdk:"targets"`
	Checker                *Checker     `tfsdk:"checker"`
}

// CheckFunc maps the resource schema data.
type CheckFunc struct {
	FuncId types.String `tfsdk:"func_id"`
	Kwargs types.String `tfsdk:"kwargs"`
}

// Checker maps the resource schema data.
type Checker struct {
	Rules []*Rule `tfsdk:"rules"`
}

// CloudDialCheck maps the resource schema data.
type CloudDialCheck struct {
	Name                   types.String `tfsdk:"name"`
	Title                  types.String `tfsdk:"title"`
	Message                types.String `tfsdk:"message"`
	Every                  types.String `tfsdk:"every"`
	Interval               types.Int64  `tfsdk:"interval"`
	RecoverNeedPeriodCount types.Int64  `tfsdk:"recover_need_period_count"`
	NoDataInterval         types.Int64  `tfsdk:"no_data_interval"`
	Targets                []*Target    `tfsdk:"targets"`
	Checker                *Checker     `tfsdk:"checker"`
}

// Condition maps the resource schema data.
type Condition struct {
	Alias    types.String   `tfsdk:"alias"`
	Operator types.String   `tfsdk:"operator"`
	Operands []types.String `tfsdk:"operands"`
}

// LoggingCheck maps the resource schema data.
type LoggingCheck struct {
	Name                   types.String `tfsdk:"name"`
	Title                  types.String `tfsdk:"title"`
	Message                types.String `tfsdk:"message"`
	Every                  types.String `tfsdk:"every"`
	Interval               types.Int64  `tfsdk:"interval"`
	RecoverNeedPeriodCount types.Int64  `tfsdk:"recover_need_period_count"`
	NoDataInterval         types.Int64  `tfsdk:"no_data_interval"`
	Targets                []*Target    `tfsdk:"targets"`
	Checker                *Checker     `tfsdk:"checker"`
}

// MonitorScript maps the resource schema data.
type MonitorScript struct {
	Type            types.String     `tfsdk:"type"`
	SimpleCheck     *SimpleCheck     `tfsdk:"simple_check"`
	SeniorCheck     *SeniorCheck     `tfsdk:"senior_check"`
	LoggingCheck    *LoggingCheck    `tfsdk:"logging_check"`
	MutationsCheck  *MutationsCheck  `tfsdk:"mutations_check"`
	WaterLevelCheck *WaterLevelCheck `tfsdk:"water_level_check"`
	RangeCheck      *RangeCheck      `tfsdk:"range_check"`
	SecurityCheck   *SecurityCheck   `tfsdk:"security_check"`
	ApmCheck        *APMCheck        `tfsdk:"apm_check"`
	RumCheck        *RUMCheck        `tfsdk:"rum_check"`
	ProcessCheck    *ProcessCheck    `tfsdk:"process_check"`
	CloudDialCheck  *CloudDialCheck  `tfsdk:"cloud_dial_check"`
}

// MutationsCheck maps the resource schema data.
type MutationsCheck struct {
	Name                   types.String `tfsdk:"name"`
	Title                  types.String `tfsdk:"title"`
	Message                types.String `tfsdk:"message"`
	Every                  types.String `tfsdk:"every"`
	Interval               types.Int64  `tfsdk:"interval"`
	RecoverNeedPeriodCount types.Int64  `tfsdk:"recover_need_period_count"`
	NoDataInterval         types.Int64  `tfsdk:"no_data_interval"`
	Targets                []*Target    `tfsdk:"targets"`
	Checker                *Checker     `tfsdk:"checker"`
}

// ProcessCheck maps the resource schema data.
type ProcessCheck struct {
	Name                   types.String `tfsdk:"name"`
	Title                  types.String `tfsdk:"title"`
	Message                types.String `tfsdk:"message"`
	Every                  types.String `tfsdk:"every"`
	Interval               types.Int64  `tfsdk:"interval"`
	RecoverNeedPeriodCount types.Int64  `tfsdk:"recover_need_period_count"`
	NoDataInterval         types.Int64  `tfsdk:"no_data_interval"`
	Targets                []*Target    `tfsdk:"targets"`
	Checker                *Checker     `tfsdk:"checker"`
}

// RUMCheck maps the resource schema data.
type RUMCheck struct {
	Name                   types.String `tfsdk:"name"`
	Title                  types.String `tfsdk:"title"`
	Message                types.String `tfsdk:"message"`
	Every                  types.String `tfsdk:"every"`
	Interval               types.Int64  `tfsdk:"interval"`
	RecoverNeedPeriodCount types.Int64  `tfsdk:"recover_need_period_count"`
	NoDataInterval         types.Int64  `tfsdk:"no_data_interval"`
	Targets                []*Target    `tfsdk:"targets"`
	Checker                *Checker     `tfsdk:"checker"`
}

// RangeCheck maps the resource schema data.
type RangeCheck struct {
	Name                   types.String `tfsdk:"name"`
	Title                  types.String `tfsdk:"title"`
	Message                types.String `tfsdk:"message"`
	Every                  types.String `tfsdk:"every"`
	Interval               types.Int64  `tfsdk:"interval"`
	RecoverNeedPeriodCount types.Int64  `tfsdk:"recover_need_period_count"`
	NoDataInterval         types.Int64  `tfsdk:"no_data_interval"`
	Targets                []*Target    `tfsdk:"targets"`
	Checker                *Checker     `tfsdk:"checker"`
}

// Rule maps the resource schema data.
type Rule struct {
	Conditions     []*Condition `tfsdk:"conditions"`
	ConditionLogic types.String `tfsdk:"condition_logic"`
	Status         types.String `tfsdk:"status"`
	Direction      types.String `tfsdk:"direction"`
	PeriodNum      types.Int64  `tfsdk:"period_num"`
	CheckPercent   types.Int64  `tfsdk:"check_percent"`
	CheckCount     types.Int64  `tfsdk:"check_count"`
	Strength       types.Int64  `tfsdk:"strength"`
}

// SecurityCheck maps the resource schema data.
type SecurityCheck struct {
	Name                   types.String `tfsdk:"name"`
	Title                  types.String `tfsdk:"title"`
	Message                types.String `tfsdk:"message"`
	Every                  types.String `tfsdk:"every"`
	Interval               types.Int64  `tfsdk:"interval"`
	RecoverNeedPeriodCount types.Int64  `tfsdk:"recover_need_period_count"`
	NoDataInterval         types.Int64  `tfsdk:"no_data_interval"`
	Targets                []*Target    `tfsdk:"targets"`
	Checker                *Checker     `tfsdk:"checker"`
}

// SeniorCheck maps the resource schema data.
type SeniorCheck struct {
	Name       types.String `tfsdk:"name"`
	Title      types.String `tfsdk:"title"`
	Message    types.String `tfsdk:"message"`
	Type       types.String `tfsdk:"type"`
	Every      types.String `tfsdk:"every"`
	CheckFuncs []*CheckFunc `tfsdk:"check_funcs"`
}

// SimpleCheck maps the resource schema data.
type SimpleCheck struct {
	Name                   types.String `tfsdk:"name"`
	Title                  types.String `tfsdk:"title"`
	Message                types.String `tfsdk:"message"`
	Every                  types.String `tfsdk:"every"`
	Interval               types.Int64  `tfsdk:"interval"`
	RecoverNeedPeriodCount types.Int64  `tfsdk:"recover_need_period_count"`
	NoDataInterval         types.Int64  `tfsdk:"no_data_interval"`
	Targets                []*Target    `tfsdk:"targets"`
	Checker                *Checker     `tfsdk:"checker"`
}

// Target maps the resource schema data.
type Target struct {
	Alias types.String `tfsdk:"alias"`
	Dql   types.String `tfsdk:"dql"`
}

// WaterLevelCheck maps the resource schema data.
type WaterLevelCheck struct {
	Name                   types.String `tfsdk:"name"`
	Title                  types.String `tfsdk:"title"`
	Message                types.String `tfsdk:"message"`
	Every                  types.String `tfsdk:"every"`
	Interval               types.Int64  `tfsdk:"interval"`
	RecoverNeedPeriodCount types.Int64  `tfsdk:"recover_need_period_count"`
	NoDataInterval         types.Int64  `tfsdk:"no_data_interval"`
	Targets                []*Target    `tfsdk:"targets"`
	Checker                *Checker     `tfsdk:"checker"`
}
