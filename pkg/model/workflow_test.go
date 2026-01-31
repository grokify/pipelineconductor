package model

import (
	"reflect"
	"testing"
)

func TestParseReusableWorkflowRef(t *testing.T) {
	tests := []struct {
		name string
		ref  string
		want *ReusableWorkflowRef
	}{
		{
			name: "full reference with path and ref",
			ref:  "owner/repo/.github/workflows/ci.yml@v1",
			want: &ReusableWorkflowRef{
				Owner:   "owner",
				Repo:    "repo",
				Path:    ".github/workflows/ci.yml",
				Ref:     "v1",
				FullRef: "owner/repo/.github/workflows/ci.yml@v1",
			},
		},
		{
			name: "reference with main branch",
			ref:  "myorg/shared-workflows/.github/workflows/build.yml@main",
			want: &ReusableWorkflowRef{
				Owner:   "myorg",
				Repo:    "shared-workflows",
				Path:    ".github/workflows/build.yml",
				Ref:     "main",
				FullRef: "myorg/shared-workflows/.github/workflows/build.yml@main",
			},
		},
		{
			name: "reference with SHA",
			ref:  "org/repo/.github/workflows/test.yml@abc123def456",
			want: &ReusableWorkflowRef{
				Owner:   "org",
				Repo:    "repo",
				Path:    ".github/workflows/test.yml",
				Ref:     "abc123def456",
				FullRef: "org/repo/.github/workflows/test.yml@abc123def456",
			},
		},
		{
			name: "local workflow reference",
			ref:  "./.github/workflows/local.yml",
			want: &ReusableWorkflowRef{
				Owner:   ".",
				Repo:    ".github",
				Path:    "workflows/local.yml",
				FullRef: "./.github/workflows/local.yml",
			},
		},
		{
			name: "reference without @ (no ref)",
			ref:  "owner/repo/.github/workflows/ci.yml",
			want: &ReusableWorkflowRef{
				Owner:   "owner",
				Repo:    "repo",
				Path:    ".github/workflows/ci.yml",
				FullRef: "owner/repo/.github/workflows/ci.yml",
			},
		},
		{
			name: "simple owner/repo format",
			ref:  "owner/repo@v1",
			want: &ReusableWorkflowRef{
				Owner:   "owner",
				Repo:    "repo",
				Ref:     "v1",
				FullRef: "owner/repo@v1",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ParseReusableWorkflowRef(tt.ref)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseReusableWorkflowRef() = %+v, want %+v", got, tt.want)
			}
		})
	}
}

func TestWorkflowUsesReusableWorkflow(t *testing.T) {
	tests := []struct {
		name     string
		workflow Workflow
		want     bool
	}{
		{
			name: "workflow using reusable workflow",
			workflow: Workflow{
				Name:                 "CI",
				UsesReusableWorkflow: true,
			},
			want: true,
		},
		{
			name: "workflow not using reusable workflow",
			workflow: Workflow{
				Name:                 "Standard CI",
				UsesReusableWorkflow: false,
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.workflow.UsesReusableWorkflow
			if got != tt.want {
				t.Errorf("Workflow.UsesReusableWorkflow = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMatrixConfigOS(t *testing.T) {
	tests := []struct {
		name   string
		matrix MatrixConfig
		wantOS []string
	}{
		{
			name: "single OS",
			matrix: MatrixConfig{
				OS: []string{"ubuntu-latest"},
			},
			wantOS: []string{"ubuntu-latest"},
		},
		{
			name: "multiple OS",
			matrix: MatrixConfig{
				OS: []string{"ubuntu-latest", "macos-latest", "windows-latest"},
			},
			wantOS: []string{"ubuntu-latest", "macos-latest", "windows-latest"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !reflect.DeepEqual(tt.matrix.OS, tt.wantOS) {
				t.Errorf("MatrixConfig.OS = %v, want %v", tt.matrix.OS, tt.wantOS)
			}
		})
	}
}
