package policy

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/grokify/pipelineconductor/pkg/model"
	"gopkg.in/yaml.v3"
)

// ProfileManager manages CI/CD profiles.
type ProfileManager struct {
	profiles map[string]*model.Profile
}

// NewProfileManager creates a new profile manager.
func NewProfileManager() *ProfileManager {
	return &ProfileManager{
		profiles: make(map[string]*model.Profile),
	}
}

// LoadFromDirectory loads all profile YAML files from a directory.
func (m *ProfileManager) LoadFromDirectory(dir string) error {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("reading directory %s: %w", dir, err)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		name := entry.Name()
		if filepath.Ext(name) != ".yaml" && filepath.Ext(name) != ".yml" {
			continue
		}

		path := filepath.Join(dir, name)
		if err := m.LoadFromFile(path); err != nil {
			return err
		}
	}

	return nil
}

// LoadFromFile loads a profile from a YAML file.
func (m *ProfileManager) LoadFromFile(path string) error {
	content, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("reading file %s: %w", path, err)
	}

	var profile model.Profile
	if err := yaml.Unmarshal(content, &profile); err != nil {
		return fmt.Errorf("parsing YAML %s: %w", path, err)
	}

	if profile.Name == "" {
		// Use filename as profile name
		base := filepath.Base(path)
		profile.Name = base[:len(base)-len(filepath.Ext(base))]
	}

	m.profiles[profile.Name] = &profile
	return nil
}

// Get returns a profile by name.
func (m *ProfileManager) Get(name string) (*model.Profile, error) {
	profile, ok := m.profiles[name]
	if !ok {
		return nil, fmt.Errorf("profile not found: %s", name)
	}
	return profile, nil
}

// GetOrDefault returns a profile by name, or the default profile if not found.
func (m *ProfileManager) GetOrDefault(name string) *model.Profile {
	if profile, ok := m.profiles[name]; ok {
		return profile
	}
	if profile, ok := m.profiles["default"]; ok {
		return profile
	}
	return DefaultProfile()
}

// List returns all profile names.
func (m *ProfileManager) List() []string {
	names := make([]string, 0, len(m.profiles))
	for name := range m.profiles {
		names = append(names, name)
	}
	return names
}

// Add adds a profile to the manager.
func (m *ProfileManager) Add(profile *model.Profile) {
	m.profiles[profile.Name] = profile
}

// LoadBuiltinProfiles loads the built-in default profiles.
func (m *ProfileManager) LoadBuiltinProfiles() {
	m.Add(DefaultProfile())
	m.Add(ModernProfile())
	m.Add(LegacyProfile())
}

// DefaultProfile returns the default CI/CD profile.
func DefaultProfile() *model.Profile {
	return &model.Profile{
		Name:        "default",
		Description: "Standard Go CI configuration for active projects",
		Go: model.ProfileGo{
			Versions: []string{"1.24", "1.25"},
		},
		OS: []string{"ubuntu-latest", "macos-latest", "windows-latest"},
		Checks: model.ProfileChecks{
			Required: []string{"test", "lint", "build"},
		},
		Lint: model.ProfileLint{
			Enabled: true,
			Tool:    "golangci-lint",
		},
		Test: model.ProfileTest{
			Coverage: true,
			Race:     true,
		},
	}
}

// ModernProfile returns a profile for modern Go projects.
func ModernProfile() *model.Profile {
	return &model.Profile{
		Name:        "modern",
		Description: "Modern Go CI for projects using latest Go features",
		Go: model.ProfileGo{
			Versions: []string{"1.25"},
		},
		OS: []string{"ubuntu-latest", "macos-latest"},
		Checks: model.ProfileChecks{
			Required: []string{"test", "lint", "build"},
		},
		Lint: model.ProfileLint{
			Enabled: true,
			Tool:    "golangci-lint",
		},
		Test: model.ProfileTest{
			Coverage: true,
			Race:     true,
		},
	}
}

// LegacyProfile returns a profile for legacy Go projects.
func LegacyProfile() *model.Profile {
	return &model.Profile{
		Name:        "legacy",
		Description: "Legacy Go CI for older projects requiring Go 1.12-1.18",
		Go: model.ProfileGo{
			Versions: []string{"1.12"},
		},
		OS: []string{"ubuntu-latest"},
		Checks: model.ProfileChecks{
			Required: []string{"test", "build"},
		},
		Lint: model.ProfileLint{
			Enabled: false,
		},
		Test: model.ProfileTest{
			Coverage: false,
			Race:     false,
		},
	}
}

// ValidateRepoAgainstProfile checks if a repo's CI config matches a profile.
func ValidateRepoAgainstProfile(ctx *model.PolicyContext, profile *model.Profile) []model.Violation {
	var violations []model.Violation

	// Check Go versions
	if len(ctx.Go.Versions) > 0 {
		for _, v := range ctx.Go.Versions {
			found := false
			for _, pv := range profile.Go.Versions {
				if v == pv {
					found = true
					break
				}
			}
			if !found {
				violations = append(violations, model.Violation{
					Policy:      "profile/go-version",
					Rule:        "allowed-versions",
					Message:     fmt.Sprintf("Go version %s not in profile %s allowed versions %v", v, profile.Name, profile.Go.Versions),
					Severity:    model.SeverityMedium,
					Remediation: fmt.Sprintf("Update go-version to one of: %v", profile.Go.Versions),
				})
			}
		}
	}

	// Check OS matrix
	if len(ctx.CI.OSMatrix) > 0 {
		profileOSSet := make(map[string]bool)
		for _, os := range profile.OS {
			profileOSSet[os] = true
		}

		for _, os := range ctx.CI.OSMatrix {
			if !profileOSSet[os] {
				violations = append(violations, model.Violation{
					Policy:      "profile/os-matrix",
					Rule:        "allowed-os",
					Message:     fmt.Sprintf("OS %s not in profile %s allowed platforms %v", os, profile.Name, profile.OS),
					Severity:    model.SeverityLow,
					Remediation: fmt.Sprintf("Update runs-on to one of: %v", profile.OS),
				})
			}
		}

		// Check for missing required OS
		for _, requiredOS := range profile.OS {
			found := false
			for _, os := range ctx.CI.OSMatrix {
				if os == requiredOS {
					found = true
					break
				}
			}
			if !found {
				violations = append(violations, model.Violation{
					Policy:      "profile/os-matrix",
					Rule:        "required-os",
					Message:     fmt.Sprintf("Profile %s requires OS %s but it's not in the matrix", profile.Name, requiredOS),
					Severity:    model.SeverityInfo,
					Remediation: fmt.Sprintf("Add %s to your OS matrix", requiredOS),
				})
			}
		}
	}

	return violations
}
