package main

import (
	"context"
	"log"
	"os"
	"time"
)

// Expire collects garbage files.
func Expire() {
	rows, err := DB.FindExpiredEntry(context.Background())
	if err != nil {
		log.Printf("Expire Error: %v", err)
		return
	}

	var cnt int
	for _, row := range rows {
		err = DB.DeleteEntry(context.Background(), row.ID)
		if err != nil {
			log.Printf("Expire Error: %v", err)
			return
		}

		err = os.Remove(row.FilePath)
		if err != nil {
			log.Printf("Expire Error: %v", err)
			return
		}

		cnt++
	}

	log.Printf("[INFO] Garbage Collected: %v Files at %v\n", cnt, time.Now())
}
