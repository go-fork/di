# di

Một container Dependency Injection (DI) hiện đại, nhẹ nhàng cho Go, lấy cảm hứng từ Laravel và các framework hiện đại. Cung cấp đăng ký dịch vụ, giải quyết phụ thuộc và quản lý vòng đời cho việc xây dựng ứng dụng Go có khả năng mở rộng.

## Tính năng
- Mẫu Service Provider với dependency management
- Tự động giải quyết phụ thuộc
- Tự động sắp xếp thứ tự provider dependencies  
- Binding singleton và transient
- Tải dịch vụ trì hoãn (deferred)
- API đơn giản, không cần tạo mã

## Cài đặt

```
go get go.fork.vn/di
```

## Sử dụng

### Cơ bản

```go
import "go.fork.vn/di"

container := di.New()
container.Bind("service", func(c di.Container) interface{} {
    return &MyService{}
})
service := container.MustMake("service").(*MyService)
```

### Singleton

```go
// Đăng ký dịch vụ singleton (chỉ khởi tạo một lần)
container.Singleton("database", func(c di.Container) interface{} {
    return database.Connect("localhost", "user", "pass")
})

// Lấy cùng một instance mỗi lần gọi
db1 := container.MustMake("database").(*Database)
db2 := container.MustMake("database").(*Database)
// db1 == db2 (cùng một instance)
```

### Instance

```go
// Đăng ký instance có sẵn
config := &Config{Debug: true}
container.Instance("config", config)

// Lấy ra instance đã đăng ký
appConfig := container.MustMake("config").(*Config)
```

### Alias

```go
// Đăng ký alias cho service
container.Singleton("logger", func(c di.Container) interface{} {
    return &Logger{}
})
container.Alias("logger", "log")

// Có thể truy cập bằng bất kỳ tên nào
log1 := container.MustMake("logger").(*Logger)
log2 := container.MustMake("log").(*Logger)
// log1 == log2
```

### Tự động Inject Dependencies

```go
// Tự động inject dependencies vào hàm
container.Singleton("userRepo", func(c di.Container) interface{} {
    return &UserRepository{}
})
container.Singleton("userService", func(c di.Container) interface{} {
    return &UserService{
        Repo: c.MustMake("userRepo").(*UserRepository),
    }
})

// Tự động resolve dependencies khi gọi hàm
container.Call(func(userService *UserService) {
    // userService được tự động inject
    userService.DoSomething()
})
```

### Service Provider Dependencies

Từ phiên bản mới, ServiceProvider hỗ trợ quản lý dependencies giữa các provider:

```go
type DatabaseProvider struct{}

func (p *DatabaseProvider) Register(app di.Application) {
    // Đăng ký database connection
    app.Singleton("db", func(c di.Container) interface{} {
        return &Database{}
    })
}

func (p *DatabaseProvider) Boot(app di.Application) {
    // Khởi tạo database
}

func (p *DatabaseProvider) Requires() []string {
    // Provider này cần config provider được đăng ký trước
    return []string{"config"}
}

func (p *DatabaseProvider) Providers() []string {
    // Provider này cung cấp các service
    return []string{"db", "db.migrator"}
}

// Sử dụng RegisterWithDependencies để tự động sắp xếp thứ tự
app.RegisterWithDependencies() // Tự động đăng ký theo đúng thứ tự dependency
```

### Mocks cho Testing

Package này cung cấp sẵn các mock cho các interface chính (trong thư mục `mocks/`):
- Application (bao gồm RegisterWithDependencies)
- ServiceProvider (bao gồm Requires, Providers)
- ServiceProviderDeferred
- ModuleLoaderContract

```go
import (
    "testing"
    "go.fork.vn/di/mocks"
    "github.com/stretchr/testify/mock"
)

func TestMyService(t *testing.T) {
    // Tạo mock cho Application
    mockApp := new(mocks.Application)
    
    // Thiết lập mock behavior
    mockApp.On("Make", "config").Return("test-config", nil)
    mockApp.On("RegisterWithDependencies").Return(nil)
    
    // Mock cho ServiceProvider với dependencies
    mockProvider := new(mocks.ServiceProvider)
    mockProvider.On("Requires").Return([]string{"config"})
    mockProvider.On("Providers").Return([]string{"database", "cache"})
    
    // Sử dụng mock trong test
    service := NewMyService(mockApp)
    result := service.DoSomething()
    
    // Kiểm tra mock được gọi đúng cách
    mockApp.AssertExpectations(t)
    mockProvider.AssertExpectations(t)
}
```

## Tài liệu

- [Tổng quan](docs/overview.md)
- [Container](docs/container.md)
- [Service Providers](docs/provider.md)
- [Application](docs/application.md)
- [Testing](docs/testing.md)
- [Migration Guide](releases/next/MIGRATION.md)
- [Release Notes](releases/next/RELEASE_NOTES.md)
- [Release Process](releases/next/RELEASE_SUMMARY.md)

Xem [doc.go](./doc.go) để biết thêm chi tiết về API và [mocks/](./mocks/) cho testing.

## Lịch sử phiên bản

Tài liệu phát triển hiện tại được lưu trong [releases/next/](releases/next/).
Tài liệu của các phiên bản đã release được lưu trữ trong thư mục [releases/](releases/).

## Đóng góp

Xem [releases/README.md](releases/README.md) để hiểu workflow quản lý release và documentation.

## Giấy phép
MIT
