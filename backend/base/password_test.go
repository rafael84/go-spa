package base

import (
	"testing"
)

func TestSaltedPasswordValid(t *testing.T) {
	secret := "qwer1234"
	var password SaltedPassword

	err := password.Encode(secret)
	if err != nil {
		t.Errorf("Unable to encode password: %s", err)
	}

	if !password.Valid(secret) {
		t.Errorf("SaltedPassword validation failed! \nSaltedPassword: %s", password)
	}

	if password.Valid(secret + "wrong") {
		t.Errorf("SaltedPassword validation failed! \nSaltedPassword: %s", password)
	}
}
