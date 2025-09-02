package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTrimPath(t *testing.T) {
	tests := []struct {
		name         string
		path         string
		expected     string
		withError    bool
		specialChars bool
	}{
		{"valid path", "/abcdefg", "abcdefg", false, false},
		{"with spaces", "/ab cde fg", "abcdefg", false, false},
		{"with more spaces", "/a b c d e f  g", "abcdefg", false, false},
		{"too long", "/abcdefgh", "", true, false},
		{"too short", "/abc", "", true, false},
		{"empty", "/", "", true, false},
		{"special characters", "/abcdef!", "", true, true},
		{"only special characters", "/..??..!", "", true, true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			path, err := trimPath(test.path)

			if test.withError {
				assert.Error(t, err)
				return
			}

			if test.specialChars {
				assert.Error(t, err)
				return
			}

			assert.Equal(t, test.expected, path)
		})
	}
}
