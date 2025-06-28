package adapter

import (
	"dbaker/pkg/model"
	"testing"
)

func TestInferPgValPlaceholders(t *testing.T) {
	tests := []struct {
		columnLen int
		expected  string
	}{
		{0, ""},
		{1, "$1"},
		{2, "$1, $2"},
		{3, "$1, $2, $3"},
		{5, "$1, $2, $3, $4, $5"},
	}

	for _, tt := range tests {
		result := inferPgValPlaceholders(tt.columnLen)
		if result != tt.expected {
			t.Errorf("inferPgValPlaceholders(%d) = %q; want %q", tt.columnLen, result, tt.expected)
		}
	}
}

func TestInferColNames(t *testing.T) {
	tests := []struct {
		name     string
		columns  []model.Column
		expected string
	}{
		{
			name:     "no columns",
			columns:  []model.Column{},
			expected: "",
		},
		{
			name: "one column",
			columns: []model.Column{
				{Name: "id"},
			},
			expected: "id",
		},
		{
			name: "multiple columns",
			columns: []model.Column{
				{Name: "id"},
				{Name: "name"},
				{Name: "email"},
			},
			expected: "id, name, email",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := inferColNames(tt.columns)
			if result != tt.expected {
				t.Errorf("inferColNames(%v) = %q; want %q", tt.columns, result, tt.expected)
			}
		})
	}
}
