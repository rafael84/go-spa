package resetpassword

import (
	"bytes"
	"fmt"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/gotk/ctx"
	"github.com/gotk/pg"
	"github.com/guregu/null"

	"github.com/rafael84/go-spa/backend/account/user"
	"github.com/rafael84/go-spa/backend/cfg"
	"github.com/rafael84/go-spa/backend/mail"
	"github.com/rafael84/go-spa/backend/random"
)

func createToken(db *pg.Session, userId int64) (*Model, error) {
	// generate key
	key, err := random.New(32)
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
	err = db.Create(resetToken)
	if err != nil {
		return nil, fmt.Errorf("Could not persist reset token: %s", err)
	}

	return resetToken, nil
}

func updateToken(db *pg.Session, token *Model) error {
	err := db.Update(token)
	if err != nil {
		return fmt.Errorf("Could not persist token: %s", err)
	}
	return nil
}

func getToken(db *pg.Session, key string) (*Model, error) {
	resetToken, err := db.FindOne(&Model{}, "key = $1", key)
	if err != nil {
		return nil, err
	}
	return resetToken.(*Model), nil
}

func sendEmail(c *ctx.Context, u *user.Model) {
	var body bytes.Buffer

	db := c.Vars["db"].(*pg.Session)

	resetToken, err := createToken(db, u.Id.NullInt64.Int64)
	if err != nil {
		log.Errorf("Unable to create a new reset token: %s", err)
		return
	}

	body.WriteString("Access this link: ")
	body.WriteString(fmt.Sprintf("%s/#/reset-password/step2/", cfg.Server.BasePath()))
	body.WriteString(resetToken.Key)

	err = mail.NewEmailAccount(
		cfg.Email.Identity,
		cfg.Email.Username,
		cfg.Email.Password,
		cfg.Email.Host,
		cfg.Email.Port,
	).Send(&mail.Message{
		From:    cfg.Email.From,
		To:      []string{u.Email},
		Subject: cfg.Email.Subject,
		Body:    body.Bytes(),
	})

	if err != nil {
		log.Errorf("Unable to send email: %s", err)
		return
	}
}
