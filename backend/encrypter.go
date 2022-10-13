package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io"
	"os"
)

// EncryptBytes encrypts pt with key using AES-GCM.
func EncryptBytes(pt []byte, key [32]byte) ([]byte, error) {
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

// DecryptBytes decrypts ct with key using AES-GCM.
func DecryptBytes(ct []byte, key [32]byte) ([]byte, error) {
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

// EncryptAndWrite encrypts the data with key, and writes it to path.
func EncryptAndWrite(data []byte, key [32]byte, path string) error {
	encryptedBytes, err := EncryptBytes(data, key)
	if err != nil {
		return err
	}

	return os.WriteFile(path, encryptedBytes, 0644)
}

// ReadAndDecrypt opens the file in path, and decrypts it with key.
func ReadAndDecrypt(key [32]byte, path string) ([]byte, error) {
	encryptedBytes, err := os.ReadFile(path)
	if err != nil {
		return []byte{}, err
	}

	return DecryptBytes(encryptedBytes, key)
}
