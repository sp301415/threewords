package main

import (
	"log"
	"os"
	"time"
)

// Expire collects garbage every 1 hour.
func Expire() {
	ExpireQuery, err := DB.Prepare("SELECT ID, Path FROM files WHERE ExpireDate < NOW()")
	if err != nil {
		panic(err)
	}

	for {
		rows, err := ExpireQuery.Query()
		if err != nil {
			log.Printf("[ERROR] Garbage Collection Error: %v\n", err)
		}

		tx, _ := DB.Begin()
		if err != nil {
			log.Printf("[ERROR] Garbage Collection Error: %v\n", err)
		}

		cnt := 0
		for rows.Next() {
			var ID, path string
			rows.Scan(&ID, &path)

			tx.Exec("DELETE FROM files WHERE ID = ?", ID)
			os.Remove(path)
			cnt++
		}

		tx.Commit()
		rows.Close()

		log.Printf("[INFO] Garbage Collected: %v Files at %v\n", cnt, time.Now())

		// Somehow ticker stops after a day or so
		// Use time.Sleep for now...
		time.Sleep(time.Hour)
	}
}
