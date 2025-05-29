# Application - Tài liệu kỹ thuật

## Tổng quan

`Application` interface định nghĩa contract cơ bản cho các ứng dụng sử dụng hệ thống Dependency Injection (DI) trong Fork framework. Đây là lớp trừu tượng cấp cao nhất, chuẩn hóa cách thức quản lý vòng đời ứng dụng, service providers, và dependency resolution.

## Định nghĩa & Vai trò

### Mục đích chính

Application interface phục vụ các mục đích sau:

- **Chuẩn hóa Contract**: Định nghĩa interface thống nhất cho các ứng dụng sử dụng DI
- **Quản lý Service Provider**: Đăng ký, boot, và quản lý vòng đời các service provider
- **Dependency Management**: Hỗ trợ binding, singleton, instance, alias, resolve, call
- **Extensibility**: Đảm bảo khả năng mở rộng và kiểm soát cấu hình các thành phần
- **Lifecycle Control**: Quản lý vòng đời và khởi tạo các thành phần của ứng dụng

### Kiến trúc và Pattern

Application interface áp dụng các pattern sau:

- **Facade Pattern**: Cung cấp interface đơn giản cho hệ thống phức tạp
- **Service Locator Pattern**: Trung tâm hóa việc locate và resolve services
- **Provider Pattern**: Quản lý các service provider theo mô-đun
- **Dependency Injection**: Tự động inject dependencies qua reflection

## Tính năng cốt lõi

### 1. Container Access
- Truy cập trực tiếp vào DI container
- Quản lý central registry của dependencies

### 2. Provider Management
- Đăng ký service providers
- Boot providers theo thứ tự phụ thuộc
- Quản lý dependencies giữa các providers

### 3. Dependency Registration
- Binding với factory functions
- Singleton instances
- Pre-built instances
- Alias mapping

### 4. Dependency Resolution
- Safe resolution với error handling
- Must-resolve với panic semantics
- Automatic function injection

## API Reference

### Container Access

#### `Container() *Container`

Trả về instance DI container của ứng dụng.

**Mục đích:**
- Cho phép truy cập trực tiếp vào dependency injection container
- Hỗ trợ đăng ký hoặc resolve các dependency

**Trả về:**
- `*Container`: Instance container hiện tại của ứng dụng

**Ví dụ:**
```go
app := NewApp()
container := app.Container()

// Truy cập trực tiếp container
container.Bind("service", factory)
service := container.MustMake("service")
```

### Service Provider Management

#### `RegisterServiceProviders() error`

Đăng ký tất cả các service provider đã cấu hình cho ứng dụng.

**Mục đích:**
- Đảm bảo các service provider được đăng ký vào container trước khi boot
- Thiết lập toàn bộ dependencies cần thiết

**Trả về:**
- `error`: Lỗi nếu provider không hợp lệ hoặc đăng ký thất bại

**Exceptions:**
- Error nếu provider nil, duplicate hoặc có lỗi khi đăng ký

**Ví dụ:**
```go
app := NewApp()
app.AddProvider(&DatabaseProvider{})
app.AddProvider(&CacheProvider{})

err := app.RegisterServiceProviders()
if err != nil {
    log.Fatal("Failed to register providers:", err)
}
```

#### `RegisterWithDependencies() error`

Đăng ký và sắp xếp providers theo thứ tự phụ thuộc.

**Mục đích:**
- Tự động sắp xếp và đăng ký providers theo thứ tự phụ thuộc
- Phát hiện circular dependencies
- Đảm bảo dependencies được đăng ký trước dependent providers

**Trả về:**
- `error`: Lỗi nếu có circular dependency hoặc không tìm thấy provider yêu cầu

**Logic thực thi:**
1. Phân tích dependency graph của providers
2. Topological sort để xác định thứ tự đăng ký
3. Phát hiện circular dependencies
4. Đăng ký providers theo thứ tự đã sắp xếp

**Ví dụ:**
```go
app := NewApp()

// Provider A yêu cầu config
app.AddProvider(&DatabaseProvider{}) // Requires: ["config"]

// Provider B yêu cầu database
app.AddProvider(&UserServiceProvider{}) // Requires: ["database"]

// Config provider không yêu cầu gì
app.AddProvider(&ConfigProvider{}) // Requires: []

// Tự động sắp xếp: Config -> Database -> UserService
err := app.RegisterWithDependencies()
if err != nil {
    log.Fatal("Dependency resolution failed:", err)
}
```

#### `BootServiceProviders() error`

Boot tất cả các service provider đã đăng ký.

**Mục đích:**
- Khởi tạo các tài nguyên, kết nối hoặc logic phụ thuộc vào provider
- Thực hiện post-registration initialization

**Trả về:**
- `error`: Lỗi nếu boot provider thất bại

**Exceptions:**
- Error nếu provider boot lỗi hoặc thiếu dependency

**Ví dụ:**
```go
// Sau khi đăng ký providers
err := app.BootServiceProviders()
if err != nil {
    log.Fatal("Failed to boot providers:", err)
}
```

#### `Register(provider ServiceProvider)`

Đăng ký một service provider vào ứng dụng.

**Tham số:**
- `provider`: ServiceProvider — provider cần đăng ký

**Mục đích:**
- Cho phép thêm động các service provider vào container
- Hỗ trợ modular architecture

**Exceptions:**
- Panic nếu provider nil hoặc không hợp lệ

**Ví dụ:**
```go
app := NewApp()

// Đăng ký provider
app.Register(&DatabaseProvider{})
app.Register(&CacheProvider{})

// Custom provider
app.Register(&CustomServiceProvider{
    ConfigPath: "/etc/app/config.yaml",
})
```

#### `Boot() error`

Khởi động tất cả các service provider đã đăng ký.

**Mục đích:**
- Thực thi logic khởi tạo của provider (Boot) sau khi đã đăng ký
- Shorthand cho BootServiceProviders()

**Trả về:**
- `error`: Lỗi nếu boot provider thất bại

**Exceptions:**
- Error nếu provider boot lỗi hoặc thiếu dependency

**Ví dụ:**
```go
app := NewApp()
app.Register(&DatabaseProvider{})

err := app.Boot()
if err != nil {
    log.Fatal("Application boot failed:", err)
}
```

### Dependency Registration

#### `Bind(abstract string, concrete BindingFunc)`

Đăng ký một binding với container.

**Tham số:**
- `abstract`: Tên abstract type
- `concrete`: Factory function tạo instance

**Mục đích:**
- Đăng ký factory function cho abstract type
- Cho phép resolve động các dependency

**Exceptions:**
- Panic nếu abstract rỗng hoặc đã tồn tại binding

**Ví dụ:**
```go
app.Bind("logger", func(c *di.Container) interface{} {
    return log.New(os.Stdout, "APP: ", log.LstdFlags)
})

app.Bind("mailer", func(c *di.Container) interface{} {
    config := c.MustMake("config").(*Config)
    return smtp.NewMailer(config.SMTPHost, config.SMTPPort)
})
```

#### `Singleton(abstract string, concrete BindingFunc)`

Đăng ký một singleton binding với container.

**Tham số:**
- `abstract`: Tên abstract type
- `concrete`: Factory function tạo instance

**Mục đích:**
- Đảm bảo chỉ tạo một instance duy nhất cho abstract type
- Trong suốt vòng đời ứng dụng

**Ví dụ:**
```go
// Database connection singleton
app.Singleton("database", func(c *di.Container) interface{} {
    config := c.MustMake("config").(*Config)
    return database.Connect(config.DatabaseURL)
})

// Cache service singleton
app.Singleton("cache", func(c *di.Container) interface{} {
    return redis.NewClient(&redis.Options{
        Addr: "localhost:6379",
    })
})
```

#### `Instance(abstract string, instance interface{})`

Đăng ký một instance đã khởi tạo sẵn vào container.

**Tham số:**
- `abstract`: Tên abstract type
- `instance`: Instance đã khởi tạo

**Mục đích:**
- Cho phép inject các instance đã tồn tại
- Thường dùng cho config, logger, pre-built objects

**Ví dụ:**
```go
// Đăng ký config object
config := &Config{
    Debug:       true,
    DatabaseURL: "postgres://user:pass@localhost/db",
    Port:        8080,
}
app.Instance("config", config)

// Đăng ký logger instance
logger := logrus.New()
logger.SetLevel(logrus.InfoLevel)
app.Instance("logger", logger)
```

#### `Alias(abstract, alias string)`

Đăng ký một alias cho abstract type.

**Tham số:**
- `abstract`: Tên abstract gốc
- `alias`: Tên alias

**Mục đích:**
- Cho phép truy cập dependency qua nhiều tên khác nhau
- Hỗ trợ backward compatibility

**Ví dụ:**
```go
app.Singleton("logger", loggerFactory)

// Tạo aliases
app.Alias("logger", "log")
app.Alias("logger", "app.logger")

// Tất cả đều trả về cùng instance
logger1 := app.MustMake("logger")
logger2 := app.MustMake("log")
logger3 := app.MustMake("app.logger")
```

### Dependency Resolution

#### `Make(abstract string) (interface{}, error)`

Resolve một dependency từ container.

**Tham số:**
- `abstract`: Tên abstract type

**Trả về:**
- `interface{}`: Instance đã resolve
- `error`: Lỗi nếu không tìm thấy hoặc binding lỗi

**Mục đích:**
- Safe resolution với error handling
- Cho phép graceful degradation

**Ví dụ:**
```go
// Safe resolution với error handling
logger, err := app.Make("logger")
if err != nil {
    // Fallback to default logger
    logger = log.New(os.Stdout, "", log.LstdFlags)
}

// Type assertion
if l, ok := logger.(*logrus.Logger); ok {
    l.Info("Logger resolved successfully")
}
```

#### `MustMake(abstract string) interface{}`

Resolve một dependency từ container, panic nếu lỗi.

**Tham số:**
- `abstract`: Tên abstract type

**Trả về:**
- `interface{}`: Instance đã resolve

**Mục đích:**
- Resolve instance, panic nếu không tìm thấy
- Dùng cho critical dependencies

**Exceptions:**
- Panic nếu không resolve được dependency

**Ví dụ:**
```go
// Critical dependencies
database := app.MustMake("database").(*sql.DB)
config := app.MustMake("config").(*Config)

// Sử dụng trong initialization
func initializeServices(app Application) {
    userRepo := &UserRepository{
        DB:     app.MustMake("database").(*sql.DB),
        Logger: app.MustMake("logger").(*logrus.Logger),
    }
    
    app.Instance("user.repository", userRepo)
}
```

#### `Call(callback interface{}, additionalParams ...interface{}) ([]interface{}, error)`

Gọi một hàm và tự động resolve các dependency.

**Tham số:**
- `callback`: Function cần gọi
- `additionalParams`: Các tham số bổ sung (ưu tiên inject)

**Trả về:**
- `[]interface{}`: Kết quả trả về của callback
- `error`: Lỗi nếu không resolve được tham số hoặc callback không hợp lệ

**Mục đích:**
- Tự động inject các dependency vào callback function qua reflection
- Hỗ trợ functional programming style

**Ví dụ:**
```go
// Tự động inject dependencies
results, err := app.Call(func(db *sql.DB, logger *logrus.Logger) error {
    logger.Info("Executing database migration")
    return runMigration(db)
})

// Với additional parameters
results, err := app.Call(func(
    userRepo *UserRepository, 
    userID string,
    action string,
) (*User, error) {
    return userRepo.PerformAction(userID, action)
}, "user123", "activate")

// Handler function injection
app.Call(func(w http.ResponseWriter, r *http.Request, userService *UserService) {
    users, err := userService.GetAllUsers()
    if err != nil {
        http.Error(w, err.Error(), 500)
        return
    }
    json.NewEncoder(w).Encode(users)
})
```

## Patterns and Best Practices

### 1. Application Structure

```go
type App struct {
    container *di.Container
    providers []ServiceProvider
    config    *Config
}

func NewApp(config *Config) *App {
    app := &App{
        container: di.New(),
        providers: make([]ServiceProvider, 0),
        config:    config,
    }
    
    // Register core instances
    app.Instance("config", config)
    app.Instance("app", app)
    
    return app
}

func (app *App) Container() *di.Container {
    return app.container
}

// Implement all Application interface methods...
```

### 2. Provider Registration Pattern

```go
func (app *App) ConfigureProviders() {
    // Core providers
    app.Register(&ConfigProvider{})
    app.Register(&LoggerProvider{})
    
    // Infrastructure providers
    app.Register(&DatabaseProvider{})
    app.Register(&CacheProvider{})
    app.Register(&QueueProvider{})
    
    // Domain providers
    app.Register(&UserServiceProvider{})
    app.Register(&AuthServiceProvider{})
    
    // Web providers (if web app)
    app.Register(&RouterProvider{})
    app.Register(&MiddlewareProvider{})
}
```

### 3. Dependency-Aware Registration

```go
func (app *App) Start() error {
    // Configure all providers
    app.ConfigureProviders()
    
    // Register with dependency resolution
    if err := app.RegisterWithDependencies(); err != nil {
        return fmt.Errorf("failed to register providers: %w", err)
    }
    
    // Boot all providers
    if err := app.Boot(); err != nil {
        return fmt.Errorf("failed to boot application: %w", err)
    }
    
    return nil
}
```

### 4. Graceful Shutdown Pattern

```go
type GracefulApp struct {
    Application
    shutdownHandlers []func() error
}

func (app *GracefulApp) RegisterShutdownHandler(handler func() error) {
    app.shutdownHandlers = append(app.shutdownHandlers, handler)
}

func (app *GracefulApp) Shutdown() error {
    var errors []error
    
    // Execute shutdown handlers in reverse order
    for i := len(app.shutdownHandlers) - 1; i >= 0; i-- {
        if err := app.shutdownHandlers[i](); err != nil {
            errors = append(errors, err)
        }
    }
    
    if len(errors) > 0 {
        return fmt.Errorf("shutdown errors: %v", errors)
    }
    
    return nil
}
```

## Testing Strategies

### 1. Mock Application

```go
func TestUserService(t *testing.T) {
    // Create mock application
    mockApp := new(mocks.Application)
    
    // Setup expectations
    mockDB := &MockDatabase{}
    mockLogger := &MockLogger{}
    
    mockApp.On("Make", "database").Return(mockDB, nil)
    mockApp.On("Make", "logger").Return(mockLogger, nil)
    
    // Test service with mocked dependencies
    service := NewUserService(mockApp)
    user, err := service.CreateUser("john@example.com")
    
    assert.NoError(t, err)
    assert.NotNil(t, user)
    mockApp.AssertExpectations(t)
}
```

### 2. Integration Testing

```go
func TestApplicationIntegration(t *testing.T) {
    config := &Config{
        DatabaseURL: "sqlite://test.db",
        Debug:       true,
    }
    
    app := NewApp(config)
    app.Register(&TestDatabaseProvider{})
    app.Register(&UserServiceProvider{})
    
    err := app.Start()
    require.NoError(t, err)
    
    // Test that services are properly registered
    userService := app.MustMake("user.service")
    assert.NotNil(t, userService)
    
    // Cleanup
    defer app.Shutdown()
}
```

### 3. Provider Testing

```go
func TestDatabaseProvider(t *testing.T) {
    app := NewApp(&Config{DatabaseURL: "sqlite://test.db"})
    provider := &DatabaseProvider{}
    
    app.Register(provider)
    err := app.Boot()
    require.NoError(t, err)
    
    // Verify database is registered
    db := app.MustMake("database")
    assert.NotNil(t, db)
    
    // Verify connection works
    _, err = app.Call(func(db *sql.DB) error {
        return db.Ping()
    })
    assert.NoError(t, err)
}
```

## Error Handling

### 1. Registration Errors

```go
func (app *App) RegisterServiceProviders() error {
    for _, provider := range app.providers {
        if provider == nil {
            return fmt.Errorf("nil provider encountered")
        }
        
        func() {
            defer func() {
                if r := recover(); r != nil {
                    err = fmt.Errorf("provider registration panic: %v", r)
                }
            }()
            
            provider.Register(app)
        }()
        
        if err != nil {
            return fmt.Errorf("failed to register provider %T: %w", provider, err)
        }
    }
    
    return nil
}
```

### 2. Boot Errors

```go
func (app *App) BootServiceProviders() error {
    for _, provider := range app.providers {
        if err := provider.Boot(app); err != nil {
            return fmt.Errorf("failed to boot provider %T: %w", provider, err)
        }
    }
    
    return nil
}
```

### 3. Dependency Resolution Errors

```go
func (app *App) Make(abstract string) (interface{}, error) {
    instance, err := app.container.Make(abstract)
    if err != nil {
        return nil, fmt.Errorf("failed to resolve %s: %w", abstract, err)
    }
    
    return instance, nil
}
```

## Performance Considerations

### 1. Lazy Loading

```go
// Implement lazy provider loading
func (app *App) RegisterDeferred(provider ServiceProviderDeferred) {
    app.deferredProviders = append(app.deferredProviders, provider)
}

func (app *App) bootDeferredProvider(name string) error {
    for _, provider := range app.deferredProviders {
        if provider.Name() == name {
            return provider.DeferredBoot(app)
        }
    }
    return fmt.Errorf("deferred provider %s not found", name)
}
```

### 2. Concurrent-Safe Operations

```go
type SafeApp struct {
    *App
    mu sync.RWMutex
}

func (app *SafeApp) Register(provider ServiceProvider) {
    app.mu.Lock()
    defer app.mu.Unlock()
    
    app.App.Register(provider)
}

func (app *SafeApp) Make(abstract string) (interface{}, error) {
    app.mu.RLock()
    defer app.mu.RUnlock()
    
    return app.App.Make(abstract)
}
```

## Migration and Compatibility

### 1. Version Compatibility

```go
type VersionedApplication interface {
    Application
    Version() string
    CompatibleWith(version string) bool
}

func (app *App) CompatibleWith(version string) bool {
    // Implement semantic version compatibility
    return semver.Compatible(app.Version(), version)
}
```

### 2. Feature Flags

```go
func (app *App) HasFeature(feature string) bool {
    if features, ok := app.MustMake("features").(FeatureFlags); ok {
        return features.IsEnabled(feature)
    }
    return false
}

func (app *App) ConditionalRegister(feature string, provider ServiceProvider) {
    if app.HasFeature(feature) {
        app.Register(provider)
    }
}
```

## Tài liệu liên quan

- [container.md](./container.md) - Chi tiết về DI Container
- [provider.md](./provider.md) - Service Provider pattern
- [binding.md](./binding.md) - Binding functions
- [deferred.md](./deferred.md) - Deferred service providers
- [examples/](../examples/) - Các ví dụ ứng dụng thực tế

## Changelog

Xem [CHANGELOG.md](../CHANGELOG.md) để biết lịch sử thay đổi của Application interface.
