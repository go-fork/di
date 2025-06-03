# Container - Tài liệu kỹ thuật

## Tổng quan

`Container` là thành phần cốt lõi ("trái tim") của hệ thống Dependency Injection (DI) trong Fork framework. Đây là hiện thực chuẩn DI hiện đại, lấy cảm hứng từ các framework như Laravel, Spring, nhưng được tối ưu hóa đặc biệt cho Go.

## Định nghĩa & Tiêu chuẩn

### Nguyên tắc thiết kế

Container tuân thủ nghiêm ngặt các nguyên tắc sau:

- **SOLID Principles**: Đặc biệt là Dependency Inversion Principle
- **Separation of Concerns**: Tách biệt logic khởi tạo, quản lý vòng đời, và resolve dependency
- **Design Patterns**: Hỗ trợ Service-Repository pattern, Adapter pattern, Service Provider pattern
- **Extensibility**: Cho phép mở rộng, kiểm soát, testability và maintainability tối đa

### Vai trò trong hệ thống

- **Registry trung tâm**: Quản lý toàn bộ dependency, service, adapter, repository
- **Lifecycle Management**: Quản lý binding, singleton, instance, alias
- **Automatic Resolution**: Tự động resolve dependency qua reflection
- **Testability**: Đảm bảo mọi thành phần có thể được inject, mock, hoặc thay thế dễ dàng
- **Foundation**: Là nền tảng cho mọi module, provider, middleware trong Fork

## Cấu trúc dữ liệu

### Container Interface

Từ version v0.1.2, `Container` là một interface thay vì struct cụ thể:

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

### Container Implementation (internal)

Hiện thực nội bộ của Container interface:

```go
type container struct {
    bindings  map[string]BindingFunc  // Factory functions cho dependency
    instances map[string]interface{}  // Singleton instances đã khởi tạo
    aliases   map[string]string       // Alias mapping
    mu        sync.RWMutex           // Concurrent safety
}
```

### BindingFunc Type

```go
type BindingFunc func(c Container) interface{}
```

`BindingFunc` là factory function được sử dụng để tạo ra instances của dependencies trong container. Đây là thành phần cốt lõi cho cơ chế dependency injection.

**Đặc điểm:**
- **Input**: Nhận container instance để có thể resolve dependencies khác
- **Output**: Trả về instance của dependency (dạng `interface{}`)
- **Factory Pattern**: Thực hiện lazy initialization - chỉ tạo instance khi cần thiết
- **Dependency Resolution**: Có thể resolve dependencies khác từ container trong quá trình khởi tạo

### Chi tiết các thành phần

#### Container Fields

| Trường | Loại | Mô tả |
|--------|------|-------|
| `bindings` | `map[string]BindingFunc` | Ánh xạ abstract type (tên logic) tới factory function khởi tạo instance |
| `instances` | `map[string]interface{}` | Lưu trữ các singleton instance đã được khởi tạo |
| `aliases` | `map[string]string` | Ánh xạ alias tới abstract type gốc, hỗ trợ truy cập đa tên |
| `mu` | `sync.RWMutex` | Đảm bảo an toàn concurrent cho mọi thao tác đăng ký/resolve |

#### BindingFunc Characteristics

| Thuộc tính | Mô tả | Ví dụ |
|------------|-------|-------|
| **Signature** | `func(c Container) interface{}` | Factory function chuẩn |
| **Container Access** | Có thể resolve dependencies khác từ container | `c.MustMake("config")` |
| **Lazy Loading** | Chỉ thực thi khi dependency được resolve | Tiết kiệm resources |
| **Flexible Return** | Trả về `interface{}` - có thể là bất kỳ type nào | Struct, interface, primitive types |
| **Error Handling** | Có thể panic hoặc return nil nếu lỗi | Tùy thuộc implementation |

#### BindingFunc Usage Patterns

```go
// 1. Simple factory - tạo instance mới
container.Bind("logger", func(c di.Container) interface{} {
    return &Logger{Level: "info"}
})

// 2. Dependency injection - resolve dependencies khác
container.Bind("userService", func(c di.Container) interface{} {
    repo := c.MustMake("userRepository").(UserRepository)
    logger := c.MustMake("logger").(Logger)
    
    return &UserService{
        Repository: repo,
        Logger:     logger,
    }
})

// 3. Configuration-based - sử dụng config để khởi tạo
container.Bind("database", func(c di.Container) interface{} {
    config := c.MustMake("config").(config.Config)
    return database.Connect(config.DatabaseURL)
})

// 4. Conditional initialization - khởi tạo theo điều kiện
container.Bind("cache", func(c di.Container) interface{} {
    config := c.MustMake("config").(config.Config)
    if config.CacheEnabled {
        return cache.NewRedisCache(config.RedisURL)
    }
    return cache.NewMemoryCache()
})

// 5. Complex initialization - khởi tạo phức tạp với multiple dependencies
container.Bind("emailService", func(c di.Container) interface{} {
    config := c.MustMake("config").(config.Config)
    logger := c.MustMake("logger").(Logger)
    templateEngine := c.MustMake("templateEngine").(TemplateEngine)
    
    service := &EmailService{
        SMTPHost:       config.SMTPHost,
        SMTPPort:       config.SMTPPort,
        Logger:         logger,
        TemplateEngine: templateEngine,
    }
    
    // Additional setup
    service.LoadTemplates()
    
    return service
})
```

## API Reference

### Constructor

#### `New() Container`

Khởi tạo một DI container rỗng, sẵn sàng cho việc đăng ký binding, instance, alias.

**Trả về:**
- `Container`: Container interface cho instance mới được khởi tạo

**Ví dụ:**
```go
container := di.New()
```

### Registration Methods

#### `Bind(abstract string, concrete BindingFunc)`

Đăng ký một binding (factory function) cho abstract type.

**Tham số:**
- `abstract`: Tên logic của dependency (thường là interface hoặc service name)
- `concrete`: Factory function nhận container, trả về instance

**Mục đích:**
- Cho phép đăng ký cách khởi tạo một dependency động
- Phục vụ cho việc resolve về sau
- Override nếu đã tồn tại

**Ví dụ:**
```go
container.Bind("logger", func(c *di.Container) interface{} {
    return &Logger{Level: "info"}
})

container.Bind("database", func(c *di.Container) interface{} {
    config := c.MustMake("config").(config.Config)
    return database.Connect(config.DSN)
})
```

#### `BindIf(abstract string, concrete BindingFunc) bool`

Đăng ký binding chỉ khi chưa tồn tại.

**Tham số:**
- `abstract`: Tên logic của dependency
- `concrete`: Factory function

**Trả về:**
- `true`: Nếu đăng ký thành công
- `false`: Nếu đã tồn tại

**Mục đích:**
- Đảm bảo không override binding đã có
- Dùng cho module mở rộng

**Ví dụ:**
```go
// Đăng ký default logger chỉ khi chưa có
success := container.BindIf("logger", func(c di.Container) interface{} {
    return &DefaultLogger{}
})
```

#### `Singleton(abstract string, concrete BindingFunc)`

Đăng ký binding singleton (chỉ tạo một instance duy nhất).

**Tham số:**
- `abstract`: Tên logic của dependency
- `concrete`: Factory function

**Mục đích:**
- Đảm bảo dependency chỉ được khởi tạo một lần duy nhất
- Factory function được wrap lại, lưu instance vào map instances khi lần đầu resolve

**Ví dụ:**
```go
container.Singleton("database", func(c di.Container) interface{} {
    return database.NewConnection("localhost:5432")
})

// Tất cả lần gọi sau sẽ trả về cùng một instance
db1 := container.MustMake("database")
db2 := container.MustMake("database")
// db1 == db2 (cùng một instance)
```

#### `Instance(abstract string, instance interface{})`

Đăng ký một instance đã khởi tạo sẵn.

**Tham số:**
- `abstract`: Tên logic
- `instance`: Giá trị đã khởi tạo

**Mục đích:**
- Cho phép inject các giá trị đã tồn tại (config, logger, ...)
- Không cần factory function

**Ví dụ:**
```go
config := &Config{Debug: true, Port: 8080}
container.Instance("config", config)

logger := log.New(os.Stdout, "APP: ", log.LstdFlags)
container.Instance("logger", logger)
```

#### `Alias(abstract, alias string)`

Đăng ký một alias cho abstract type.

**Tham số:**
- `abstract`: Tên gốc
- `alias`: Tên alias

**Mục đích:**
- Cho phép truy cập dependency qua nhiều tên khác nhau
- Resolve alias sẽ trả về instance của abstract gốc

**Ví dụ:**
```go
container.Singleton("logger", func(c di.Container) interface{} {
    return &Logger{}
})
container.Alias("logger", "log")

// Có thể truy cập bằng cả hai tên
logger1 := container.MustMake("logger")
logger2 := container.MustMake("log")
// logger1 == logger2
```

### Resolution Methods

#### `Make(abstract string) (interface{}, error)`

Resolve một dependency từ container.

**Tham số:**
- `abstract`: Tên logic của dependency

**Trả về:**
- `interface{}`: Instance đã resolve
- `error`: Nếu không tìm thấy hoặc binding lỗi

**Logic:**
1. Ưu tiên alias (nếu có)
2. Kiểm tra instance đã có
3. Resolve từ binding
4. Trả về error nếu không tìm thấy

**Ví dụ:**
```go
logger, err := container.Make("logger")
if err != nil {
    log.Fatal("Cannot resolve logger:", err)
}
```

#### `MustMake(abstract string) interface{}`

Resolve một dependency, panic nếu lỗi.

**Tham số:**
- `abstract`: Tên logic của dependency

**Trả về:**
- `interface{}`: Instance đã resolve

**Mục đích:**
- Resolve instance, panic nếu không tìm thấy
- Dùng cho critical dependency

**Ví dụ:**
```go
logger := container.MustMake("logger").(*Logger)
database := container.MustMake("database").(*Database)
```

### Utility Methods

#### `Bound(abstract string) bool`

Kiểm tra một abstract đã được đăng ký binding/instance/alias chưa.

**Tham số:**
- `abstract`: Tên logic cần kiểm tra

**Trả về:**
- `true`: Nếu đã đăng ký
- `false`: Nếu chưa đăng ký

**Mục đích:**
- Hỗ trợ kiểm tra trạng thái container
- Phục vụ cho module động

**Ví dụ:**
```go
if container.Bound("logger") {
    logger := container.MustMake("logger")
    // Sử dụng logger
}

// Đăng ký conditional
if !container.Bound("cache") {
    container.Singleton("cache", func(c di.Container) interface{} {
        return cache.NewMemoryCache()
    })
}
```

#### `Reset()`

Xóa toàn bộ binding, instance, alias khỏi container.

**Mục đích:**
- Làm sạch container
- Thường dùng cho test hoặc reload

**Ví dụ:**
```go
func TestSomething(t *testing.T) {
    container := di.New()
    
    // Setup test dependencies
    container.Instance("config", testConfig)
    
    // Run test
    // ...
    
    // Cleanup
    container.Reset()
}
```

### Advanced Methods

#### `Call(callback interface{}, additionalParams ...interface{}) ([]interface{}, error)`

Gọi một hàm và tự động resolve các dependency qua reflection.

**Tham số:**
- `callback`: Function cần gọi
- `additionalParams`: Các tham số bổ sung (ưu tiên inject)

**Trả về:**
- `[]interface{}`: Kết quả trả về của callback
- `error`: Nếu không resolve được tham số hoặc callback không hợp lệ

**Logic:**
1. Phân tích các tham số của callback
2. Resolve từ container hoặc lấy từ additionalParams
3. Gọi callback với các tham số đã resolve

**Ví dụ:**
```go
// Tự động inject dependencies
container.Call(func(logger *Logger, db *Database) error {
    logger.Info("Starting application")
    return db.Ping()
})

// Với additional parameters
container.Call(func(logger *Logger, userID string) {
    logger.Info("User logged in:", userID)
}, "user123")
```

## Concurrent Safety

Container được thiết kế để an toàn với concurrent access:

- **Read-Write Mutex**: Sử dụng `sync.RWMutex` để bảo vệ các thao tác
- **Read Operations**: `Make`, `MustMake`, `Bound` sử dụng read lock
- **Write Operations**: `Bind`, `Instance`, `Alias`, `Reset` sử dụng write lock
- **Singleton Resolution**: Có cơ chế đặc biệt để tránh race condition khi khởi tạo singleton

```go
// An toàn khi sử dụng từ nhiều goroutine
go func() {
    service := container.MustMake("service")
    // Sử dụng service
}()

go func() {
    anotherService := container.MustMake("another-service")
    // Sử dụng anotherService
}()
```

## Best Practices

### 1. Naming Convention

```go
// Sử dụng interface names cho abstract types
container.Bind("logger", func(c di.Container) interface{} {
    return &StdLogger{}
})

// Hoặc service names rõ ràng
container.Bind("user.repository", func(c di.Container) interface{} {
    return &UserRepository{}
})
```

### 2. Dependency Resolution

```go
// Resolve dependencies trong factory function
container.Singleton("user.service", func(c di.Container) interface{} {
    repo := c.MustMake("user.repository").(UserRepository)
    logger := c.MustMake("logger").(Logger)
    
    return &UserService{
        Repository: repo,
        Logger:     logger,
    }
})
```

### 3. Error Handling

```go
// Sử dụng Make cho non-critical dependencies
cache, err := container.Make("cache")
if err != nil {
    // Fallback to no-cache mode
    cache = &NoOpCache{}
}

// Sử dụng MustMake cho critical dependencies
database := container.MustMake("database").(*Database)
```

### 4. Testing

```go
func TestUserService(t *testing.T) {
    container := di.New()
    
    // Mock dependencies
    mockRepo := &MockUserRepository{}
    mockLogger := &MockLogger{}
    
    container.Instance("user.repository", mockRepo)
    container.Instance("logger", mockLogger)
    
    // Test target
    container.Singleton("user.service", func(c di.Container) interface{} {
        return &UserService{
            Repository: c.MustMake("user.repository").(UserRepository),
            Logger:     c.MustMake("logger").(Logger),
        }
    })
    
    service := container.MustMake("user.service").(*UserService)
    
    // Run tests
    // ...
}
```

## Lưu ý quan trọng

### 1. Type Assertion

Container trả về `interface{}`, cần type assertion khi sử dụng:

```go
// Cần type assertion
logger := container.MustMake("logger").(*Logger)

// Hoặc kiểm tra type
if logger, ok := container.Make("logger"); ok {
    if l, ok := logger.(*Logger); ok {
        l.Info("Hello")
    }
}
```

### 2. Circular Dependencies

Tránh circular dependencies trong factory functions:

```go
// ❌ Sai - circular dependency
container.Bind("service-a", func(c di.Container) interface{} {
    serviceB := c.MustMake("service-b")
    return &ServiceA{ServiceB: serviceB}
})

container.Bind("service-b", func(c di.Container) interface{} {
    serviceA := c.MustMake("service-a") // Deadlock!
    return &ServiceB{ServiceA: serviceA}
})
```

### 3. Memory Management

Singleton instances được giữ trong memory suốt vòng đời container:

```go
// Cân nhắc memory usage với large objects
container.Singleton("large-cache", func(c di.Container) interface{} {
    return cache.NewLargeCache(1000000) // Sẽ tồn tại suốt đời app
})
```

## Tích hợp với các thành phần khác

### Service Providers

```go
type DatabaseProvider struct{}

func (p *DatabaseProvider) Register(app Application) {
    container := app.Container()
    container.Singleton("database", func(c di.Container) interface{} {
        config := c.MustMake("config").(config.Config)
        return database.Connect(config.DatabaseURL)
    })
}
```

### Application Integration

```go
type App struct {
    container Container
}

func (app *App) Container() Container {
    return app.container
}

func (app *App) Make(abstract string) (interface{}, error) {
    return app.container.Make(abstract)
}
```

## Tài liệu liên quan

- [binding.md](./binding.md) - Chi tiết về BindingFunc
- [provider.md](./provider.md) - Service Provider pattern
- [application.md](./application.md) - Application interface
- [examples/](../examples/) - Các ví dụ sử dụng thực tế

## Changelog

Xem [CHANGELOG.md](../CHANGELOG.md) để biết lịch sử thay đổi của Container API.
