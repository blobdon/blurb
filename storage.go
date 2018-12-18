package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/boltdb/bolt"
)

func initBolt(dbfilepath string) (*bolt.DB, error) {
	// Open the .db data file at given filepath.
	// It will be created if it doesn't exist.
	// TODO: need to add Timeout option to avoid indefinite wait if already opened
	// db, err := bolt.Open(dbfilepath, 0600, &bolt.Options{Timeout: 1 * time.Second})
	db, err := bolt.Open(dbfilepath, 0600, nil)
	if err != nil {
		log.Fatal(err)
		return db, err
	}
	// defer db.Close()

	// create any needed buckets
	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("books"))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		_, err = tx.CreateBucketIfNotExists([]byte("reviewIndex"))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		return nil
	})
	if err != nil {
		log.Println(err)
		return db, err
	}
	return db, err
}

// bookToBolt saves b to the db.
func bookToBolt(b book, db *bolt.DB) error {
	// Marshal book data into bytes. Done outside db func to limit tx time, esp for writes
	buf, err := json.Marshal(b)
	if err != nil {
		return err
	}
	return db.Update(func(tx *bolt.Tx) error {
		bkt := tx.Bucket([]byte("books"))
		return bkt.Put([]byte(b.ISBN13), buf)
	})
}

// bookFromBolt gets a book from bolt, currentlly keyed by ISBN13
func bookFromBolt(isbn13 string, db *bolt.DB) (book, error) {
	b := book{}
	err := db.View(func(tx *bolt.Tx) error {
		bkt := tx.Bucket([]byte("books"))
		j := bkt.Get([]byte(isbn13))
		return json.Unmarshal(j, &b)
	})
	return b, err
}
