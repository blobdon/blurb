package main

import (
	"log"
	"sort"
	"strings"
	"unicode"

	"github.com/boltdb/bolt"
)

func saveTokensToIndex(tokens []string, isbn string, db *bolt.DB) {
	for _, t := range tokens {

		err := db.Update(func(tx *bolt.Tx) error {
			ibkt := tx.Bucket([]byte("reviewIndex"))
			// create bucket for the token
			tbkt, err := ibkt.CreateBucketIfNotExists([]byte(t))
			if err != nil {
				return err
			}
			// add isbn to the token bucket, blank value
			return tbkt.Put([]byte(isbn), []byte(""))
		})
		if err != nil {
			log.Printf("Error saving index %s: %s - %s", isbn, t, err)
		}
	}
}

// tokenize returns a slice of string tokens from text
func tokenize(text string) ([]string, error) {

	ns := normString(text)
	f := strings.Fields(ns)
	// remove duplicates, based on example at https://github.com/golang/go/wiki/SliceTricks
	sort.Strings(f)
	j := 0
	for i := 1; i < len(f); i++ {
		if f[j] == f[i] {
			continue
		}
		j++
		f[i], f[j] = f[j], f[i]
	}
	return f[:j+1], nil
}

// clean text = lowercase, remove non-alphanumeric - replacing with space for tokenization
func normString(s string) string {
	low := strings.ToLower(s)
	return toAlphanumeric(low)
}

// remove all non-alphanumeric characters from string, leaving spaces for tokenizing
func toAlphanumeric(s string) string {
	return strings.Map(
		func(r rune) rune {
			if unicode.IsLetter(r) || unicode.IsNumber(r) || unicode.IsSpace(r) {
				return r
			}
			return ' ' // replace with space
		},
		s,
	)
}
