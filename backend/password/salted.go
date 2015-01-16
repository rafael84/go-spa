package password

import (
	"bytes"
	"database/sql/driver"
	"encoding/hex"
	"strings"
)

type Salted string

func (p *Salted) Encode(secret string) error {
	salt, err := generateSalt([]byte(secret))
	if err != nil {
		return err
	}

	hash, err := generateHash(salt, []byte(secret))
	if err != nil {
		return err
	}
	*p = Salted(hex.EncodeToString(salt) + ":" + hex.EncodeToString(hash.Sum(nil)))
	return nil
}

func (p *Salted) Valid(secret string) bool {
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

func (p *Salted) String() string {
	return string(*p)
}

func (p *Salted) Value() (driver.Value, error) {
	return string(*p), nil
}

func (p *Salted) Scan(val interface{}) error {
	*p = Salted(val.([]byte))
	return nil
}
