// Package schema provides embedded JSON schemas for validation.
package schema

import (
	_ "embed"
)

// CheckResultSchema is the JSON Schema for CheckResult, embedded at compile time.
//
//go:embed check_result.schema.json
var CheckResultSchema []byte
