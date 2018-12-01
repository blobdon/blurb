package main

import (
	"log"
	"strings"

	"github.com/boltdb/bolt"
)

type bookView struct {
	ISBN13        string
	Title         string
	Author        string
	ReviewAuthors map[string]interface{} //[name]=blankfornow(formatted for author page url)
}

type authorView struct {
	Name          string
	BooksAuthored map[string]string // [ISBN]=title
	BooksReviewed map[string]string // [ISBN]=title
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
		bv := bookView{ReviewAuthors: map[string]interface{}{}}
		bv, err := bookViewFromBolt(i, db)
		if err != nil {
			log.Println("Error getting bookview from bolt for isbn ", i, err)
		}

		err = saveBookView(bv, dirpath)
		if err != nil {
			log.Println("Error saving bookview for isbn ", i, err)
		}

		// If author already in avmap, update with new info, otherwise initialize maps and add new info
		// Make sure maps are initilised whether author is added to map first as author of reviewer
		// For author
		if av, ok := avMap[bv.Author]; ok {
			av.BooksAuthored[bv.ISBN13] = bv.Title
			avMap[bv.Author] = av
		} else {
			// add name and initilize map fields
			av.Name = bv.Author
			av.BooksAuthored = map[string]string{}
			av.BooksReviewed = map[string]string{}
			av.BooksAuthored[bv.ISBN13] = bv.Title // add book authored
			avMap[bv.Author] = av
		}
		// For each reviewer
		for ra := range bv.ReviewAuthors {
			if av, ok := avMap[ra]; ok {
				av.BooksReviewed[bv.ISBN13] = bv.Title
				avMap[ra] = av
			} else {
				// add name and initilize map fields
				av.Name = ra
				av.BooksAuthored = map[string]string{}
				av.BooksReviewed = map[string]string{}
				av.BooksReviewed[bv.ISBN13] = bv.Title // add book reviewed
				avMap[ra] = av
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
