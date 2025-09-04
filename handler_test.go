package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTrimPath(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		expected string
		errorMsg string
	}{
		{"valid path", "/abcdefg", "abcdefg", ""},
		{"with spaces", "/ab cde fg", "abcdefg", ""},
		{"with more spaces", "/a b c d e f  g", "abcdefg", ""},
		{"too long", "/abcdefgh", "", "must be exactly 7 chars"},
		{"too short", "/abc", "", "must be exactly 7 chars"},
		{"empty", "/", "", "must be exactly 7 chars"},
		{"special characters", "/abcdef!", "", "only regular chars are valid"},
		{"only special characters", "/..??..!", "", "only regular chars are valid"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			path, err := trimPath(test.path)

			if err != nil {
				assert.Equal(t, test.errorMsg, err.Error())
				return
			}

			assert.Equal(t, test.expected, path)
		})
	}
}
