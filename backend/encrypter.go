package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io"
)

func encryptFile(pt []byte, key [32]byte) ([]byte, error) {
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

func decryptFile(ct []byte, key [32]byte) ([]byte, error) {
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
