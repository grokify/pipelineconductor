package report

import (
	"bytes"
	"encoding/csv"

	"github.com/grokify/pipelineconductor/pkg/model"
)

// CSVFormatter generates CSV reports.
type CSVFormatter struct{}

// Format generates a CSV report.
func (f *CSVFormatter) Format(result *model.ComplianceResult) ([]byte, error) {
	var buf bytes.Buffer
	w := csv.NewWriter(&buf)

	// Write header
	header := []string{
		"repo",
		"org",
		"compliant",
		"violation_count",
		"warning_count",
		"error",
		"skipped",
		"skip_reason",
		"scan_time_ms",
	}
	if err := w.Write(header); err != nil {
		return nil, err
	}

	// Write repo summary rows
	for _, repo := range result.Repos {
		compliant := "true"
		if !repo.Compliant {
			compliant = "false"
		}
		skipped := "false"
		if repo.Skipped {
			skipped = "true"
		}

		row := []string{
			repo.Repo.FullName,
			repo.Repo.Owner,
			compliant,
			itoa(len(repo.Violations)),
			itoa(len(repo.Warnings)),
			repo.Error,
			skipped,
			repo.SkipReason,
			itoa64(repo.ScanTimeMs),
		}
		if err := w.Write(row); err != nil {
			return nil, err
		}
	}

	w.Flush()
	if err := w.Error(); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// CSVViolationsFormatter generates detailed CSV with one row per violation.
type CSVViolationsFormatter struct{}

// Format generates a detailed CSV report with violation details.
func (f *CSVViolationsFormatter) Format(result *model.ComplianceResult) ([]byte, error) {
	var buf bytes.Buffer
	w := csv.NewWriter(&buf)

	// Write header
	header := []string{
		"repo",
		"org",
		"policy",
		"rule",
		"severity",
		"message",
		"remediation",
		"file",
		"line",
	}
	if err := w.Write(header); err != nil {
		return nil, err
	}

	// Write violation rows
	for _, repo := range result.Repos {
		if len(repo.Violations) == 0 {
			// Write a row even for compliant repos (no violations)
			row := []string{
				repo.Repo.FullName,
				repo.Repo.Owner,
				"",
				"",
				"",
				"",
				"",
				"",
				"",
			}
			if err := w.Write(row); err != nil {
				return nil, err
			}
			continue
		}

		for _, v := range repo.Violations {
			row := []string{
				repo.Repo.FullName,
				repo.Repo.Owner,
				v.Policy,
				v.Rule,
				string(v.Severity),
				v.Message,
				v.Remediation,
				v.File,
				itoa(v.Line),
			}
			if err := w.Write(row); err != nil {
				return nil, err
			}
		}
	}

	w.Flush()
	if err := w.Error(); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func itoa(i int) string {
	if i == 0 {
		return "0"
	}
	return intToString(i)
}

func itoa64(i int64) string {
	if i == 0 {
		return "0"
	}
	return int64ToString(i)
}

func intToString(n int) string {
	if n == 0 {
		return "0"
	}
	neg := false
	if n < 0 {
		neg = true
		n = -n
	}
	var digits []byte
	for n > 0 {
		digits = append([]byte{byte('0' + n%10)}, digits...)
		n /= 10
	}
	if neg {
		digits = append([]byte{'-'}, digits...)
	}
	return string(digits)
}

func int64ToString(n int64) string {
	if n == 0 {
		return "0"
	}
	neg := false
	if n < 0 {
		neg = true
		n = -n
	}
	var digits []byte
	for n > 0 {
		digits = append([]byte{byte('0' + n%10)}, digits...)
		n /= 10
	}
	if neg {
		digits = append([]byte{'-'}, digits...)
	}
	return string(digits)
}
