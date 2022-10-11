package main

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

var DownloadQuery *sql.Stmt
var UploadQuery *sql.Stmt
var ExpireQuery *sql.Stmt

// Initialize DB Connection
func init() {
	var err error
	DB, err = sql.Open("mysql", "root:@tcp(localhost)/threewords")
	if err != nil {
		panic(err)
	}

	DownloadQuery, err = DB.Prepare("SELECT Path, OriginalName FROM files WHERE ID = ?")
	if err != nil {
		panic(err)
	}

	UploadQuery, err = DB.Prepare("INSERT INTO files VALUES (?, ?, ?, NOW() + INTERVAL 24 HOUR)")
	if err != nil {
		panic(err)
	}

	ExpireQuery, err = DB.Prepare("SELECT ID, Path FROM files WHERE ExpireDate < NOW()")
	if err != nil {
		panic(err)
	}
}

func main() {
	defer DownloadQuery.Close()
	defer UploadQuery.Close()
	defer ExpireQuery.Close()
	defer DB.Close()

	expireTicker := time.NewTicker(24 * time.Hour)
	go func() {
		for t := range expireTicker.C {
			log.Printf("[INFO] Garbage collected at %v", t)
			expire()
		}
	}()

	http.HandleFunc("/download", DownloadHandler)
	http.HandleFunc("/upload", UploadHandler)
	http.ListenAndServe(":8000", nil)
}
