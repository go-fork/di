# ServiceProvider - Tài liệu Kỹ thuật

## Tổng quan

`ServiceProvider` là interface cốt lõi định nghĩa contract cho các service provider trong hệ thống Dependency Injection của Fork framework. Đây là thành phần quan trọng cho phép các module hoặc package đăng ký dịch vụ một cách mô-đun và tách biệt logic khởi tạo khỏi phần core của ứng dụng.

### Vai trò và Mục đích

- **Modular Service Registration**: Cho phép đăng ký các binding, singleton, instance vào container theo module
- **Separation of Concerns**: Tách biệt logic khởi tạo và cấu hình dịch vụ khỏi application logic
- **Dependency Management**: Quản lý dependencies giữa các service providers
- **Service Discovery**: Cung cấp metadata về các service được đăng ký
- **Bootstrap Coordination**: Điều phối quá trình khởi tạo các service theo đúng thứ tự

### Đặc điểm Chính

- **Idempotent Operations**: Register method có thể được gọi nhiều lần mà không gây lỗi
- **Two-Phase Initialization**: Tách biệt phase Register và Boot để đảm bảo dependency resolution
- **Dependency Declaration**: Khai báo rõ ràng dependencies và provided services
- **Framework Integration**: Tích hợp seamlessly với DI container và application lifecycle

## Architecture Pattern

ServiceProvider áp dụng các design pattern sau:

### 1. Service Locator Pattern
Provider đăng ký services vào container trung tâm để các component khác có thể locate và sử dụng.

### 2. Factory Pattern
Provider hoạt động như factory để tạo và cấu hình các service instance.

### 3. Dependency Injection Pattern
Provider inject dependencies cần thiết vào services được tạo.

### 4. Bootstrap Pattern
Two-phase initialization (Register → Boot) đảm bảo proper dependency resolution.

## API Reference

### Register

```go
Register(app Application)
```

Đăng ký các bindings vào DI container.

**Mô tả**:
- Phase đầu tiên của service provider lifecycle
- Đăng ký tất cả services mà provider cung cấp vào container
- Không nên thực hiện initialization logic phức tạp
- Phải idempotent (có thể gọi nhiều lần an toàn)
- Type-safe interface, không cần type assertion

**Parameters**:
- `app Application`: Application instance với typed interface cho DI container access

**Behavior**:
- Đăng ký bindings, singletons, instances vào container
- Không thực hiện actual service initialization
- Có thể panic nếu binding không hợp lệ
- Nên kiểm tra xem service đã được đăng ký chưa trước khi đăng ký lại

**Ví dụ Implementation**:
```go
type DatabaseProvider struct {}

func (p *DatabaseProvider) Register(app Application) {
    container := app.Container()
    
    // Đăng ký database connection factory
    container.Singleton("db.connection", func() interface{} {
        config := container.MustMake("config").(Config)
        return NewDatabaseConnection(config.Database)
    })
    
    // Đăng ký repository pattern
    container.Bind("user.repository", func() interface{} {
        db := container.MustMake("db.connection").(Database)
        return NewUserRepository(db)
    })
    
    // Đăng ký migration service
    container.Singleton("db.migrator", func() interface{} {
        db := container.MustMake("db.connection").(Database)
        return NewMigrator(db)
    })
}
```

### Boot

```go
Boot(app Application)
```

Được gọi sau khi tất cả các service provider đã được đăng ký.

**Mô tả**:
- Phase thứ hai của service provider lifecycle
- Thực hiện initialization logic phức tạp
- Có thể sử dụng các services đã được đăng ký bởi providers khác
- Thích hợp cho việc setup middleware, event listeners, scheduled tasks
- Type-safe interface, không cần type assertion

**Parameters**:
- `app Application`: Application instance với đầy đủ services đã được đăng ký

**Timing**:
- Được gọi sau khi tất cả providers đã hoàn thành Register phase
- Được gọi theo thứ tự dependency resolution
- Providers với dependencies được boot sau providers mà chúng depend on

**Use Cases**:
- Database migration và seeding
- Event listener registration
- Middleware setup
- Background job scheduling
- Cache warming
- External service connection establishment

**Ví dụ Implementation**:
```go
func (p *DatabaseProvider) Boot(app Application) {
    container := app.Container()
    
    // Run database migrations
    migrator := container.MustMake("db.migrator").(Migrator)
    if err := migrator.RunMigrations(); err != nil {
        panic(fmt.Sprintf("Database migration failed: %v", err))
    }
    
    // Setup database health checks
    db := container.MustMake("db.connection").(Database)
    healthChecker := container.MustMake("health.checker").(HealthChecker)
    healthChecker.Register("database", func() error {
        return db.Ping()
    })
    
    // Register database events
    eventBus := container.MustMake("event.bus").(EventBus)
    eventBus.Subscribe("user.created", func(event Event) {
        // Handle user creation event
        log.Printf("User created: %v", event.Data)
    })
}
```

### Requires

```go
Requires() []string
```

Trả về danh sách các provider mà provider này phụ thuộc vào.

**Mô tả**:
- Khai báo dependencies với các providers khác
- Framework sử dụng thông tin này để sắp xếp thứ tự boot
- Đảm bảo dependency providers được khởi tạo trước current provider

**Return Values**:
- `[]string`: Mảng tên các providers mà provider này yêu cầu

**Dependency Resolution**:
- Framework thực hiện topological sort dựa trên dependency graph
- Circular dependencies sẽ được detect và báo lỗi
- Providers không có dependencies được boot trước

**Best Practices**:
- Chỉ khai báo direct dependencies, không cần transitive dependencies
- Sử dụng consistent naming convention cho provider names
- Tránh over-specification dependencies không thực sự cần thiết

**Ví dụ Implementation**:
```go
type WebProvider struct {}

func (p *WebProvider) Requires() []string {
    return []string{
        "config",           // Cần configuration để setup routes
        "logger",           // Cần logging service
        "database",         // Cần database connection
        "cache",            // Cần cache service cho session
    }
}

type AuthProvider struct {}

func (p *AuthProvider) Requires() []string {
    return []string{
        "web",              // Cần web framework để register routes
        "user.repository",  // Cần user repository
        "session",          // Cần session management
    }
}
```

### Providers

```go
Providers() []string
```

Trả về danh sách các service mà provider này đăng ký.

**Mô tả**:
- Metadata về các services mà provider cung cấp
- Hỗ trợ debugging và service discovery
- Có thể được sử dụng để auto-generate documentation

**Return Values**:
- `[]string`: Mảng tên các services mà provider đăng ký vào container

**Use Cases**:
- Service discovery và introspection
- Debugging dependency issues
- Auto-generating API documentation
- Service health monitoring
- Container inspection tools

**Naming Conventions**:
- Sử dụng dot notation cho namespace: `"database.connection"`
- Group related services: `"cache.redis"`, `"cache.memory"`
- Use consistent prefixes: `"auth.user"`, `"auth.session"`

**Ví dụ Implementation**:
```go
type DatabaseProvider struct {}

func (p *DatabaseProvider) Providers() []string {
    return []string{
        "db.connection",    // Main database connection
        "db.migrator",      // Database migration service
        "user.repository",  // User data repository
        "post.repository",  // Post data repository
    }
}

type CacheProvider struct {}

func (p *CacheProvider) Providers() []string {
    return []string{
        "cache.redis",      // Redis cache implementation
        "cache.memory",     // In-memory cache implementation
        "cache.manager",    // Cache manager with fallback logic
    }
}

type WebProvider struct {}

func (p *WebProvider) Providers() []string {
    return []string{
        "web.app",          // Web application instance
        "web.router",       // HTTP router
        "web.middleware",   // Middleware stack
        "web.template",     // Template engine
    }
}
```

## Implementation Patterns

### 1. Basic Service Provider

```go
type ConfigProvider struct {
    configPath string
}

func NewConfigProvider(configPath string) *ConfigProvider {
    return &ConfigProvider{
        configPath: configPath,
    }
}

func (p *ConfigProvider) Register(app Application) {
    container := app.Container()
    
    // Register config as singleton
    container.Singleton("config", func() interface{} {
        config, err := LoadConfigFromFile(p.configPath)
        if err != nil {
            panic(fmt.Sprintf("Failed to load config: %v", err))
        }
        return config
    })
    
    // Register environment-specific config
    container.Bind("config.env", func() interface{} {
        config := container.MustMake("config").(Config)
        return config.Environment
    })
}

func (p *ConfigProvider) Boot(app Application) {
    // Validate configuration
    container := app.Container()
    config := container.MustMake("config").(Config)
    
    if err := config.Validate(); err != nil {
        panic(fmt.Sprintf("Invalid configuration: %v", err))
    }
    
    log.Printf("Configuration loaded successfully for environment: %s", config.Environment)
}

func (p *ConfigProvider) Requires() []string {
    return []string{} // Config provider typically has no dependencies
}

func (p *ConfigProvider) Providers() []string {
    return []string{
        "config",
        "config.env",
    }
}
```

### 2. Advanced Provider với Factory Pattern

```go
type DatabaseProvider struct {
    drivers map[string]DriverFactory
}

type DriverFactory func(config DatabaseConfig) Database

func NewDatabaseProvider() *DatabaseProvider {
    return &DatabaseProvider{
        drivers: map[string]DriverFactory{
            "mysql":    NewMySQLDriver,
            "postgres": NewPostgresDriver,
            "sqlite":   NewSQLiteDriver,
        },
    }
}

func (p *DatabaseProvider) Register(app Application) {
    container := app.Container()
    
    // Register database connection factory
    container.Singleton("db.connection", func() interface{} {
        config := container.MustMake("config").(Config)
        
        factory, exists := p.drivers[config.Database.Driver]
        if !exists {
            panic(fmt.Sprintf("Unsupported database driver: %s", config.Database.Driver))
        }
        
        db := factory(config.Database)
        
        // Setup connection pool
        db.SetMaxOpenConns(config.Database.MaxOpenConns)
        db.SetMaxIdleConns(config.Database.MaxIdleConns)
        
        return db
    })
    
    // Register query builder
    container.Bind("db.query", func() interface{} {
        db := container.MustMake("db.connection").(Database)
        return NewQueryBuilder(db)
    })
    
    // Register transaction manager
    container.Bind("db.transaction", func() interface{} {
        db := container.MustMake("db.connection").(Database)
        return NewTransactionManager(db)
    })
}

func (p *DatabaseProvider) Boot(app Application) {
    container := app.Container()
    
    // Test database connection
    db := container.MustMake("db.connection").(Database)
    if err := db.Ping(); err != nil {
        panic(fmt.Sprintf("Database connection failed: %v", err))
    }
    
    // Run migrations if enabled
    config := container.MustMake("config").(Config)
    if config.Database.AutoMigrate {
        migrator := NewMigrator(db)
        if err := migrator.RunMigrations(); err != nil {
            panic(fmt.Sprintf("Database migration failed: %v", err))
        }
    }
}

func (p *DatabaseProvider) Requires() []string {
    return []string{"config"}
}

func (p *DatabaseProvider) Providers() []string {
    return []string{
        "db.connection",
        "db.query",
        "db.transaction",
    }
}
```

### 3. Conditional Provider

```go
type CacheProvider struct {
    enableRedis  bool
    enableMemory bool
}

func NewCacheProvider(enableRedis, enableMemory bool) *CacheProvider {
    return &CacheProvider{
        enableRedis:  enableRedis,
        enableMemory: enableMemory,
    }
}

func (p *CacheProvider) Register(app Application) {
    container := app.Container()
    
    // Register Redis cache if enabled
    if p.enableRedis {
        container.Singleton("cache.redis", func() interface{} {
            config := container.MustMake("config").(Config)
            return NewRedisCache(config.Redis)
        })
    }
    
    // Register memory cache if enabled
    if p.enableMemory {
        container.Singleton("cache.memory", func() interface{} {
            return NewMemoryCache()
        })
    }
    
    // Register cache manager với fallback logic
    container.Singleton("cache.manager", func() interface{} {
        manager := NewCacheManager()
        
        if p.enableRedis {
            redis := container.MustMake("cache.redis").(Cache)
            manager.AddCache("redis", redis, 100) // High priority
        }
        
        if p.enableMemory {
            memory := container.MustMake("cache.memory").(Cache)
            manager.AddCache("memory", memory, 50) // Medium priority
        }
        
        return manager
    })
    
    // Register default cache alias
    container.Alias("cache", "cache.manager")
}

func (p *CacheProvider) Boot(app Application) {
    container := app.Container()
    cacheManager := container.MustMake("cache.manager").(CacheManager)
    
    // Test cache connections
    if err := cacheManager.HealthCheck(); err != nil {
        log.Printf("Cache health check failed: %v", err)
    }
    
    // Warm up cache with essential data
    p.warmUpCache(cacheManager)
}

func (p *CacheProvider) warmUpCache(cache CacheManager) {
    // Load essential configuration into cache
    go func() {
        essentialKeys := []string{"app.settings", "feature.flags"}
        for _, key := range essentialKeys {
            // Load and cache essential data
            if data, err := p.loadEssentialData(key); err == nil {
                cache.Set(key, data, time.Hour)
            }
        }
    }()
}

func (p *CacheProvider) Requires() []string {
    deps := []string{"config"}
    if p.enableRedis {
        deps = append(deps, "logger") // Redis cache might need logging
    }
    return deps
}

func (p *CacheProvider) Providers() []string {
    providers := []string{"cache.manager", "cache"}
    
    if p.enableRedis {
        providers = append(providers, "cache.redis")
    }
    
    if p.enableMemory {
        providers = append(providers, "cache.memory")
    }
    
    return providers
}
```

### 4. Event-Driven Provider

```go
type EventProvider struct {
    subscribers map[string][]EventHandler
}

type EventHandler func(Event) error

func NewEventProvider() *EventProvider {
    return &EventProvider{
        subscribers: make(map[string][]EventHandler),
    }
}

func (p *EventProvider) Register(app Application) {
    container := app.Container()
    
    // Register event bus
    container.Singleton("event.bus", func() interface{} {
        return NewEventBus()
    })
    
    // Register event dispatcher
    container.Singleton("event.dispatcher", func() interface{} {
        bus := container.MustMake("event.bus").(EventBus)
        return NewEventDispatcher(bus)
    })
    
    // Register async event processor
    container.Singleton("event.async", func() interface{} {
        dispatcher := container.MustMake("event.dispatcher").(EventDispatcher)
        return NewAsyncEventProcessor(dispatcher, 10) // 10 workers
    })
}

func (p *EventProvider) Boot(app Application) {
    container := app.Container()
    bus := container.MustMake("event.bus").(EventBus)
    
    // Register pre-defined event handlers
    for eventName, handlers := range p.subscribers {
        for _, handler := range handlers {
            bus.Subscribe(eventName, handler)
        }
    }
    
    // Start async event processor
    asyncProcessor := container.MustMake("event.async").(AsyncEventProcessor)
    asyncProcessor.Start()
    
    // Register shutdown hook để gracefully stop async processor
    if shutdownHook, ok := app.(ShutdownHook); ok {
        shutdownHook.OnShutdown(func() {
            asyncProcessor.Stop()
        })
    }
}

func (p *EventProvider) Subscribe(event string, handler EventHandler) {
    if p.subscribers[event] == nil {
        p.subscribers[event] = make([]EventHandler, 0)
    }
    p.subscribers[event] = append(p.subscribers[event], handler)
}

func (p *EventProvider) Requires() []string {
    return []string{"logger"} // Event system needs logging
}

func (p *EventProvider) Providers() []string {
    return []string{
        "event.bus",
        "event.dispatcher", 
        "event.async",
    }
}
```

## Error Handling Patterns

### 1. Graceful Degradation Provider

```go
type ResilientProvider struct {
    fallbackEnabled bool
    healthCheck     HealthChecker
}

func (p *ResilientProvider) Register(app Application) {
    container := app.Container()
    
    // Primary service with fallback
    container.Singleton("service.primary", func() interface{} {
        defer func() {
            if r := recover(); r != nil {
                log.Printf("Primary service registration failed: %v", r)
                if p.fallbackEnabled {
                    // Register fallback service instead
                    container.Singleton("service.fallback", func() interface{} {
                        return NewFallbackService()
                    })
                    container.Alias("service.primary", "service.fallback")
                }
            }
        }()
        
        config := container.MustMake("config").(Config)
        return NewPrimaryService(config)
    })
}

func (p *ResilientProvider) Boot(app Application) {
    container := app.Container()
    
    // Health check with graceful handling
    if service, err := container.Make("service.primary"); err == nil {
        if healthChecker, ok := service.(HealthChecker); ok {
            p.healthCheck = healthChecker
            
            // Start health monitoring
            go p.monitorHealth()
        }
    }
}

func (p *ResilientProvider) monitorHealth() {
    ticker := time.NewTicker(30 * time.Second)
    defer ticker.Stop()
    
    for range ticker.C {
        if err := p.healthCheck.HealthCheck(); err != nil {
            log.Printf("Service health check failed: %v", err)
            // Trigger recovery mechanism if needed
        }
    }
}
```

### 2. Validation Provider

```go
type ValidatedProvider struct {
    validators map[string]Validator
}

type Validator func(interface{}) error

func (p *ValidatedProvider) Register(app Application) {
    container := app.Container()
    
    container.Singleton("validated.service", func() interface{} {
        config := container.MustMake("config").(Config)
        
        // Validate configuration before creating service
        if validator, exists := p.validators["config"]; exists {
            if err := validator(config); err != nil {
                panic(fmt.Sprintf("Configuration validation failed: %v", err))
            }
        }
        
        service := NewService(config)
        
        // Validate service after creation
        if validator, exists := p.validators["service"]; exists {
            if err := validator(service); err != nil {
                panic(fmt.Sprintf("Service validation failed: %v", err))
            }
        }
        
        return service
    })
}

func (p *ValidatedProvider) AddValidator(name string, validator Validator) {
    if p.validators == nil {
        p.validators = make(map[string]Validator)
    }
    p.validators[name] = validator
}
```

## Testing Strategies

### 1. Mock Provider for Testing

```go
type MockProvider struct {
    mockServices map[string]interface{}
}

func NewMockProvider() *MockProvider {
    return &MockProvider{
        mockServices: make(map[string]interface{}),
    }
}

func (p *MockProvider) Register(app Application) {
    container := app.Container()
    
    for serviceName, mockService := range p.mockServices {
        container.Singleton(serviceName, func() interface{} {
            return mockService
        })
    }
}

func (p *MockProvider) Boot(app Application) {
    // No-op for mock provider
}

func (p *MockProvider) AddMockService(name string, mock interface{}) {
    p.mockServices[name] = mock
}

func (p *MockProvider) Requires() []string {
    return []string{}
}

func (p *MockProvider) Providers() []string {
    services := make([]string, 0, len(p.mockServices))
    for serviceName := range p.mockServices {
        services = append(services, serviceName)
    }
    return services
}

// Test usage
func TestServiceWithMockProvider(t *testing.T) {
    app := NewTestApplication()
    
    mockProvider := NewMockProvider()
    mockProvider.AddMockService("database", &MockDatabase{})
    mockProvider.AddMockService("cache", &MockCache{})
    
    app.Register(mockProvider)
    app.Boot()
    
    // Test service that depends on database and cache
    service := app.MustMake("my.service")
    assert.NotNil(t, service)
}
```

### 2. Provider Testing

```go
func TestDatabaseProvider(t *testing.T) {
    tests := []struct {
        name     string
        config   Config
        wantErr  bool
        errMsg   string
    }{
        {
            name: "valid mysql config",
            config: Config{
                Database: DatabaseConfig{
                    Driver: "mysql",
                    DSN:    "user:pass@tcp(localhost:3306)/testdb",
                },
            },
            wantErr: false,
        },
        {
            name: "invalid driver",
            config: Config{
                Database: DatabaseConfig{
                    Driver: "invalid",
                    DSN:    "some-dsn",
                },
            },
            wantErr: true,
            errMsg:  "Unsupported database driver: invalid",
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            app := NewTestApplication()
            app.Container().Singleton("config", func() interface{} {
                return tt.config
            })
            
            provider := NewDatabaseProvider()
            
            if tt.wantErr {
                assert.Panics(t, func() {
                    provider.Register(app)
                    app.Boot()
                }, tt.errMsg)
            } else {
                assert.NotPanics(t, func() {
                    provider.Register(app)
                    provider.Boot(app)
                })
                
                // Verify services are registered
                db := app.MustMake("db.connection")
                assert.NotNil(t, db)
            }
        })
    }
}
```

### 3. Integration Testing

```go
func TestProviderIntegration(t *testing.T) {
    app := NewApplication()
    
    // Register providers in dependency order
    providers := []ServiceProvider{
        NewConfigProvider("config/test.yaml"),
        NewLoggerProvider(),
        NewDatabaseProvider(),
        NewCacheProvider(true, true), // Enable both Redis and Memory
        NewWebProvider(),
    }
    
    // Register all providers
    for _, provider := range providers {
        app.Register(provider)
    }
    
    // Boot application
    err := app.Boot()
    require.NoError(t, err)
    
    // Verify critical services are available
    criticalServices := []string{
        "config",
        "logger", 
        "db.connection",
        "cache.manager",
        "web.app",
    }
    
    for _, serviceName := range criticalServices {
        service := app.MustMake(serviceName)
        assert.NotNil(t, service, "Service %s should be available", serviceName)
    }
    
    // Test service interactions
    web := app.MustMake("web.app").(WebApp)
    assert.NotNil(t, web.Database())
    assert.NotNil(t, web.Cache())
}
```

## Performance Considerations

### 1. Lazy Loading Provider

```go
type LazyProvider struct {
    factories map[string]func() interface{}
    loaded    map[string]bool
    mu        sync.RWMutex
}

func NewLazyProvider() *LazyProvider {
    return &LazyProvider{
        factories: make(map[string]func() interface{}),
        loaded:    make(map[string]bool),
    }
}

func (p *LazyProvider) Register(app Application) {
    container := app.Container()
    
    for serviceName, factory := range p.factories {
        // Create lazy wrapper
        container.Singleton(serviceName, func() interface{} {
            p.mu.Lock()
            defer p.mu.Unlock()
            
            if !p.loaded[serviceName] {
                service := factory()
                p.loaded[serviceName] = true
                return service
            }
            
            return factory()
        })
    }
}

func (p *LazyProvider) AddLazyService(name string, factory func() interface{}) {
    p.factories[name] = factory
}
```

### 2. Cached Provider

```go
type CachedProvider struct {
    cache    map[string]interface{}
    cacheTTL map[string]time.Time
    mu       sync.RWMutex
}

func (p *CachedProvider) Register(app Application) {
    container := app.Container()
    
    container.Singleton("expensive.service", func() interface{} {
        p.mu.RLock()
        if service, exists := p.cache["expensive.service"]; exists {
            if time.Now().Before(p.cacheTTL["expensive.service"]) {
                p.mu.RUnlock()
                return service
            }
        }
        p.mu.RUnlock()
        
        // Create expensive service
        service := p.createExpensiveService()
        
        p.mu.Lock()
        p.cache["expensive.service"] = service
        p.cacheTTL["expensive.service"] = time.Now().Add(5 * time.Minute)
        p.mu.Unlock()
        
        return service
    })
}
```

## Best Practices

### 1. Provider Naming và Organization

```go
// Good: Organize providers by domain
type DatabaseProvider struct {}   // Handles all database-related services
type CacheProvider struct {}      // Handles all cache-related services
type WebProvider struct {}        // Handles all web-related services

// Good: Use descriptive service names
func (p *DatabaseProvider) Providers() []string {
    return []string{
        "db.connection.master",     // Clear purpose
        "db.connection.replica",    // Clear purpose
        "db.migrator.up",          // Clear action
        "db.migrator.down",        // Clear action
    }
}

// Avoid: Generic or unclear names
// "db", "cache", "service1", "helper"
```

### 2. Error Handling Best Practices

```go
type RobustProvider struct {}

func (p *RobustProvider) Register(app Application) {
    container := app.Container()
    
    container.Singleton("robust.service", func() interface{} {
        // Validate dependencies first
        config := container.MustMake("config").(Config)
        if config == nil {
            panic("config service is required but not available")
        }
        
        // Handle service creation errors gracefully
        service, err := NewRobustService(config)
        if err != nil {
            // Log error với context
            logger := container.MustMake("logger").(Logger)
            logger.Error("Failed to create robust service", map[string]interface{}{
                "error": err.Error(),
                "config": config,
            })
            
            panic(fmt.Sprintf("Failed to create robust service: %v", err))
        }
        
        return service
    })
}

func (p *RobustProvider) Boot(app Application) {
    // Always use defensive programming
    defer func() {
        if r := recover(); r != nil {
            log.Printf("Provider boot failed: %v", r)
            // Có thể trigger fallback logic ở đây
        }
    }()
    
    container := app.Container()
    service := container.MustMake("robust.service").(RobustService)
    
    // Validate service health before marking as ready
    if err := service.HealthCheck(); err != nil {
        panic(fmt.Sprintf("Service health check failed: %v", err))
    }
}
```

### 3. Configuration-driven Provider

```go
type ConfigurableProvider struct {
    config ProviderConfig
}

type ProviderConfig struct {
    Enabled          bool                   `yaml:"enabled"`
    Services         []ServiceConfig        `yaml:"services"`
    Dependencies     []string               `yaml:"dependencies"`
    BootTimeout      time.Duration          `yaml:"boot_timeout"`
    HealthCheck      HealthCheckConfig      `yaml:"health_check"`
    RetryPolicy      RetryPolicyConfig      `yaml:"retry_policy"`
}

type ServiceConfig struct {
    Name     string                 `yaml:"name"`
    Type     string                 `yaml:"type"`
    Enabled  bool                   `yaml:"enabled"`
    Config   map[string]interface{} `yaml:"config"`
}

func NewConfigurableProvider(config ProviderConfig) *ConfigurableProvider {
    return &ConfigurableProvider{config: config}
}

func (p *ConfigurableProvider) Register(app Application) {
    if !p.config.Enabled {
        return // Skip registration if disabled
    }
    
    container := app.Container()
    
    for _, serviceConfig := range p.config.Services {
        if !serviceConfig.Enabled {
            continue
        }
        
        // Register service based on configuration
        p.registerService(container, serviceConfig)
    }
}

func (p *ConfigurableProvider) registerService(container Container, config ServiceConfig) {
    switch config.Type {
    case "singleton":
        container.Singleton(config.Name, func() interface{} {
            return p.createService(config)
        })
    case "transient":
        container.Bind(config.Name, func() interface{} {
            return p.createService(config)
        })
    default:
        panic(fmt.Sprintf("Unknown service type: %s", config.Type))
    }
}
```

### 4. Versioned Provider

```go
type VersionedProvider struct {
    version     string
    migrations  []Migration
    compatible  []string
}

type Migration struct {
    From    string
    To      string
    Migrate func(Container) error
}

func (p *VersionedProvider) Register(app Application) {
    container := app.Container()
    
    // Check version compatibility
    if appVersion := p.getAppVersion(container); appVersion != "" {
        if !p.isCompatible(appVersion) {
            panic(fmt.Sprintf("Provider version %s not compatible with app version %s", 
                p.version, appVersion))
        }
    }
    
    // Register versioned services
    container.Singleton("service.v1", func() interface{} {
        return NewServiceV1()
    })
    
    // Register version alias
    container.Alias("service", "service.v1")
}

func (p *VersionedProvider) isCompatible(appVersion string) bool {
    for _, compatible := range p.compatible {
        if compatible == appVersion {
            return true
        }
    }
    return false
}
```

## Advanced Patterns

### 1. Middleware Provider

```go
type MiddlewareProvider struct {
    middlewares []MiddlewareDefinition
}

type MiddlewareDefinition struct {
    Name     string
    Priority int
    Factory  func(Container) Middleware
}

type Middleware interface {
    Handle(Context, func(Context)) error
}

func (p *MiddlewareProvider) Register(app Application) {
    container := app.Container()
    
    // Register middleware factory
    container.Singleton("middleware.stack", func() interface{} {
        stack := NewMiddlewareStack()
        
        // Sort middlewares by priority
        sort.Slice(p.middlewares, func(i, j int) bool {
            return p.middlewares[i].Priority < p.middlewares[j].Priority
        })
        
        // Register middlewares in order
        for _, mw := range p.middlewares {
            middleware := mw.Factory(container)
            stack.Add(mw.Name, middleware)
        }
        
        return stack
    })
}

func (p *MiddlewareProvider) AddMiddleware(name string, priority int, factory func(Container) Middleware) {
    p.middlewares = append(p.middlewares, MiddlewareDefinition{
        Name:     name,
        Priority: priority,
        Factory:  factory,
    })
}
```

### 2. Plugin Provider System

```go
type PluginProvider struct {
    plugins     []Plugin
    loader      PluginLoader
    sandbox     Sandbox
}

type Plugin interface {
    Name() string
    Version() string
    Dependencies() []string
    Register(Container) error
    Unload() error
}

type PluginLoader interface {
    LoadPlugin(path string) (Plugin, error)
    UnloadPlugin(Plugin) error
}

type Sandbox interface {
    Execute(Plugin, func()) error
}

func (p *PluginProvider) Register(app Application) {
    container := app.Container()
    
    // Register plugin manager
    container.Singleton("plugin.manager", func() interface{} {
        return NewPluginManager(p.loader, p.sandbox)
    })
    
    // Register each plugin in sandbox
    for _, plugin := range p.plugins {
        err := p.sandbox.Execute(plugin, func() {
            if err := plugin.Register(container); err != nil {
                panic(fmt.Sprintf("Plugin %s registration failed: %v", plugin.Name(), err))
            }
        })
        
        if err != nil {
            log.Printf("Plugin %s failed to register: %v", plugin.Name(), err)
        }
    }
}
```

## Monitoring và Observability

### 1. Metrics Provider

```go
type MetricsProvider struct {
    registry MetricsRegistry
}

func (p *MetricsProvider) Register(app Application) {
    container := app.Container()
    
    // Register metrics registry
    container.Singleton("metrics.registry", func() interface{} {
        return p.registry
    })
    
    // Register metrics collector
    container.Singleton("metrics.collector", func() interface{} {
        registry := container.MustMake("metrics.registry").(MetricsRegistry)
        return NewMetricsCollector(registry)
    })
    
    // Register service metrics wrappers
    container.Extend("db.connection", func(db interface{}) interface{} {
        collector := container.MustMake("metrics.collector").(MetricsCollector)
        return NewInstrumentedDatabase(db.(Database), collector)
    })
}

func (p *MetricsProvider) Boot(app Application) {
    container := app.Container()
    collector := container.MustMake("metrics.collector").(MetricsCollector)
    
    // Start metrics collection
    go collector.Start()
    
    // Register HTTP metrics endpoint
    if webApp, err := container.Make("web.app"); err == nil {
        web := webApp.(WebApp)
        web.GET("/metrics", func(c Context) {
            metrics := collector.GetMetrics()
            c.JSON(200, metrics)
        })
    }
}
```

### 2. Tracing Provider

```go
type TracingProvider struct {
    tracer Tracer
    config TracingConfig
}

func (p *TracingProvider) Register(app Application) {
    container := app.Container()
    
    container.Singleton("tracer", func() interface{} {
        return p.tracer
    })
    
    // Wrap services with tracing
    services := []string{"db.connection", "cache.manager", "web.router"}
    for _, serviceName := range services {
        container.Extend(serviceName, func(service interface{}) interface{} {
            return p.wrapWithTracing(serviceName, service)
        })
    }
}

func (p *TracingProvider) wrapWithTracing(name string, service interface{}) interface{} {
    return &TracedService{
        name:    name,
        service: service,
        tracer:  p.tracer,
    }
}
```

## Liên kết

- [Container Documentation](./container.md) - Tài liệu về DI Container
- [Application Documentation](./application.md) - Tài liệu về Application interface
- [Loader Documentation](./loader.md) - Tài liệu về ModuleLoader
- [Deferred Documentation](./deferred.md) - Tài liệu về Deferred Provider (nếu có)
- [Best Practices Guide](./best-practices.md) - Hướng dẫn best practices (nếu có)

---

*Tài liệu này được tạo tự động dựa trên source code và doctype của `provider.go`. Vui lòng cập nhật khi có thay đổi trong implementation.*
