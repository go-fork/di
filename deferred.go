// Package di cung cấp dependency injection container cho framework Fork.
package di

// ServiceProviderDeferred mở rộng ServiceProvider cho các dịch vụ cần deferred boot.
//
// Mục đích:
//   - Cho phép boot dịch vụ sau khi xử lý HTTP request (ví dụ: cleanup, flush, async job).
//
// Tính năng:
//   - DeferredBoot được gọi sau khi request kết thúc.
//
// Tham số:
//   - app: interface{} — ứng dụng chính hoặc struct có Container().
//
// Lưu ý:
//   - DeferredBoot chỉ nên dùng cho các tác vụ không blocking request.
type ServiceProviderDeferred interface {
	ServiceProvider

	// DeferredBoot được gọi sau khi HTTP request được xử lý.
	//
	// app: interface{} — ứng dụng chính hoặc struct có Container().
	// Có thể panic nếu cleanup thất bại.
	DeferredBoot(app interface{})
}
