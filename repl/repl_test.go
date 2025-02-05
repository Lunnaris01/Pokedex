package main

import (
	"testing"
)

func TestCleanInput(t *testing.T) {
	cases := []struct {
		input    string
		expected []string
		}{
			{
				input:    "  ",
				expected: []string{},
			},
			{
				input:    "  hello  ",
				expected: []string{"hello"},
			},
			{
				input:    "  hello  world  ",
				expected: []string{"hello", "world"},
			},
			{
				input:    "  HellO  World  ",
				expected: []string{"hello", "world"},
			},
		}
	
	for _, c := range cases {
		actual := cleanInput(c.input)
		if len(c.expected) != len(actual){
			t.Errorf("Length missmatch")
		}
		for i := range actual {
			word := actual[i]
			expectedWord := c.expected[i]
			if expectedWord != word {
				t.Errorf("Expected: %v\nActual: %v",expectedWord,word)
			}
		}
	}
	
}
