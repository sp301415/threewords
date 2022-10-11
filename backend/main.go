package main

import (
	"database/sql"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

var DownloadQuery *sql.Stmt
var UploadQuery *sql.Stmt
var CheckExpireQuery *sql.Stmt
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

	CheckExpireQuery, err = DB.Prepare("SELECT ID, Path FROM files WHERE ExpireDate < NOW()")
	if err != nil {
		panic(err)
	}

	ExpireQuery, err = DB.Prepare("DELETE FROM files WHERE ID = ?")
	if err != nil {
		panic(err)
	}
}

func main() {
	defer DownloadQuery.Close()
	defer UploadQuery.Close()
	defer CheckExpireQuery.Close()
	defer ExpireQuery.Close()
	defer DB.Close()

	expireTicker := time.NewTicker(24 * time.Hour)
	go func() {
		for t := range expireTicker.C {
			log.Printf("[INFO] Garbage collected at %v", t)
			Expire()
		}
	}()

	gin.SetMode(gin.ReleaseMode)

	r := gin.Default()
	r.SetTrustedProxies([]string{"127.0.0.1"})
	r.MaxMultipartMemory = 100 * (1 << 20) // 100MB

	r.POST("/upload", UploadHandler)
	r.POST("/download", DownloadHandler)
	r.Run(":8000")
}
