package requests

type RequestSignUp struct {
	FullName string `json:"fullName,omitempty" validate:"required"`
	Email    string `json:"email,omitempty" validate:"required"`
	Password string `json:"password,omitempty" validate:"pwd"`
}
