package resetpassword

import (
	"fmt"
	"time"

	"github.com/gotk/pg"
	"github.com/guregu/null"

	"github.com/rafael84/go-spa/backend/base"
)

type Service struct {
	session *pg.Session
}

func NewService(session *pg.Session) *Service {
	return &Service{session}
}

func (r *Service) Create(userId int64) (*Model, error) {
	// generate key
	key, err := base.Random(32)
	if err != nil {
		return nil, err
	}

	// create new reset token structure
	resetToken := &Model{
		Id:         null.NewInt(0, false),
		State:      ResetTokenActive,
		Key:        key,
		Expiration: time.Now().Add(time.Minute * 10),
		UserId:     userId,
	}

	// create new user in database
	err = r.session.Create(resetToken)
	if err != nil {
		return nil, fmt.Errorf("Could not persist reset token: %s", err)
	}

	return resetToken, nil
}

func (r *Service) Update(token *Model) error {
	err := r.session.Update(token)
	if err != nil {
		return fmt.Errorf("Could not persist token: %s", err)
	}
	return nil
}

func (r *Service) GetByKey(key string) (*Model, error) {
	resetToken, err := r.session.FindOne(&Model{}, "key = $1", key)
	if err != nil {
		return nil, err
	}
	return resetToken.(*Model), nil
}
