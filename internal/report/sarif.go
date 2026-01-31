package report

import (
	"encoding/json"
	"fmt"

	"github.com/grokify/pipelineconductor/pkg/model"
)

// SARIFFormatter generates SARIF reports for GitHub Security integration.
type SARIFFormatter struct{}

// SARIF represents a SARIF log file.
type SARIF struct {
	Version string     `json:"version"`
	Schema  string     `json:"$schema"`
	Runs    []SARIFRun `json:"runs"`
}

// SARIFRun represents a single run in the SARIF log.
type SARIFRun struct {
	Tool    SARIFTool     `json:"tool"`
	Results []SARIFResult `json:"results"`
	Rules   []SARIFRule   `json:"rules,omitempty"`
}

// SARIFTool describes the analysis tool.
type SARIFTool struct {
	Driver SARIFDriver `json:"driver"`
}

// SARIFDriver describes the tool driver.
type SARIFDriver struct {
	Name           string      `json:"name"`
	Version        string      `json:"version"`
	InformationURI string      `json:"informationUri"`
	Rules          []SARIFRule `json:"rules,omitempty"`
}

// SARIFRule describes a rule/policy.
type SARIFRule struct {
	ID               string                   `json:"id"`
	Name             string                   `json:"name"`
	ShortDescription SARIFMultiformatMessage  `json:"shortDescription"`
	FullDescription  *SARIFMultiformatMessage `json:"fullDescription,omitempty"`
	DefaultConfig    *SARIFRuleConfig         `json:"defaultConfiguration,omitempty"`
	HelpURI          string                   `json:"helpUri,omitempty"`
}

// SARIFRuleConfig describes rule configuration.
type SARIFRuleConfig struct {
	Level string `json:"level"`
}

// SARIFMultiformatMessage represents a message in multiple formats.
type SARIFMultiformatMessage struct {
	Text string `json:"text"`
}

// SARIFResult represents a single result/finding.
type SARIFResult struct {
	RuleID    string          `json:"ruleId"`
	RuleIndex int             `json:"ruleIndex"`
	Level     string          `json:"level"`
	Message   SARIFMessage    `json:"message"`
	Locations []SARIFLocation `json:"locations,omitempty"`
	Fixes     []SARIFFix      `json:"fixes,omitempty"`
}

// SARIFMessage represents a result message.
type SARIFMessage struct {
	Text string `json:"text"`
}

// SARIFLocation represents a location in source.
type SARIFLocation struct {
	PhysicalLocation *SARIFPhysicalLocation `json:"physicalLocation,omitempty"`
	LogicalLocations []SARIFLogicalLocation `json:"logicalLocations,omitempty"`
}

// SARIFPhysicalLocation represents a physical file location.
type SARIFPhysicalLocation struct {
	ArtifactLocation SARIFArtifactLocation `json:"artifactLocation"`
	Region           *SARIFRegion          `json:"region,omitempty"`
}

// SARIFArtifactLocation represents the artifact (file) location.
type SARIFArtifactLocation struct {
	URI       string `json:"uri"`
	URIBaseID string `json:"uriBaseId,omitempty"`
}

// SARIFRegion represents a region within a file.
type SARIFRegion struct {
	StartLine   int `json:"startLine,omitempty"`
	StartColumn int `json:"startColumn,omitempty"`
	EndLine     int `json:"endLine,omitempty"`
	EndColumn   int `json:"endColumn,omitempty"`
}

// SARIFLogicalLocation represents a logical location (repo, component).
type SARIFLogicalLocation struct {
	Name               string `json:"name"`
	FullyQualifiedName string `json:"fullyQualifiedName"`
	Kind               string `json:"kind"`
}

// SARIFFix represents a suggested fix.
type SARIFFix struct {
	Description SARIFMessage `json:"description"`
}

// Format generates a SARIF report.
func (f *SARIFFormatter) Format(result *model.ComplianceResult) ([]byte, error) {
	sarif := &SARIF{
		Version: "2.1.0",
		Schema:  "https://raw.githubusercontent.com/oasis-tcs/sarif-spec/master/Schemata/sarif-schema-2.1.0.json",
		Runs:    []SARIFRun{},
	}

	// Collect unique rules across all results
	rulesMap := make(map[string]*SARIFRule)
	ruleIndex := make(map[string]int)
	var results []SARIFResult

	for _, repo := range result.Repos {
		for _, v := range repo.Violations {
			ruleID := v.Policy
			if v.Rule != "" {
				ruleID = fmt.Sprintf("%s/%s", v.Policy, v.Rule)
			}

			// Add rule if not seen
			if _, ok := rulesMap[ruleID]; !ok {
				rule := &SARIFRule{
					ID:   ruleID,
					Name: v.Policy,
					ShortDescription: SARIFMultiformatMessage{
						Text: v.Message,
					},
					DefaultConfig: &SARIFRuleConfig{
						Level: severityToSARIFLevel(v.Severity),
					},
				}
				rulesMap[ruleID] = rule
				ruleIndex[ruleID] = len(rulesMap) - 1
			}

			// Create result
			sarifResult := SARIFResult{
				RuleID:    ruleID,
				RuleIndex: ruleIndex[ruleID],
				Level:     severityToSARIFLevel(v.Severity),
				Message: SARIFMessage{
					Text: fmt.Sprintf("[%s] %s", repo.Repo.FullName, v.Message),
				},
				Locations: []SARIFLocation{
					{
						LogicalLocations: []SARIFLogicalLocation{
							{
								Name:               repo.Repo.Name,
								FullyQualifiedName: repo.Repo.FullName,
								Kind:               "repository",
							},
						},
					},
				},
			}

			// Add file location if available
			if v.File != "" {
				sarifResult.Locations[0].PhysicalLocation = &SARIFPhysicalLocation{
					ArtifactLocation: SARIFArtifactLocation{
						URI: v.File,
					},
				}
				if v.Line > 0 {
					sarifResult.Locations[0].PhysicalLocation.Region = &SARIFRegion{
						StartLine: v.Line,
					}
				}
			}

			// Add fix suggestion if remediation available
			if v.Remediation != "" {
				sarifResult.Fixes = []SARIFFix{
					{
						Description: SARIFMessage{
							Text: v.Remediation,
						},
					},
				}
			}

			results = append(results, sarifResult)
		}
	}

	// Convert rules map to slice
	rules := make([]SARIFRule, len(rulesMap))
	for ruleID, rule := range rulesMap {
		rules[ruleIndex[ruleID]] = *rule
	}

	// Create run
	run := SARIFRun{
		Tool: SARIFTool{
			Driver: SARIFDriver{
				Name:           "PipelineConductor",
				Version:        "0.1.0",
				InformationURI: "https://github.com/grokify/pipelineconductor",
				Rules:          rules,
			},
		},
		Results: results,
	}

	sarif.Runs = append(sarif.Runs, run)

	return json.MarshalIndent(sarif, "", "  ")
}

func severityToSARIFLevel(s model.Severity) string {
	switch s {
	case model.SeverityCritical, model.SeverityHigh:
		return "error"
	case model.SeverityMedium:
		return "warning"
	case model.SeverityLow, model.SeverityInfo:
		return "note"
	default:
		return "note"
	}
}
