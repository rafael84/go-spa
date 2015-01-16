package password

import (
	"crypto/rand"
	"crypto/sha1"
	"hash"
	"io"
)

const saltSize = 32

func generateSalt(secret []byte) ([]byte, error) {
	buf := make([]byte, saltSize, saltSize+sha1.Size)
	_, err := io.ReadFull(rand.Reader, buf)
	if err != nil {
		return nil, err
	}

	hash := sha1.New()
	hash.Write(buf)
	hash.Write(secret)

	return hash.Sum(buf), nil
}

func generateHash(salt, secret []byte) (hash.Hash, error) {
	combination := string(salt) + string(secret)
	passwordHash := sha1.New()
	io.WriteString(passwordHash, combination)

	return passwordHash, nil
}
