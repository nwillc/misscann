package main

import (
	"github.com/nwillc/genfuncs/container"
	"testing"
)

func Test_allowable(t *testing.T) {
	type args struct {
		group container.GSlice[person]
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "none",
			args: args{
				group: container.GSlice[person]{},
			},
			want: true,
		},
		{
			name: "missionaries",
			args: args{
				group: container.GSlice[person]{MISSIONARY, MISSIONARY},
			},
			want: true,
		},
		{
			name: "cannibals",
			args: args{
				group: container.GSlice[person]{CANNIBAL, CANNIBAL},
			},
			want: true,
		},
		{
			name: "more missionaries",
			args: args{
				group: container.GSlice[person]{MISSIONARY, MISSIONARY, CANNIBAL},
			},
			want: true,
		},
		{
			name: "more missionaries",
			args: args{
				group: container.GSlice[person]{MISSIONARY, MISSIONARY, CANNIBAL},
			},
			want: true,
		},
		{
			name: "more cannibals",
			args: args{
				group: container.GSlice[person]{MISSIONARY, CANNIBAL, CANNIBAL},
			},
			want: false,
		},
		{
			name: "equal",
			args: args{
				group: container.GSlice[person]{MISSIONARY, CANNIBAL},
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := allowable(tt.args.group); got != tt.want {
				t.Errorf("allowable() = %v, want %v", got, tt.want)
			}
		})
	}
}
