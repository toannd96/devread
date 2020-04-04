package req

type ReqtSignUp struct {
	Email    string `json:"email,omitempty" validate:"required"`
	FullName string `json:"full_name,omitempty" validate:"required"`
	Password string `json:"password,omitempty" validate:"pwd"`
}
