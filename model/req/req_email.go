package req

type ReqEmail struct {
	Email string `json:"email,omitempty" validate:"required,email"`
}
