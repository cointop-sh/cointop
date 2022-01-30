package cointop

import "testing"

func Test_getStructHash(t *testing.T) {
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
				str1: struct {
					Name       string
					Properties struct {
						P7D  int
						P10D int
					}
				}{},
				str2: struct {
					Name       string
					Properties struct {
						P7D  int
						P10D int
					}
				}{},
			},
			want: true,
		},
		{
			name: "different structs but have similar fields and different field type",
			args: args{
				str1: struct {
					Name       string
					Properties struct {
						P7D  int
						P10D int
					}
				}{},
				str2: struct {
					Name       rune
					Properties struct {
						P7D  int
						P10D int
					}
				}{},
			},
			want: false,
		},
		{
			name: "different structs and different fields",
			args: args{
				str1: struct {
					Name       string
					Properties struct {
						P7D  int
						P10D int
					}
				}{},
				str2: struct {
					Name       string
					Age        int
					Properties struct {
						P7D  int
						P10D int
					}
				}{},
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
