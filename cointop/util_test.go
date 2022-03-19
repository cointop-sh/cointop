package cointop

import "testing"

func Test_getStructHash(t *testing.T) {
	type args struct {
		str1 interface{}
		str2 interface{}
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
		want    bool
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
				str2: &struct {
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
		{
			name: "error occurs at str1 when struct is nil",
			args: args{
				str1: nil,
				str2: struct {
					Name       string
					Age        int
					Properties struct {
						P7D  int
						P10D int
					}
				}{},
			},
			wantErr: true,
			want:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hash1, err1 := getStructHash(tt.args.str1)
			hash2, _ := getStructHash(tt.args.str2)

			if err1 != nil && !tt.wantErr {
				t.Errorf("getStructHash() error = %v, wantErr %v", err1, tt.wantErr)
				return
			}

			if cp := hash1 == hash2; cp != tt.want {
				t.Errorf("getStructHash() = %v, want %v", cp, tt.want)
			}
		})
	}
}
