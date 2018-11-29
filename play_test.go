package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"testing"
)

// To see the println results use
// go test -v -timeout 30s github.com/blobdon/blurb -run ^TestBig$
func TestBig(t *testing.T) {
	// get list of test json files
	directory := "/Users/blobdon/blurbdata/test/json/"
	files, err := ioutil.ReadDir(directory)
	if err != nil {
		log.Fatal(err)
	}

	matches := 0
	empty := struct{}{}
	isbns := map[string]struct{}{}
	authors := map[string]struct{}{}
	books := map[string]book{}

	for _, file := range files {

		//get boook from json file
		jsonbytes, err := ioutil.ReadFile(directory + file.Name())
		if err != nil {
			fmt.Println("Error opening file", err)
			continue
		}
		var newb book
		err = json.Unmarshal(jsonbytes, &newb)
		if err != nil {
			fmt.Println("Error unmarshalling json", err)
			continue
		}

		// New book's reviews contain any of all other authors?
		for a := range authors {
			// don't check reviews for a book's own author
			if a == newb.Author {
				continue
			}
			match, err := containsAuthor(newb.ReviewText, a)
			if err != nil {
				t.Errorf("Error matching book: %s author: %s - %s", newb.ISBN13, a, err)
			}
			if match {
				matches++
				fmt.Printf("Found old author %s in new reviews of %s: %s.\n", a, newb.ISBN13, newb.Title)
			}
		}

		// All other books reviews contain the new author?
		for _, oldb := range books {
			// don't check reviews for a book's own author
			if oldb.Author == newb.Author {
				continue
			}

			match, err := containsAuthor(oldb.ReviewText, newb.Author)
			if err != nil {
				t.Errorf("Error matching book: %s author: %s - %s", oldb.ISBN13, newb.Author, err)
			}
			if match {
				matches++
				fmt.Printf("Found new author %s in old reviews of %s: %s.\n", newb.Author, oldb.ISBN13, oldb.Title)
			}
		}

		// add book, isbn, author, to maps
		// this is done at the end to help avoid self-comparisons
		books[newb.ISBN13] = newb
		isbns[newb.ISBN13] = empty
		authors[newb.Author] = empty

	}
	fmt.Println(matches, "matches found")
}

// func TestBuildJson(t *testing.T) {
// 	// get list of test html files
// 	directory := "/Users/blobdon/blurbdata/test/html/"
// 	files, err := ioutil.ReadDir(directory)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	// loop over files, unmarshal json to book struct, assign ID and save book to db
// 	for _, file := range files {
// 		var b book
// 		b, err := htmlToBookBN(directory+file.Name(), b)
// 		if err != nil {
// 			fmt.Println("html error", err)
// 		}
// 		bjson := fmt.Sprintf("/Users/blobdon/blurbdata/test/json/%s.json", b.ISBN13)
// 		err = bookToFile(b, bjson)
// 		if err != nil {
// 			fmt.Println("json error", err)
// 		}
// 	}
// }
