package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

type book struct {
	ISBN13     string
	Title      string
	Author     string
	ReviewText string // if needed for performance later, dont marshal this one to json, save separately
	// ReviewAuthors as map will break current tests b/c cant compare books with non-comparable types (ie maps)
	ReviewAuthors map[string]interface{}
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)

	}

}
func run() error {
	if len(os.Args) < 2 {
		return errors.New("Please provide the full path to a directory of new html files to ingest")
	}
	// check directory for new html files
	inHTMLs := []string{}
	inDir := os.Args[1]
	inFiles, err := ioutil.ReadDir(inDir)
	if err != nil {
		exiterr := fmt.Sprintf("Error accessing given filepath - exiting: %s", err)
		return errors.New(exiterr)
	}
	for _, f := range inFiles {
		fp := inDir + f.Name()
		if filepath.Ext(fp) == ".html" {
			inHTMLs = append(inHTMLs, fp)
		}
	}
	if len(inHTMLs) < 1 {
		return errors.New("No HTML files in given directory")
	}

	// open and set up db
	db, err := initBolt("/Users/blobdon/blurbdata/blurb.db")
	if err != nil {
		return err
	}
	defer db.Close()

	authorMap := map[string]interface{}{}
	isbn13Map := map[string]interface{}{}
	inchecks := 0
	outchecks := 0
	matches := 0
	newbooks := 0
	// newauthors := 0

	authorMap, err = newAuthorMap(db)
	if err != nil {
		exiterr := fmt.Sprintf("Error creating authorMap: %s", err)
		return errors.New(exiterr)
	}
	isbn13Map, err = newISBN13Map(db)
	if err != nil {
		exiterr := fmt.Sprintf("Error creating isbn13Map: %s", err)
		return errors.New(exiterr)
	}

	for i, fp := range inHTMLs {
		log.Println("Ingesting html file # ", i+1)
		// turn html into book
		newB := book{ReviewAuthors: map[string]interface{}{}}
		newB, err := htmlToBookBN(fp, newB)
		if err != nil {
			// TODO Move to error directory?
			fmt.Println("htmltoBook error", err)
			continue
		}
		// check if isbn already in list - if yes, discard, continue
		if _, ok := isbn13Map[newB.ISBN13]; ok {
			// TODO: move file to duplicates directory
			log.Printf("Book: %s already exists, moved to duplicates directory", newB.ISBN13)
			continue
		}

		// save book to db
		err = bookToBolt(newB, db)
		if err != nil {
			// TODO move file to error directory
			log.Printf("Error saving book: %s to db, moved to errors directory", newB.ISBN13)
		}
		newbooks++
		// add author/isbn to maps
		authorMap[newB.Author] = struct{}{}
		isbn13Map[newB.ISBN13] = struct{}{}

		log.Println("Checking reviews in/out for html file # ", i+1)
		// check/process incoming reviews of text
		// New book's reviews contain any of all other authors?
		for a := range authorMap {
			// don't check reviews for a book's own author
			if a == newB.Author {
				continue
			}
			inchecks++

			match, err := containsAuthor(newB.ReviewText, a)
			if err != nil {
				log.Printf("Error matching book: %s author: %s - %s", newB.ISBN13, a, err)
			}
			if match {
				// fmt.Printf("Found old author %s in new reviews of %s: %s.\n", a, newB.ISBN13, newB.Title)
				newB.ReviewAuthors[a] = struct{}{}
				err = bookToBolt(newB, db)
				if err != nil {
					// TODO move file to error directory
					log.Printf("Error saving book: %s to db after adding new ReviewAuthor, moved html to errors directory", newB.ISBN13)
				}
				matches++
			}
		}

		// check/process outgoing reviews by author
		// All other books reviews contain the new author?
		for i := range isbn13Map {
			// get old book
			oldB, err := bookFromBolt(i, db)
			if err != nil {
				// TODO move file to error directory
				log.Printf("Error getting book: %s from db, moved html to errors directory", i)
			}
			// don't check reviews for a book's own author
			if oldB.Author == newB.Author {
				continue
			}
			outchecks++

			match, err := containsAuthor(oldB.ReviewText, newB.Author)
			if err != nil {
				log.Printf("Error matching book: %s author: %s - %s", oldB.ISBN13, newB.Author, err)
			}
			if match {
				// fmt.Printf("Found new author %s in old reviews of %s: %s.\n", newB.Author, oldB.ISBN13, oldB.Title)
				oldB.ReviewAuthors[newB.Author] = struct{}{}
				err = bookToBolt(oldB, db)
				if err != nil {
					// TODO move file to error directory
					log.Printf("Error saving book: %s to db after adding new ReviewAuthor, moved html to errors directory", oldB.ISBN13)
				}
				matches++
			}
		}

		// TODO: move html file to completed directory
	}
	log.Printf("New Books: %v\t Checked In: %v\t Out: %v\tMatches: %v", newbooks, inchecks, outchecks, matches)
	return nil
}
