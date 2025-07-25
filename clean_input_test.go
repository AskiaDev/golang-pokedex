package main

import (
	"testing"
)


func TestCleanInput(t *testing.T){
	cases := []struct {
		input string
		expected []string
	} {
		{
			input : " hello world ", 
			expected: []string{"hello", "world"},
		},
	}

	for _, c := range cases {
		actual := cleanInput(c.input)

		for i := range c.expected {
			word := actual[i]
			expectedWord := c.expected[i]

			if word != expectedWord {
				t.Errorf("expected %q but got %q", expectedWord, word)
			}
		}
	}
}