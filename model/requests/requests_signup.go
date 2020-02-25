package requests

type ReqSignUp struct {
	FullName string `validate:"required"`
	Email    string `validate:"required"`
	Password string `validate:"required"`
}
