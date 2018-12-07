package main

import (
	"log"
	"strings"

	"github.com/boltdb/bolt"
)

type review struct {
	Reviewer string
	Title    string
	Author   string
	ISBN13   string
}

type bookView struct {
	ISBN13    string
	Title     string
	Author    string
	ReviewsIn []review // name, blankfornow(formatted for author page url)
}

type authorView struct {
	Name       string
	ReviewsOut []review // isbn, title
	ReviewsIn  []review
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
	aMap, err := newAuthorMap(db)
	if err != nil {
		return err
	}

	// avMap will be used to store authorView for each author to be updated during range over books
	// this will allow all av's to be made in one pass of books db
	// Initialize each av's Review slices with empty slice literals - this ensure that json.Marshal will
	// produce ReviewsIn: [] (empty slice) instead of ReviewsIn: null (nil slice)
	avMap := map[string]authorView{}
	for a := range aMap {
		avMap[a] = authorView{
			Name:       a,
			ReviewsIn:  []review{},
			ReviewsOut: []review{},
		}
	}

	// build book/author views
	// - may benefit from concurrency with large # of views being saved to files, but benchmark first
	for i := range isbn13Map {
		b, err := bookFromBolt(i, db)
		if err != nil {
			log.Println("Error getting bookview from bolt for isbn ", i, err)
		}
		// copy book details to bookview
		bv := bookView{
			Title:     b.Title,
			ISBN13:    b.ISBN13,
			Author:    b.Author,
			ReviewsIn: []review{},
		}

		// For each reviewer
		for ra := range b.ReviewAuthors {
			r := review{
				Reviewer: ra,
				Title:    b.Title,
				ISBN13:   b.ISBN13,
				Author:   b.Author,
			}

			// In for current book
			bv.ReviewsIn = append(bv.ReviewsIn, r)

			// In for current author
			av := avMap[b.Author]
			av.ReviewsIn = append(av.ReviewsIn, r)
			avMap[b.Author] = av

			// Out for the review author
			av = avMap[ra]
			av.ReviewsOut = append(av.ReviewsOut, r)
			avMap[ra] = av
		}

		err = saveBookView(bv, dirpath)
		if err != nil {
			log.Println("Error saving bookview for isbn ", i, err)
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
