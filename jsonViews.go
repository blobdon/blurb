package main

import (
	"encoding/json"
	"log"
	"strings"

	"github.com/boltdb/bolt"
)

type link struct {
	Key  string
	Text string
}

type bookView struct {
	ISBN13        string
	Title         string
	Author        string
	ReviewAuthors []link // name, blankfornow(formatted for author page url)
}

type authorView struct {
	Name          string
	BooksAuthored []link // isbn, title
	BooksReviewed []link // isbn, title
}

// buildJSONViews will build a json directory, at dirpath, of the
// blurb data as viewed in the website/webapp, i.e. individual books and authors
// The resulting json directory is a potential static-file API for a webapp
func buildJSONViews(db *bolt.DB, dirpath string) error {
	// add final slash to dirpath if not provided (more sophisticated = use filepath, etc)
	if dirpath[len(dirpath)-1:] != "/" {
		dirpath += "/"
	}

	isbn13Map, err := newISBN13Map(db)
	if err != nil {
		return err
	}

	// avMap will be used to store authorView for each author to be updated during range over books
	// this will allow all av's to be made in one pass of books db
	// - map fields in authorViews must be initialised before any additions
	avMap := map[string]authorView{}

	// build book/author views
	// - may benefit from concurrency with large # of views being saved to files, but benchmark first
	for i := range isbn13Map {
		// get bookview from bolt, using isbn of book
		bv := bookView{}
		bv, err := bookViewFromBolt(i, db)
		if err != nil {
			log.Println("Error getting bookview from bolt for isbn ", i, err)
		}

		err = saveBookView(bv, dirpath)
		if err != nil {
			log.Println("Error saving bookview for isbn ", i, err)
		}

		// If author already in avmap, update with new info, otherwise add name and new info
		// For author
		if av, ok := avMap[bv.Author]; ok {
			av.BooksAuthored = append(av.BooksAuthored, link{bv.ISBN13, bv.Title})
			avMap[bv.Author] = av
		} else {
			av.Name = bv.Author
			av.BooksAuthored = append(av.BooksAuthored, link{bv.ISBN13, bv.Title})
			avMap[bv.Author] = av
		}
		// For each reviewer
		for _, ra := range bv.ReviewAuthors {
			if av, ok := avMap[ra.Key]; ok {
				av.BooksReviewed = append(av.BooksReviewed, link{bv.ISBN13, bv.Title})
				avMap[ra.Key] = av
			} else {
				av.Name = ra.Key
				av.BooksReviewed = append(av.BooksReviewed, link{bv.ISBN13, bv.Title})
				avMap[ra.Key] = av
			}
		}
	}
	for _, av := range avMap {
		err = saveAuthorView(av, dirpath)
		if err != nil {
			log.Println("Error saving authorview for name ", av.Name, err)
		}
	}
	return nil
}

func saveAuthorView(av authorView, dirpath string) error {
	// does this filename need to be built using filepath pkg? probably yes
	filename := dirpath + "authors/" + strings.Replace(av.Name, " ", "-", -1) + ".json"
	return structToFile(av, filename)
}

func saveBookView(bv bookView, dirpath string) error {
	// does this filename need to be built using filepath pkg? probably yes
	filename := dirpath + "books/" + bv.ISBN13 + ".json"
	return structToFile(bv, filename)
}

// bookViewFromBolt gets a bookView from bolt, using the stored book keyed by ISBN13
// this is possible because bookview is an assignable subset of fields of book
func bookViewFromBolt(isbn13 string, db *bolt.DB) (bookView, error) {
	b := book{}
	err := db.View(func(tx *bolt.Tx) error {
		bkt := tx.Bucket([]byte("books"))
		j := bkt.Get([]byte(isbn13))
		return json.Unmarshal(j, &b)
	})

	// copy book details to bookview
	bv := bookView{
		Title:  b.Title,
		ISBN13: b.ISBN13,
		Author: b.Author,
	}
	for ra := range b.ReviewAuthors {
		bv.ReviewAuthors = append(bv.ReviewAuthors, link{Key: ra})
	}

	return bv, err
}
