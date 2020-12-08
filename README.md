# DevReading - Ứng dụng tổng hợp bài viết kỹ thuật chất lượng từ các blog IT

## Bài toán
- Cập nhập kiến thức thường xuyên trong lĩnh vực IT là việc hết sức quan trọng, hiện nay có rất nhiều các blog IT được lập nên để chia sẻ kiến thức cho mọi người ví dụ viblo, tôi đi code dạo, the full snack, kipalog,... (chỉ nêu tên một số).
- Tuy nhiên có nhược điểm là người dùng phải truy cập các trang web khác nhau để cập nhật những bài viết mới.
- Vì thế dự án này được lập ra với ý tưởng ban đầu là gom các bài viết hay nhất từ các blog IT của cộng đồng công nghệ Việt Nam từ nhiều nguồn lại với nhau, hiển thị trên một giao diện đơn giản để theo dõi nhanh chóng hơn nhằm đem lại nhiều nguồn kiến thức chất lượng tới mọi người.

## Mô tả ứng dụng
- Ứng dụng sẽ tự động thu thập các thông tin về bài viết gồm **tên bài viết**, **thể loại**, **link liên kết bài viết đến trang nguồn** nằm trong top những bài viết hay nhất về và cập nhật theo thời gian thực (đảm bảo không thu thập nội dung bài viết, vi phạm nội dung bản quyền tác giả).Tệp dữ liệu mà ứng dụng thu thập, xem tại [đây](https://github.com/dactoankmapydev/tech_posts_trending/blob/master/tong_hop_bai_viet.csv)
- Thay vì mò mẫm vào từng trang blog thì nay chỉ cần vào 1 trang duy nhất để xem những bài đăng hấp dẫn từ các trang blog it như:
  - [Viblo](https://viblo.asia/) mục [trending](https://viblo.asia/trending), [series](https://viblo.asia/series), một phần mục [newest](https://viblo.asia/newest)
  - [Toidicodedao](https://toidicodedao.com/) mục [chuyện coding](https://toidicodedao.com/category/chuyen-coding/)
  - [Yellowcodebooks](https://yellowcodebooks.com/) mục [java](https://yellowcodebooks.com/category/lap-trinh-java/), [android](https://yellowcodebooks.com/category/lap-trinh-android/)
  - [Quan-cam](https://quan-cam.com/) mục [programming](https://quan-cam.com/tags/programming)
  - [Codeaholicguy](https://codeaholicguy.com/) mục [chuyện coding](https://codeaholicguy.com/category/chuyen-coding)
  - [Thefullsnack](https://thefullsnack.com/)
  ...
  
- Người dùng là khách ghé qua ứng dụng có thể xem thông tin, đọc các bài viết (chuyển tiếp link tới bài viết gốc).
- Người dùng có thể tìm kiếm các bài viết theo thể loại (tags) của bài viết.
- Khi người dùng muốn lưu lại các bài viết mình thấy hay (vì bài viết đó hay nhưng chưa có thời gian đọc ngay, lưu lại để lần sau vào đỡ mất công tìm lại) thì người dùng cần đăng ký tài khoản trên ứng dụng để sử dụng.
  - Ứng dụng sẽ có hình thức đăng ký tài khoản là **đăng ký thông thường**, chức năng **đăng ký qua ứng dụng khác** đang được phát triển.
  - Vì việc đăng ký chỉ để sử dụng chức năng cơ bản nhất định nên thông tin đăng ký sẽ bao gồm **email**, **tên tài khoản**, **mật khẩu**.
- Khi hoàn thành thông tin đăng ký sẽ có một email gửi tới địa chỉ email người dùng để xác thực email, người dùng tìm kiếm trong phần **thư rác**, nhấp link đó người dùng sẽ phải nhập lại mật khẩu để hoàn thành xác thực tài khoản.
- Sau khi đăng ký xong tiến hành đăng nhập vào ứng dụng để sử dụng, người dùng có thể xem thông tin tài khoản của mình, cập nhật thông tin tài khoản, mỗi người dùng sẽ có 1 kho lưu trữ để quản lý các bài viết mình đánh dấu trước đó, khi người dùng thích bài viết nào thì thêm bài viết đấy vào kho lưu trữ, khi không cần nữa có thể bỏ đánh dấu (xóa) khỏi kho lưu trữ.
- Trong trường hợp người dùng quên mật khẩu đăng nhập, người dùng có thể tạo lại mật khẩu bằng cách sử dụng chức năng quên mật khẩu, điền vào thông tin email đăng ký đã xác thực, một email sẽ được gửi tới email đó (trong hộp thư rác), người dùng nhấp vào link và cập nhập lại mật khẩu của mình.

## Công nghệ sử dụng
- Backend: [echo go web framework](https://echo.labstack.com/), [colly scraping framework](http://go-colly.org/)

## Chạy ứng dụng
- Sử dụng [swagger](https://swagger.io/) để test API
```
  go run main.go
```
- Ứng dụng chạy trên ```http://localhost:3000/swagger/index.html```

## Người phát triển dự án
1. Nguyễn Đắc Toàn
- Email: nguyendactoankma@gmail.com
- Facebook: https://www.facebook.com/toan.nguyen.31392410/
- Linkedin: https://www.linkedin.com/in/dac-toan-nguyen-a94b2a146/

## Tham khảo
- Ứng dụng được xây dựng dựa trên nền tảng của khóa học [golang-flutter](https://www.code4func.com/course/golang-flutter) tại [code4func](https://www.code4func.com/).
