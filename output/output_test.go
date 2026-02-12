// Copyright 2025 Adobe. All rights reserved.
// This file is licensed to you under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License. You may obtain a copy
// of the License at http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software distributed under
// the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR REPRESENTATIONS
// OF ANY KIND, either express or implied. See the License for the specific language
// governing permissions and limitations under the License.

package output

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"testing"
)

// captureStdout captures the output of a function that writes to stdout.
func captureStdout(t *testing.T, fn func()) string {
	t.Helper()
	old := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("failed to create pipe: %v", err)
	}
	os.Stdout = w

	fn()

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	if _, err := io.Copy(&buf, r); err != nil {
		t.Fatalf("failed to read pipe: %v", err)
	}
	return buf.String()
}

// loadExpected reads the expected output from a golden file in testdata/.
func loadExpected(t *testing.T, filename string) string {
	t.Helper()
	data, err := os.ReadFile(filepath.Join("testdata", filename))
	if err != nil {
		t.Fatalf("failed to read golden file %s: %v", filename, err)
	}
	return string(data)
}

func TestPrintPrettyJSON(t *testing.T) {
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
		{name: "non-JSON is printed as-is", input: "this is not JSON", file: "non_json.txt"},
		{name: "empty string is printed as-is", input: "", file: "empty_string.txt"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expected := loadExpected(t, tt.file)

			got := captureStdout(t, func() {
				PrintPrettyJSON(tt.input)
			})

			if got != expected {
				t.Errorf("PrintPrettyJSON(%q)\ngot:\n%s\nexpected output is in testdata/%s", tt.input, got, tt.file)
			}
		})
	}
}
