package main

type bookView struct {
	ISBN13        string
	Title         string
	Author        string
	ReviewAuthors map[string]string //[name]=blankfornow
}

type authorView struct {
	Name          string
	BooksAuthored map[string]string // [ISBN]=title
	BooksReviewed map[string]string // [ISBN]=title
}

// buildJSONViews will build a json directory, at dirpath, of the
// blurb data as viewed in the website/webapp, i.e. individual books and authors
// The resulting json directory is a potential static-file API for a webapp
// func buildJSONViews(db *bolt.DB, dirpath string) error {
// 	authorsMap, err := newAuthorMap(db)
// 	if err != nil {
// 		return err
// 	}
// 	isbn13Map, err := newISBN13Map(db)
// 	if err != nil {
// 		return err
// 	}
// 	for _, i := range isbn13Map {
// 		//TODO get bookView from bolt by isbn
// 		//TODO save bookView
// 	}

// 	return nil
// }
