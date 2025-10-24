package schema

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/samuraidays/macinsight/pkg/types"
)

// JSONSchemaValidator validates JSON against a schema
type JSONSchemaValidator struct {
	schema map[string]interface{}
}

// NewValidator creates a new validator with the report schema
func NewValidator() (*JSONSchemaValidator, error) {
	generator := &JSONSchemaGenerator{}
	schema, err := generator.GenerateReportSchema()
	if err != nil {
		return nil, err
	}

	return &JSONSchemaValidator{schema: schema}, nil
}

// ValidateFile validates a JSON file against the schema
func (v *JSONSchemaValidator) ValidateFile(filename string) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("failed to read file %s: %w", filename, err)
	}

	return v.ValidateJSON(data)
}

// ValidateJSON validates JSON data against the schema
func (v *JSONSchemaValidator) ValidateJSON(data []byte) error {
	var report types.Report
	if err := json.Unmarshal(data, &report); err != nil {
		return fmt.Errorf("invalid JSON: %w", err)
	}

	// Use the generator's validation logic
	generator := &JSONSchemaGenerator{}
	return generator.ValidateReport(report)
}

// ValidateReader validates JSON from a reader against the schema
func (v *JSONSchemaValidator) ValidateReader(r io.Reader) error {
	data, err := io.ReadAll(r)
	if err != nil {
		return fmt.Errorf("failed to read data: %w", err)
	}

	return v.ValidateJSON(data)
}
