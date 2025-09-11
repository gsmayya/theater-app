package utils

import (
	"os"
	"testing"
)

func TestGetEnvOrDefault(t *testing.T) {
	tests := []struct {
		name         string
		key          string
		defaultValue string
		envValue     string
		expected     string
	}{
		{
			name:         "environment variable exists and not empty",
			key:          "TEST_VAR",
			defaultValue: "default",
			envValue:     "env_value",
			expected:     "env_value",
		},
		{
			name:         "environment variable exists but empty",
			key:          "TEST_VAR_EMPTY",
			defaultValue: "default",
			envValue:     "",
			expected:     "default",
		},
		{
			name:         "environment variable does not exist",
			key:          "TEST_VAR_MISSING",
			defaultValue: "default",
			envValue:     "",
			expected:     "default",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set up environment
			if tt.envValue != "" {
				os.Setenv(tt.key, tt.envValue)
			} else {
				os.Unsetenv(tt.key)
			}

			// Clean up after test
			defer os.Unsetenv(tt.key)

			result := GetEnvOrDefault(tt.key, tt.defaultValue)
			if result != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, result)
			}
		})
	}
}

func TestGetInt32(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expected    int32
		expectError bool
	}{
		{
			name:        "valid positive integer",
			input:       "123",
			expected:    123,
			expectError: false,
		},
		{
			name:        "valid negative integer",
			input:       "-456",
			expected:    -456,
			expectError: false,
		},
		{
			name:        "zero",
			input:       "0",
			expected:    0,
			expectError: false,
		},
		{
			name:        "large number",
			input:       "2147483647",
			expected:    2147483647,
			expectError: false,
		},
		{
			name:        "invalid string",
			input:       "not_a_number",
			expected:    0,
			expectError: true,
		},
		{
			name:        "empty string",
			input:       "",
			expected:    0,
			expectError: true,
		},
		{
			name:        "float number",
			input:       "123.45",
			expected:    0,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := GetInt32(tt.input)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error for input %s", tt.input)
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error for input %s: %v", tt.input, err)
				}
				if result != tt.expected {
					t.Errorf("Expected %d, got %d", tt.expected, result)
				}
			}
		})
	}
}

func TestGetInt32EdgeCases(t *testing.T) {
	// Test overflow
	t.Run("overflow", func(t *testing.T) {
		_, err := GetInt32("2147483648") // Max int32 + 1
		if err == nil {
			t.Error("Expected error for overflow")
		}
	})

	// Test underflow
	t.Run("underflow", func(t *testing.T) {
		_, err := GetInt32("-2147483649") // Min int32 - 1
		if err == nil {
			t.Error("Expected error for underflow")
		}
	})
}

func TestGetInt32OrPanic(t *testing.T) {
	// Test valid case
	t.Run("valid positive integer", func(t *testing.T) {
		result := GetInt32OrPanic("123")
		if result != 123 {
			t.Errorf("Expected 123, got %d", result)
		}
	})

	// Note: Panic test removed as it causes log.Fatalf to exit the process
	// In production, this function should be avoided in favor of GetInt32
}
