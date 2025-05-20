// Package di cung cấp dependency injection container cho framework Fork.
package di

// ServiceProvider định nghĩa contract cho các service provider trong hệ thống DI.
//
// Mục đích:
//   - Cho phép module hoặc package đăng ký các dịch vụ, binding vào container một cách mô-đun.
//   - Tách biệt logic khởi tạo, cấu hình dịch vụ khỏi phần còn lại của ứng dụng.
//
// Tính năng:
//   - Đăng ký các binding, singleton, instance vào container.
//   - Thực hiện các thao tác khởi tạo bổ sung sau khi đăng ký (Boot).
//
// Tham số:
//   - app: interface{} — Thường là ứng dụng chính hoặc struct có method Container() *Container.
//
// Trả về:
//   - Không trả về giá trị, nhưng có thể panic nếu đăng ký binding lỗi.
//
// Lưu ý:
//   - Register phải idempotent (gọi nhiều lần không gây lỗi trạng thái).
//   - Boot có thể dùng để khởi tạo tài nguyên phụ thuộc vào các binding đã đăng ký.
type ServiceProvider interface {
	// Register đăng ký các bindings vào container.
	//
	// app: interface{} — ứng dụng chính hoặc struct có Container().
	// Có thể panic nếu binding không hợp lệ.
	Register(app interface{})

	// Boot được gọi sau khi tất cả các service provider đã được đăng ký.
	//
	// app: interface{} — ứng dụng chính hoặc struct có Container().
	// Có thể panic nếu khởi tạo thất bại.
	Boot(app interface{})
}
