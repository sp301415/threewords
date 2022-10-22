package main

import (
	"bytes"
	"database/sql"
	"encoding/base64"
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strings"
	"threewords/internal/db"
	"threewords/threewords"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const STORE_BASE_DIR = "files"

const (
	fileError       = "파일을 읽거나 쓸 수 없습니다."
	encryptionError = "암호화/복호화하는 도중 문제가 생겼습니다."
	dbError         = "데이터베이스를 처리하는 도중 문제가 생겼습니다."
	keyError        = "해당 키는 잘못되었거나 존재하지 않습니다."
)

// UploadHandler handles /upload API.
// It mainly reads the file uploaded via multipart-form, and saves the encrypted file, and assigns new threeword.
func UploadHandler(c *gin.Context) {
	// Set CORS header
	c.Header("Access-Control-Allow-Origin", "https://threewords.sp301415.com")

	// Create unique threeword
	var words threewords.ThreeWords
	for {
		words = threewords.Generate()

		_, err := DB.CheckID(c, words.ID())
		if errors.Is(err, sql.ErrNoRows) {
			break
		} else if err != nil {
			c.String(http.StatusInternalServerError, dbError)
			return
		}
	}

	// Read uploaded file
	fileHeader, err := c.FormFile("upload")
	if err != nil {
		c.String(http.StatusBadRequest, fileError)
		return
	}

	file, err := fileHeader.Open()
	if err != nil {
		c.String(http.StatusInternalServerError, fileError)
		return
	}
	defer file.Close()

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		c.String(http.StatusInternalServerError, fileError)
		return
	}

	filePath := filepath.Join(STORE_BASE_DIR, uuid.NewString())

	// Encrypt and write file
	err = EncryptAndWrite(fileBytes, words.Key(), filePath)
	if err != nil {
		c.String(http.StatusInternalServerError, encryptionError)
		return
	}

	// Encrypt original file name
	encryptedName, err := EncryptBytes([]byte(fileHeader.Filename), words.Key())
	if err != nil {
		c.String(http.StatusInternalServerError, encryptionError)
		return
	}
	encryptedNameBase64 := base64.StdEncoding.EncodeToString(encryptedName)

	// Write to database
	err = DB.CreateEntry(c, db.CreateEntryParams{
		ID:           words.ID(),
		FilePath:     filePath,
		OriginalName: encryptedNameBase64,
	})
	if err != nil {
		c.String(http.StatusInternalServerError, dbError)
		return
	}

	c.String(http.StatusOK, words.String())
}

// DownloadHandler handles /donwload API.
// It validates threeword sent by user, and returns the associated file.
func DownloadHandler(c *gin.Context) {
	// Set CORS header
	c.Header("Access-Control-Allow-Origin", "https://threewords.sp301415.com")

	// Get threeword from user
	words := threewords.ThreeWords{
		strings.TrimSpace(c.PostForm("word0")),
		strings.TrimSpace(c.PostForm("word1")),
		strings.TrimSpace(c.PostForm("word2")),
	}

	// Validate threeword
	if !threewords.Validate(words) {
		c.String(http.StatusBadRequest, keyError)
		return
	}

	// Read from database
	row, err := DB.ReadEntry(c, words.ID())
	if errors.Is(err, sql.ErrNoRows) {
		c.String(http.StatusBadRequest, keyError)
		return
	} else if err != nil {
		c.String(http.StatusInternalServerError, dbError)
		return
	}

	// Decrypt file and originalName
	fileBytes, err := ReadAndDecrypt(words.Key(), row.FilePath)
	if err != nil {
		c.String(http.StatusInternalServerError, encryptionError)
		return
	}

	encryptedName, err := base64.StdEncoding.DecodeString(row.OriginalName)
	if err != nil {
		c.String(http.StatusInternalServerError, fileError)
		return
	}

	originalName, err := DecryptBytes(encryptedName, words.Key())
	if err != nil {
		c.String(http.StatusInternalServerError, encryptionError)
		return
	}

	// Send as multipart/form-data
	var formResponse bytes.Buffer
	formWriter := multipart.NewWriter(&formResponse)
	fileWriter, err := formWriter.CreateFormFile("file", base64.StdEncoding.EncodeToString(originalName))
	if err != nil {
		c.String(http.StatusInternalServerError, fileError)
		return
	}

	_, err = fileWriter.Write(fileBytes)
	if err != nil {
		c.String(http.StatusInternalServerError, encryptionError)
		return
	}

	err = formWriter.Close()
	if err != nil {
		c.String(http.StatusInternalServerError, encryptionError)
		return
	}

	c.Data(http.StatusOK, formWriter.FormDataContentType(), formResponse.Bytes())
}
