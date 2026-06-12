package tfconvert

import (
	"encoding/json"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

func StringsToTypes(values []string) []types.String {
	result := make([]types.String, 0, len(values))
	for _, value := range values {
		result = append(result, types.StringValue(value))
	}
	return result
}

func TypesToStrings(values []types.String) []string {
	result := make([]string, 0, len(values))
	for _, value := range values {
		if value.IsNull() || value.IsUnknown() {
			continue
		}
		result = append(result, value.ValueString())
	}
	return result
}

func StringValueOrNull(value string) types.String {
	if value == "" {
		return types.StringNull()
	}
	return types.StringValue(value)
}

func CanonicalJSONFromString(value string) (any, string, error) {
	var decoded any
	if err := json.Unmarshal([]byte(value), &decoded); err != nil {
		return nil, "", err
	}
	canonical, err := CanonicalJSONFromValue(decoded)
	if err != nil {
		return nil, "", err
	}
	return decoded, canonical, nil
}

func CanonicalJSONFromValue(value any) (string, error) {
	encoded, err := json.Marshal(value)
	if err != nil {
		return "", err
	}
	return string(encoded), nil
}
