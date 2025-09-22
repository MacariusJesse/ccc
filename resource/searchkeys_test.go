package resource

import (
	"testing"

	"github.com/cccteam/ccc/accesstypes"
	"github.com/google/go-cmp/cmp"
)

type testResource struct{}

func (t testResource) Resource() accesstypes.Resource {
	return accesstypes.Resource("test_resource")
}

func (t testResource) DefaultConfig() Config {
	return Config{DBType: SpannerDBType}
}

type testRequest struct {
	Name        string `substring:"name,title" fulltext:"content" ngram:"description"`
	Email       string `substring:"email"`
	Content     string `fulltext:"content"`
	Description string `ngram:"description"`
	NoTags      string
}

type testRequestPostgres struct {
	Name        string `substring:"name,title"`
	Email       string `substring:"email"`
	Content     string `fulltext:"content"`
	Description string `ngram:"description"`
}

func (t testRequestPostgres) Resource() accesstypes.Resource {
	return accesstypes.Resource("test_resource_postgres")
}

func (t testRequestPostgres) DefaultConfig() Config {
	return Config{DBType: PostgresDBType}
}

func TestNewSearchKeys(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		res  Resourcer
		want *SearchKeys
	}{
		{
			name: "Spanner resource with all search types",
			res:  testResource{},
			want: &SearchKeys{
				keys: map[SearchKey]SearchType{
					"name":        SubString,
					"title":       SubString,
					"email":       SubString,
					"content":     FullText,
					"description": Ngram,
				},
			},
		},
		{
			name: "PostgreSQL resource - no search types",
			res:  testRequestPostgres{},
			want: &SearchKeys{
				keys: map[SearchKey]SearchType{},
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := NewSearchKeys[testRequest](tt.res)
			if len(got.keys) != len(tt.want.keys) {
				t.Errorf("NewSearchKeys() keys length = %v, want %v", len(got.keys), len(tt.want.keys))
			}
			for key, searchType := range tt.want.keys {
				if got.keys[key] != searchType {
					t.Errorf("NewSearchKeys() keys[%v] = %v, want %v", key, got.keys[key], searchType)
				}
			}
		})
	}
}

func TestSplitSplitKeys(t *testing.T) {
	t.Parallel()

	type args struct {
		keys string
	}
	tests := []struct {
		name string
		args args
		want []SearchKey
	}{
		{
			name: "Single key",
			args: args{keys: "name"},
			want: []SearchKey{"name"},
		},
		{
			name: "Multiple keys",
			args: args{keys: "name,title,description"},
			want: []SearchKey{"name", "title", "description"},
		},
		{
			name: "Empty string",
			args: args{keys: ""},
			want: []SearchKey{""},
		},
		{
			name: "Keys with spaces",
			args: args{keys: "name, title, description"},
			want: []SearchKey{"name", " title", " description"},
		},
		{
			name: "Single key with comma",
			args: args{keys: "name,"},
			want: []SearchKey{"name", ""},
		},
		{
			name: "Empty keys",
			args: args{keys: ","},
			want: []SearchKey{"", ""},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := splitSplitKeys(tt.args.keys)
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("splitSplitKeys() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestSearchKeys_Integration(t *testing.T) {
	t.Parallel()

	type testStruct struct {
		Field1 string `substring:"field1,field1_alt"`
		Field2 string `fulltext:"field2"`
		Field3 string `ngram:"field3"`
		Field4 string `substring:"field4"`
		Field5 string
	}

	res := testResource{}

	searchKeys := NewSearchKeys[testStruct](res)

	expectedKeys := map[SearchKey]SearchType{
		"field1":     SubString,
		"field1_alt": SubString,
		"field2":     FullText,
		"field3":     Ngram,
		"field4":     SubString,
	}

	if len(searchKeys.keys) != len(expectedKeys) {
		t.Errorf("SearchKeys length = %v, want %v", len(searchKeys.keys), len(expectedKeys))
	}

	for key, searchType := range expectedKeys {
		if searchKeys.keys[key] != searchType {
			t.Errorf("SearchKeys[%v] = %v, want %v", key, searchKeys.keys[key], searchType)
		}
	}
}

func TestSearchKeys_EmptyTags(t *testing.T) {
	t.Parallel()

	type testStruct struct {
		Field1 string `substring:""`
		Field2 string `fulltext:""`
		Field3 string `ngram:""`
		Field4 string
	}

	res := testResource{}

	searchKeys := NewSearchKeys[testStruct](res)

	if len(searchKeys.keys) != 0 {
		t.Errorf("SearchKeys should be empty for empty tags, got %v keys", len(searchKeys.keys))
	}
}

func TestSearchKeys_MixedTags(t *testing.T) {
	t.Parallel()

	type testStruct struct {
		Field1 string `substring:"field1" fulltext:"field1" ngram:"field1"`
		Field2 string `substring:"field2,field2_alt" fulltext:"field2_content"`
		Field3 string `ngram:"field3"`
	}

	res := testResource{}

	searchKeys := NewSearchKeys[testStruct](res)

	expectedKeys := map[SearchKey]SearchType{
		"field1":         SubString,
		"field2":         SubString,
		"field2_alt":     SubString,
		"field2_content": FullText,
		"field3":         Ngram,
	}

	if len(searchKeys.keys) != len(expectedKeys) {
		t.Errorf("SearchKeys length = %v, want %v", len(searchKeys.keys), len(expectedKeys))
	}

	for key, searchType := range expectedKeys {
		if searchKeys.keys[key] != searchType {
			t.Errorf("SearchKeys[%v] = %v, want %v", key, searchKeys.keys[key], searchType)
		}
	}
}

func TestSearchKeys_PostgresIntegration(t *testing.T) {
	t.Parallel()

	type testStruct struct {
		Field1 string `substring:"field1"`
		Field2 string `fulltext:"field2"`
		Field3 string `ngram:"field3"`
	}

	res := testRequestPostgres{}

	searchKeys := NewSearchKeys[testStruct](res)

	if len(searchKeys.keys) != 0 {
		t.Errorf("SearchKeys should be empty for PostgreSQL, got %v keys", len(searchKeys.keys))
	}
}