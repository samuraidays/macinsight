package schema

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/samuraidays/macinsight/pkg/types"
)

func TestJSONSchemaGenerator_GenerateReportSchema(t *testing.T) {
	generator := &JSONSchemaGenerator{}
	schema, err := generator.GenerateReportSchema()
	if err != nil {
		t.Fatalf("GenerateReportSchema failed: %v", err)
	}

	// Basic schema structure validation
	if schema["$schema"] == nil {
		t.Error("Schema missing $schema field")
	}
	if schema["title"] == nil {
		t.Error("Schema missing title field")
	}
	if schema["type"] != "object" {
		t.Error("Schema type should be object")
	}

	properties, ok := schema["properties"].(map[string]interface{})
	if !ok {
		t.Fatal("Schema properties should be a map")
	}

	// Check required fields
	requiredFields := []string{"version", "host", "score", "checks"}
	for _, field := range requiredFields {
		if properties[field] == nil {
			t.Errorf("Schema missing required field: %s", field)
		}
	}
}

func TestJSONSchemaGenerator_ValidateReport(t *testing.T) {
	generator := &JSONSchemaGenerator{}

	// Valid report
	validReport := types.Report{
		Version: "v1.0.0",
		Host: types.HostInfo{
			Hostname: "test-host",
			OS: types.OSInfo{
				Product: "macOS",
				Version: "14.2.1",
				Build:   "23C71",
			},
		},
		Score: 85,
		Checks: []types.CheckResult{
			{
				ID:     "sip",
				Title:  "System Integrity Protection enabled",
				Status: "pass",
				Score:  20,
			},
		},
	}

	if err := generator.ValidateReport(validReport); err != nil {
		t.Errorf("Valid report should pass validation: %v", err)
	}

	// Invalid report - invalid score
	invalidReport := validReport
	invalidReport.Score = 150
	if err := generator.ValidateReport(invalidReport); err == nil {
		t.Error("Invalid score should fail validation")
	}

	// Invalid report - invalid check ID
	invalidReport = validReport
	invalidReport.Checks[0].ID = "invalid"
	if err := generator.ValidateReport(invalidReport); err == nil {
		t.Error("Invalid check ID should fail validation")
	}
}

func TestJSONSchemaGenerator_WriteSchema(t *testing.T) {
	generator := &JSONSchemaGenerator{}
	var buf bytes.Buffer

	if err := generator.WriteSchema(&buf); err != nil {
		t.Fatalf("WriteSchema failed: %v", err)
	}

	// Check that output is valid JSON
	var result map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &result); err != nil {
		t.Fatalf("Generated schema is not valid JSON: %v", err)
	}

	// Check basic structure
	if result["$schema"] == nil {
		t.Error("Generated schema missing $schema")
	}
}

func TestJSONSchemaValidator_ValidateJSON(t *testing.T) {
	validator, err := NewValidator()
	if err != nil {
		t.Fatalf("NewValidator failed: %v", err)
	}

	// Valid JSON
	validJSON := `{
		"version": "v1.0.0",
		"host": {
			"hostname": "test-host",
			"os": {
				"product": "macOS",
				"version": "14.2.1",
				"build": "23C71"
			}
		},
		"score": 85,
		"checks": [
			{
				"id": "sip",
				"title": "System Integrity Protection enabled",
				"status": "pass",
				"score": 20
			}
		]
	}`

	if err := validator.ValidateJSON([]byte(validJSON)); err != nil {
		t.Errorf("Valid JSON should pass validation: %v", err)
	}

	// Invalid JSON - missing required field
	invalidJSON := `{"version": "v1.0.0"}`
	if err := validator.ValidateJSON([]byte(invalidJSON)); err == nil {
		t.Error("Invalid JSON should fail validation")
	}
}
