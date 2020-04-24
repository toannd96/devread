package req

type PasswordSubmit struct {
	Password string `json:"password,omitempty" validate:"required,pwd"`
	Confirm  string `json:"confirm,omitempty" validate:"required,pwd"`
}
