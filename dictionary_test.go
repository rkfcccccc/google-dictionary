package dictionary

import (
	"context"
	"errors"
	"testing"
)

func TestGetWordEntry(t *testing.T) {
	testCases := []struct {
		word string
		err  error
	}{
		{"run", nil}, {"equality", nil}, {"assuming", nil}, {"test", nil},
		{"unknown", nil}, {"error", nil}, {"entry", nil}, {"word", nil}, {"messy';:", nil},
		{"dajshd9as", ErrNoDefinitionsFound}, {";-=;=p-=1-2=3p]]", ErrNoDefinitionsFound},
		{"309das-§∞¶¢£∞•ªº–dsd;a';;:a…¬˚∆˙©¥¨†®§•ª˙º∆ˆ", ErrNoDefinitionsFound},
	}

	for i, testCase := range testCases {
		_, err := GetWordEntry(context.Background(), "en", testCase.word)

		if !errors.Is(err, testCase.err) {
			t.Logf("test case %d (%s): error should be %q got %q", i, testCase.word, testCase.err, err)
			t.Fail()
		}
	}
}
