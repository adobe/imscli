// Copyright 2025 Adobe. All rights reserved.
// This file is licensed to you under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License. You may obtain a copy
// of the License at http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software distributed under
// the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR REPRESENTATIONS
// OF ANY KIND, either express or implied. See the License for the specific language
// governing permissions and limitations under the License.

package pretty

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func loadExpected(t *testing.T, filename string) string {
	t.Helper()
	data, err := os.ReadFile(filepath.Join("testdata", filename))
	if err != nil {
		t.Fatalf("failed to read golden file %s: %v", filename, err)
	}
	return strings.TrimSuffix(string(data), "\n")
}

func TestJSON(t *testing.T) {
	tests := []struct {
		name  string
		input string
		file  string
	}{
		{name: "compact object is prettified", input: `{"name":"John","age":30}`, file: "compact_object.json"},
		{name: "already indented JSON is normalized", input: "{\n    \"key\": \"value\"\n}", file: "already_indented.json"},
		{name: "compact array is prettified", input: `[{"id":1},{"id":2}]`, file: "compact_array.json"},
		{name: "empty object", input: `{}`, file: "empty_object.json"},
		{name: "empty array", input: `[]`, file: "empty_array.json"},
		{name: "nested objects", input: `{"a":{"b":{"c":1}}}`, file: "nested_objects.json"},
		{name: "non-JSON is returned as-is", input: "this is not JSON", file: "non_json.txt"},
		{name: "empty string is returned as-is", input: "", file: "empty_string.txt"},
		{name: "unicode characters", input: `{"greeting":"こんにちは","emoji":"🎉","accented":"café"}`, file: "unicode.json"},
		{name: "special characters", input: `{"html":"\u003cscript\u003ealert('xss')\u003c/script\u003e","newlines":"line1\nline2","tabs":"col1\tcol2","quotes":"she said \"hello\""}`, file: "special_chars.json"},
		{name: "null and bool values", input: `{"value":null,"enabled":true,"disabled":false}`, file: "null_and_bool.json"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expected := loadExpected(t, tt.file)
			got := JSON(tt.input)
			if got != expected {
				t.Errorf("JSON(%q)\ngot:\n%s\nexpected:\n%s", tt.input, got, expected)
			}
		})
	}
}
