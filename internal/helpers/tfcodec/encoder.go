package tfcodec

import (
	"encoding/json"
	"reflect"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	stringType  = reflect.TypeOf(types.StringValue(""))
	int64Type   = reflect.TypeOf(types.Int64Value(0))
	boolType    = reflect.TypeOf(types.BoolValue(false))
	float64Type = reflect.TypeOf(types.Float64Value(0.0))
)

func Encode(obj interface{}) map[string]interface{} {
	objValue := reflect.ValueOf(obj)
	objType := objValue.Type()

	if objType.Kind() == reflect.Ptr {
		objValue = objValue.Elem()
		objType = objValue.Type()
	}

	data := make(map[string]interface{})

	for i := 0; i < objType.NumField(); i++ {
		fieldType := objType.Field(i)
		fieldValue := objValue.Field(i)

		// skip unexported field
		if !fieldType.IsExported() {
			continue
		}

		fieldValue, ok := resolvePointer(fieldValue)
		if !ok {
			continue
		}

		if fieldValue.IsZero() {
			continue
		}

		tagName := fieldType.Tag.Get("tfsdk")

		switch fieldValue.Type().Kind() {
		case reflect.Struct:
			data[tagName] = encodeStruct(fieldValue)
		case reflect.Array, reflect.Slice:
			var slice []interface{}
			for i := 0; i < fieldValue.Len(); i++ {
				itemValue := fieldValue.Index(i)
				itemValue, ok := resolvePointer(itemValue)
				if !ok {
					continue
				}
				if itemValue.IsZero() {
					continue
				}
				switch itemValue.Type().Kind() {
				case reflect.Struct:
					slice = append(slice, encodeStruct(itemValue))
				default:
					slice = append(slice, itemValue.Interface())
				}
			}
			data[tagName] = slice
		default:
			data[tagName] = fieldValue.Interface()
		}
	}
	return data
}

func resolvePointer(fieldValue reflect.Value) (reflect.Value, bool) {
	for fieldValue.Kind() == reflect.Ptr {
		if fieldValue.IsNil() {
			break
		}
		fieldValue = fieldValue.Elem()
	}

	// skip nil pointer
	if fieldValue.Kind() == reflect.Ptr && fieldValue.IsNil() {
		return reflect.Value{}, false
	}
	return fieldValue, true
}

func encodeStruct(fieldValue reflect.Value) interface{} {
	switch {
	case reflect.DeepEqual(fieldValue.Type(), int64Type):
		return fieldValue.Interface().(types.Int64).ValueInt64()
	case reflect.DeepEqual(fieldValue.Type(), boolType):
		return fieldValue.Interface().(types.Bool).ValueBool()
	case reflect.DeepEqual(fieldValue.Type(), float64Type):
		return fieldValue.Interface().(types.Float64).ValueFloat64()
	case reflect.DeepEqual(fieldValue.Type(), stringType):
		return fieldValue.Interface().(types.String).ValueString()
	default:
		return Encode(fieldValue.Interface())
	}
}

func EncodeJSON(in interface{}) ([]byte, error) {
	m := Encode(in)
	return json.Marshal(m)
}
