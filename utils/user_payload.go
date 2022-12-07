package utils

// UserSignUpPayload
// json payload for user request
type UserSignUpPayload struct {
	FirstName string `json:"first_name" validate:"required"`
	LastName  string `json:"last_name" validate:"required"`
	Username  string `json:"username" validate:"required"`
	Email     string `json:"email" validate:"required,email"`
	Password1 string `json:"password1" validate:"required"`
	Password2 string `json:"password2" validate:"required"`
}

// UserSignInPayload
// json payload for user request
type UserSignInPayload struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}
