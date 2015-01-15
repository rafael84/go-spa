package cfg

import (
	"os"
)

type email struct {
	From     string
	Subject  string
	Username string
	Password string
}

var Email email

func init() {
	Email.From = os.Getenv("EMAIL_FROM")
	Email.Subject = os.Getenv("EMAIL_SUBJECT")
	Email.Username = os.Getenv("EMAIL_USERNAME")
	Email.Password = os.Getenv("EMAIL_PASSWORD")
}
