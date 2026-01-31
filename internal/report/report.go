// Package report provides compliance report generation in multiple formats.
package report

import (
	"fmt"
	"io"

	"github.com/grokify/pipelineconductor/pkg/model"
)

// Format represents the output format for reports.
type Format string

// Supported report formats.
const (
	FormatJSON     Format = "json"
	FormatMarkdown Format = "markdown"
	FormatSARIF    Format = "sarif"
	FormatCSV      Format = "csv"
)

// ParseFormat parses a format string into a Format.
func ParseFormat(s string) (Format, error) {
	switch s {
	case "json":
		return FormatJSON, nil
	case "markdown", "md":
		return FormatMarkdown, nil
	case "sarif":
		return FormatSARIF, nil
	case "csv":
		return FormatCSV, nil
	default:
		return "", fmt.Errorf("unknown format: %s (supported: json, markdown, sarif, csv)", s)
	}
}

// Formatter generates reports in a specific format.
type Formatter interface {
	Format(result *model.ComplianceResult) ([]byte, error)
}

// Builder creates reports from compliance results.
type Builder struct {
	formatters map[Format]Formatter
}

// NewBuilder creates a new report builder with all formatters.
func NewBuilder() *Builder {
	return &Builder{
		formatters: map[Format]Formatter{
			FormatJSON:     &JSONFormatter{},
			FormatMarkdown: &MarkdownFormatter{},
			FormatSARIF:    &SARIFFormatter{},
			FormatCSV:      &CSVFormatter{},
		},
	}
}

// Generate generates a report in the specified format.
func (b *Builder) Generate(result *model.ComplianceResult, format Format) ([]byte, error) {
	formatter, ok := b.formatters[format]
	if !ok {
		return nil, fmt.Errorf("unsupported format: %s", format)
	}
	return formatter.Format(result)
}

// Write generates a report and writes it to a writer.
func (b *Builder) Write(result *model.ComplianceResult, format Format, w io.Writer) error {
	data, err := b.Generate(result, format)
	if err != nil {
		return err
	}
	_, err = w.Write(data)
	return err
}
