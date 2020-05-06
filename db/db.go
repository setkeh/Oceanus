package db

import (
	"encoding/json"
	"log"
	"time"

	"github.com/setkeh/Oceanus/models"
	bolt "go.etcd.io/bbolt"
)

type dbClient struct {
	// Filename to the BoltDB database.
	Path string

	// Returns the current time.
	Now func() time.Time

	db *bolt.DB
}

func (c *dbClient) dbOpen() error {
	// Open database file.
	db, err := bolt.Open(c.Path, 0666, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return err
	}
	c.db = db

	// Start writable transaction.
	tx, err := c.db.Begin(true)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Initialize top-level buckets.
	if _, err := tx.CreateBucketIfNotExists([]byte("Images")); err != nil {
		return err
	}

	// Save transaction to disk.
	return tx.Commit()
}

func (c *dbClient) dbClose() error {
	if c.db != nil {
		return c.db.Close()
	}
	return nil
}

func (c *dbClient) insertImage(i models.Image) {
	//db, err := bolt.Open("images.db", 0600, &bolt.Options{Timeout: 1 * time.Second})
	tx, err := c.db.Begin(true)
	if err != nil {
		log.Fatal(d.Err)
	}
	//defer d.Db.Close()

	tx.
	d.Db.Update(func(tx *bolt.Tx) error {
		tx.CreateBucketIfNotExists([]byte("Images"))
		b := tx.Bucket([]byte("Images"))

		buf, err := json.Marshal(i)
		if err != nil {
			return err
		}

		return b.Put([]byte(i.ID), buf)
	})
}

func (d *db) getImage(id string) ([]byte, error) {
	//db, err := bolt.Open("images.db", 0600, &bolt.Options{Timeout: 100 * time.Second})
	if d.Err != nil {
		return nil, d.Err
	}
	//defer db.Close()

	var ret []byte

	d.Db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("Images")).Get([]byte(id))

		//fmt.Println(id)
		//fmt.Println(b)

		ret = append(ret, b...)

		return nil //b.Get([]byte(id))
	})

	return ret, nil
}

func (d *db) getImageList() ([]models.Photo, error) {
	//db, err := bolt.Open("images.db", 0600, &bolt.Options{Timeout: 100 * time.Second})
	if d.Err != nil {
		return nil, d.Err
	}
	//defer db.Close()

	var ret []models.Photo

	d.Db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("Images"))

		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			var pic models.Photo
			json.Unmarshal(v, &pic)
			ret = append(ret, pic)
		}
		return nil
	})

	return ret, nil
}
