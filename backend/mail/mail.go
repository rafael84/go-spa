package mail

import (
	"bytes"
	"net/smtp"
	"strconv"
)

type Message struct {
	From    string
	To      []string
	Subject string
	Body    []byte
}

func (m *Message) Bytes() []byte {
	var msg bytes.Buffer
	msg.WriteString("From: " + m.From + "\n")
	msg.WriteString("Subject: " + m.Subject + "\n")
	msg.WriteString("\n")
	msg.Write(m.Body)
	return msg.Bytes()
}

type account struct {
	Identity string
	Username string
	Password string
	Host     string
	Port     int

	auth smtp.Auth
}

func (a *account) Send(message *Message) error {
	return smtp.SendMail(a.Host+":"+strconv.Itoa(a.Port), a.auth, a.Username, message.To, message.Bytes())
}

func NewEmailAccount(identity, username, password, host string, port int) *account {
	return &account{
		Identity: identity,
		Username: username,
		Password: password,
		Host:     host,
		Port:     port,

		auth: smtp.PlainAuth(
			identity,
			username,
			password,
			host,
		),
	}
}
