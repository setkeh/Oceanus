package db

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	//	guuid "github.com/google/uuid"
	_ "github.com/go-sql-driver/mysql"
)

type Mysql struct {
	DB *sql.DB
}

func (d *Mysql) Init() {
	db, err := sql.Open("mysql", fmt.Sprintf("%s", os.Getenv("DB_CONNECTION_STRING")))
	if err != nil {
		panic(err)
	}

	// See "Important settings" section.
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	d.DB = db
}

func (d *Mysql) Close() error {
	if d.DB != nil {
		return d.DB.Close()
	}
	return nil
}

func (d *Mysql) Insert(guid string, path string, imageName string) (int64, error) {
	//	query := fmt.Sprintf("INSERT INTO Images VALUES(%s, %s, %s)", guid, path, imageName)

	in, err := d.DB.Prepare("INSERT INTO Images VALUES(?, ?, ?)")

	defer in.Close()

	res, err := in.Exec(guid, path, imageName)
	if err != nil {
		return 0, err
	}

	return res.RowsAffected()
}

func (d *Mysql) GetByID(guid string) (string, error) {
	out, err := d.DB.Prepare(fmt.Sprintf("SELECT imagePath FROM Images WHERE id = %s", guid))
	if err != nil {
		return "", err
	}

	defer out.Close()

	var path string
	out.QueryRow(1).Scan(&path)
	return path, nil
}
