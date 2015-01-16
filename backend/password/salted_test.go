package password

import (
	"testing"
)

func TestSaltedPasswordValid(t *testing.T) {
	secret := "qwer1234"
	var pw Salted

	err := pw.Encode(secret)
	if err != nil {
		t.Errorf("Unable to encode password: %s", err)
	}

	if !pw.Valid(secret) {
		t.Errorf("SaltedPassword validation failed! \nSaltedPassword: %s", pw)
	}

	if pw.Valid(secret + "wrong") {
		t.Errorf("SaltedPassword validation failed! \nSaltedPassword: %s", pw)
	}
}
