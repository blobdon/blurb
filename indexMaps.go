package main

import (
	"encoding/json"
	"fmt"

	"github.com/boltdb/bolt"
)

func newAuthorMap(db *bolt.DB) (map[string]interface{}, error) {
	m := map[string]interface{}{}
	err := db.View(func(tx *bolt.Tx) error {
		// Assume bucket exists and has keys
		bkt := tx.Bucket([]byte("books"))

		bkt.ForEach(func(k, v []byte) error {
			b := book{}
			err := json.Unmarshal(v, &b)
			if err != nil {
				fmt.Println("Error unmarshalling json", err) // to log once implemented
			}
			m[b.Author] = struct{}{}
			return nil
		})
		return nil
	})
	return m, err
}

func newISBN13Map(db *bolt.DB) (map[string]interface{}, error) {
	m := map[string]interface{}{}
	err := db.View(func(tx *bolt.Tx) error {
		// Assume bucket exists and has keys
		bkt := tx.Bucket([]byte("books"))

		bkt.ForEach(func(k, v []byte) error {
			b := book{}
			err := json.Unmarshal(v, &b)
			if err != nil {
				fmt.Println("Error unmarshalling json", err) // to log once implemented
			}
			m[b.ISBN13] = struct{}{}
			return nil
		})
		return nil
	})
	return m, err
}
