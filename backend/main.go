package main

import (
	"database/sql"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

// Initialize DB Connection
func init() {
	var err error
	DB, err = sql.Open("sqlite3", "threewords.db")
	if err != nil {
		panic(err)
	}
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

	r.POST("/upload", UploadHandler)
	r.POST("/download", DownloadHandler)
	r.Run(":8000")
}
