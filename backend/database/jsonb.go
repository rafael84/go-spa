package database

import (
	"encoding/json"
	"errors"

	"github.com/guregu/null"
)

type JSONB null.String

func (j *JSONB) Set(src interface{}) error {
	jsonData, err := json.Marshal(src)
	if err != nil {
		return err
	}
	*j = JSONB(null.NewString(string(jsonData), true))
	return nil
}

func (j *JSONB) Get(dst interface{}) error {
	if !j.NullString.Valid {
		return errors.New("Empty JSON data")
	}
	return json.Unmarshal([]byte(j.NullString.String), dst)
}

// MarshalJSON implements the `json.Marshaller` interface
func (j *JSONB) MarshalJSON() ([]byte, error) {
	if j.NullString.Valid {
		return []byte(j.NullString.String), nil
	}
	return []byte("null"), nil
}
