package resource

import (
	"context"
	"testing"
)

func TestUserEvent(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		ctx  context.Context
		wantPanic bool
	}{
		{
			name: "Context without session info",
			ctx:  context.Background(),
			wantPanic: true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if tt.wantPanic {
				defer func() {
					if r := recover(); r == nil {
						t.Error("Expected panic but got none")
					}
				}()
				UserEvent(tt.ctx)
			}
		})
	}
}

func TestProcessEvent(t *testing.T) {
	t.Parallel()

	type args struct {
		processName string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Process with name",
			args: args{processName: "data_export"},
			want: "Process data_export",
		},
		{
			name: "Process with empty name",
			args: args{processName: ""},
			want: "Process ",
		},
		{
			name: "Process with special characters",
			args: args{processName: "process-with_special.chars"},
			want: "Process process-with_special.chars",
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := ProcessEvent(tt.args.processName)
			if got != tt.want {
				t.Errorf("ProcessEvent() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUserProcessEvent(t *testing.T) {
	t.Parallel()

	type args struct {
		ctx         context.Context
		processName string
	}
	tests := []struct {
		name string
		args args
		wantPanic bool
	}{
		{
			name: "Context without session info and process",
			args: args{
				ctx:         context.Background(),
				processName: "cleanup",
			},
			wantPanic: true,
		},
		{
			name: "Context without session info and empty process",
			args: args{
				ctx:         context.Background(),
				processName: "",
			},
			wantPanic: true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if tt.wantPanic {
				defer func() {
					if r := recover(); r == nil {
						t.Error("Expected panic but got none")
					}
				}()
				UserProcessEvent(tt.args.ctx, tt.args.processName)
			}
		})
	}
}