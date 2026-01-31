package report

import (
	"encoding/json"

	"github.com/grokify/pipelineconductor/pkg/model"
)

// JSONFormatter generates JSON reports.
type JSONFormatter struct {
	Indent bool
}

// Format generates a JSON report.
func (f *JSONFormatter) Format(result *model.ComplianceResult) ([]byte, error) {
	if f.Indent {
		return json.MarshalIndent(result, "", "  ")
	}
	return json.MarshalIndent(result, "", "  ")
}
