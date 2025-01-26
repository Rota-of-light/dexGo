package main

import "testing"

func TestCleanInput(t *testing.T) {
	cases := []struct {
		input	string
		expected []string
	}{
		{
			input:	"  hello  world  ",
			expected:	[]string{"hello", "world"},
		},
		{
			input:	"  Hello  World  ",
			expected:	[]string{"hello", "world"},
		},
		{
			input:	"  gOOdByE  wORLd  ",
			expected:	[]string{"goodbye", "world"},
		},
	}
	for _, c := range cases {
		actual := cleanInput(c.input)
		if len(actual) != len(c.expected) {
			t.Errorf("Expected length %d but got %d for input %q", len(c.expected), len(actual), c.input)
			continue
		}
		for i := range actual {
			if actual[i] != c.expected[i] {
				t.Errorf("Expected %q but got %q at index %d for input %q", c.expected[i], actual[i], i, c.input)
			}
		}
	}
}