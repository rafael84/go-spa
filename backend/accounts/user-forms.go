package accounts

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
