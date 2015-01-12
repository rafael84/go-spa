package account

import "github.com/guregu/null"

type SignUpForm struct {
	FirstName     string `json:"firstName"`
	LastName      string `json:"lastName"`
	Email         string `json:"email"`
	Password      string `json:"password"`
	PasswordAgain string `json:"passwordAgain"`
}

type SignInForm struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type ResetPasswordForm struct {
	Email string `json:"email"`
}

type MeForm struct {
	Id       null.Int     `json:"id"`
	Email    string       `json:"email"`
	JsonData UserJsonData `json:"jsonData,omitempty"`
}
