package mute

import (
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/GuanceCloud/terraform-provider-guance/internal/api"
)

type MuteRange = muteRange
type NotifyTarget = notifyTarget
type RepeatCrontabSet = repeatCrontabSet

func StringPointerValue(value string) types.String {
	return stringPointerValue(value)
}

func MuteRangesFromContent(values []api.MuteRange, existing []MuteRange, present bool) []MuteRange {
	return muteRangesFromContent(values, existing, present)
}

func NotifyTargetsFromContent(values []api.MuteNotifyTarget, existing []NotifyTarget, present bool) []NotifyTarget {
	return notifyTargetsFromContent(values, existing, present)
}

func RepeatCrontabSetFromContent(value *api.RepeatCrontabSet) *RepeatCrontabSet {
	return repeatCrontabSetFromContent(value)
}

func RepeatTimeSetFromContent(content *api.MuteContent) int {
	return repeatTimeSetFromContent(content)
}

func RepeatExpireTimeValueOrExisting(value string, existing types.String) types.String {
	return repeatExpireTimeValueOrExisting(value, existing)
}

func DeclarationFromContent(values map[string]any) map[string]string {
	return declarationFromContent(values)
}

func MuteEnabledFromStatus(status int, existing types.Bool) types.Bool {
	return muteEnabledFromStatus(status, existing)
}
