package tfcodec

import (
	"encoding/json"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
)

type MyStruct struct {
	Name    string        `tfsdk:"name"`
	Age     int           `tfsdk:"age"`
	Address string        `tfsdk:"address"`
	Score   float64       `tfsdk:"score"`
	MyInt   types.Int64   `tfsdk:"my_int"`
	MyArr   []types.Int64 `tfsdk:"my_arr"`
	Structs *MyStruct     `tfsdk:"structs"`
	Nested  []*MyStruct   `tfsdk:"nested"`
}

func TestEncode(t *testing.T) {
	type args struct {
		obj interface{}
	}
	tests := []struct {
		name string
		args args
		want map[string]interface{}
	}{
		{
			name: "ok",
			args: args{obj: MyStruct{
				Name:    "John",
				Age:     30,
				Address: "123 Main St",
				Score:   8.5,
				MyInt:   types.Int64Value(42),
				MyArr: []types.Int64{
					types.Int64Value(42),
				},
				Structs: &MyStruct{
					Name:    "John",
					Age:     30,
					Address: "123 Main St",
					Score:   8.5,
					MyInt:   types.Int64Value(42),
					MyArr: []types.Int64{
						types.Int64Value(42),
					},
				},
				Nested: []*MyStruct{
					{
						MyInt: types.Int64Value(42),
					},
				},
			}},
			want: map[string]interface{}{
				"address": "123 Main St",
				"age":     30,
				"name":    "John",
				"score":   8.5,
				"my_int":  42,
				"my_arr":  []interface{}{42},
				"structs": map[string]interface{}{
					"address": "123 Main St",
					"age":     30,
					"name":    "John",
					"score":   8.5,
					"my_int":  42,
					"my_arr":  []interface{}{42},
				},
				"nested": []map[string]interface{}{
					{"my_int": 42},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wantJSON, err := json.Marshal(tt.want)
			if !assert.NoError(t, err) {
				t.FailNow()
			}
			gotJSON, err := json.Marshal(Encode(tt.args.obj))
			if !assert.NoError(t, err) {
				t.FailNow()
			}
			assert.JSONEq(t, string(wantJSON), string(gotJSON))
		})
	}
}
