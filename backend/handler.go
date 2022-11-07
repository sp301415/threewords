package main

import (
	"bytes"
	"database/sql"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"strings"
	"threewords/threewords"

	"github.com/gin-gonic/gin"
)

// UploadPostHandler handles /upload POST API.
// It mainly reads the file uploaded via multipart-form, and saves the encrypted file, and assigns new threeword.
func UploadPostHandler(c *gin.Context) {
	// Set CORS header
	c.Header("Access-Control-Allow-Origin", "https://threewords.sp301415.com")

	// Create unique threeword
	var words threewords.ThreeWords
	for {
		words = threewords.Generate()

		_, err := DB.CheckID(c, words.ID())
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				break
			}
			c.String(http.StatusInternalServerError, ErrDBOperation.Error())
			return
		}
	}

	// Read uploaded file
	fileHeader, err := c.FormFile("upload")
	if err != nil {
		c.String(http.StatusBadRequest, ErrFileOperation.Error())
		return
	}

	file, err := fileHeader.Open()
	if err != nil {
		c.String(http.StatusInternalServerError, ErrFileOperation.Error())
		return
	}
	defer file.Close()

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		c.String(http.StatusInternalServerError, ErrFileOperation.Error())
		return
	}

	err = Add(words, File{Name: fileHeader.Filename, Data: fileBytes})
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.String(http.StatusOK, words.String())
}

// QueryEscape escapes the string so it can be safely placed inside a URL query.
// The Go stdlib offers url.QueryEscape, but it encodes spaces to `+`; this function fixes it.
func QueryEscape(s string) string {
	u := &url.URL{Path: s}
	return strings.ReplaceAll(u.String(), "+", "%2B")
}

// DownloadPostHandler handles /download POST API.
// It validates threeword sent by user, and returns the associated file.
func DownloadPostHandler(c *gin.Context) {
	// Set CORS header
	c.Header("Access-Control-Allow-Origin", "https://threewords.sp301415.com")

	// Get threeword from user
	words := threewords.ThreeWords{
		strings.TrimSpace(c.PostForm("word0")),
		strings.TrimSpace(c.PostForm("word1")),
		strings.TrimSpace(c.PostForm("word2")),
	}

	file, err := Get(words)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	// Send as multipart/form-data
	var formResponse bytes.Buffer
	formWriter := multipart.NewWriter(&formResponse)
	fileWriter, err := formWriter.CreateFormFile("file", QueryEscape(file.Name))
	if err != nil {
		c.String(http.StatusInternalServerError, ErrFileOperation.Error())
		return
	}

	_, err = fileWriter.Write(file.Data)
	if err != nil {
		c.String(http.StatusInternalServerError, ErrFileOperation.Error())
		return
	}

	err = formWriter.Close()
	if err != nil {
		c.String(http.StatusInternalServerError, ErrFileOperation.Error())
		return
	}

	c.Data(http.StatusOK, formWriter.FormDataContentType(), formResponse.Bytes())
}

// DownloadGetHandler handles POST API of the form /download/%s-%s-%s.
func DownloadGetHandler(c *gin.Context) {
	words, ok := threewords.FromString(c.Param("words"))
	if !ok {
		c.String(http.StatusBadRequest, ErrKeyNotFound.Error())
		return
	}

	file, err := Get(words)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename*=utf-8''%s", QueryEscape(file.Name)))
	c.Data(http.StatusOK, "application/octet-stream", file.Data)
}
