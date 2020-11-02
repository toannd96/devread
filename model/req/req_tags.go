package req

type ReqTags struct {
	Tags string `json:"tags,omitempty" validate:"required"`
}
