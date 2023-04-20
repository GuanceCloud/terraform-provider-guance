package tfcodec

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
)

func TestDecode(t *testing.T) {
	type args struct {
		in  map[string]interface{}
		out interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    interface{}
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "ok",
			args: args{
				in: map[string]interface{}{
					"address": "123 Main St",
					"age":     30,
					"my_int":  42,
					"name":    "John",
					"score":   8.5,
					"structs": map[string]interface{}{
						"address": "123 Main St",
						"age":     30,
						"my_int":  42,
						"name":    "John",
						"score":   8.5,
					},
				},
				out: &MyStruct{},
			},
			want: &MyStruct{
				Name:    "John",
				Age:     30,
				Address: "123 Main St",
				Score:   8.5,
				MyInt:   types.Int64Value(42),
				Structs: &MyStruct{
					Name:    "John",
					Age:     30,
					Address: "123 Main St",
					Score:   8.5,
					MyInt:   types.Int64Value(42),
				},
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.wantErr(t, Decode(tt.args.in, tt.args.out))
			assert.Equal(t, tt.want, tt.args.out)
		})
	}
}
