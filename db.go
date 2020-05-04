package main

import (
	"encoding/json"
	"log"
	"time"

	bolt "go.etcd.io/bbolt"
)

type image struct {
	ID  string `json:"id"`
	Src string `json:"src"`
	URL string `json:"url"`
	B64 string `json:"b64"`
}

func insertImage(i image) {
	db, err := bolt.Open("images.db", 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	db.Update(func(tx *bolt.Tx) error {
		tx.CreateBucketIfNotExists([]byte("Images"))
		b := tx.Bucket([]byte("Images"))

		buf, err := json.Marshal(i)
		if err != nil {
			return err
		}

		return b.Put([]byte(i.ID), buf)
	})
}

func getImage(id string) ([]byte, error) {
	db, err := bolt.Open("images.db", 0600, &bolt.Options{Timeout: 100 * time.Second})
	if err != nil {
		return nil, err
	}
	defer db.Close()

	var ret []byte

	db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("Images")).Get([]byte(id))

		//fmt.Println(id)
		//fmt.Println(b)

		ret = b

		return nil //b.Get([]byte(id))
	})

	return ret, nil
}
