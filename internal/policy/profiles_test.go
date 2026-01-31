package policy

import (
	"testing"
)

func TestNewProfileManager(t *testing.T) {
	pm := NewProfileManager()
	if pm == nil {
		t.Fatal("NewProfileManager() returned nil")
	}
	if pm.profiles == nil {
		t.Error("NewProfileManager() profiles map is nil")
	}
}

func TestProfileManagerLoadBuiltinProfiles(t *testing.T) {
	pm := NewProfileManager()
	pm.LoadBuiltinProfiles()

	// Should have default, modern, and legacy profiles
	expectedProfiles := []string{"default", "modern", "legacy"}
	for _, name := range expectedProfiles {
		profile, err := pm.Get(name)
		if err != nil {
			t.Errorf("LoadBuiltinProfiles() error getting profile %s: %v", name, err)
		}
		if profile == nil {
			t.Errorf("LoadBuiltinProfiles() missing profile: %s", name)
		}
	}
}

func TestProfileManagerGetOrDefault(t *testing.T) {
	pm := NewProfileManager()
	pm.LoadBuiltinProfiles()

	// Get existing profile
	profile := pm.GetOrDefault("modern")
	if profile == nil {
		t.Fatal("GetOrDefault('modern') returned nil")
	}
	if profile.Name != "modern" {
		t.Errorf("GetOrDefault('modern').Name = %s, want modern", profile.Name)
	}

	// Get non-existent profile returns default
	profile = pm.GetOrDefault("nonexistent")
	if profile == nil {
		t.Fatal("GetOrDefault('nonexistent') returned nil")
	}
	if profile.Name != "default" {
		t.Errorf("GetOrDefault('nonexistent').Name = %s, want default", profile.Name)
	}

	// Get with empty string returns default
	profile = pm.GetOrDefault("")
	if profile == nil {
		t.Fatal("GetOrDefault('') returned nil")
	}
	if profile.Name != "default" {
		t.Errorf("GetOrDefault('').Name = %s, want default", profile.Name)
	}
}

func TestProfileManagerList(t *testing.T) {
	pm := NewProfileManager()
	pm.LoadBuiltinProfiles()

	profiles := pm.List()
	if len(profiles) != 3 {
		t.Errorf("List() returned %d profiles, want 3", len(profiles))
	}
}

func TestProfileManagerGet(t *testing.T) {
	pm := NewProfileManager()
	pm.LoadBuiltinProfiles()

	// Existing profile
	profile, err := pm.Get("default")
	if err != nil {
		t.Errorf("Get('default') error: %v", err)
	}
	if profile == nil {
		t.Error("Get('default') returned nil profile")
	}

	// Non-existent profile
	profile, err = pm.Get("nonexistent")
	if err == nil {
		t.Error("Get('nonexistent') should return error")
	}
	if profile != nil {
		t.Error("Get('nonexistent') should return nil profile")
	}
}

func TestDefaultProfile(t *testing.T) {
	profile := DefaultProfile()
	if profile == nil {
		t.Fatal("DefaultProfile() returned nil")
	}
	if profile.Name != "default" {
		t.Errorf("DefaultProfile().Name = %s, want default", profile.Name)
	}
	if len(profile.Go.Versions) == 0 {
		t.Error("DefaultProfile().Go.Versions should not be empty")
	}
	if len(profile.OS) == 0 {
		t.Error("DefaultProfile().OS should not be empty")
	}
	if len(profile.Checks.Required) == 0 {
		t.Error("DefaultProfile().Checks.Required should not be empty")
	}
}

func TestModernProfile(t *testing.T) {
	profile := ModernProfile()
	if profile == nil {
		t.Fatal("ModernProfile() returned nil")
	}
	if profile.Name != "modern" {
		t.Errorf("ModernProfile().Name = %s, want modern", profile.Name)
	}
	if len(profile.Go.Versions) == 0 {
		t.Error("ModernProfile().Go.Versions should not be empty")
	}
	if len(profile.OS) == 0 {
		t.Error("ModernProfile().OS should not be empty")
	}
}

func TestLegacyProfile(t *testing.T) {
	profile := LegacyProfile()
	if profile == nil {
		t.Fatal("LegacyProfile() returned nil")
	}
	if profile.Name != "legacy" {
		t.Errorf("LegacyProfile().Name = %s, want legacy", profile.Name)
	}
	// Legacy should not require lint
	if profile.Lint.Enabled {
		t.Error("LegacyProfile().Lint.Enabled should be false")
	}
	// Legacy should not require race detector or coverage
	if profile.Test.Race {
		t.Error("LegacyProfile().Test.Race should be false")
	}
	if profile.Test.Coverage {
		t.Error("LegacyProfile().Test.Coverage should be false")
	}
}
