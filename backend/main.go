package main

import (
	"database/sql"
	"threewords/internal/db"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)

var DB *db.Queries

// Initialize DB Connection
func init() {
	sqlDB, err := sql.Open("sqlite3", "sql/threewords.db")
	if err != nil {
		panic(err)
	}

	DB = db.New(sqlDB)
}

func main() {
	// Run garbage collector every hour.
	go func() {
		for {
			Expire()
			time.Sleep(time.Hour)
		}
	}()

	gin.SetMode(gin.ReleaseMode)

	r := gin.Default()
	r.SetTrustedProxies([]string{"127.0.0.1"})
	r.MaxMultipartMemory = 100 * (1 << 20) // 100MB

	r.POST("/upload", UploadPostHandler)
	r.POST("/download", DownloadPostHandler)
	r.GET("/download/:words", DownloadGetHandler)
	r.Run(":8000")
}
