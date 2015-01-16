package cfg

import (
	"os"
	"strconv"
)

type email struct {
	Identity string
	Host     string
	Port     int
	Username string
	Password string

	From    string
	Subject string

	Regex string
}

var Email email

func init() {
	Email.Identity = os.Getenv("EMAIL_IDENTITY")
	Email.Host = os.Getenv("EMAIL_HOST")
	Email.Port, _ = strconv.Atoi(os.Getenv("EMAIL_PORT"))
	Email.Username = os.Getenv("EMAIL_USERNAME")
	Email.Password = os.Getenv("EMAIL_PASSWORD")

	Email.From = os.Getenv("EMAIL_FROM")
	Email.Subject = os.Getenv("EMAIL_SUBJECT")

	Email.Regex = `\b[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,4}\b`
}
