package action

import (
	"dbaker/pkg/model"
	"encoding/json"
	"os"
	"testing"
)

func TestSplitTableName(t *testing.T) {
	testCases := []struct {
		name           string
		input          string
		expectedName   string
		expectedSchema string
	}{
		{
			name:           "valid schema and table name",
			input:          "public.users",
			expectedName:   "users",
			expectedSchema: "public",
		},
		{
			name:           "missing schema",
			input:          "users",
			expectedName:   "",
			expectedSchema: "",
		},
		{
			name:           "too many parts",
			input:          "db.public.users",
			expectedName:   "",
			expectedSchema: "",
		},
		{
			name:           "empty input string",
			input:          "",
			expectedName:   "",
			expectedSchema: "",
		},
		{
			name:           "leading dot",
			input:          ".users",
			expectedName:   "users",
			expectedSchema: "",
		},
		{
			name:           "trailing dot",
			input:          "public.",
			expectedName:   "",
			expectedSchema: "public",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			name, schema := splitTableName(tc.input)
			if name != tc.expectedName || schema != tc.expectedSchema {
				t.Errorf("splitTableName(%q) = %q, %q; want %q, %q", tc.input, name, schema, tc.expectedName, tc.expectedSchema)
			}
		})
	}
}

func TestWriteJson(t *testing.T) {
	testCases := []struct {
		name    string
		tables  []*model.Table
		wantErr bool
	}{
		{
			name:    "empty tables slice",
			tables:  []*model.Table{},
			wantErr: false,
		},
		{
			name: "single table",
			tables: []*model.Table{
				{
					Name:   "users",
					Schema: "public",
					Columns: []model.Column{
						{Name: "id", Typ: model.Int, IsGenerated: true},
						{Name: "name", Typ: model.Varchar, MaxLength: 255, IsNullable: false},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "multiple tables",
			tables: []*model.Table{
				{
					Name:   "users",
					Schema: "public",
					Columns: []model.Column{
						{Name: "id", Typ: model.Int, IsGenerated: true},
						{Name: "name", Typ: model.Varchar, MaxLength: 255, IsNullable: false},
					},
				},
				{
					Name:   "products",
					Schema: "public",
					Columns: []model.Column{
						{Name: "product_id", Typ: model.UUID, IsUnique: true},
						{Name: "price", Typ: model.Decimal, IsNullable: false},
					},
				},
			},
			wantErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tmpfile, err := os.CreateTemp("", "test_recipe_*.json")
			if err != nil {
				t.Fatalf("failed to create temporary file: %v", err)
			}
			defer os.Remove(tmpfile.Name())
			tmpfile.Close() // Close the file handle so writeJson can open it

			err = writeJson(tmpfile.Name(), tc.tables)
			if (err != nil) != tc.wantErr {
				t.Errorf("writeJson() error = %v, wantErr %v", err, tc.wantErr)
				return
			}

			if !tc.wantErr {
				contents, err := os.ReadFile(tmpfile.Name())
				if err != nil {
					t.Fatalf("failed to read written file: %v", err)
				}

				var readTables []*model.Table
				err = json.Unmarshal(contents, &readTables)
				if err != nil {
					t.Fatalf("failed to unmarshal JSON from file: %v", err)
				}

				// Marshal both original and read tables to canonical JSON for comparison
				// Using json.Marshal (without indent) ensures consistent comparison regardless of formatting
				originalJSON, err := json.Marshal(tc.tables)
				if err != nil {
					t.Fatalf("failed to marshal original tables for comparison: %v", err)
				}
				readJSON, err := json.Marshal(readTables)
				if err != nil {
					t.Fatalf("failed to marshal read tables for comparison: %v", err)
				}

				if string(originalJSON) != string(readJSON) {
					t.Errorf("written JSON does not match original.\nExpected: %s\nGot: %s", string(originalJSON), string(readJSON))
				}
			}
		})
	}
}
