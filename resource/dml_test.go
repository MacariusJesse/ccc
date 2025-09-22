package resource

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestConfig_SetDBType(t *testing.T) {
	t.Parallel()

	type args struct {
		dbType DBType
	}
	tests := []struct {
		name string
		c    Config
		args args
		want Config
	}{
		{
			name: "Set SpannerDBType",
			c:    Config{},
			args: args{dbType: SpannerDBType},
			want: Config{DBType: SpannerDBType},
		},
		{
			name: "Set PostgresDBType",
			c:    Config{},
			args: args{dbType: PostgresDBType},
			want: Config{DBType: PostgresDBType},
		},
		{
			name: "Override existing DBType",
			c:    Config{DBType: SpannerDBType},
			args: args{dbType: PostgresDBType},
			want: Config{DBType: PostgresDBType},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := tt.c.SetDBType(tt.args.dbType)
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("Config.SetDBType() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestConfig_SetChangeTrackingTable(t *testing.T) {
	t.Parallel()

	type args struct {
		changeTrackingTable string
	}
	tests := []struct {
		name string
		c    Config
		args args
		want Config
	}{
		{
			name: "Set change tracking table",
			c:    Config{},
			args: args{changeTrackingTable: "change_tracking"},
			want: Config{ChangeTrackingTable: "change_tracking"},
		},
		{
			name: "Set empty change tracking table",
			c:    Config{},
			args: args{changeTrackingTable: ""},
			want: Config{ChangeTrackingTable: ""},
		},
		{
			name: "Override existing change tracking table",
			c:    Config{ChangeTrackingTable: "old_table"},
			args: args{changeTrackingTable: "new_table"},
			want: Config{ChangeTrackingTable: "new_table"},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := tt.c.SetChangeTrackingTable(tt.args.changeTrackingTable)
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("Config.SetChangeTrackingTable() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestConfig_SetTrackChanges(t *testing.T) {
	t.Parallel()

	type args struct {
		trackChanges bool
	}
	tests := []struct {
		name string
		c    Config
		args args
		want Config
	}{
		{
			name: "Set track changes to true",
			c:    Config{},
			args: args{trackChanges: true},
			want: Config{TrackChanges: true},
		},
		{
			name: "Set track changes to false",
			c:    Config{},
			args: args{trackChanges: false},
			want: Config{TrackChanges: false},
		},
		{
			name: "Override existing track changes",
			c:    Config{TrackChanges: true},
			args: args{trackChanges: false},
			want: Config{TrackChanges: false},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := tt.c.SetTrackChanges(tt.args.trackChanges)
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("Config.SetTrackChanges() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestDBType(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		db   DBType
		want string
	}{
		{
			name: "SpannerDBType string",
			db:   SpannerDBType,
			want: "spanner",
		},
		{
			name: "PostgresDBType string",
			db:   PostgresDBType,
			want: "postgres",
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := string(tt.db)
			if got != tt.want {
				t.Errorf("DBType.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStatement(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		stmt Statement
	}{
		{
			name: "Empty statement",
			stmt: Statement{},
		},
		{
			name: "Statement with SQL only",
			stmt: Statement{
				SQL: "SELECT * FROM users",
			},
		},
		{
			name: "Statement with Spanner params",
			stmt: Statement{
				SQL: "SELECT * FROM users WHERE id = @id",
				SpannerParams: map[string]any{
					"id": "user123",
				},
			},
		},
		{
			name: "Statement with PostgreSQL params",
			stmt: Statement{
				SQL: "SELECT * FROM users WHERE id = $1",
				PostgreSQLParams: []any{"user123"},
			},
		},
		{
			name: "Statement with both param types",
			stmt: Statement{
				SQL: "SELECT * FROM users WHERE id = @id OR id = $1",
				SpannerParams: map[string]any{
					"id": "user123",
				},
				PostgreSQLParams: []any{"user456"},
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if tt.stmt.SQL == "" && len(tt.stmt.SpannerParams) == 0 && len(tt.stmt.PostgreSQLParams) == 0 {
				return
			}
			if tt.stmt.SQL == "" {
				t.Error("Statement should have SQL")
			}
		})
	}
}