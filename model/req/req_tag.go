package req

type ReqTag struct {
	Tag string `json:"tag,omitempty" validate:"required"`
}
