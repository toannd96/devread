package custom_error

import "errors"

var (
	UserConflict = errors.New("Người dùng đã tồn tại ")
	SignUpFail   = errors.New("Đăng ký thất bại")
)
