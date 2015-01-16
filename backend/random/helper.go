package random

import (
	"crypto/rand"
	"encoding/hex"
	"io"
)

func New(size int) (string, error) {
	buf := make([]byte, size, size)
	_, err := io.ReadFull(rand.Reader, buf)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(buf), nil
}
