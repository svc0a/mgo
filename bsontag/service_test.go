package bsontag

import (
	"testing"
)

type Doc struct {
	ID int64 `bson:"_id"`
}

type Entity struct {
	Data string `bson:"data,omitempty"`
}

type Member struct {
	Doc    `bson:",inline"`
	Entity Entity `bson:"entity"`
	Name   string `bson:"name"`
}

var member = Register[Member]()

func TestGetBsonTag(t *testing.T) {
	type args struct {
		field any
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "1",
			args: args{
				field: &member.ID,
			},
			want: "_id",
		},
		{
			name: "2",
			args: args{
				field: &member.Entity.Data,
			},
			want: "entity.data",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Get(tt.args.field); got != tt.want {
				t.Errorf("GetBsonTag() = %v, want %v", got, tt.want)
			}
		})
	}
}
