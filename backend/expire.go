package main

import (
	"context"
	"log"
	"os"
	"time"
)

const (
	expireDBError   = "[ERROR] Garbage Collection Error (Database): %v\n"
	expireFileError = "[ERROR] Garbage Collectoin Error (os.Remove): %v\n"
)

// Expire collects garbage files.
func Expire() {
	rows, err := DB.FindExpiredEntry(context.Background())
	if err != nil {
		log.Printf(expireDBError, err)
		return
	}

	var cnt int
	for _, row := range rows {
		err = DB.DeleteEntry(context.Background(), row.ID)
		if err != nil {
			log.Printf(expireDBError, err)
			return
		}

		err = os.Remove(row.FilePath)
		if err != nil {
			log.Printf(expireFileError, err)
			return
		}

		cnt++
	}

	log.Printf("[INFO] Garbage Collected: %v Files at %v\n", cnt, time.Now())
}
