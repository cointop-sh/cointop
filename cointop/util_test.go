package cointop

import "testing"

func Test_getStructHash(t *testing.T) {
	type SCoin struct {
		Name       string
		Properties struct {
			P7D  int
			P10D int
		}
	}
	type SCoin1 struct {
		Name       string
		Properties struct {
			P7D  int
			P10D int
		}
	}
	type SCoin2 struct {
		Name       interface{}
		Properties struct {
			P7D  int
			P10D int
		}
	}
	type SCoin3 struct {
		Name       string
		Age        int
		Properties *struct {
			P7D  int
			P10D int
		}
	}
	type args struct {
		str1 interface{}
		str2 interface{}
	}

	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "the same structs",
			args: args{
				str1: SCoin{},
				str2: SCoin{},
			},
			want: true,
		},
		{
			name: "different structs but have similar fields",
			args: args{
				str1: SCoin{},
				str2: SCoin1{},
			},
			want: true,
		},
		{
			name: "different structs but have similar fields and different field type",
			args: args{
				str1: SCoin{},
				str2: SCoin2{},
			},
			want: false,
		},
		{
			name: "different structs and different fields",
			args: args{
				str1: SCoin{},
				str2: SCoin3{},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cp := getStructHash(tt.args.str1) == getStructHash(tt.args.str2)
			if cp != tt.want {
				t.Errorf("getStructHash() = %v, want %v", cp, tt.want)
			}
		})
	}
}
