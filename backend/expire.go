package main

import "os"

func expire() {
	rows, _ := ExpireQuery.Query()

	query, _ := DB.Begin()
	for rows.Next() {
		var ID, path string
		rows.Scan(&ID, &path)

		query.Exec("DELETE FROM files WHERE ID=?", ID)
		os.Remove(path)
	}
	query.Commit()
}
