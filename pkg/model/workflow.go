package model

import (
	"strings"
	"time"
)

// Workflow represents a CI/CD workflow configuration.
type Workflow struct {
	Name                 string                `json:"name"`
	Path                 string                `json:"path"`
	Content              string                `json:"content,omitempty"`
	Triggers             []string              `json:"triggers"`
	Jobs                 []WorkflowJob         `json:"jobs"`
	UsesReusableWorkflow bool                  `json:"usesReusableWorkflow"`
	ReusableWorkflowRefs []ReusableWorkflowRef `json:"reusableWorkflowRefs,omitempty"`
	State                string                `json:"state"`
	CreatedAt            time.Time             `json:"createdAt,omitempty"`
	UpdatedAt            time.Time             `json:"updatedAt,omitempty"`
}

// WorkflowJob represents a job within a workflow.
type WorkflowJob struct {
	ID                   string               `json:"id"`
	Name                 string               `json:"name"`
	RunsOn               []string             `json:"runsOn"`
	Steps                []WorkflowStep       `json:"steps,omitempty"`
	Needs                []string             `json:"needs,omitempty"`
	Matrix               *MatrixConfig        `json:"matrix,omitempty"`
	UsesReusableWorkflow bool                 `json:"usesReusableWorkflow"`
	ReusableWorkflowRef  *ReusableWorkflowRef `json:"reusableWorkflowRef,omitempty"`
}

// WorkflowStep represents a step within a job.
type WorkflowStep struct {
	Name string            `json:"name,omitempty"`
	Uses string            `json:"uses,omitempty"`
	Run  string            `json:"run,omitempty"`
	With map[string]string `json:"with,omitempty"`
	Env  map[string]string `json:"env,omitempty"`
}

// MatrixConfig represents a build matrix configuration.
type MatrixConfig struct {
	OS            []string            `json:"os,omitempty"`
	GoVersion     []string            `json:"goVersion,omitempty"`
	PythonVersion []string            `json:"pythonVersion,omitempty"`
	NodeVersion   []string            `json:"nodeVersion,omitempty"`
	Include       []map[string]string `json:"include,omitempty"`
	Exclude       []map[string]string `json:"exclude,omitempty"`
	FailFast      bool                `json:"failFast"`
}

// ReusableWorkflowRef represents a reference to a reusable workflow.
type ReusableWorkflowRef struct {
	Owner   string `json:"owner"`
	Repo    string `json:"repo"`
	Path    string `json:"path"`
	Ref     string `json:"ref"`
	FullRef string `json:"fullRef"`
}

// WorkflowRun represents a workflow execution.
type WorkflowRun struct {
	ID         int64     `json:"id"`
	WorkflowID int64     `json:"workflowId"`
	Name       string    `json:"name"`
	Status     string    `json:"status"`
	Conclusion string    `json:"conclusion"`
	Branch     string    `json:"branch"`
	HeadSHA    string    `json:"headSha"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
	HTMLURL    string    `json:"htmlUrl"`
}

// ParseReusableWorkflowRef parses a reusable workflow reference string.
// Format: "owner/repo/.github/workflows/workflow.yml@ref"
func ParseReusableWorkflowRef(ref string) *ReusableWorkflowRef {
	result := &ReusableWorkflowRef{FullRef: ref}

	parts := strings.SplitN(ref, "@", 2)
	if len(parts) == 2 {
		result.Ref = parts[1]
		ref = parts[0]
	}

	pathParts := strings.SplitN(ref, "/", 3)
	if len(pathParts) >= 1 {
		result.Owner = pathParts[0]
	}
	if len(pathParts) >= 2 {
		result.Repo = pathParts[1]
	}
	if len(pathParts) >= 3 {
		result.Path = pathParts[2]
	}

	return result
}
