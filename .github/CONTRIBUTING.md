# Contributing to go-fork/di

Cảm ơn bạn đã quan tâm đến việc đóng góp cho dự án go-fork/di! Đây là một số hướng dẫn để giúp quá trình đóng góp trở nên suôn sẻ và hiệu quả.

## Quy trình đóng góp

1. **Fork repository** và tạo branch từ `develop`.
2. **Viết code** và đảm bảo tuân thủ tiêu chuẩn code của dự án.
3. **Viết tests** cho code mới của bạn.
4. **Đảm bảo tất cả tests đều pass**.
5. **Cập nhật tài liệu** nếu cần thiết.
6. **Tạo pull request** tới branch `develop`.

## Tiêu chuẩn code

- Tuân thủ [Effective Go](https://golang.org/doc/effective_go.html)
- Sử dụng `go fmt` trước khi commit
- Đảm bảo code của bạn pass `go vet` và `golangci-lint`
- Viết comment theo chuẩn [godoc](https://blog.golang.org/godoc-documenting-go-code)
- Đặt tên biến, hàm, package rõ ràng và có ý nghĩa
- Kiểm tra lỗi và xử lý lỗi theo chuẩn Go

## Viết test

- Viết unit tests cho tất cả các public functions/methods
- Đảm bảo code coverage tối thiểu 70%
- Test phải thể hiện rõ ràng cách sử dụng và các use case

## Tạo Pull Request

- Mô tả rõ ràng những gì PR của bạn làm
- Tham chiếu đến issues liên quan (nếu có)
- Đảm bảo tất cả CI checks đều pass

## Quy tắc commit

Tuân thủ [Conventional Commits](https://www.conventionalcommits.org/):

- `feat`: Thêm tính năng mới
- `fix`: Sửa lỗi
- `docs`: Thay đổi tài liệu
- `style`: Thay đổi không ảnh hưởng đến nghĩa của code
- `refactor`: Thay đổi code không thêm tính năng hoặc sửa lỗi
- `perf`: Cải thiện hiệu suất
- `test`: Thêm hoặc sửa tests
- `chore`: Thay đổi công cụ build, package manager, etc.

Ví dụ: `feat(container): thêm phương thức Reset()`

## Các tài nguyên hữu ích

- [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- [Effective Go](https://golang.org/doc/effective_go.html)
- [Go Doc Comments](https://tip.golang.org/doc/comment)

Cảm ơn bạn đã đóng góp cho dự án!
