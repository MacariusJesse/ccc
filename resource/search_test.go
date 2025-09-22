package resource

import (
	"testing"

	"github.com/go-playground/errors/v5"
)

func TestSearchKey_String(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		f    SearchKey
		want string
	}{
		{
			name: "Simple search key",
			f:    SearchKey("name"),
			want: "name",
		},
		{
			name: "Empty search key",
			f:    SearchKey(""),
			want: "",
		},
		{
			name: "Search key with special characters",
			f:    SearchKey("user_email"),
			want: "user_email",
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := tt.f.String()
			if got != tt.want {
				t.Errorf("SearchKey.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSearchType(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		st   SearchType
		want string
	}{
		{
			name: "SubString type",
			st:   SubString,
			want: "substring",
		},
		{
			name: "FullText type",
			st:   FullText,
			want: "fulltext",
		},
		{
			name: "Ngram type",
			st:   Ngram,
			want: "ngram",
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := string(tt.st)
			if got != tt.want {
				t.Errorf("SearchType.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewSearch(t *testing.T) {
	t.Parallel()

	type args struct {
		typ    SearchType
		values map[SearchKey]string
	}
	tests := []struct {
		name string
		args args
		want *Search
	}{
		{
			name: "New substring search",
			args: args{
				typ: SubString,
				values: map[SearchKey]string{
					"name": "John Doe",
				},
			},
			want: &Search{
				typ: SubString,
				values: map[SearchKey]string{
					"name": "John Doe",
				},
			},
		},
		{
			name: "New fulltext search",
			args: args{
				typ: FullText,
				values: map[SearchKey]string{
					"content": "search text",
				},
			},
			want: &Search{
				typ: FullText,
				values: map[SearchKey]string{
					"content": "search text",
				},
			},
		},
		{
			name: "New ngram search",
			args: args{
				typ: Ngram,
				values: map[SearchKey]string{
					"description": "ngram search",
				},
			},
			want: &Search{
				typ: Ngram,
				values: map[SearchKey]string{
					"description": "ngram search",
				},
			},
		},
		{
			name: "Empty values",
			args: args{
				typ:    SubString,
				values: map[SearchKey]string{},
			},
			want: &Search{
				typ:    SubString,
				values: map[SearchKey]string{},
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := NewSearch(tt.args.typ, tt.args.values)
			if got.typ != tt.want.typ {
				t.Errorf("NewSearch().typ = %v, want %v", got.typ, tt.want.typ)
			}
			if len(got.values) != len(tt.want.values) {
				t.Errorf("NewSearch().values length = %v, want %v", len(got.values), len(tt.want.values))
			}
			for key, value := range tt.want.values {
				if got.values[key] != value {
					t.Errorf("NewSearch().values[%v] = %v, want %v", key, got.values[key], value)
				}
			}
		})
	}
}

func TestSearch_SpannerStmt(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		s       Search
		want    Statement
		wantErr bool
	}{
		{
			name: "SubString search with single term",
			s: Search{
				typ: SubString,
				values: map[SearchKey]string{
					"name": "John",
				},
			},
			want: Statement{
				SQL: "WHERE SEARCH_SUBSTRING(name, @searchsubstringterm0)\nORDER BY SCORE_NGRAMS(name, @ngramscoreterm0)",
				SpannerParams: map[string]any{
					"searchsubstringterm0": "John",
					"ngramscoreterm0":      "John",
				},
			},
			wantErr: false,
		},
		{
			name: "SubString search with multiple terms",
			s: Search{
				typ: SubString,
				values: map[SearchKey]string{
					"content": "John Doe",
				},
			},
			want: Statement{
				SQL: "WHERE SEARCH_SUBSTRING(content, @searchsubstringterm0) OR SEARCH_SUBSTRING(content, @searchsubstringterm1)\nORDER BY SCORE_NGRAMS(content, @ngramscoreterm0) + SCORE_NGRAMS(content, @ngramscoreterm1)",
				SpannerParams: map[string]any{
					"searchsubstringterm0": "John",
					"searchsubstringterm1": "Doe",
					"ngramscoreterm0":      "John",
					"ngramscoreterm1":      "Doe",
				},
			},
			wantErr: false,
		},
		{
			name: "SubString search with multiple key-value pairs",
			s: Search{
				typ: SubString,
				values: map[SearchKey]string{
					"name":    "John",
					"content": "search",
				},
			},
			want:    Statement{},
			wantErr: true,
		},
		{
			name: "FullText search - not implemented",
			s: Search{
				typ: FullText,
				values: map[SearchKey]string{
					"content": "search text",
				},
			},
			want:    Statement{},
			wantErr: true,
		},
		{
			name: "Ngram search - not implemented",
			s: Search{
				typ: Ngram,
				values: map[SearchKey]string{
					"content": "ngram search",
				},
			},
			want:    Statement{},
			wantErr: true,
		},
		{
			name: "Invalid search type",
			s: Search{
				typ: SearchType("invalid"),
				values: map[SearchKey]string{
					"name": "test",
				},
			},
			want:    Statement{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := tt.s.spannerStmt()
			if (err != nil) != tt.wantErr {
				t.Errorf("Search.spannerStmt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				return
			}
			if got.SQL != tt.want.SQL {
				t.Errorf("Search.spannerStmt().SQL = %v, want %v", got.SQL, tt.want.SQL)
			}
			if len(got.SpannerParams) != len(tt.want.SpannerParams) {
				t.Errorf("Search.spannerStmt().SpannerParams length = %v, want %v", len(got.SpannerParams), len(tt.want.SpannerParams))
			}
			for key, value := range tt.want.SpannerParams {
				if got.SpannerParams[key] != value {
					t.Errorf("Search.spannerStmt().SpannerParams[%v] = %v, want %v", key, got.SpannerParams[key], value)
				}
			}
		})
	}
}

func TestSearch_ParseToSearchSubstring(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		s       Search
		want    Statement
		wantErr bool
	}{
		{
			name: "Single term",
			s: Search{
				typ: SubString,
				values: map[SearchKey]string{
					"name": "John",
				},
			},
			want: Statement{
				SQL: "SEARCH_SUBSTRING(name, @searchsubstringterm0)",
				SpannerParams: map[string]any{
					"searchsubstringterm0": "John",
				},
			},
			wantErr: false,
		},
		{
			name: "Multiple terms",
			s: Search{
				typ: SubString,
				values: map[SearchKey]string{
					"content": "John Doe",
				},
			},
			want: Statement{
				SQL: "SEARCH_SUBSTRING(content, @searchsubstringterm0) OR SEARCH_SUBSTRING(content, @searchsubstringterm1)",
				SpannerParams: map[string]any{
					"searchsubstringterm0": "John",
					"searchsubstringterm1": "Doe",
				},
			},
			wantErr: false,
		},
		{
			name: "Multiple key-value pairs",
			s: Search{
				typ: SubString,
				values: map[SearchKey]string{
					"name":    "John",
					"content": "search",
				},
			},
			want:    Statement{},
			wantErr: true,
		},
		{
			name: "Empty values",
			s: Search{
				typ:    SubString,
				values: map[SearchKey]string{},
			},
			want:    Statement{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := tt.s.parseToSearchSubstring()
			if (err != nil) != tt.wantErr {
				t.Errorf("Search.parseToSearchSubstring() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				return
			}
			if got.SQL != tt.want.SQL {
				t.Errorf("Search.parseToSearchSubstring().SQL = %v, want %v", got.SQL, tt.want.SQL)
			}
			if len(got.SpannerParams) != len(tt.want.SpannerParams) {
				t.Errorf("Search.parseToSearchSubstring().SpannerParams length = %v, want %v", len(got.SpannerParams), len(tt.want.SpannerParams))
			}
			for key, value := range tt.want.SpannerParams {
				if got.SpannerParams[key] != value {
					t.Errorf("Search.parseToSearchSubstring().SpannerParams[%v] = %v, want %v", key, got.SpannerParams[key], value)
				}
			}
		})
	}
}

func TestSearch_ParseToNgramScore(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		s       Search
		want    Statement
		wantErr bool
	}{
		{
			name: "Single term",
			s: Search{
				typ: SubString,
				values: map[SearchKey]string{
					"name": "John",
				},
			},
			want: Statement{
				SQL: "SCORE_NGRAMS(name, @ngramscoreterm0)",
				SpannerParams: map[string]any{
					"ngramscoreterm0": "John",
				},
			},
			wantErr: false,
		},
		{
			name: "Multiple terms",
			s: Search{
				typ: SubString,
				values: map[SearchKey]string{
					"content": "John Doe",
				},
			},
			want: Statement{
				SQL: "SCORE_NGRAMS(content, @ngramscoreterm0) + SCORE_NGRAMS(content, @ngramscoreterm1)",
				SpannerParams: map[string]any{
					"ngramscoreterm0": "John",
					"ngramscoreterm1": "Doe",
				},
			},
			wantErr: false,
		},
		{
			name: "Multiple key-value pairs",
			s: Search{
				typ: SubString,
				values: map[SearchKey]string{
					"name":    "John",
					"content": "search",
				},
			},
			want:    Statement{},
			wantErr: true,
		},
		{
			name: "Empty values",
			s: Search{
				typ:    SubString,
				values: map[SearchKey]string{},
			},
			want:    Statement{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := tt.s.parseToNgramScore()
			if (err != nil) != tt.wantErr {
				t.Errorf("Search.parseToNgramScore() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				return
			}
			if got.SQL != tt.want.SQL {
				t.Errorf("Search.parseToNgramScore().SQL = %v, want %v", got.SQL, tt.want.SQL)
			}
			if len(got.SpannerParams) != len(tt.want.SpannerParams) {
				t.Errorf("Search.parseToNgramScore().SpannerParams length = %v, want %v", len(got.SpannerParams), len(tt.want.SpannerParams))
			}
			for key, value := range tt.want.SpannerParams {
				if got.SpannerParams[key] != value {
					t.Errorf("Search.parseToNgramScore().SpannerParams[%v] = %v, want %v", key, got.SpannerParams[key], value)
				}
			}
		})
	}
}

func TestSearch_ErrorTypes(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		s         Search
		wantErr   bool
		errString string
	}{
		{
			name: "FullText search error",
			s: Search{
				typ: FullText,
				values: map[SearchKey]string{
					"content": "search text",
				},
			},
			wantErr:   true,
			errString: "fulltext search is not yet implemented",
		},
		{
			name: "Ngram search error",
			s: Search{
				typ: Ngram,
				values: map[SearchKey]string{
					"content": "ngram search",
				},
			},
			wantErr:   true,
			errString: "ngram search is not yet implemented",
		},
		{
			name: "Invalid search type error",
			s: Search{
				typ: SearchType("invalid"),
				values: map[SearchKey]string{
					"name": "test",
				},
			},
			wantErr:   true,
			errString: "invalid search type not supported",
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			_, err := tt.s.spannerStmt()
			if (err != nil) != tt.wantErr {
				t.Errorf("Search.spannerStmt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && err != nil {
				if !errors.Is(err, errors.New(tt.errString)) && err.Error() != tt.errString {
					t.Errorf("Search.spannerStmt() error = %v, want error containing %v", err, tt.errString)
				}
			}
		})
	}
}