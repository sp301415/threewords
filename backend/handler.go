package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"threewords/threewords"

	"github.com/google/uuid"
)

const STORE_DIR = "files"

// DownloadHandler handles /download API.
func DownloadHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "https://threewords.sp301415.com")

	r.ParseForm()

	words := threewords.ThreeWords{
		strings.TrimSpace(r.PostFormValue("word0")),
		strings.TrimSpace(r.PostFormValue("word1")),
		strings.TrimSpace(r.PostFormValue("word2")),
	}

	if !threewords.Validate(words) {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "잘못된 키입니다.")
		log.Printf("[ERROR] User requested wrong threewords: %v", words)
		return
	}

	row, err := DownloadQuery.Query(words.ID())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "DB 쿼리를 실행할 수 없습니다.")
		log.Printf("[ERROR] Cannot run DownloadQuery.Exec: %v\n", err)
		return
	}

	if !row.Next() {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "키를 찾을 수 없습니다.")
		log.Printf("[ERROR] Cannot find files associated with threewords %v\n", words)
		return
	}

	var path, originalName string
	row.Scan(&path, &originalName)

	encryptedBytes, err := os.ReadFile(path)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "파일을 열 수 없습니다.")
		log.Printf("[ERROR] Cannot read file: %v\n", err)
		return
	}

	fileBytes, err := decryptFile(encryptedBytes, words.Key())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "파일을 복호화할 수 없습니다.")
		log.Printf("[ERROR] Cannot decrypt file: %v\n", err)
		return
	}

	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%v\"", originalName))
	w.Write(fileBytes)

	log.Printf("[INFO] Threewords %v accessed!\n", words)
}

// UploadHandler handles /upload API.
func UploadHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "https://threewords.sp301415.com")

	words := threewords.Generate()

	uploadedFile, header, err := r.FormFile("upload")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "HTML Form을 파싱할 수 없습니다.")
		log.Printf("[ERROR] FormFile error: %v\n", err)
		return
	}

	fileName := uuid.NewString()
	filePath := filepath.Join(STORE_DIR, fileName)
	originalName := header.Filename

	content, err := io.ReadAll(uploadedFile)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "업로드한 파일을 읽을 수 없습니다.")
		log.Printf("[ERROR] Cannot read uploaded file: %v\n", err)
		return
	}

	encryptedContent, err := encryptFile(content, words.Key())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "파일을 암호화할 수 없습니다.")
		log.Printf("[ERROR] Cannot encrypt file: %v\n", err)
		return
	}

	err = os.WriteFile(filePath, encryptedContent, 0644)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "파일을 쓸 수 없습니다.")
		log.Printf("[ERROR] Cannot write file: %v\n", err)
		return
	}

	_, err = UploadQuery.Exec(words.ID(), filePath, originalName)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "DB 쿼리를 실행할 수 없습니다.")
		log.Printf("[ERROR] Cannot run UploadQuery.Exec: %v\n", err)
		return
	}

	fmt.Fprint(w, words)

	log.Printf("[INFO] File %v uploaded!\n", originalName)
}
