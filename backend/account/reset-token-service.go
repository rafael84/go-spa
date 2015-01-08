package account

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"time"

	"github.com/guregu/null"
	"github.com/rafael84/go-spa/backend/database"
)

type resetTokenService struct {
	session *database.Session
}

func NewResetTokenService(session *database.Session) *resetTokenService {
	return &resetTokenService{session}
}

func (r *resetTokenService) Create(userId int64) (*ResetToken, error) {

	// generate key
	buf := make([]byte, 32, 32)
	_, err := io.ReadFull(rand.Reader, buf)
	if err != nil {
		return nil, err
	}

	// create new reset token structure
	resetToken := &ResetToken{
		Id:         null.NewInt(0, false),
		State:      0,
		Key:        hex.EncodeToString(buf),
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

func (r *resetTokenService) GetByKey(key string) (*ResetToken, error) {
	resetToken, err := r.session.One(&ResetToken{}, "key = $1", key)
	if err != nil {
		return nil, err
	}
	return resetToken.(*ResetToken), nil
}
