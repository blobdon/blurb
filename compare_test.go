package main

import (
	"testing"
)

func TestContainsAuthor(t *testing.T) {
	testfile := "/Users/blobdon/blurbdata/test/json/9780374533557.json"
	b := book{}
	b, err := fileToBook(testfile)
	if err != nil {
		t.Errorf("containsAuthor fail - can't open json: %s", testfile)
	}

	authorY := "Richard Thaler"
	authorN := "Susan Cain"
	yes, err := containsAuthor(b.ReviewText, authorY)
	if !yes {
		t.Errorf("containsAuthor fail, failed to find %q", authorY)
	}
	if err != nil {
		t.Errorf("containsAuthor fail, error with author %q: %q", authorY, err)
	}
	no, err := containsAuthor(b.ReviewText, authorN)
	if no {
		t.Errorf("containsAuthor fail, incorrectly found %q", authorN)
	}
	if err != nil {
		t.Errorf("containsAuthor fail, error with author %q: %q", authorN, err)
	}
}
