package req

type ReqUpdateUser struct {
	FullName string `json:"full_name,omitempty"`
	Password string `json:"password,omitempty"`
	Confirm  string `json:"confirm,omitempty"`
}
