# 🚀 go.fork.vn/di v0.1.2 - Container Interface Update

Chúng tôi vui mừng giới thiệu **go.fork.vn/di v0.1.2** - phiên bản cập nhật với Container Interface và mocks mới!

## 🌟 Highlights

### 🔄 Container Interface
- **Interface-first Design**: Container giờ đây là interface thay vì struct cụ thể
- **Better Testability**: Dễ dàng mock cho unit testing
- **Extensibility**: Cho phép tạo nhiều implementation tùy chỉnh

### 🧪 Regenerated Mocks
- **Updated Mock Objects**: Tất cả mock objects được tạo lại để hỗ trợ Container interface mới
- **Testing Support**: Cải thiện khả năng testing với mockable interfaces
- **Interface Consistency**: Đảm bảo tất cả mocks tuân thủ các interface mới

### 📚 Updated Documentation
- **Interface-Based**: Cập nhật tài liệu để phản ánh Container interface mới
- **Testing Patterns**: Bổ sung các mẫu test với Container interface
- **Architecture Guide**: Cập nhật hướng dẫn kiến trúc với best practices mới

## 📦 Tính năng mới

### Container Interface

```go
// Container là interface của hệ thống Dependency Injection (DI) trong Fork framework.
type Container interface {
    // Bind đăng ký một binding (factory function) cho abstract type.
    Bind(abstract string, concrete BindingFunc)
    
    // BindIf đăng ký binding chỉ khi chưa tồn tại.
    BindIf(abstract string, concrete BindingFunc) bool
    
    // Singleton đăng ký binding singleton (chỉ tạo một instance duy nhất).
    Singleton(abstract string, concrete BindingFunc)
    
    // Instance đăng ký một instance đã khởi tạo sẵn.
    Instance(abstract string, instance interface{})
    
    // Alias đăng ký một alias cho abstract type.
    Alias(abstract, alias string)
    
    // Make resolve một dependency từ container.
    Make(abstract string) (interface{}, error)
    
    // MustMake resolve một dependency, panic nếu lỗi.
    MustMake(abstract string) interface{}
    
    // Bound kiểm tra một abstract đã được đăng ký binding/instance/alias chưa.
    Bound(abstract string) bool
    
    // Reset xóa toàn bộ binding, instance, alias khỏi container.
    Reset()
    
    // Call gọi một hàm và tự động resolve các dependency qua reflection.
    Call(callback interface{}, additionalParams ...interface{}) ([]interface{}, error)
}
```

### Tạo container mới

```go
// Vẫn sử dụng di.New() nhưng giờ đây trả về Container interface
container := di.New()

// Sử dụng các phương thức bình thường
container.Bind("service", func(c di.Container) interface{} {
    return &MyService{}
})

service := container.MustMake("service").(*MyService)
```

### Testing với Container interface mới

```go
import (
    "testing"
    "github.com/stretchr/testify/mock"
    "go.fork.vn/di"
    "go.fork.vn/di/mocks"
)

func TestService(t *testing.T) {
    // Tạo mock cho Container interface
    mockContainer := new(mocks.Container)
    
    // Thiết lập expectations
    mockContainer.On("Make", "config").Return(&Config{}, nil)
    mockContainer.On("MustMake", "logger").Return(&Logger{})
    
    // Sử dụng mock container trong tests
    service := NewService(mockContainer)
    
    // Test...
}
```

## 🚦 Khả năng tương thích

Phiên bản v0.1.2 không giới thiệu breaking changes cho public API. Tất cả code hiện tại sẽ tiếp tục hoạt động bình thường. Chỉ code truy cập trực tiếp vào kiểu `container` nội bộ mới cần điều chỉnh.

## 📋 Migration Guide

Chi tiết về cách migrate lên v0.1.2 có thể xem tại [MIGRATION_v0.1.2.md](MIGRATION_v0.1.2.md)

## 🧪 Testability Improvements

Việc chuyển sang Container interface mang lại nhiều lợi ích cho testing:

- **Mock Injections**: Dễ dàng inject mock container vào components
- **Flexible Testing**: Test với nhiều container implementations khác nhau  
- **Test Isolation**: Tách biệt dependency khi testing
- **Cleaner Tests**: Giảm sự phụ thuộc vào cấu trúc nội bộ của DI container

## 🔍 Advanced Use Cases

### Custom Container Implementations

```go
type LoggingContainer struct {
    di.Container
    logger Logger
}

func (lc *LoggingContainer) Make(abstract string) (interface{}, error) {
    lc.logger.Log("Resolving: " + abstract)
    return lc.Container.Make(abstract)
}

// Tạo container tùy chỉnh
baseContainer := di.New()
loggingContainer := &LoggingContainer{
    Container: baseContainer,
    logger: &MyLogger{},
}
```

### Middleware Pattern

```go
type ContainerMiddleware func(di.Container) di.Container

func MetricsMiddleware(next di.Container) di.Container {
    return &metricsContainer{Container: next}
}

func ValidationMiddleware(next di.Container) di.Container {
    return &validationContainer{Container: next}
}

// Áp dụng middleware
container := ValidationMiddleware(MetricsMiddleware(di.New()))
```

## 🛠️ Installation

```bash
go get go.fork.vn/di@v0.1.2
```

## 🔄 Migration Steps

1. Cập nhật package
```bash
go get -u go.fork.vn/di@v0.1.2
```

2. Cập nhật code sử dụng container struct trực tiếp (nếu có)
```go
// Trước
c := di.New().(*di.container) // Không còn hoạt động và không được khuyến nghị

// Sau
c := di.New() // Sử dụng Container interface
```

3. Cập nhật mocks (nếu đang sử dụng)
```go
mockContainer := new(mocks.Container) // Đã được cập nhật để implement Container interface mới
```

## 📊 Performance

Việc chuyển sang Container interface không ảnh hưởng đến hiệu suất vì Go compiler rất hiệu quả trong việc tối ưu hóa các lời gọi interface.

## 📚 Tài liệu

Tất cả tài liệu đã được cập nhật để phản ánh Container interface mới:

- **[Container](docs/container.md)** - Tài liệu Container interface cập nhật
- **[Testing](docs/testing.md)** - Mẫu testing với Container interface mới
- **[Application](docs/application.md)** - Container integration với Application
- **[Mocks](docs/mocks.md)** - Hướng dẫn sử dụng mocks mới

---

## 🎉 Cảm ơn!

Cảm ơn community đã support Fork framework! Phiên bản v0.1.2 là một bước tiến quan trọng trong việc cải thiện khả năng kiểm thử và tính mở rộng của framework.

**Happy Coding với go.fork.vn/di!** 🚀
