package resource

import (
	"testing"

	"github.com/cccteam/ccc/accesstypes"
	"github.com/google/go-cmp/cmp"
)

func TestNewFieldSet(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		want *fieldSet
	}{
		{
			name: "New",
			want: &fieldSet{
				data:   make(map[accesstypes.Field]any),
				fields: nil,
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := newFieldSet()
			if diff := cmp.Diff(tt.want, got, cmp.AllowUnexported(fieldSet{})); diff != "" {
				t.Errorf("newFieldSet() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestFieldSet_Set(t *testing.T) {
	t.Parallel()

	type args struct {
		field accesstypes.Field
		value any
	}
	tests := []struct {
		name string
		args args
		want *fieldSet
	}{
		{
			name: "Set single field",
			args: args{
				field: "Name",
				value: "John",
			},
			want: &fieldSet{
				data: map[accesstypes.Field]any{
					"Name": "John",
				},
				fields: []accesstypes.Field{"Name"},
			},
		},
	}
	
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			fs := newFieldSet()
			fs.Set(tt.args.field, tt.args.value)
			if diff := cmp.Diff(tt.want, fs, cmp.AllowUnexported(fieldSet{})); diff != "" {
				t.Errorf("fieldSet.Set() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestFieldSet_SetMultiple(t *testing.T) {
	t.Parallel()

	fs := newFieldSet()
	fs.Set("Name", "John")
	fs.Set("Age", 30)

	want := &fieldSet{
		data: map[accesstypes.Field]any{
			"Name": "John",
			"Age":  30,
		},
		fields: []accesstypes.Field{"Name", "Age"},
	}

	if diff := cmp.Diff(want, fs, cmp.AllowUnexported(fieldSet{})); diff != "" {
		t.Errorf("fieldSet.Set() multiple fields mismatch (-want +got):\n%s", diff)
	}
}

func TestFieldSet_SetExisting(t *testing.T) {
	t.Parallel()

	fs := newFieldSet()
	fs.Set("Name", "John")
	fs.Set("Age", 30)
	fs.Set("Name", "Jane")

	want := &fieldSet{
		data: map[accesstypes.Field]any{
			"Name": "Jane",
			"Age":  30,
		},
		fields: []accesstypes.Field{"Name", "Age"},
	}

	if diff := cmp.Diff(want, fs, cmp.AllowUnexported(fieldSet{})); diff != "" {
		t.Errorf("fieldSet.Set() existing field mismatch (-want +got):\n%s", diff)
	}
}

func TestFieldSet_Get(t *testing.T) {
	t.Parallel()

	type args struct {
		field accesstypes.Field
	}
	tests := []struct {
		name string
		args args
		want any
	}{
		{
			name: "Get existing field",
			args: args{field: "Name"},
			want: "Jane",
		},
		{
			name: "Get another existing field",
			args: args{field: "Age"},
			want: 30,
		},
		{
			name: "Get non-existing field",
			args: args{field: "Email"},
			want: nil,
		},
	}

	fs := newFieldSet()
	fs.Set("Name", "Jane")
	fs.Set("Age", 30)

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := fs.Get(tt.args.field)
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("fieldSet.Get() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestFieldSet_IsSet(t *testing.T) {
	t.Parallel()

	type args struct {
		field accesstypes.Field
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Check existing field",
			args: args{field: "Name"},
			want: true,
		},
		{
			name: "Check another existing field",
			args: args{field: "Age"},
			want: true,
		},
		{
			name: "Check non-existing field",
			args: args{field: "Email"},
			want: false,
		},
	}

	fs := newFieldSet()
	fs.Set("Name", "Jane")
	fs.Set("Age", 30)

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := fs.IsSet(tt.args.field)
			if got != tt.want {
				t.Errorf("fieldSet.IsSet() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFieldSet_KeySet(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		want KeySet
	}{
		{
			name: "KeySet",
			want: KeySet{
				keyParts: []KeyPart{
					{Key: "Name", Value: "Jane"},
					{Key: "Age", Value: 30},
				},
			},
		},
	}

	fs := newFieldSet()
	fs.Set("Name", "Jane")
	fs.Set("Age", 30)

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := fs.KeySet()
			if diff := cmp.Diff(tt.want, got, cmp.AllowUnexported(KeySet{}, KeyPart{})); diff != "" {
				t.Errorf("fieldSet.KeySet() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}