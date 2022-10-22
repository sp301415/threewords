package main

import (
	"database/sql"
	"log"
	"os"
	"time"
)

const (
	expireDBError   = "[ERROR] Garbage Collection Error (Database): %v\n"
	expireFileError = "[ERROR] Garbage Collectoin Error (os.Remove): %v\n"
)

var ExpireQuery *sql.Stmt

func init() {
	var err error
	ExpireQuery, err = DB.Prepare("SELECT ID, Path FROM files WHERE ExpireDate < NOW()")
	if err != nil {
		panic(err)
	}
}

// Expire collects garbage files.
func Expire() {
	rows, err := ExpireQuery.Query()
	if err != nil {
		log.Printf(expireDBError, err)
	}
	defer rows.Close()

	tx, err := DB.Begin()
	if err != nil {
		log.Printf(expireDBError, err)
	}

	var cnt int
	var ID, path string
	for rows.Next() {
		err := rows.Scan(&ID, &path)
		if err != nil {
			log.Printf(expireDBError, err)
		}

		_, err = tx.Exec("DELETE FROM files WHERE ID = ?", ID)
		if err != nil {
			log.Printf(expireDBError, err)
		}

		err = os.Remove(path)
		if err != nil {
			log.Printf(expireFileError, err)
		}

		cnt++
	}

	err = tx.Commit()
	if err != nil {
		log.Printf(expireDBError, err)
	}

	log.Printf("[INFO] Garbage Collected: %v Files at %v\n", cnt, time.Now())
}
