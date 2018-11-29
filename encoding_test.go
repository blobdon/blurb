package main

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestBookJSON(t *testing.T) {
	tmpfile, err := ioutil.TempFile("/Users/blobdon/blurbdata/test", "testbook*.json")
	if err != nil {
		t.Errorf("create tempfile: %s", err)
	}

	tmpfilename := tmpfile.Name()
	defer os.Remove(tmpfile.Name()) // clean up

	tmpfile.Close()
	if err != nil {
		t.Errorf("close tempfile: %s", err)
	}

	err = bookToFile(startbook, tmpfilename)
	if err != nil {
		t.Errorf("book to json file: %s", err)
	}

	endbook := book{}
	endbook, err = fileToBook(tmpfilename)
	if err != nil {
		t.Errorf("book from json file: %s", err)
	}

	// TODO replace this with test for incomparable types
	if endbook.Title != startbook.Title {
		t.Errorf("endbook != startbook\n wanted %q got %q\nError: %q", startbook, endbook, err)
	}
}
func TestStructJSON(t *testing.T) {
	tmpfile, err := ioutil.TempFile("/Users/blobdon/blurbdata/test", "testbook*.json")
	if err != nil {
		t.Errorf("create tempfile: %s", err)
	}

	tmpfilename := tmpfile.Name()
	defer os.Remove(tmpfile.Name()) // clean up

	tmpfile.Close()
	if err != nil {
		t.Errorf("close tempfile: %s", err)
	}

	err = structToFile(startbook, tmpfilename)
	if err != nil {
		t.Errorf("book to json file: %s", err)
	}

	endbook := book{}
	endbook, err = fileToBook(tmpfilename)
	if err != nil {
		t.Errorf("book from json file: %s", err)
	}

	// TODO replace this with test for incomparable types
	if endbook.Title != startbook.Title {
		t.Errorf("endbook != startbook\n wanted %q got %q\nError: %q", startbook, endbook, err)
	}
}
func TestMapJSON(t *testing.T) {
	tmpfile, err := ioutil.TempFile("/Users/blobdon/blurbdata/test", "testmap*.json")
	if err != nil {
		t.Errorf("create tempfile: %s", err)
	}

	tmpfilename := tmpfile.Name()
	defer os.Remove(tmpfile.Name()) // clean up

	tmpfile.Close()
	if err != nil {
		t.Errorf("close tempfile: %s", err)
	}

	err = mapToFile(startmap, tmpfilename)
	if err != nil {
		t.Errorf("map to json file: %s", err)
	}

	endmap := map[string]struct{}{}
	endmap, err = fileToMap(tmpfilename)
	if err != nil {
		t.Errorf("map from json file: %s", err)
	}

	if endmap["testingtesting"] != struct{}{} {
		t.Errorf("endmap != startmap\n wanted %q got %q\nError: %q", startmap, endmap, err)
	}

}

var startbook = book{
	Title:      "Thinking, Fast and Slow",
	Author:     "Daniel Kahneman",
	ISBN13:     "9780374533557",
	ReviewText: "Good Book -me",
}

var startmap = map[string]struct{}{
	"testingtesting": struct{}{},
}
