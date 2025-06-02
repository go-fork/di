package di

// ServiceProviderDeferred mở rộng ServiceProvider cho các dịch vụ cần deferred boot.
//
// Mục đích:
//   - Cho phép boot dịch vụ sau khi xử lý HTTP request (ví dụ: cleanup, flush, async job).
//   - Tuân thủ SOLID principles với type-safe contract.
//
// Tính năng:
//   - DeferredBoot được gọi sau khi request kết thúc.
//   - Type-safe interface, không cần type assertion.
//
// Tham số:
//   - app: Application — ứng dụng với DI container và lifecycle management.
//
// Lưu ý:
//   - DeferredBoot chỉ nên dùng cho các tác vụ không blocking request.
//   - Tuân thủ Dependency Inversion Principle với typed contract.
type ServiceProviderDeferred interface {
	ServiceProvider

	// DeferredBoot được gọi sau khi HTTP request được xử lý.
	//
	// app: Application — ứng dụng với typed interface, không cần casting.
	// Có thể panic nếu cleanup thất bại.
	DeferredBoot(app Application)
}
