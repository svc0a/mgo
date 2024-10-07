package tagx

import (
	"github.com/sirupsen/logrus"
	"testing"
)

func TestToLowerCase(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "test",
			args: args{
				s: "updatedAt",
			},
			want: "updatedAt",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ToLowerCase(tt.args.s)
			if got != tt.want {
				t.Errorf("ToLowerCase() = %v, want %v", got, tt.want)
			}
			logrus.Info(got)
		})
	}
}
