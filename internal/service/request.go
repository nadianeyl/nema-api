package service

import "github.com/nadianeyl/nema-api/internal/validator"

type RegisterUserRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func ValidateEmail(v *validator.Validator, email string) {
	v.Check(email != "", "email", "email is required")
}

func ValidatePasswordPlaintext(v *validator.Validator, password string) {
	v.Check(password != "", "password", "password is required")
	v.Check(len(password) >= 8, "password", "password should be at least 8 bytes long")
	v.Check(len(password) <= 72, "password", "password should not be more than 72 bytes long")
}

func ValidateRegisterUserReq(v *validator.Validator, req *RegisterUserRequest) {
	v.Check(len(req.Name) <= 500, "name", "name should not be more than 500 bytes long")

	ValidateEmail(v, req.Email)

	ValidatePasswordPlaintext(v, req.Password)
}
