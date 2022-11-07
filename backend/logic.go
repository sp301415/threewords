package main

import (
	"context"
	"database/sql"
	"encoding/base64"
	"errors"
	"path"
	"threewords/internal/db"
	"threewords/threewords"

	"github.com/google/uuid"
)

const ROOT = "files"

type File struct {
	Name string
	Data []byte
}

// Add encrypts adds file to the database, with given threewords.
func Add(words threewords.ThreeWords, file File) error {
	id := words.ID()
	key := words.Key()

	// Encrypt and write file
	savePath := path.Join(ROOT, uuid.NewString())
	err := EncryptAndWrite(file.Data, key, savePath)
	if err != nil {
		return ErrFileOperation
	}

	// Encrypt original file name
	encryptedName, err := EncryptBytes([]byte(file.Name), key)
	if err != nil {
		return ErrFileOperation
	}
	encryptedNameBase64 := base64.StdEncoding.EncodeToString(encryptedName)

	// Write to database
	err = DB.CreateEntry(context.Background(), db.CreateEntryParams{
		ID:           id,
		FilePath:     savePath,
		OriginalName: encryptedNameBase64,
	})
	if err != nil {
		return ErrDBOperation
	}

	return nil
}

// Get returns the file data from database using given threewords.
func Get(words threewords.ThreeWords) (File, error) {
	if !threewords.Validate(words) {
		return File{}, ErrKeyNotFound
	}

	id := words.ID()
	key := words.Key()

	// Read from database
	row, err := DB.ReadEntry(context.Background(), id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return File{}, ErrKeyNotFound
		}
		return File{}, ErrDBOperation
	}

	// Decrypt file and encryptedName
	data, err := ReadAndDecrypt(key, row.FilePath)
	if err != nil {
		return File{}, ErrFileOperation
	}

	encryptedName, err := base64.StdEncoding.DecodeString(row.OriginalName)
	if err != nil {
		return File{}, ErrFileOperation
	}

	name, err := DecryptBytes(encryptedName, key)
	if err != nil {
		return File{}, ErrFileOperation
	}

	return File{Name: string(name), Data: data}, nil
}
