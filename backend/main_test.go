package main

import (
	"bytes"
	"testing"
	"threewords/threewords"
)

func TestEncDec(t *testing.T) {
	words := threewords.Generate()
	pt := []byte(words.ID())

	ct, _ := encryptFile(pt, words.Key())
	ptct, _ := decryptFile(ct, words.Key())

	if !bytes.Equal(pt, ptct) {
		t.Fail()
	}
}
