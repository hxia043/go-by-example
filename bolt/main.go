package main

import (
	"fmt"
	"log"

	"github.com/boltdb/bolt"
)

func main() {
	db, err := bolt.Open("my.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte("MyBucket"))
		if err != nil {
			return err
		}

		err = b.Put([]byte("answer1"), []byte("42"))
		if err != nil {
			return err
		}

		err = b.Put([]byte("answer2"), []byte("43"))
		return err
	})

	var answer1 string
	var answer2 string

	db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("MyBucket"))
		v := b.Get([]byte("answer1"))
		answer1 = string(v)

		v = b.Get([]byte("answer2"))
		answer2 = string(v)
		return nil
	})

	fmt.Println(answer1)
	fmt.Println(answer2)
}
