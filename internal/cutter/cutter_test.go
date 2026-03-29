package cutter

import (
	"testing"
)

func TestParseFields(t *testing.T) {
	testCases := []struct {
		name  string
		input string
		want  []int
	}{
		{"simple", "2", []int{2}},
		{"range", "1-3", []int{1, 2, 3}},
		{"multiple", "1,3", []int{1, 3}},
		{"single-range", "5-5", []int{5}},
		{"large-single", "10", []int{10}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := parseFields(tc.input)
			if !equalIntSlices(got, tc.want) {
				t.Errorf("parseFields(%q) = %v, want %v", tc.input, got, tc.want)
			}
		})
	}
}

func TestExtractColumns(t *testing.T) {
	c := &Cutter{}
	testCases := []struct {
		fields   []string
		columns  []int
		flagS    bool
		expected []string
		name     string
	}{
		{
			fields:   []string{"a", "b", "c"},
			columns:  []int{2},
			flagS:    false,
			expected: []string{"b"},
			name:     "basic",
		},
		{
			fields:   []string{"a"},
			columns:  []int{2},
			flagS:    false,
			expected: []string{""},
			name:     "out-of-range-no-s",
		},
		{
			fields:   []string{"a"},
			columns:  []int{2},
			flagS:    true,
			expected: nil,
			name:     "out-of-range-with-s",
		},
		{
			fields:   []string{},
			columns:  []int{1},
			flagS:    false,
			expected: []string{""},
			name:     "empty-fields-no-s",
		},
		{
			fields:   []string{"a"},
			columns:  []int{1},
			flagS:    false,
			expected: []string{"a"},
			name:     "exact-len-no-s",
		},
		{
			fields:   []string{"a", "b"},
			columns:  []int{3, 4},
			flagS:    false,
			expected: []string{"", ""},
			name:     "all-missing-multi-no-s",
		},
		{
			fields:   []string{},
			columns:  []int{1},
			flagS:    true,
			expected: nil,
			name:     "empty-with-s",
		},
		{
			fields:   []string{"a1", "b2", "c3", "d4", "e5", "f6", "g7", "h8", "i9", "j10"},
			columns:  []int{5, 8},
			flagS:    false,
			expected: []string{"e5", "h8"},
			name:     "wide-line",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			out := c.extractColumns(tc.fields, tc.columns, tc.flagS)
			if (tc.expected == nil && out != nil) || (tc.expected != nil && !equalSlices(out, tc.expected)) {
				t.Errorf("want %v got %v", tc.expected, out)
			}
		})
	}
}

func equalSlices(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func equalIntSlices(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
