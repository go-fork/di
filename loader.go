package di

// ModuleLoaderContract định nghĩa contract cho các loader chịu trách nhiệm nạp module/service provider vào ứng dụng.
//
// Mục đích:
//   - Chuẩn hóa quy trình nạp, đăng ký, boot các module hoặc service provider vào ứng dụng sử dụng DI container.
//   - Tách biệt logic khởi tạo module khỏi phần còn lại của app, tăng khả năng mở rộng và kiểm soát vòng đời module.
//
// Tính năng:
//   - Đăng ký các core service provider (RegisterCoreProviders).
//   - Khởi tạo ứng dụng với các cấu hình và provider cần thiết (BootstrapApplication).
//   - Tải module đơn lẻ hoặc nhiều module (LoadModule, LoadModules).
//
// Tham số:
//   - module: interface{} — module cần nạp, có thể là service provider hoặc các kiểu module khác.
//   - modules: ...interface{} — danh sách module cần nạp.
//
// Trả về:
//   - error: trả về lỗi nếu quá trình nạp, đăng ký, boot module thất bại.
//
// Exceptions/Lưu ý:
//   - Nếu module không hợp lệ hoặc không phải service provider, trả về error rõ ràng.
//   - BootstrapApplication sẽ trả về error nếu đăng ký hoặc boot provider thất bại.
//   - LoadModules dừng lại khi gặp lỗi ở bất kỳ module nào.
//
// Ví dụ sử dụng:
//
//	loader := NewModuleLoader(app)
//	err := loader.BootstrapApplication()
//	err := loader.LoadModules(moduleA, moduleB)
type ModuleLoaderContract interface {
	// BootstrapApplication khởi tạo ứng dụng với các cấu hình và provider cần thiết.
	// Trả về error nếu đăng ký hoặc boot provider thất bại.
	BootstrapApplication() error

	// RegisterCoreProviders đăng ký các core service provider.
	// Trả về error nếu đăng ký thất bại.
	RegisterCoreProviders() error

	// LoadModule tải một module vào ứng dụng.
	// Trả về error nếu module không hợp lệ hoặc nạp thất bại.
	LoadModule(module interface{}) error

	// LoadModules tải nhiều module vào ứng dụng.
	// Trả về error nếu bất kỳ module nào nạp thất bại.
	LoadModules(modules ...interface{}) error
}
