package tfcodec

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/mitchellh/mapstructure"
)

var jsonNumberType = reflect.ValueOf(json.Number("42")).Kind()

// StringToTimeHookFunc returns a DecodeHookFunc that converts
// strings to time.Time.
func tfDecoderHookFunc() mapstructure.DecodeHookFunc {
	return func(
		f reflect.Type,
		t reflect.Type,
		data interface{}) (interface{}, error) {
		switch {
		case t == stringType:
			if f.Kind() != reflect.String {
				return nil, fmt.Errorf("expect type %q, got %q", t, f)
			}
			return types.StringValue(data.(string)), nil
		case t == boolType:
			if f.Kind() != reflect.Bool {
				return nil, fmt.Errorf("expect type %q, got %q", t, f)
			}
			return types.BoolValue(data.(bool)), nil
		case t == int64Type:
			switch f.Kind() {
			case jsonNumberType, reflect.Int64, reflect.Int32, reflect.Int16, reflect.Int8, reflect.Int:
				v, _ := strconv.ParseInt(fmt.Sprint(data), 10, 64)
				return types.Int64Value(v), nil
			default:
				return nil, fmt.Errorf("expect type %q, got %q", t, f)
			}
		case t == float64Type:
			switch f.Kind() {
			case jsonNumberType, reflect.Float64, reflect.Float32:
				v, _ := strconv.ParseFloat(fmt.Sprint(data), 64)
				return types.Float64Value(v), nil
			default:
				return nil, fmt.Errorf("expect type %q, got %q", t, f)
			}
		default:
			return data, nil
		}
	}
}

func Decode(in map[string]interface{}, out interface{}) error {
	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		DecodeHook:       tfDecoderHookFunc(),
		Result:           out,
		TagName:          "tfsdk",
		WeaklyTypedInput: true,
	})
	if err != nil {
		return fmt.Errorf("new decoder failed: %w", err)
	}
	if err := decoder.Decode(in); err != nil {
		return fmt.Errorf("decode map failed: %w", err)
	}
	return decoder.Decode(in)
}

func DecodeJSON(body []byte, out interface{}) error {
	var props map[string]interface{}
	decoder := json.NewDecoder(bytes.NewReader(body))
	decoder.UseNumber()
	if err := decoder.Decode(&props); err != nil {
		return fmt.Errorf("failed to unmarshal properties: %w", err)
	}
	if err := Decode(props, out); err != nil {
		return fmt.Errorf("failed to decode properties: %w", err)
	}
	return nil
}
