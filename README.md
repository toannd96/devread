Xây dựng ứng dụng web thu thập thông tin bài viết trên [blog viblo mục trending](https://viblo.asia/trending) sử dụng Echo Framework (Golang)

## Mô tả
- Ứng dụng (web) sẽ tự động crawl các thông tin về bài viết nằm trong top trending về và cập nhật hàng ngày.
- Người dùng là khách ghé qua ứng dụng có thể xem thông tin, đọc các bài viết (chuyển tiếp link tới bài viết gốc).
- Khi người dùng muốn đánh dấu các bài viết mình đã xem (vì bài viết đó hay, để lần sau vào đỡ mất công tìm lại) thì người dùng cần đăng ký tài khoản trên ứng dụng để sử dụng.
  - Ứng dụng sẽ có 2 hình thức đăng ký tài khoản là **đăng ký thông thường** và **đăng ký thông qua ứng dụng khác như facebook**
  - Vì việc đăng ký chỉ để sử dụng chức năng cơ bản nhất định nên thông tin đăng ký sẽ bao gồm **tên tài khoản**, **mật khẩu**.
- Sau khi có đăng ký xong tiến hành đăng nhập vào ứng dụng để sử dụng, người dùng có thể xem thông tin tài khoản của mình, cập nhật thông tin tài khoản, mỗi người dùng sẽ có 1 kho lưu trữ để quản lý các bài viết mình quan tâm, khi người dùng thích bài viết nào thì đánh dấu bài viết đấy vào kho lưu trữ, khi không cần nữa có thể xóa bài viết khỏi kho lưu trữ đi.
- Về vấn đề xác thực người dùng, ứng dụng sử dụng JWT để xác thực như mô tả phía dưới.

## Xác thực
Sử dụng JWT để xác thực:

- Hai cặp mã thông báo được tạo, mã thông báo truy cập (access token) tồn tại ngắn (15 phút) và mã thông báo làm mới (refresh token) tồn tại lâu (1 ngày).
- Mỗi JWT được tạo bao gồm UUID và ID người dùng.UUID cho phép có các mã thông báo khác nhau cho cùng một người dùng trên các thiết bị khác nhau. Nó cũng giúp vô hiệu hóa một JWT theo ý muốn. Một ví dụ là khi người dùng đăng xuất, JWT ngay lập tức bị vô hiệu vì UUID được sử dụng bị loại bỏ.
- Redis được sử dụng để lưu dữ liệu của JWT (UUID và ID người dùng). Vì Redis là kho lưu trữ khóa-giá trị, UUID đóng vai trò là khóa trong khi ID của người dùng là giá trị. Vì vậy, khi các cặp mã thông báo được tạo (truy cập và làm mới mã thông báo), dữ liệu mã thông báo được lưu trong redis.
- Đối với mỗi yêu cầu được xác thực, mã thông báo truy cập được trích xuất, sau đó dữ liệu của mã thông báo đó được tìm kiếm trong redis. Nếu tìm thấy, yêu cầu được cấp.
- Vì redis có tính năng hết hạn tài liệu, dữ liệu JWT có thể bị xóa khỏi redis sau khi thời gian hết hạn.
- Khi người dùng thực hiện yêu cầu đăng xuất, dữ liệu của JWT được cung cấp sẽ bị xóa khỏi redis.
- Bản chất của mã thông báo làm mới là khi mã thông báo truy cập của người dùng hết hạn, mã thông báo làm mới sẽ được gửi trong yêu cầu tạo bộ mã thông báo truy cập mới và mã thông báo làm mới.


## Cập nhật
- API rate limiting
- Oauth2
