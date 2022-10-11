package main

import "os"

func expire() {
	rows, _ := CheckExpireQuery.Query()

	tx, _ := DB.Begin()
	for rows.Next() {
		var ID, path string
		rows.Scan(&ID, &path)
		tx.Stmt(ExpireQuery).Exec(ID)
		os.Remove(path)
	}
	tx.Commit()
}
