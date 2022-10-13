package main

import (
	"log"
	"os"
	"time"
)

// Expire collects garbage every 1 hour.
func Expire() {
	expireTicker := time.NewTicker(time.Hour)

	for t := range expireTicker.C {
		log.Printf("[INFO] Garbage Collected: %v\n", t)

		rows, _ := DB.Query("SELECT ID, Path FROM files WHERE ExpireDate < NOW()")

		tx, _ := DB.Begin()
		for rows.Next() {
			var ID, path string
			rows.Scan(&ID, &path)

			tx.Exec("DELETE FROM files WHERE ID = ?", ID)
			os.Remove(path)
		}
		tx.Commit()
	}
}
