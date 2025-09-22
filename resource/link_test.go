package resource

import (
	"encoding/json"
	"testing"

	"github.com/cccteam/ccc"
	"github.com/google/go-cmp/cmp"
)

func TestLink_EncodeSpanner(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		l       Link
		want    any
		wantErr bool
	}{
		{
			name: "Valid link",
			l: Link{
				ID:       ccc.MustParseUUID("550e8400-e29b-41d4-a716-446655440000"),
				Resource: "users",
				Text:     "User Profile",
			},
			want:    []byte(`{"id":"550e8400-e29b-41d4-a716-446655440000","resource":"users","text":"User Profile"}`),
			wantErr: false,
		},
		{
			name: "Empty link",
			l: Link{
				ID:       ccc.UUID{},
				Resource: "",
				Text:     "",
			},
			want:    []byte(`{"id":"00000000-0000-0000-0000-000000000000","resource":"","text":""}`),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := tt.l.EncodeSpanner()
			if (err != nil) != tt.wantErr {
				t.Errorf("Link.EncodeSpanner() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !cmp.Equal(got, tt.want) {
				t.Errorf("Link.EncodeSpanner() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLink_DecodeSpanner(t *testing.T) {
	t.Parallel()

	type args struct {
		val any
	}
	tests := []struct {
		name    string
		l       *Link
		args    args
		want    Link
		wantErr bool
	}{
		{
			name: "Decode valid JSON string",
			l:    &Link{},
			args: args{
				val: `{"id":"550e8400-e29b-41d4-a716-446655440000","resource":"users","text":"User Profile"}`,
			},
			want: Link{
				ID:       ccc.MustParseUUID("550e8400-e29b-41d4-a716-446655440000"),
				Resource: "users",
				Text:     "User Profile",
			},
			wantErr: false,
		},
		{
			name: "Decode empty JSON",
			l:    &Link{},
			args: args{
				val: `{"id":"00000000-0000-0000-0000-000000000000","resource":"","text":""}`,
			},
			want: Link{
				ID:       ccc.UUID{},
				Resource: "",
				Text:     "",
			},
			wantErr: false,
		},
		{
			name: "Invalid type - int",
			l:    &Link{},
			args: args{
				val: 123,
			},
			want:    Link{},
			wantErr: true,
		},
		{
			name: "Invalid type - nil",
			l:    &Link{},
			args: args{
				val: nil,
			},
			want:    Link{},
			wantErr: true,
		},
		{
			name: "Invalid JSON string",
			l:    &Link{},
			args: args{
				val: "invalid json",
			},
			want:    Link{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := tt.l.DecodeSpanner(tt.args.val)
			if (err != nil) != tt.wantErr {
				t.Errorf("Link.DecodeSpanner() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if diff := cmp.Diff(tt.want, *tt.l); diff != "" {
					t.Errorf("Link.DecodeSpanner() mismatch (-want +got):\n%s", diff)
				}
			}
		})
	}
}

func TestLink_MarshalJSON(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		l       Link
		want    []byte
		wantErr bool
	}{
		{
			name: "Valid link",
			l: Link{
				ID:       ccc.MustParseUUID("550e8400-e29b-41d4-a716-446655440000"),
				Resource: "users",
				Text:     "User Profile",
			},
			want:    []byte(`{"id":"550e8400-e29b-41d4-a716-446655440000","resource":"users","text":"User Profile"}`),
			wantErr: false,
		},
		{
			name: "Empty link",
			l: Link{
				ID:       ccc.UUID{},
				Resource: "",
				Text:     "",
			},
			want:    []byte(`{"id":"00000000-0000-0000-0000-000000000000","resource":"","text":""}`),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := tt.l.MarshalJSON()
			if (err != nil) != tt.wantErr {
				t.Errorf("Link.MarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !cmp.Equal(got, tt.want) {
				t.Errorf("Link.MarshalJSON() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLink_UnmarshalJSON(t *testing.T) {
	t.Parallel()

	type args struct {
		data []byte
	}
	tests := []struct {
		name    string
		l       *Link
		args    args
		want    Link
		wantErr bool
	}{
		{
			name: "Valid JSON",
			l:    &Link{},
			args: args{
				data: []byte(`{"id":"550e8400-e29b-41d4-a716-446655440000","resource":"users","text":"User Profile"}`),
			},
			want: Link{
				ID:       ccc.MustParseUUID("550e8400-e29b-41d4-a716-446655440000"),
				Resource: "users",
				Text:     "User Profile",
			},
			wantErr: false,
		},
		{
			name: "Null data",
			l:    &Link{},
			args: args{
				data: nil,
			},
			want:    Link{},
			wantErr: false,
		},
		{
			name: "Null string",
			l:    &Link{},
			args: args{
				data: []byte("null"),
			},
			want:    Link{},
			wantErr: false,
		},
		{
			name: "Empty JSON object",
			l:    &Link{},
			args: args{
				data: []byte("{}"),
			},
			want: Link{
				ID:       ccc.UUID{},
				Resource: "",
				Text:     "",
			},
			wantErr: false,
		},
		{
			name: "Invalid JSON",
			l:    &Link{},
			args: args{
				data: []byte(`{"id":"invalid"}`),
			},
			want:    Link{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := tt.l.UnmarshalJSON(tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("Link.UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if diff := cmp.Diff(tt.want, *tt.l); diff != "" {
					t.Errorf("Link.UnmarshalJSON() mismatch (-want +got):\n%s", diff)
				}
			}
		})
	}
}

func TestLink_IsNull(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		l    Link
		want bool
	}{
		{
			name: "Null link",
			l: Link{
				ID:       ccc.UUID{},
				Resource: "",
				Text:     "",
			},
			want: true,
		},
		{
			name: "Non-null link",
			l: Link{
				ID:       ccc.MustParseUUID("550e8400-e29b-41d4-a716-446655440000"),
				Resource: "users",
				Text:     "User Profile",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := tt.l.IsNull()
			if got != tt.want {
				t.Errorf("Link.IsNull() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNullLink_EncodeSpanner(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		nl      NullLink
		want    any
		wantErr bool
	}{
		{
			name: "Valid null link",
			nl: NullLink{
				Link: Link{
					ID:       ccc.MustParseUUID("550e8400-e29b-41d4-a716-446655440000"),
					Resource: "users",
					Text:     "User Profile",
				},
				Valid: true,
			},
			want:    []byte(`{"id":"550e8400-e29b-41d4-a716-446655440000","resource":"users","text":"User Profile"}`),
			wantErr: false,
		},
		{
			name: "Invalid null link",
			nl: NullLink{
				Link:  Link{},
				Valid: false,
			},
			want:    nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := tt.nl.EncodeSpanner()
			if (err != nil) != tt.wantErr {
				t.Errorf("NullLink.EncodeSpanner() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !cmp.Equal(got, tt.want) {
				t.Errorf("NullLink.EncodeSpanner() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNullLink_DecodeSpanner(t *testing.T) {
	t.Parallel()

	type args struct {
		val any
	}
	tests := []struct {
		name    string
		nl      *NullLink
		args    args
		want    NullLink
		wantErr bool
	}{
		{
			name: "Decode valid string",
			nl:   &NullLink{},
			args: args{
				val: `{"id":"550e8400-e29b-41d4-a716-446655440000","resource":"users","text":"User Profile"}`,
			},
			want: NullLink{
				Link: Link{
					ID:       ccc.MustParseUUID("550e8400-e29b-41d4-a716-446655440000"),
					Resource: "users",
					Text:     "User Profile",
				},
				Valid: true,
			},
			wantErr: false,
		},
		{
			name: "Decode nil string pointer",
			nl:   &NullLink{},
			args: args{
				val: (*string)(nil),
			},
			want: NullLink{
				Link:  Link{},
				Valid: false,
			},
			wantErr: false,
		},
		{
			name: "Decode nil value",
			nl:   &NullLink{},
			args: args{
				val: nil,
			},
			want: NullLink{
				Link:  Link{},
				Valid: false,
			},
			wantErr: false,
		},
		{
			name: "Invalid type",
			nl:   &NullLink{},
			args: args{
				val: 123,
			},
			want:    NullLink{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := tt.nl.DecodeSpanner(tt.args.val)
			if (err != nil) != tt.wantErr {
				t.Errorf("NullLink.DecodeSpanner() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if diff := cmp.Diff(tt.want, *tt.nl); diff != "" {
					t.Errorf("NullLink.DecodeSpanner() mismatch (-want +got):\n%s", diff)
				}
			}
		})
	}
}

func TestNullLink_MarshalJSON(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		nl      NullLink
		want    []byte
		wantErr bool
	}{
		{
			name: "Valid null link",
			nl: NullLink{
				Link: Link{
					ID:       ccc.MustParseUUID("550e8400-e29b-41d4-a716-446655440000"),
					Resource: "users",
					Text:     "User Profile",
				},
				Valid: true,
			},
			want:    []byte(`{"id":"550e8400-e29b-41d4-a716-446655440000","resource":"users","text":"User Profile"}`),
			wantErr: false,
		},
		{
			name: "Invalid null link",
			nl: NullLink{
				Link:  Link{},
				Valid: false,
			},
			want:    []byte("null"),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := tt.nl.MarshalJSON()
			if (err != nil) != tt.wantErr {
				t.Errorf("NullLink.MarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !cmp.Equal(got, tt.want) {
				t.Errorf("NullLink.MarshalJSON() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNullLink_UnmarshalJSON(t *testing.T) {
	t.Parallel()

	type args struct {
		data []byte
	}
	tests := []struct {
		name    string
		nl      *NullLink
		args    args
		want    NullLink
		wantErr bool
	}{
		{
			name: "Valid JSON",
			nl:   &NullLink{},
			args: args{
				data: []byte(`{"id":"550e8400-e29b-41d4-a716-446655440000","resource":"users","text":"User Profile"}`),
			},
			want: NullLink{
				Link: Link{
					ID:       ccc.MustParseUUID("550e8400-e29b-41d4-a716-446655440000"),
					Resource: "users",
					Text:     "User Profile",
				},
				Valid: true,
			},
			wantErr: false,
		},
		{
			name: "Null data",
			nl:   &NullLink{},
			args: args{
				data: nil,
			},
			want: NullLink{
				Link:  Link{},
				Valid: false,
			},
			wantErr: false,
		},
		{
			name: "Null string",
			nl:   &NullLink{},
			args: args{
				data: []byte("null"),
			},
			want: NullLink{
				Link:  Link{},
				Valid: false,
			},
			wantErr: false,
		},
		{
			name: "Invalid JSON",
			nl:   &NullLink{},
			args: args{
				data: []byte(`{"id":"invalid"}`),
			},
			want:    NullLink{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := tt.nl.UnmarshalJSON(tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("NullLink.UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if diff := cmp.Diff(tt.want, *tt.nl); diff != "" {
					t.Errorf("NullLink.UnmarshalJSON() mismatch (-want +got):\n%s", diff)
				}
			}
		})
	}
}