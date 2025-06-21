package internal

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

type Case struct {
	fileInput string
	expected  string
}

func TestDecoder(t *testing.T) {
	testCases := []Case{
		{
			fileInput: "files/listing1",
			expected:  "files/listing1_expected",
		},
		{
			fileInput: "files/listing39",
			expected:  "files/listing39_expected",
		},
		{
			fileInput: "files/listing40",
			expected:  "files/listing40_expected",
		},
		{
			fileInput: "files/listing41",
			expected:  "files/listing41_expected",
		},
	}

	for _, test := range testCases {
		t.Run(test.fileInput, func(t *testing.T) {
			decoder := NewDecoder(test.fileInput)
			decoder.Decode(false, false)
			b, err := os.ReadFile(test.expected)
			if err != nil {
				t.Error(err)
			}
			expected := string(b)

			output := bytes.Buffer{}
			if err := decoder.Disassemble(&output); err != nil {
				t.Errorf("error found: %v", err)
			}

			if diff := cmp.Diff(strings.TrimSpace(expected), strings.TrimSpace(output.String())); diff != "" {
				t.Errorf("Mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
