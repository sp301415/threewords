package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io"
	"os"
)

// encryptBytes encrypts pt with key using AES-GCM.
func encryptBytes(pt []byte, key [32]byte) ([]byte, error) {
	block, err := aes.NewCipher(key[:])
	if err != nil {
		return nil, err
	}

	AESGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, AESGCM.NonceSize())
	_, err = io.ReadFull(rand.Reader, nonce)
	if err != nil {
		return nil, err
	}

	return AESGCM.Seal(nonce, nonce, pt, nil), nil
}

// decryptBytes decrypts ct with key using AES-GCM.
func decryptBytes(ct []byte, key [32]byte) ([]byte, error) {
	block, err := aes.NewCipher(key[:])
	if err != nil {
		return nil, err
	}

	AESGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce, ct := ct[:AESGCM.NonceSize()], ct[AESGCM.NonceSize():]
	return AESGCM.Open(nil, nonce, ct, nil)
}

// encryptAndWrite encrypts the data with key, and writes it to path.
func encryptAndWrite(data []byte, key [32]byte, path string) error {
	encryptedBytes, err := encryptBytes(data, key)
	if err != nil {
		return err
	}

	return os.WriteFile(path, encryptedBytes, 0644)
}

// readAndDecrypt opens the file in path, and decrypts it with key.
func readAndDecrypt(key [32]byte, path string) ([]byte, error) {
	encryptedBytes, err := os.ReadFile(path)
	if err != nil {
		return []byte{}, err
	}

	return decryptBytes(encryptedBytes, key)
}
