package user

import (
	"errors"
	"fmt"

	"github.com/gotk/pg"
	"github.com/guregu/null"
	"github.com/rafael84/go-spa/backend/base"
)

type userService struct {
	session *pg.Session
}

func NewUserService(session *pg.Session) *userService {
	return &userService{session}
}

func (us *userService) Create(email, password string, userJsonData *UserJsonData) (*Model, error) {
	// encode password
	var saltedPassword base.SaltedPassword
	err := saltedPassword.Encode(password)
	if err != nil {
		return nil, errors.New("Could not encode password")
	}

	// create new user structure
	user := &Model{
		Id:       null.NewInt(0, false),
		State:    UserStateActive,
		Email:    email,
		Password: saltedPassword,
	}

	// fill user structure with additional data
	err = user.JsonData.Encode(userJsonData)
	if err != nil {
		return nil, fmt.Errorf("Failed to marshal json data: %s", err)
	}

	// create new user in database
	err = us.session.Create(user)
	if err != nil {
		return nil, fmt.Errorf("Could not persist user: %s", err)
	}

	return user, nil
}

func (us *userService) Update(user *Model) error {
	err := us.session.Update(user)
	if err != nil {
		return fmt.Errorf("Could not persist user: %s", err)
	}
	return nil
}

func (us *userService) GetById(id int64) (*Model, error) {
	user, err := us.session.FindOne(&Model{}, "id = $1", id)
	if err != nil {
		return nil, err
	}
	return user.(*Model), nil
}

func (us *userService) GetByEmail(email string) (*Model, error) {
	user, err := us.session.FindOne(&Model{}, "email = $1", email)
	if err != nil {
		return nil, err
	}
	return user.(*Model), nil
}
