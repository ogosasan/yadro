package comics

import (
	"testing"
)

func TestIsStopWord(t *testing.T) {
	testTable := []struct {
		word     string
		expected bool
	}{
		{word: "you",
			expected: true,
		},
		{
			word:     "what",
			expected: true,
		},
		{
			word:     "hello",
			expected: false,
		},
	}

	for _, testCase := range testTable {
		result := IsStopWord(testCase.word)

		if result != testCase.expected {
			t.Errorf("Incorrect result. Expect %t, got %t", testCase.expected, result)
		}
	}
}
