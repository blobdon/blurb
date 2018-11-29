package main

import (
	"encoding/json"
	"fmt"
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
