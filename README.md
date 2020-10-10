# Ứng dụng tổng hợp bài viết kỹ thuật chất lượng từ các blog IT

## Mở đầu
- Cập nhập kiến thức thường xuyên trong lĩnh vực IT là việc hết sức quan trọng, hiện nay có rất nhiều các blog IT được lập nên để chia sẻ kiến thức cho mọi người ví dụ viblo, tôi đi code dạo, the full snack, kipalog,... (chỉ nêu tên một số). Tuy nhiên có nhược điểm là người dùng phải truy cập các trang web khác nhau để cập nhật những bài viết mới, vì thế dự án này được lập ra với ý tưởng ban đầu là gom các bài viết hay nhất (trending) từ các blog IT của cộng đồng công nghệ Việt Nam từ nhiều nguồn lại với nhau, hiển thị trên một giao diện đơn giản để theo dõi nhanh chóng hơn nhằm đem lại nhiều nguồn kiến thức chất lượng tới mọi người.

## Mô tả
- Ứng dụng web sẽ tự động thu thập các thông tin về bài viết gồm bài viết và link liên kết bài viết đến trang nguồn nằm trong top bài viết hot nhất về và cập nhật theo thời gian thực (đảm bảo không thu thập nội dung bài viết, vi phạm nội dung bản quyền tác giả).
- Người dùng là khách ghé qua ứng dụng có thể xem thông tin, đọc các bài viết (chuyển tiếp link tới bài viết gốc).
- Khi người dùng muốn đánh dấu các bài viết mình đã xem (vì bài viết đó hay, để lần sau vào đỡ mất công tìm lại) thì người dùng cần đăng ký tài khoản trên ứng dụng để sử dụng.
  - Ứng dụng sẽ có hình thức đăng ký tài khoản là **đăng ký thông thường**, chức năng **đăng ký qua ứng dụng khác** đang được phát triển.
  - Vì việc đăng ký chỉ để sử dụng chức năng cơ bản nhất định nên thông tin đăng ký sẽ bao gồm **email**, **tên tài khoản**, **mật khẩu**.
- Khi hoàn thành thông tin đăng ký sẽ có một email gửi tới địa chỉ email người dùng để xác thực email, người dùng tìm kiếm trong phần **thư rác**, nhấp link đó người dùng sẽ phải nhập lại mật khẩu để hoàn thành xác thực tài khoản.
- Sau khi đăng ký xong tiến hành đăng nhập vào ứng dụng để sử dụng, người dùng có thể xem thông tin tài khoản của mình, cập nhật thông tin tài khoản, mỗi người dùng sẽ có 1 kho lưu trữ để quản lý các bài viết mình quan tâm, khi người dùng thích bài viết nào thì đánh dấu bài viết đấy vào kho lưu trữ, khi không cần nữa có thể xóa bài viết khỏi kho lưu trữ đi.
- Trong trường hợp người dùng quên mật khẩu đăng nhập, người dùng có thể tạo lại mật khẩu bằng cách sử dụng chức năng quên mật khẩu, điền vào thông tin email đăng ký đã xác thực, một email sẽ được gửi tới email đó (trong hộp thư rác), người dùng nhấp vào link và cập nhập lại mật khẩu của mình.

## Người phát triển dự án
1. Nguyễn Đắc Toàn, Email: nguyendactoankma@gmail.com
   Facebook: https://www.facebook.com/toan.nguyen.31392410/

## Tham khảo
- Ứng dụng được xây dựng dựa trên nền tảng của khóa học [golang-flutter](https://www.code4func.com/course/golang-flutter) tại [code4func](https://www.code4func.com/).
