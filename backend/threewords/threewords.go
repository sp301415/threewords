// Package threewords provides low level API for ThreeWords.
package threewords

import (
	"bufio"
	"crypto/md5"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math/big"
	"os"
	"strings"
)

// type ThreeWords represents three korean words.
// ThreeWords SHOULD NOT be saved in any form; use ID() method to get unique ID.
type ThreeWords [3]string

// wordList saves possible words, initialized at init().
var wordList = struct {
	slice        []string
	set          map[string]struct{} // for fast validation
	lengthBigInt *big.Int
}{}

const wordCount = 3000 // Approximate number of words

func init() {
	wordList.slice = make([]string, 0, wordCount)
	wordList.set = make(map[string]struct{}, wordCount)

	file, err := os.Open("threewords/words.txt")
	if err != nil {
		panic(err)
	}

	fileScanner := bufio.NewScanner(file)
	for fileScanner.Scan() {
		word := strings.TrimSpace(fileScanner.Text())
		wordList.slice = append(wordList.slice, word)
		wordList.set[word] = struct{}{}
	}
	wordList.lengthBigInt = big.NewInt(int64(len(wordList.slice)))
}

// Generate creates a new ThreeWords, chosen from given word list
// with cryptographically secure RNG.
func Generate() ThreeWords {
	var i0, i1, i2 int
	for {
		i0Big, _ := rand.Int(rand.Reader, wordList.lengthBigInt)
		i1Big, _ := rand.Int(rand.Reader, wordList.lengthBigInt)
		i2Big, _ := rand.Int(rand.Reader, wordList.lengthBigInt)

		i0 = int(i0Big.Int64())
		i1 = int(i1Big.Int64())
		i2 = int(i2Big.Int64())

		if i0 != i1 && i1 != i2 && i2 != i0 {
			break
		}
	}

	return ThreeWords{wordList.slice[i0], wordList.slice[i1], wordList.slice[i2]}
}

// FromString creates new threewords from the string of form `%s-%s-%s`.
func FromString(s string) (ThreeWords, bool) {
	w := strings.Split(s, "-")
	if len(w) != 3 {
		return ThreeWords{}, false
	}

	words := ThreeWords{w[0], w[1], w[2]}
	if !Validate(words) {
		return ThreeWords{}, false
	}
	return words, true
}

// Validate checks if given ThreeWords is inside wordList.
func Validate(words ThreeWords) bool {
	for _, word := range words {
		if _, ok := wordList.set[word]; !ok {
			return false
		}
	}
	return true
}

// String implements fmt.Stringer interface. Default form is `%s, %s, %s`.
func (w ThreeWords) String() string {
	return fmt.Sprintf("%s, %s, %s", w[0], w[1], w[2])
}

// ID returns the unique ID of this words, which is MD5(words).
// words itself SHOULD NOT be saved in any form.
func (w ThreeWords) ID() string {
	h := md5.Sum([]byte(w.String()))
	return hex.EncodeToString(h[:])
}

// Key returns a secret key derived from this words to sign the file, which is SHA256(words).
// Key SHOULD NOT be saved in any form.
func (w ThreeWords) Key() [32]byte {
	return sha256.Sum256([]byte(w.String()))
}
