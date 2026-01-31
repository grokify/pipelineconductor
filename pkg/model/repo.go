// Package model provides core data structures for PipelineConductor.
package model

import "time"

// Repo represents a repository with its metadata.
type Repo struct {
	Owner           string    `json:"owner"`
	Name            string    `json:"name"`
	FullName        string    `json:"fullName"`
	DefaultBranch   string    `json:"defaultBranch"`
	Languages       []string  `json:"languages"`
	PrimaryLanguage string    `json:"primaryLanguage"`
	Topics          []string  `json:"topics"`
	Visibility      string    `json:"visibility"`
	Archived        bool      `json:"archived"`
	Fork            bool      `json:"fork"`
	CreatedAt       time.Time `json:"createdAt"`
	UpdatedAt       time.Time `json:"updatedAt"`
	PushedAt        time.Time `json:"pushedAt"`
	HTMLURL         string    `json:"htmlUrl"`
	CloneURL        string    `json:"cloneUrl"`
}

// RepoFilter defines criteria for filtering repositories.
type RepoFilter struct {
	IncludeLanguages []string `json:"includeLanguages,omitempty"`
	ExcludeLanguages []string `json:"excludeLanguages,omitempty"`
	IncludeTopics    []string `json:"includeTopics,omitempty"`
	ExcludeTopics    []string `json:"excludeTopics,omitempty"`
	IncludeArchived  bool     `json:"includeArchived"`
	IncludeForks     bool     `json:"includeForks"`
	VisibilityFilter []string `json:"visibilityFilter,omitempty"`
	NamePattern      string   `json:"namePattern,omitempty"`
}

// BranchProtection represents branch protection settings.
type BranchProtection struct {
	Branch               string   `json:"branch"`
	Enabled              bool     `json:"enabled"`
	RequireReviews       bool     `json:"requireReviews"`
	RequiredReviewers    int      `json:"requiredReviewers"`
	RequireStatusChecks  bool     `json:"requireStatusChecks"`
	RequiredStatusChecks []string `json:"requiredStatusChecks"`
	EnforceAdmins        bool     `json:"enforceAdmins"`
	RequireSignedCommits bool     `json:"requireSignedCommits"`
	AllowForcePushes     bool     `json:"allowForcePushes"`
	AllowDeletions       bool     `json:"allowDeletions"`
}

// Matches returns true if the repo matches the filter criteria.
func (r *Repo) Matches(filter RepoFilter) bool {
	if r.Archived && !filter.IncludeArchived {
		return false
	}
	if r.Fork && !filter.IncludeForks {
		return false
	}
	if len(filter.VisibilityFilter) > 0 && !contains(filter.VisibilityFilter, r.Visibility) {
		return false
	}
	if len(filter.IncludeLanguages) > 0 && !hasAny(r.Languages, filter.IncludeLanguages) {
		return false
	}
	if hasAny(r.Languages, filter.ExcludeLanguages) {
		return false
	}
	if len(filter.IncludeTopics) > 0 && !hasAny(r.Topics, filter.IncludeTopics) {
		return false
	}
	if hasAny(r.Topics, filter.ExcludeTopics) {
		return false
	}
	return true
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func hasAny(slice, items []string) bool {
	for _, item := range items {
		if contains(slice, item) {
			return true
		}
	}
	return false
}
