package user

type LoginDto struct {
	Email    string `json:"email"    validate:"required, is_string, is_email"`
	Password string `json:"password" validate:"required, is_string"`
}
