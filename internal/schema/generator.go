package schema

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/samuraidays/macinsight/pkg/types"
)

// JSONSchemaGenerator generates JSON Schema from Go structs
type JSONSchemaGenerator struct{}

// GenerateReportSchema generates JSON Schema for types.Report
func (g *JSONSchemaGenerator) GenerateReportSchema() (map[string]interface{}, error) {
	schema := map[string]interface{}{
		"$schema":     "https://json-schema.org/draft/2020-12/schema",
		"$id":         "https://github.com/samuraidays/macinsight/schema/report.json",
		"title":       "macinsight Security Audit Report",
		"description": "JSON schema for macinsight security audit report output",
		"type":        "object",
		"properties": map[string]interface{}{
			"version": map[string]interface{}{
				"type":        "string",
				"description": "macinsight version",
				"pattern":     "^v[0-9]+\\.[0-9]+\\.[0-9]+(-[a-zA-Z0-9]+)?$",
			},
			"host": map[string]interface{}{
				"type":        "object",
				"description": "Host information",
				"properties": map[string]interface{}{
					"hostname": map[string]interface{}{
						"type":        "string",
						"description": "Hostname of the audited system",
					},
					"os": map[string]interface{}{
						"type":        "object",
						"description": "Operating system information",
						"properties": map[string]interface{}{
							"product": map[string]interface{}{
								"type":        "string",
								"description": "OS product name",
								"const":       "macOS",
							},
							"version": map[string]interface{}{
								"type":        "string",
								"description": "OS version",
								"pattern":     "^[0-9]+\\.[0-9]+\\.[0-9]+$",
							},
							"build": map[string]interface{}{
								"type":        "string",
								"description": "OS build number",
								"pattern":     "^[0-9]+[A-Z][0-9]+[A-Z]?[0-9]*$",
							},
						},
						"required": []string{"product", "version", "build"},
					},
				},
				"required": []string{"hostname", "os"},
			},
			"score": map[string]interface{}{
				"type":        "integer",
				"description": "Total security score (0-100)",
				"minimum":     0,
				"maximum":     100,
			},
			"checks": map[string]interface{}{
				"type":        "array",
				"description": "Security check results",
				"items": map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"id": map[string]interface{}{
							"type":        "string",
							"description": "Check identifier",
							"enum":        []string{"sip", "gatekeeper", "filevault", "firewall", "autologin", "osupdate"},
						},
						"title": map[string]interface{}{
							"type":        "string",
							"description": "Human-readable check title",
						},
						"status": map[string]interface{}{
							"type":        "string",
							"description": "Check result status",
							"enum":        []string{"pass", "fail", "warn", "unknown"},
						},
						"score": map[string]interface{}{
							"type":        "integer",
							"description": "Points awarded for this check",
							"minimum":     0,
							"maximum":     20,
						},
						"evidence": map[string]interface{}{
							"type":        "object",
							"description": "Evidence data from the check",
							"additionalProperties": map[string]interface{}{
								"type": "string",
							},
						},
						"recommendation": map[string]interface{}{
							"type":        "string",
							"description": "Recommendation for improvement",
						},
					},
					"required": []string{"id", "title", "status", "score"},
				},
			},
		},
		"required": []string{"version", "host", "score", "checks"},
	}

	return schema, nil
}

// WriteSchema writes the JSON schema to a writer
func (g *JSONSchemaGenerator) WriteSchema(w io.Writer) error {
	schema, err := g.GenerateReportSchema()
	if err != nil {
		return err
	}

	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(schema)
}

// ValidateReport validates a Report against the schema
func (g *JSONSchemaGenerator) ValidateReport(report types.Report) error {
	// Basic validation
	if report.Version == "" {
		return fmt.Errorf("version is required")
	}
	if report.Host.Hostname == "" {
		return fmt.Errorf("hostname is required")
	}
	if report.Score < 0 || report.Score > 100 {
		return fmt.Errorf("score must be between 0 and 100, got %d", report.Score)
	}

	// Validate checks
	validStatuses := map[string]bool{"pass": true, "fail": true, "warn": true, "unknown": true}
	validIDs := map[string]bool{"sip": true, "gatekeeper": true, "filevault": true, "firewall": true, "autologin": true, "osupdate": true}

	for _, check := range report.Checks {
		if !validIDs[check.ID] {
			return fmt.Errorf("invalid check ID: %s", check.ID)
		}
		if !validStatuses[check.Status] {
			return fmt.Errorf("invalid status: %s", check.Status)
		}
		if check.Score < 0 || check.Score > 20 {
			return fmt.Errorf("check score must be between 0 and 20, got %d for %s", check.Score, check.ID)
		}
	}

	return nil
}
