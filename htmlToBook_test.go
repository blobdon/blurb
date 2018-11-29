package main

import (
	"testing"
)

func TestHtmlToBookBN(t *testing.T) {
	// filename := "/Users/blobdon/blurbdata/test/"
	// filename := "/Users/blobdon/blurbdata/html/BNSubjects/psychology/books/9780374533557.html"
	filename := "/Users/blobdon/blurbdata/test/html/9780374533557.html"
	wantb := book{
		ISBN13: "9780374533557",
		Title:  "Thinking, Fast and Slow",
		Author: "Daniel Kahneman",
	}
	gotb := book{}
	gotb, err := htmlToBookBN(filename, gotb)
	if gotb.Title != wantb.Title {
		t.Errorf("htmlToBookBN title fail, wanted %q got %q\nError: %q", wantb.Title, gotb.Title, err)
	}
	if gotb.Author != wantb.Author {
		t.Errorf("htmlToBookBN author fail, wanted %q got %q\nError: %q", wantb.Author, gotb.Author, err)
	}
	if gotb.ISBN13 != wantb.ISBN13 {
		t.Errorf("htmlToBookBN ISBN13 fail, wanted %q got %q\nError: %q", wantb.ISBN13, gotb.ISBN13, err)
	}
}
