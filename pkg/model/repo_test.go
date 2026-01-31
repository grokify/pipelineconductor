package model

import (
	"testing"
)

func TestRepoMatches(t *testing.T) {
	tests := []struct {
		name   string
		repo   Repo
		filter RepoFilter
		want   bool
	}{
		{
			name: "empty filter matches all",
			repo: Repo{
				Name:            "test-repo",
				FullName:        "org/test-repo",
				PrimaryLanguage: "Go",
			},
			filter: RepoFilter{},
			want:   true,
		},
		{
			name: "archived repo excluded by default",
			repo: Repo{
				Name:     "archived-repo",
				FullName: "org/archived-repo",
				Archived: true,
			},
			filter: RepoFilter{
				IncludeArchived: false,
			},
			want: false,
		},
		{
			name: "archived repo included when flag set",
			repo: Repo{
				Name:     "archived-repo",
				FullName: "org/archived-repo",
				Archived: true,
			},
			filter: RepoFilter{
				IncludeArchived: true,
			},
			want: true,
		},
		{
			name: "fork excluded by default",
			repo: Repo{
				Name:     "forked-repo",
				FullName: "org/forked-repo",
				Fork:     true,
			},
			filter: RepoFilter{
				IncludeForks: false,
			},
			want: false,
		},
		{
			name: "fork included when flag set",
			repo: Repo{
				Name:     "forked-repo",
				FullName: "org/forked-repo",
				Fork:     true,
			},
			filter: RepoFilter{
				IncludeForks: true,
			},
			want: true,
		},
		{
			name: "language filter matches",
			repo: Repo{
				Name:      "go-repo",
				FullName:  "org/go-repo",
				Languages: []string{"Go", "Shell"},
			},
			filter: RepoFilter{
				IncludeLanguages: []string{"Go", "Python"},
			},
			want: true,
		},
		{
			name: "language filter does not match",
			repo: Repo{
				Name:      "java-repo",
				FullName:  "org/java-repo",
				Languages: []string{"Java"},
			},
			filter: RepoFilter{
				IncludeLanguages: []string{"Go", "Python"},
			},
			want: false,
		},
		{
			name: "topic filter matches",
			repo: Repo{
				Name:     "api-repo",
				FullName: "org/api-repo",
				Topics:   []string{"api", "rest"},
			},
			filter: RepoFilter{
				IncludeTopics: []string{"api"},
			},
			want: true,
		},
		{
			name: "topic filter does not match",
			repo: Repo{
				Name:     "api-repo",
				FullName: "org/api-repo",
				Topics:   []string{"api", "rest"},
			},
			filter: RepoFilter{
				IncludeTopics: []string{"grpc"},
			},
			want: false,
		},
		{
			name: "visibility filter matches",
			repo: Repo{
				Name:       "public-repo",
				FullName:   "org/public-repo",
				Visibility: "public",
			},
			filter: RepoFilter{
				VisibilityFilter: []string{"public", "internal"},
			},
			want: true,
		},
		{
			name: "visibility filter does not match",
			repo: Repo{
				Name:       "private-repo",
				FullName:   "org/private-repo",
				Visibility: "private",
			},
			filter: RepoFilter{
				VisibilityFilter: []string{"public"},
			},
			want: false,
		},
		{
			name: "exclude language filter",
			repo: Repo{
				Name:      "vendor-repo",
				FullName:  "org/vendor-repo",
				Languages: []string{"Go", "C"},
			},
			filter: RepoFilter{
				ExcludeLanguages: []string{"C"},
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.repo.Matches(tt.filter)
			if got != tt.want {
				t.Errorf("Repo.Matches() = %v, want %v", got, tt.want)
			}
		})
	}
}
