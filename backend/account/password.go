package account

import (
	"bytes"
	"crypto/rand"
	"crypto/sha1"
	"database/sql/driver"
	"encoding/hex"
	"hash"
	"io"
	"strings"
)

const saltSize = 32

type SaltedPassword string

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

func (p *SaltedPassword) Encode(secret string) error {
	salt, err := generateSalt([]byte(secret))
	if err != nil {
		return err
	}

	hash, err := generateHash(salt, []byte(secret))
	if err != nil {
		return err
	}
	*p = SaltedPassword(hex.EncodeToString(salt) + ":" + hex.EncodeToString(hash.Sum(nil)))
	return nil
}

func (p *SaltedPassword) Valid(secret string) bool {
	hexPassword := strings.Split(p.String(), ":")

	salt, err := hex.DecodeString(hexPassword[0])
	if err != nil {
		return false
	}

	saltPlusSecret, err := hex.DecodeString(hexPassword[1])
	if err != nil {
		return false
	}

	hash, err := generateHash(salt, []byte(secret))
	if err != nil {
		return false
	}

	return bytes.Equal(hash.Sum(nil), saltPlusSecret)
}

func (p *SaltedPassword) String() string {
	return string(*p)
}

func (p *SaltedPassword) Value() (driver.Value, error) {
	return string(*p), nil
}

func (p *SaltedPassword) Scan(val interface{}) error {
	*p = SaltedPassword(val.([]byte))
	return nil
}
