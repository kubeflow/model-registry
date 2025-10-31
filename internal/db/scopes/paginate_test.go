package scopes

import (
	"encoding/base64"
	"testing"

	"github.com/kubeflow/model-registry/internal/db/models"
	"github.com/stretchr/testify/assert"
)

// TestInputValidation ensures input validation works correctly
func TestInputValidation(t *testing.T) {
	tests := []struct {
		name        string
		tablePrefix string
		expected    bool
		description string
	}{
		{
			name:        "Valid table prefix",
			tablePrefix: "TestTable",
			expected:    true,
			description: "Valid alphanumeric table name should be accepted",
		},
		{
			name:        "Valid table prefix with underscores",
			tablePrefix: "test_table_name",
			expected:    true,
			description: "Valid table name with underscores should be accepted",
		},
		{
			name:        "Empty table prefix",
			tablePrefix: "",
			expected:    true,
			description: "Empty table prefix should be valid",
		},
		{
			name:        "Invalid table prefix with special chars",
			tablePrefix: "test-table!",
			expected:    false,
			description: "Table prefix with special characters should be rejected",
		},
		{
			name:        "Invalid table prefix starting with number",
			tablePrefix: "123table",
			expected:    false,
			description: "Table prefix starting with number should be rejected",
		},
		{
			name:        "SQL injection attempt in table prefix",
			tablePrefix: "table'; DROP TABLE users; --",
			expected:    false,
			description: "SQL injection attempt should be rejected",
		},
		{
			name:        "SQL injection attempt with semicolon",
			tablePrefix: "test; DELETE FROM users",
			expected:    false,
			description: "Table prefix with semicolon should be rejected",
		},
		{
			name:        "SQL injection attempt with quotes",
			tablePrefix: "test'table",
			expected:    false,
			description: "Table prefix with quotes should be rejected",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isValidTablePrefix(tt.tablePrefix)
			assert.Equal(t, tt.expected, result, tt.description)
		})
	}
}

// TestColumnValidation ensures orderBy column validation works correctly
func TestColumnValidation(t *testing.T) {
	tests := []struct {
		name        string
		orderBy     string
		expected    string
		description string
	}{
		{
			name:        "Valid CREATE_TIME column",
			orderBy:     "CREATE_TIME",
			expected:    "create_time_since_epoch",
			description: "CREATE_TIME should map to create_time_since_epoch",
		},
		{
			name:        "Valid LAST_UPDATE_TIME column",
			orderBy:     "LAST_UPDATE_TIME",
			expected:    "last_update_time_since_epoch",
			description: "LAST_UPDATE_TIME should map to last_update_time_since_epoch",
		},
		{
			name:        "Valid ID column",
			orderBy:     "ID",
			expected:    "id",
			description: "ID should map to id",
		},
		{
			name:        "Valid id column (lowercase)",
			orderBy:     "id",
			expected:    "id",
			description: "id should map to id",
		},
		{
			name:        "Invalid column",
			orderBy:     "invalid_column",
			expected:    "id",
			description: "Invalid column should fallback to default (id)",
		},
		{
			name:        "SQL injection attempt",
			orderBy:     "'; DROP TABLE users; --",
			expected:    "id",
			description: "SQL injection attempt should fallback to default",
		},
		{
			name:        "SQL injection with UNION",
			orderBy:     "id UNION SELECT * FROM users",
			expected:    "id",
			description: "SQL injection with UNION should fallback to default",
		},
		{
			name:        "Empty orderBy",
			orderBy:     "",
			expected:    "id",
			description: "Empty orderBy should fallback to default",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, ok := allowedOrderByColumns[tt.orderBy]
			if !ok {
				result = models.DefaultOrderBy
			}
			assert.Equal(t, tt.expected, result, tt.description)
		})
	}
}

// TestSortOrderValidation ensures sortOrder validation works correctly
func TestSortOrderValidation(t *testing.T) {
	tests := []struct {
		name        string
		sortOrder   string
		expected    string
		description string
	}{
		{
			name:        "Valid ASC sort order",
			sortOrder:   "ASC",
			expected:    "ASC",
			description: "ASC should be valid",
		},
		{
			name:        "Valid DESC sort order",
			sortOrder:   "DESC",
			expected:    "DESC",
			description: "DESC should be valid",
		},
		{
			name:        "Invalid sort order",
			sortOrder:   "INVALID",
			expected:    "",
			description: "INVALID should be rejected and default to empty",
		},
		{
			name:        "SQL injection in sort order",
			sortOrder:   "ASC; DROP TABLE users; --",
			expected:    "",
			description: "SQL injection attempt should be rejected and default to empty",
		},
		{
			name:        "Lowercase asc",
			sortOrder:   "asc",
			expected:    "",
			description: "Lowercase asc should be rejected and default to empty",
		},
		{
			name:        "Empty sort order",
			sortOrder:   "",
			expected:    "",
			description: "Empty sort order should be rejected and default to empty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := allowedSortOrders[tt.sortOrder]
			assert.Equal(t, tt.expected, result, tt.description)
		})
	}
}

// TestCursorDecoding ensures cursor decoding is safe
func TestCursorDecoding(t *testing.T) {
	tests := []struct {
		name        string
		token       string
		expectError bool
		description string
	}{
		{
			name:        "Valid cursor",
			token:       createValidCursor(),
			expectError: false,
			description: "Valid cursor should decode successfully",
		},
		{
			name:        "Invalid base64",
			token:       "invalid-base64!",
			expectError: true,
			description: "Invalid base64 should return error",
		},
		{
			name:        "Malicious cursor with SQL injection",
			token:       createMaliciousCursor(),
			expectError: false,
			description: "Malicious cursor should decode but be safely handled",
		},
		{
			name:        "Cursor with invalid format",
			token:       base64.StdEncoding.EncodeToString([]byte("invalid:format:extra")),
			expectError: true,
			description: "Cursor with wrong format should return error",
		},
		{
			name:        "Cursor with non-numeric ID",
			token:       base64.StdEncoding.EncodeToString([]byte("notanumber:value")),
			expectError: true,
			description: "Cursor with non-numeric ID should return error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cursor, err := DecodeCursor(tt.token)
			if tt.expectError {
				assert.Error(t, err, tt.description)
				assert.Nil(t, cursor, "Cursor should be nil on error")
			} else {
				assert.NoError(t, err, tt.description)
				assert.NotNil(t, cursor, "Cursor should not be nil on success")
			}
		})
	}
}

// createValidCursor creates a valid cursor for testing
func createValidCursor() string {
	cursor := "123:test_value"
	return base64.StdEncoding.EncodeToString([]byte(cursor))
}

// createMaliciousCursor creates a cursor with malicious SQL injection payload
func createMaliciousCursor() string {
	maliciousValue := "'; DROP TABLE test_table; --"
	cursor := "1:" + maliciousValue
	return base64.StdEncoding.EncodeToString([]byte(cursor))
}
