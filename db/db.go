package db

import (
	"encoding/json"
	"log"
	"time"

	"github.com/setkeh/Oceanus/models"
	bolt "go.etcd.io/bbolt"
)

type DbClient struct {
	// Filename to the BoltDB database.
	Path string

	// Returns the current time.
	Now func() time.Time

	db *bolt.DB
}

func (c *DbClient) DbOpen() error {
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

func (c *DbClient) dbClose() error {
	if c.db != nil {
		return c.db.Close()
	}
	return nil
}

func (c *DbClient) InsertImage(i models.Image) {
	tx, err := c.db.Begin(true)
	if err != nil {
		log.Fatal(err)
	}
	//defer d.Db.Close()

	tx.DB().Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("Images"))

		buf, err := json.Marshal(i)
		if err != nil {
			return err
		}

		return b.Put([]byte(i.ID), buf)
	})
}

func (c *DbClient) GetImage(id string) ([]byte, error) {
	tx, err := c.db.Begin(true)
	if err != nil {
		return nil, err
	}
	//defer db.Close()

	var ret []byte

	tx.DB().View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("Images")).Get([]byte(id))

		//fmt.Println(id)
		//fmt.Println(b)

		ret = append(ret, b...)

		return nil //b.Get([]byte(id))
	})

	return ret, nil
}

func (c *DbClient) GetImageList() ([]models.Photo, error) {
	tx, err := c.db.Begin(true)
	if err != nil {
		return nil, err
	}
	//defer db.Close()

	var ret []models.Photo

	tx.DB().View(func(tx *bolt.Tx) error {
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
