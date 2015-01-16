package user

import (
	"errors"
	"fmt"

	"github.com/gotk/pg"
	"github.com/guregu/null"

	"github.com/rafael84/go-spa/backend/password"
)

func Create(db *pg.Session, email, pw string, role int, userJsonData *UserJsonData) (*Model, error) {
	// encode password
	var saltedPassword password.Salted
	err := saltedPassword.Encode(pw)
	if err != nil {
		return nil, errors.New("Could not encode password")
	}

	// create new user structure
	user := &Model{
		Id:       null.NewInt(0, false),
		State:    UserStateActive,
		Email:    email,
		Password: saltedPassword,
		Role:     role,
	}

	// fill user structure with additional data
	err = user.JsonData.Encode(userJsonData)
	if err != nil {
		return nil, fmt.Errorf("Failed to marshal json data: %s", err)
	}

	// create new user in database
	err = db.Create(user)
	if err != nil {
		return nil, fmt.Errorf("Could not persist user: %s", err)
	}

	return user, nil
}

func Update(db *pg.Session, user *Model) error {
	err := db.Update(user)
	if err != nil {
		return fmt.Errorf("Could not persist user: %s", err)
	}
	return nil
}

func GetById(db *pg.Session, id int64) (*Model, error) {
	user, err := db.FindOne(&Model{}, "id = $1", id)
	if err != nil {
		return nil, err
	}
	return user.(*Model), nil
}

func GetByEmail(db *pg.Session, email string) (*Model, error) {
	user, err := db.FindOne(&Model{}, "email = $1", email)
	if err != nil {
		return nil, err
	}
	return user.(*Model), nil
}
