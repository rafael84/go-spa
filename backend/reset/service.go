package reset

import (
	"fmt"
	"time"

	"github.com/gotk/pg"
	"github.com/guregu/null"
	"github.com/rafael84/go-spa/backend/base"
)

type resetTokenService struct {
	session *pg.Session
}

func NewResetTokenService(session *pg.Session) *resetTokenService {
	return &resetTokenService{session}
}

func (r *resetTokenService) Create(userId int64) (*ResetToken, error) {

	// generate key
	key, err := base.Random(32)
	if err != nil {
		return nil, err
	}

	// create new reset token structure
	resetToken := &ResetToken{
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

func (r *resetTokenService) Update(token *ResetToken) error {
	err := r.session.Update(token)
	if err != nil {
		return fmt.Errorf("Could not persist token: %s", err)
	}
	return nil
}

func (r *resetTokenService) GetByKey(key string) (*ResetToken, error) {
	resetToken, err := r.session.FindOne(&ResetToken{}, "key = $1", key)
	if err != nil {
		return nil, err
	}
	return resetToken.(*ResetToken), nil
}
