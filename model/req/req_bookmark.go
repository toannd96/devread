package req

type ReqBookmark struct {
	PostName string `json:"post,omitempty" validate:"required"`
}
