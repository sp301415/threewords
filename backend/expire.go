package main

import (
	"context"
	"log"
	"os"
	"threewords/internal/db"
	"time"
)

const (
	expireDBError   = "[ERROR] Garbage Collection Error (Database): %v\n"
	expireFileError = "[ERROR] Garbage Collectoin Error (os.Remove): %v\n"
)

// Expire collects garbage files.
func Expire() {
	q := db.New(DB)

	rows, err := q.FindExpiredEntry(context.Background())
	if err != nil {
		log.Printf(expireDBError, err)
	}

	var cnt int
	for _, row := range rows {
		err = q.DeleteEntry(context.Background(), row.ID)
		if err != nil {
			log.Printf(expireDBError, err)
		}

		err = os.Remove(row.FilePath)
		if err != nil {
			log.Printf(expireFileError, err)
		}

		cnt++
	}

	log.Printf("[INFO] Garbage Collected: %v Files at %v\n", cnt, time.Now())
}
