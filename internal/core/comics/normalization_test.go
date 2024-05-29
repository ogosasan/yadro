package comics

import (
	"reflect"
	"testing"
)

func TestNormalization(t *testing.T) {
	testTable := []struct {
		sentence string
		expected []string
	}{
		{sentence: "i'll follow you as long as you are following me",
			expected: []string{"follow", "long"},
		},
		{
			sentence: "follower brings bunch of questions",
			expected: []string{"follow", "bring", "bunch", "question"},
		},
	}

	for _, testCase := range testTable {
		result := Normalization(testCase.sentence)

		if reflect.DeepEqual(result, testCase.expected) == false {
			t.Errorf("Incorrect result. Expect %s, got %s", testCase.expected, result)
		}
	}
}
