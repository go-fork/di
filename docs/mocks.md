# Mocks - Tài liệu Kỹ thuật

## Tổng quan

Tài liệu này mô tả việc sử dụng và tạo lại mock objects cho các interface chính của gói `go.fork.vn/di`. Mock objects này được thiết kế để hỗ trợ viết tests cho các components sử dụng DI container và các interface liên quan.

## Mock Objects

### 1. Container Mock

`Container` mock là hiện thực giả lập của `Container` interface, cho phép bạn kiểm soát hoàn toàn hành vi của DI container trong tests.

```go
import (
    "testing"
    "go.fork.vn/di/mocks"
    "github.com/stretchr/testify/mock"
)

func TestWithContainerMock(t *testing.T) {
    mockContainer := new(mocks.Container)
    
    // Cấu hình hành vi cho mock
    mockContainer.On("Bound", "config").Return(true)
    mockContainer.On("MustMake", "config").Return(&Config{Debug: true})
    mockContainer.On("Make", "logger").Return(&Logger{}, nil)
    
    // Cấu hình cho binding function
    mockContainer.On("Bind", mock.AnythingOfType("string"), mock.AnythingOfType("di.BindingFunc")).Return()
    
    // Sử dụng trong component cần test
    component := NewComponent(mockContainer)
    result := component.Process()
    
    // Kiểm tra kết quả
    assert.NotNil(t, result)
    
    // Verify expectations
    mockContainer.AssertExpectations(t)
}
```

### 2. Application Mock

`Application` mock cho phép test các components tương tác với application interface:

```go
func TestWithApplicationMock(t *testing.T) {
    mockApp := new(mocks.Application)
    mockContainer := new(mocks.Container)
    
    // Cấu hình Application
    mockApp.On("Container").Return(mockContainer)
    
    // Cấu hình Container
    mockContainer.On("MustMake", "config").Return(&Config{})
    
    // Sử dụng trong provider
    provider := NewMyProvider()
    provider.Register(mockApp)
    
    // Verify expectations
    mockApp.AssertExpectations(t)
    mockContainer.AssertExpectations(t)
}
```

### 3. ServiceProvider Mock

`ServiceProvider` mock cho phép test các components quản lý service providers:

```go
func TestWithProviderMock(t *testing.T) {
    mockProvider := new(mocks.ServiceProvider)
    
    // Cấu hình Provider
    mockProvider.On("Register", mock.AnythingOfType("*app.Application")).Return()
    mockProvider.On("Boot", mock.AnythingOfType("*app.Application")).Return()
    
    // Sử dụng trong application
    app := NewApplication()
    app.Register(mockProvider)
    app.Boot()
    
    // Verify expectations
    mockProvider.AssertExpectations(t)
}
```

### 4. ServiceProviderDeferred Mock

`ServiceProviderDeferred` mock cho phép test deferred operations:

```go
func TestWithDeferredProviderMock(t *testing.T) {
    mockDeferredProvider := new(mocks.ServiceProviderDeferred)
    
    // Cấu hình Provider
    mockDeferredProvider.On("Register", mock.AnythingOfType("*app.Application")).Return()
    mockDeferredProvider.On("Boot", mock.AnythingOfType("*app.Application")).Return()
    mockDeferredProvider.On("DeferredBoot", mock.AnythingOfType("*app.Application")).Return()
    
    // Sử dụng trong application
    app := NewApplication()
    app.Register(mockDeferredProvider)
    app.Boot()
    app.ExecuteDeferred()
    
    // Verify expectations
    mockDeferredProvider.AssertExpectations(t)
}
```

### 5. ModuleLoaderContract Mock

`ModuleLoaderContract` mock cho phép test module loading:

```go
func TestWithModuleLoaderMock(t *testing.T) {
    mockLoader := new(mocks.ModuleLoaderContract)
    
    // Cấu hình Loader
    mockLoader.On("LoadModules", mock.AnythingOfType("*app.Application"), mock.AnythingOfType("[]string")).Return()
    
    // Sử dụng trong application
    app := NewApplication()
    modules := []string{"module1", "module2"}
    mockLoader.LoadModules(app, modules)
    
    // Verify expectations
    mockLoader.AssertExpectations(t)
}
```

## Tạo lại Mock Objects

Các mock objects được tạo tự động bằng công cụ [mockery](https://github.com/vektra/mockery). Để tạo lại các mocks sau khi có thay đổi về interface, thực hiện các bước sau:

### Cài đặt Mockery

```bash
go install github.com/vektra/mockery/v2@latest
```

### Tạo lại Mocks

```bash
cd /path/to/di
mockery --name Container --output ./mocks --outpkg mocks
mockery --name Application --output ./mocks --outpkg mocks
mockery --name ServiceProvider --output ./mocks --outpkg mocks
mockery --name ServiceProviderDeferred --output ./mocks --outpkg mocks
mockery --name ModuleLoaderContract --output ./mocks --outpkg mocks
```

## Best Practices

### 1. Cấu hình trước khi sử dụng

Luôn cấu hình đầy đủ các expectations cho mock trước khi sử dụng:

```go
mockContainer.On("MustMake", "config").Return(&Config{})
```

### 2. Verify Expectations

Luôn verify rằng mọi expectation đều được gọi:

```go
mockContainer.AssertExpectations(t)
```

### 3. Specificity

Cấu hình expectations càng cụ thể càng tốt:

```go
// Tốt - cụ thể về tham số
mockContainer.On("MustMake", "config").Return(&Config{})

// Kém cụ thể hơn - chấp nhận bất kỳ string nào
mockContainer.On("MustMake", mock.AnythingOfType("string")).Return(&Config{})
```

### 4. Mock Chaining

Khi mock nhiều interfaces liên quan, kết nối chúng với nhau:

```go
mockApp.On("Container").Return(mockContainer)
mockContainer.On("MustMake", "logger").Return(mockLogger)
```

### 5. State Verification

Sử dụng custom arguments matchers khi cần thiết:

```go
mockContainer.On("Bind", mock.MatchedBy(func(s string) bool {
    return strings.HasPrefix(s, "service.")
}), mock.Anything).Return()
```

## Lưu ý

- Các mock objects là dynamic, có thể thay đổi hành vi trong quá trình test
- Mỗi mock method trả về mock.Call, cho phép thiết lập hành vi phức tạp hơn
- Có thể sử dụng `mock.Anything` khi không quan tâm đến tham số cụ thể
- Sử dụng `Times(n)` để chỉ định số lần một method được gọi

## Tham khảo

- [Mocks Source Code](../mocks/)
- [Testify Mock Documentation](https://pkg.go.dev/github.com/stretchr/testify/mock)
- [Container Interface Documentation](container.md)
- [Application Interface Documentation](application.md)
- [Testing Guide](testing.md)
