package main

import (
	"regexp"
	"testing"
)

func Test_buildURL(t *testing.T) {
	type test struct {
		board       string
		postID      string
		expectedURL string
	}

	tests := []test{
		test{
			board:       "lit",
			postID:      "123",
			expectedURL: "https://boards.4channel.org/lit/thread/123",
		},
	}

	for _, tc := range tests {
		got := buildURL(tc.board, tc.postID)
		if got != tc.expectedURL {
			t.Errorf("got=[%s], expected=[%s]", got, tc.expectedURL)
		}
	}
}

func Test_postOccurrencesCount(t *testing.T) {
	type test struct {
		html          string
		postID        string
		expectedCount int
	}

	tests := []test{
		test{
			html:          "9876\nabc9876\n@<h2>ok</h2>9876",
			postID:        "9876",
			expectedCount: 3,
		},
	}

	for _, tc := range tests {
		rgx := regexp.MustCompile(tc.postID)
		got := postOccurrencesCount(rgx, tc.html)
		if got != tc.expectedCount {
			t.Errorf("got=[%d], expected=[%d]", got, tc.expectedCount)
		}
	}
}
