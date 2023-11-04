package cmd

import (
	"fmt"
	"testing"
)

func TestExtractStudentNameAndIDRegex(t *testing.T) {
	tests := []struct {
		name         string
		tidy         string
		expectedID   string
		expectedName string
		expectedErr  error
	}{
		{
			name:         "Valid input",
			tidy:         "John Doe 12345",
			expectedID:   "12345",
			expectedName: "John Doe",
			expectedErr:  nil,
		},
		{
			name:         "Input with no name",
			tidy:         "12345",
			expectedID:   "12345",
			expectedName: "",
			expectedErr:  fmt.Errorf("illegal input: 12345, lack of name and id"),
		},
		{
			name:         "Input with no name and no ID",
			tidy:         "220301033刘徐明",
			expectedID:   "220301033",
			expectedName: "刘徐明",
			expectedErr:  nil,
		},
		{
			name:         "Input with valid characters",
			tidy:         "李四 12345",
			expectedID:   "12345",
			expectedName: "李四",
			expectedErr:  nil,
		},
	}

	for _, test := range tests {
		name, id, err := extractStudentNameAndIDRegex(test.tidy)

		if err != nil && err != test.expectedErr {
			t.Errorf("Expected error %v, but got %v", test.expectedErr, err)
		} else {
			if id != test.expectedID {
				t.Errorf("Expected ID %s, but got %s", test.expectedID, id)
			}

			if name != test.expectedName {
				t.Errorf("Expected name %s, but got %s", test.expectedName, name)
			}
		}
	}
}
