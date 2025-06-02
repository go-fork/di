package di

// Application định nghĩa interface cơ bản cho ứng dụng sử dụng DI container.
//
// Mục đích:
//   - Chuẩn hóa contract cho các ứng dụng sử dụng DI, hỗ trợ quản lý service provider, binding, instance, alias, resolve, call.
//   - Đảm bảo khả năng mở rộng, kiểm soát vòng đời và cấu hình các thành phần của ứng dụng.
//
// Tính năng:
//   - Truy cập container DI.
//   - Đăng ký, boot, quản lý service provider.
//   - Đăng ký binding, singleton, instance, alias.
//   - Resolve dependency, gọi hàm tự động inject dependency.
//
// Các phương thức:
//
//   - Container() *Container
//     Trả về DI container của ứng dụng.
//
//   - RegisterServiceProviders() error
//     Đăng ký tất cả service provider. Trả về error nếu có lỗi khi đăng ký.
//
//   - RegisterWithDependencies() error
//     Đăng ký và sắp xếp providers theo thứ tự phụ thuộc. Trả về error nếu có circular dependency.
//
//   - BootServiceProviders() error
//     Boot tất cả service provider. Trả về error nếu có lỗi khi boot.
//
//   - Register(provider ServiceProvider)
//     Đăng ký một service provider vào ứng dụng. Panic nếu provider nil hoặc không hợp lệ.
//
//   - Boot() error
//     Khởi động tất cả các service provider đã đăng ký. Trả về error nếu boot thất bại.
//
//   - Bind(abstract string, concrete BindingFunc)
//     Đăng ký binding cho abstract type với factory function. Panic nếu abstract rỗng hoặc trùng.
//
//   - Singleton(abstract string, concrete BindingFunc)
//     Đăng ký singleton binding. Chỉ tạo một instance duy nhất cho abstract type.
//
//   - Instance(abstract string, instance interface{})
//     Đăng ký instance đã khởi tạo sẵn vào container.
//
//   - Alias(abstract, alias string)
//     Đăng ký alias cho abstract type. Resolve alias sẽ trả về instance của abstract gốc.
//
//   - Make(abstract string) (interface{}, error)
//     Resolve dependency theo tên. Trả về error nếu không tìm thấy hoặc binding lỗi.
//
//   - MustMake(abstract string) interface{}
//     Resolve dependency, panic nếu không tìm thấy hoặc binding lỗi.
//
//   - Call(callback interface{}, additionalParams ...interface{}) ([]interface{}, error)
//     Gọi hàm và tự động resolve các dependency qua reflection. Trả về error nếu không resolve được tham số.
//
// Exceptions/Lưu ý:
//   - Các method resolve (Make, MustMake, Call) có thể panic hoặc trả về error nếu dependency không tồn tại hoặc binding lỗi.
//   - Nên implement interface này cho app chính để tận dụng hệ sinh thái provider của Fork.
type Application interface {
	// Container trả về DI container của ứng dụng.
	//
	// Mục đích:
	//   - Cho phép truy cập trực tiếp vào dependency injection container để đăng ký hoặc resolve các dependency.
	// Trả về:
	//   - Container: instance container hiện tại của ứng dụng.
	Container() Container

	// RegisterServiceProviders đăng ký tất cả các service provider đã cấu hình cho ứng dụng.
	//
	// Mục đích:
	//   - Đảm bảo các service provider được đăng ký vào container trước khi boot.
	// Trả về:
	//   - error: lỗi nếu provider không hợp lệ hoặc đăng ký thất bại.
	// Exceptions:
	//   - Trả về error nếu provider nil, duplicate hoặc có lỗi khi đăng ký.
	RegisterServiceProviders() error

	// RegisterWithDependencies đăng ký và sắp xếp providers theo thứ tự phụ thuộc
	//
	// Mục đích:
	//   - Tự động sắp xếp và đăng ký providers theo thứ tự phụ thuộc
	//   - Phát hiện circular dependencies
	// Trả về:
	//   - error: Lỗi nếu có circular dependency hoặc không tìm thấy provider yêu cầu
	RegisterWithDependencies() error

	// BootServiceProviders boot tất cả các service provider đã đăng ký.
	//
	// Mục đích:
	//   - Khởi tạo các tài nguyên, kết nối hoặc logic phụ thuộc vào provider sau khi đã đăng ký.
	// Trả về:
	//   - error: lỗi nếu boot provider thất bại.
	// Exceptions:
	//   - Trả về error nếu provider boot lỗi hoặc thiếu dependency.
	BootServiceProviders() error

	// Register đăng ký một service provider vào ứng dụng.
	//
	// Mục đích:
	//   - Cho phép thêm động các service provider vào container.
	// Tham số:
	//   - provider: ServiceProvider — provider cần đăng ký.
	// Exceptions:
	//   - Panic nếu provider nil hoặc không hợp lệ.
	Register(provider ServiceProvider)

	// Boot khởi động tất cả các service provider đã đăng ký.
	//
	// Mục đích:
	//   - Thực thi logic khởi tạo của provider (Boot) sau khi đã đăng ký.
	// Trả về:
	//   - error: lỗi nếu boot provider thất bại.
	// Exceptions:
	//   - Trả về error nếu provider boot lỗi hoặc thiếu dependency.
	Boot() error

	// Bind đăng ký một binding với container.
	//
	// Mục đích:
	//   - Đăng ký factory function cho abstract type, cho phép resolve động các dependency.
	// Tham số:
	//   - abstract: string — tên abstract type.
	//   - concrete: BindingFunc — factory function tạo instance.
	// Exceptions:
	//   - Panic nếu abstract rỗng hoặc đã tồn tại binding.
	Bind(abstract string, concrete BindingFunc)

	// Singleton đăng ký một singleton binding với container.
	//
	// Mục đích:
	//   - Đảm bảo chỉ tạo một instance duy nhất cho abstract type trong suốt vòng đời app.
	// Tham số:
	//   - abstract: string — tên abstract type.
	//   - concrete: BindingFunc — factory function tạo instance.
	Singleton(abstract string, concrete BindingFunc)

	// Instance đăng ký một instance đã khởi tạo sẵn vào container.
	//
	// Mục đích:
	//   - Cho phép inject các instance đã tồn tại (ví dụ: config, logger).
	// Tham số:
	//   - abstract: string — tên abstract type.
	//   - instance: interface{} — instance đã khởi tạo.
	Instance(abstract string, instance interface{})

	// Alias đăng ký một alias cho abstract type.
	//
	// Mục đích:
	//   - Cho phép truy cập dependency qua nhiều tên khác nhau.
	// Tham số:
	//   - abstract: string — tên abstract gốc.
	//   - alias: string — tên alias.
	Alias(abstract, alias string)

	// Make resolve một dependency từ container.
	//
	// Mục đích:
	//   - Resolve instance theo abstract type đã đăng ký.
	// Tham số:
	//   - abstract: string — tên abstract type.
	// Trả về:
	//   - interface{}: instance đã resolve.
	//   - error: lỗi nếu không tìm thấy hoặc binding lỗi.
	Make(abstract string) (interface{}, error)

	// MustMake resolve một dependency từ container, panic nếu lỗi.
	//
	// Mục đích:
	//   - Resolve instance, panic nếu không tìm thấy hoặc binding lỗi.
	// Tham số:
	//   - abstract: string — tên abstract type.
	// Trả về:
	//   - interface{}: instance đã resolve.
	// Exceptions:
	//   - Panic nếu không resolve được dependency.
	MustMake(abstract string) interface{}

	// Call gọi một hàm và tự động resolve các dependency.
	//
	// Mục đích:
	//   - Tự động inject các dependency vào callback function qua reflection.
	// Tham số:
	//   - callback: interface{} — function cần gọi.
	//   - additionalParams: ...interface{} — các tham số bổ sung (ưu tiên inject).
	// Trả về:
	//   - []interface{}: kết quả trả về của callback.
	//   - error: lỗi nếu không resolve được tham số hoặc callback không hợp lệ.
	Call(callback interface{}, additionalParams ...interface{}) ([]interface{}, error)
}
