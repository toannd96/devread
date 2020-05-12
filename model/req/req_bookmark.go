package req

type ReqBookmark struct {
	RepoName string `json:"repo,omitempty" validate:"required"`
}
