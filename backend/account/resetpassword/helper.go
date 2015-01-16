package resetpassword

import (
	"bytes"
	"fmt"

	log "github.com/Sirupsen/logrus"
	"github.com/gotk/ctx"
	"github.com/gotk/pg"

	"github.com/rafael84/go-spa/backend/account/user"
	"github.com/rafael84/go-spa/backend/cfg"
	"github.com/rafael84/go-spa/backend/mail"
)

func sendEmail(c *ctx.Context, u *user.Model) {
	var body bytes.Buffer

	resetTokenService := NewService(c.Vars["db"].(*pg.Session))

	resetToken, err := resetTokenService.Create(u.Id.NullInt64.Int64)
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
