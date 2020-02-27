package custom_error

import "errors"

var (
	UserConflict = errors.New("Người dùng đã tồn tại")
	UserNotFound = errors.New("Người dùng không tồn tại")
	SignUpFail   = errors.New("Đăng ký thất bại")
)
