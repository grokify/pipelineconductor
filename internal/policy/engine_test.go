package policy

import (
	"testing"

	"github.com/grokify/pipelineconductor/pkg/model"
)

func TestNewEngine(t *testing.T) {
	engine := NewEngine()
	if engine == nil {
		t.Fatal("NewEngine() returned nil")
	}
	if engine.policySet == nil {
		t.Error("NewEngine() policySet is nil")
	}
}

func TestEngineEvaluate(t *testing.T) {
	engine := NewEngine()

	ctx := &model.PolicyContext{
		Repo: model.RepoContext{
			FullName: "org/test-repo",
			Name:     "test-repo",
			Org:      "org",
		},
	}

	result := engine.Evaluate(ctx, "merge")
	// With no policies loaded, default should be deny
	if result.Action != "merge" {
		t.Errorf("Evaluate() action = %s, want merge", result.Action)
	}
	if result.RepoName != "org/test-repo" {
		t.Errorf("Evaluate() repoName = %s, want org/test-repo", result.RepoName)
	}
}

func TestEngineEvaluateAll(t *testing.T) {
	engine := NewEngine()

	ctx := &model.PolicyContext{
		Repo: model.RepoContext{
			FullName: "org/test-repo",
		},
	}

	results := engine.EvaluateAll(ctx)
	// Should have results for build, test, lint, merge actions
	if len(results) != 4 {
		t.Errorf("EvaluateAll() returned %d results, want 4", len(results))
	}

	// Verify action types
	expectedActions := map[string]bool{
		ActionBuild: false,
		ActionTest:  false,
		ActionLint:  false,
		ActionMerge: false,
	}
	for _, result := range results {
		expectedActions[result.Action] = true
	}
	for action, found := range expectedActions {
		if !found {
			t.Errorf("EvaluateAll() missing action: %s", action)
		}
	}
}

func TestEvaluationResultToViolation(t *testing.T) {
	tests := []struct {
		name    string
		result  EvaluationResult
		wantNil bool
	}{
		{
			name: "allowed result returns nil",
			result: EvaluationResult{
				Action:  "merge",
				Allowed: true,
			},
			wantNil: true,
		},
		{
			name: "denied result returns violation",
			result: EvaluationResult{
				Action:  "merge",
				Allowed: false,
				Reasons: []string{"branch-protection-policy"},
			},
			wantNil: false,
		},
		{
			name: "denied build result returns medium severity",
			result: EvaluationResult{
				Action:  "build",
				Allowed: false,
			},
			wantNil: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			violation := tt.result.ToViolation()
			if tt.wantNil && violation != nil {
				t.Errorf("ToViolation() = %v, want nil", violation)
			}
			if !tt.wantNil && violation == nil {
				t.Error("ToViolation() = nil, want violation")
			}
			if !tt.wantNil && violation != nil {
				if violation.Policy != "cedar/"+tt.result.Action {
					t.Errorf("ToViolation().Policy = %s, want cedar/%s", violation.Policy, tt.result.Action)
				}
				// Merge and deploy should be high severity
				if (tt.result.Action == "merge" || tt.result.Action == "deploy") &&
					violation.Severity != model.SeverityHigh {
					t.Errorf("ToViolation().Severity = %s, want high for %s action",
						violation.Severity, tt.result.Action)
				}
			}
		})
	}
}

func TestEngineAddPolicy(t *testing.T) {
	engine := NewEngine()

	// Valid policy
	validPolicy := []byte(`permit(principal, action, resource);`)
	err := engine.AddPolicy("test-permit", validPolicy)
	if err != nil {
		t.Errorf("AddPolicy() with valid policy returned error: %v", err)
	}

	// Invalid policy
	invalidPolicy := []byte(`this is not valid cedar`)
	err = engine.AddPolicy("test-invalid", invalidPolicy)
	if err == nil {
		t.Error("AddPolicy() with invalid policy should return error")
	}
}
