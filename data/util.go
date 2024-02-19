package data

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	rando "math/rand"
	"time"
)

func GetRandGameCode(n int) string {
	var letters = []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rando.Intn(len(letters))]
	}
	return string(b)
}

// TimeToString converts a time.Time to RFC3339 string
func TimeToString(t time.Time) string {
	return t.Format(time.RFC3339)
}

// StringToTime used to convert a RFC3339 to time.Time
func StringToTime(s string) (time.Time, error) {
	return time.Parse(time.RFC3339, s) //returns both a time and an error so it can be returned directly.
}

const timestamp = "20060102150405"

// NewUniqueID Returns a unique id that fills 12 bytes
func NewUniqueID() string {
	timestamp := time.Now().Format(timestamp)
	b := make([]byte, 12)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		log.Println(err)
		return ""
	}
	return fmt.Sprintf("%s%s%s", base64.URLEncoding.EncodeToString(b), "&tt&", timestamp)
}
