package main

import "os"

// Expire collects garbage every 24 hours.
func Expire() {
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
