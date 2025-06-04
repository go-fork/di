// Package di cung cấp hệ thống dependency injection (DI) hiện đại, linh hoạt cho Fork framework.
//
// # Tổng quan
//
// Package này hiện thực DI Container chuẩn, lấy cảm hứng từ các framework lớn (Laravel, Spring), tối ưu cho Go:
//   - Tuân thủ SOLID principles, separation of concerns, dependency inversion.
//   - Hỗ trợ Service-Repository pattern, Adapter pattern, Service Provider pattern.
//   - Type-safe interfaces, loại bỏ type assertion và runtime casting.
//   - Cho phép mở rộng, kiểm soát, testability và maintainability tối đa cho ứng dụng Go.
//
// # Thành phần chính
//
//   - container.go: Định nghĩa struct Container, các phương thức quản lý dependency (Bind, Singleton, Instance, Alias, Make, MustMake, Bound, Reset, Call).
//   - binding.go: Định nghĩa BindingFunc (factory function cho dependency).
//   - provider.go: Định nghĩa interface ServiceProvider, ServiceProviderDeferred, Application (chuẩn hóa contract cho module/service provider/app).
//   - deferred.go: Hỗ trợ deferred service provider (DeferredBoot sau HTTP request).
//   - application.go: Chuẩn hóa interface Application cho app sử dụng DI, quản lý provider, binding, resolve, call.
//   - loader.go: Định nghĩa contract và thực thi cho ModuleLoader, hỗ trợ nạp module/service provider động.
//
// # Chuẩn & Best Practice
//
//   - Tách biệt logic khởi tạo, quản lý vòng đời, resolve dependency.
//   - Đăng ký binding qua factory function, hỗ trợ singleton, instance, alias.
//   - ServiceProvider pattern cho phép module hóa, mở rộng hệ thống.
//   - Hỗ trợ test dễ dàng (inject mock, reset container, ...).
//   - Tài liệu, comment theo chuẩn godoc, mỗi file/module đều có doc.go riêng.
//
// # Sơ đồ sử dụng
//
//  1. Khởi tạo container: container := di.New()
//  2. Đăng ký binding/singleton/instance/alias.
//  3. Đăng ký ServiceProvider (thường qua app.Register hoặc loader).
//  4. Resolve dependency qua Make/MustMake hoặc inject tự động qua Call.
//  5. Có thể reset hoặc kiểm tra trạng thái container qua Bound/Reset.
//
// # Ví dụ sử dụng
//
//	container := di.New()
//	container.Bind("logger", func(c di.Container) interface{} { return log.NewLogger() })
//	container.Singleton("db", func(c di.Container) interface{} { return database.New("dsn") })
//	container.Instance("config", configManager)
//	container.Alias("logger", "log")
//	logger := container.MustMake("logger").(log.Logger)
//	result, err := container.Call(func(l log.Logger, db database.DB) error { l.Info("Hello"); return db.Ping() })
//
// # Service Provider Pattern
//
//	type MyProvider struct{}
//	func (p *MyProvider) Register(app Application) { ... }
//	func (p *MyProvider) Boot(app Application)     { ... }
//
// # Lưu ý
//
// - Đảm bảo tuân thủ SOLID, separation of concerns, best practices Go.
// - Sử dụng DI để giảm coupling, tăng khả năng test và mở rộng.
// - ServiceProvider nên dùng cho các module lớn hoặc tích hợp bên ngoài.
// - Có thể sử dụng mocks cho unit test các service phụ thuộc DI.
// - Mỗi file thành phần nên có doc.go riêng, mô tả rõ vai trò và contract.
//
// # Xem thêm
//
//   - container.go: Định nghĩa Container, các phương thức quản lý dependency
//   - binding.go: Định nghĩa BindingFunc
//   - provider.go: Định nghĩa ServiceProvider, ServiceProviderDeferred, Application
//   - deferred.go: Deferred service provider
//   - application.go: Chuẩn hóa interface Application
//   - loader.go: ModuleLoader contract và thực thi
//   - mocks/: Chứa các mock objects cho các interface chính, được tạo bằng mockery
//
// Package này là nền tảng cho cấu trúc hệ thống Fork, cho phép xây dựng các ứng dụng theo kiến trúc mô-đun với khả năng mở rộng cao.
package di
