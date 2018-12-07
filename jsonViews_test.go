package main

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"
)

// Test whether json from type book can be unmarshalled into bookView
func TestBookToBookView(t *testing.T) {
	ra := map[string]interface{}{}
	ra["Other Author"] = struct{}{}
	bv := bookView{}
	b := book{
		ISBN13:        "1234567890123",
		Title:         "Testbook Title",
		Author:        "Test Author",
		ReviewAuthors: ra,
		ReviewText:    "Vrey good book -Author",
	}
	j, err := json.Marshal(b)
	if err != nil {
		t.Errorf("Error marshalling book to json: %s", err)
	}
	err = json.Unmarshal(j, &bv)
	if err != nil {
		t.Errorf("Error unmarshalling json to bookView: %s", err)
	}
	fmt.Printf("Book: %v \n BookView: %v \n", b, bv)
}

// TestBuildJSONViews can currently be used to build test views JSON
func TestBuildJSONViews(t *testing.T) {
	// clean test directory
	dir := "/Users/blobdon/blurbdata/test/views/"
	err := os.RemoveAll(dir)
	if err != nil {
		t.Fatalf("Failed to clear test views directory: %s", err)
	}
	for _, d := range [2]string{"books", "authors"} {
		err = os.MkdirAll(dir+d, 0755)
		if err != nil {
			t.Fatalf("Failed to recreate test views directory for %s: %s", d, err)
		}
	}

	db, err := initBolt("/Users/blobdon/blurbdata/test/test.db")
	if err != nil {
		t.Errorf("Error on init bolt db: %s", err)
	}
	defer db.Close()
	err = buildJSONViews(db, dir)
	if err != nil {
		t.Errorf("Error within buildJSONViews: %s", err)
	}
	// add want/get comparisons to check output
}
