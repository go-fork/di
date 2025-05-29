# ModuleLoader - Tài liệu kỹ thuật

## Tổng quan

`ModuleLoaderContract` là interface định nghĩa contract cho các loader chịu trách nhiệm nạp module/service provider vào ứng dụng trong hệ thống Fork framework. Đây là thành phần cốt lõi hỗ trợ kiến trúc modular, cho phép tách biệt logic khởi tạo module khỏi phần còn lại của ứng dụng.

## Định nghĩa & Vai trò

### Mục đích chính

ModuleLoader interface phục vụ các mục đích sau:

- **Modular Architecture**: Chuẩn hóa quy trình nạp, đăng ký, boot các module
- **Separation of Concerns**: Tách biệt logic khởi tạo module khỏi application core
- **Extensibility**: Tăng khả năng mở rộng và kiểm soát vòng đời module
- **Bootstrap Management**: Quản lý quá trình khởi tạo ứng dụng theo phases
- **Dynamic Loading**: Hỗ trợ tải module động theo nhu cầu

### Kiến trúc và Pattern

ModuleLoader áp dụng các pattern sau:

- **Builder Pattern**: Xây dựng ứng dụng theo từng bước có thứ tự
- **Plugin Architecture**: Hỗ trợ tải module như các plugin
- **Bootstrap Pattern**: Quản lý vòng đời khởi tạo theo phases
- **Factory Pattern**: Tạo và cấu hình modules
- **Command Pattern**: Encapsulate các thao tác loading

## Tính năng cốt lõi

### 1. Bootstrap Management
- Khởi tạo ứng dụng với cấu hình cần thiết
- Đăng ký core service providers
- Thiết lập foundation dependencies

### 2. Module Loading
- Tải module đơn lẻ hoặc batch
- Validation và type checking
- Error handling và rollback

### 3. Provider Management
- Đăng ký core vs. application providers
- Dependency resolution
- Boot sequence management

### 4. Lifecycle Control
- Phân chia rõ ràng các phase khởi tạo
- Ordered execution
- Error recovery và cleanup

## API Reference

### Bootstrap Methods

#### `BootstrapApplication() error`

Khởi tạo ứng dụng với các cấu hình và provider cần thiết.

**Mục đích:**
- Thực hiện toàn bộ quá trình bootstrap ứng dụng
- Đăng ký core providers và boot chúng
- Thiết lập foundation cho việc load modules

**Trả về:**
- `error`: Lỗi nếu đăng ký hoặc boot provider thất bại

**Workflow:**
1. Đăng ký core service providers
2. Register dependencies với dependency resolution
3. Boot tất cả registered providers
4. Validate application state

**Ví dụ:**
```go
type AppLoader struct {
    app Application
}

func NewAppLoader(app Application) *AppLoader {
    return &AppLoader{app: app}
}

func (loader *AppLoader) BootstrapApplication() error {
    // Phase 1: Register core providers
    if err := loader.RegisterCoreProviders(); err != nil {
        return fmt.Errorf("failed to register core providers: %w", err)
    }
    
    // Phase 2: Register with dependencies
    if err := loader.app.RegisterWithDependencies(); err != nil {
        return fmt.Errorf("failed to resolve dependencies: %w", err)
    }
    
    // Phase 3: Boot all providers
    if err := loader.app.Boot(); err != nil {
        return fmt.Errorf("failed to boot application: %w", err)
    }
    
    return nil
}
```

#### `RegisterCoreProviders() error`

Đăng ký các core service provider.

**Mục đích:**
- Đăng ký các provider cơ bản cần thiết cho mọi ứng dụng
- Thiết lập infrastructure dependencies
- Chuẩn bị foundation cho application providers

**Trả về:**
- `error`: Lỗi nếu đăng ký thất bại

**Core Providers thường bao gồm:**
- Configuration Provider
- Logging Provider  
- Error Handling Provider
- Event System Provider
- Cache Provider (optional)

**Ví dụ:**
```go
func (loader *AppLoader) RegisterCoreProviders() error {
    coreProviders := []ServiceProvider{
        &ConfigProvider{},
        &LoggerProvider{},
        &EventProvider{},
        &ErrorHandlerProvider{},
    }
    
    for _, provider := range coreProviders {
        if provider == nil {
            return fmt.Errorf("nil core provider encountered")
        }
        
        loader.app.Register(provider)
    }
    
    return nil
}
```

### Module Loading Methods

#### `LoadModule(module interface{}) error`

Tải một module vào ứng dụng.

**Tham số:**
- `module`: Module cần nạp, có thể là service provider hoặc các kiểu module khác

**Trả về:**
- `error`: Lỗi nếu module không hợp lệ hoặc nạp thất bại

**Mục đích:**
- Tải và đăng ký một module đơn lẻ
- Validation module trước khi load
- Graceful error handling

**Module Types được hỗ trợ:**
- `ServiceProvider`: Standard service provider
- `ServiceProviderDeferred`: Deferred service provider  
- `ModuleConfig`: Configuration-based module
- `PluginModule`: Plugin-style module

**Ví dụ:**
```go
func (loader *AppLoader) LoadModule(module interface{}) error {
    if module == nil {
        return fmt.Errorf("module cannot be nil")
    }
    
    switch m := module.(type) {
    case ServiceProvider:
        return loader.loadServiceProvider(m)
    case ServiceProviderDeferred:
        return loader.loadDeferredProvider(m)
    case *ModuleConfig:
        return loader.loadConfigModule(m)
    case PluginModule:
        return loader.loadPluginModule(m)
    default:
        return fmt.Errorf("unsupported module type: %T", module)
    }
}

func (loader *AppLoader) loadServiceProvider(provider ServiceProvider) error {
    // Validate provider
    if err := loader.validateProvider(provider); err != nil {
        return fmt.Errorf("provider validation failed: %w", err)
    }
    
    // Register provider
    loader.app.Register(provider)
    
    // Boot immediately if app is already bootstrapped
    if loader.isBootstrapped() {
        if err := provider.Boot(loader.app); err != nil {
            return fmt.Errorf("failed to boot provider %T: %w", provider, err)
        }
    }
    
    return nil
}
```

#### `LoadModules(modules ...interface{}) error`

Tải nhiều module vào ứng dụng.

**Tham số:**
- `modules`: Danh sách module cần nạp

**Trả về:**
- `error`: Lỗi nếu bất kỳ module nào nạp thất bại

**Mục đích:**
- Batch loading nhiều modules
- Atomic operation (all or nothing)
- Optimized dependency resolution

**Strategies:**
- **Fail-Fast**: Dừng lại khi gặp lỗi đầu tiên
- **Dependency-Aware**: Sắp xếp modules theo dependencies
- **Rollback**: Rollback khi có lỗi

**Ví dụ:**
```go
func (loader *AppLoader) LoadModules(modules ...interface{}) error {
    if len(modules) == 0 {
        return nil
    }
    
    // Phase 1: Validate all modules
    for i, module := range modules {
        if module == nil {
            return fmt.Errorf("module at index %d is nil", i)
        }
        
        if err := loader.validateModule(module); err != nil {
            return fmt.Errorf("module %d validation failed: %w", i, err)
        }
    }
    
    // Phase 2: Load modules in dependency order
    orderedModules, err := loader.resolveDependencyOrder(modules)
    if err != nil {
        return fmt.Errorf("dependency resolution failed: %w", err)
    }
    
    // Phase 3: Load each module
    var loadedModules []interface{}
    for _, module := range orderedModules {
        if err := loader.LoadModule(module); err != nil {
            // Rollback loaded modules
            loader.rollbackModules(loadedModules)
            return fmt.Errorf("failed to load module %T: %w", module, err)
        }
        loadedModules = append(loadedModules, module)
    }
    
    return nil
}
```

## Implementation Patterns

### 1. Basic Module Loader

```go
type BasicModuleLoader struct {
    app         Application
    bootstrapped bool
    coreProviders []ServiceProvider
    modules      []interface{}
    mu          sync.RWMutex
}

func NewBasicModuleLoader(app Application) *BasicModuleLoader {
    return &BasicModuleLoader{
        app:           app,
        bootstrapped:  false,
        coreProviders: make([]ServiceProvider, 0),
        modules:       make([]interface{}, 0),
    }
}

func (loader *BasicModuleLoader) BootstrapApplication() error {
    loader.mu.Lock()
    defer loader.mu.Unlock()
    
    if loader.bootstrapped {
        return fmt.Errorf("application already bootstrapped")
    }
    
    if err := loader.RegisterCoreProviders(); err != nil {
        return err
    }
    
    if err := loader.app.RegisterWithDependencies(); err != nil {
        return err
    }
    
    if err := loader.app.Boot(); err != nil {
        return err
    }
    
    loader.bootstrapped = true
    return nil
}

func (loader *BasicModuleLoader) RegisterCoreProviders() error {
    for _, provider := range loader.coreProviders {
        loader.app.Register(provider)
    }
    return nil
}

func (loader *BasicModuleLoader) isBootstrapped() bool {
    loader.mu.RLock()
    defer loader.mu.RUnlock()
    return loader.bootstrapped
}
```

### 2. Configuration-Driven Loader

```go
type ConfigModuleLoader struct {
    BasicModuleLoader
    config *LoaderConfig
}

type LoaderConfig struct {
    CoreProviders []string          `yaml:"core_providers"`
    Modules       []ModuleConfig    `yaml:"modules"`
    LoadOrder     []string          `yaml:"load_order"`
    Environment   string            `yaml:"environment"`
}

type ModuleConfig struct {
    Name        string                 `yaml:"name"`
    Type        string                 `yaml:"type"`
    Package     string                 `yaml:"package"`
    Config      map[string]interface{} `yaml:"config"`
    Dependencies []string              `yaml:"dependencies"`
    Environment []string               `yaml:"environment"`
}

func NewConfigModuleLoader(app Application, config *LoaderConfig) *ConfigModuleLoader {
    return &ConfigModuleLoader{
        BasicModuleLoader: *NewBasicModuleLoader(app),
        config:           config,
    }
}

func (loader *ConfigModuleLoader) RegisterCoreProviders() error {
    for _, providerName := range loader.config.CoreProviders {
        provider, err := loader.createProvider(providerName)
        if err != nil {
            return fmt.Errorf("failed to create core provider %s: %w", providerName, err)
        }
        
        loader.app.Register(provider)
    }
    
    return nil
}

func (loader *ConfigModuleLoader) LoadModulesFromConfig() error {
    modules := make([]interface{}, 0, len(loader.config.Modules))
    
    for _, moduleConfig := range loader.config.Modules {
        // Check environment
        if !loader.isEnvironmentMatch(moduleConfig.Environment) {
            continue
        }
        
        module, err := loader.createModule(moduleConfig)
        if err != nil {
            return fmt.Errorf("failed to create module %s: %w", moduleConfig.Name, err)
        }
        
        modules = append(modules, module)
    }
    
    return loader.LoadModules(modules...)
}
```

### 3. Plugin-Based Loader

```go
type PluginModuleLoader struct {
    BasicModuleLoader
    pluginDir   string
    plugins     map[string]*plugin.Plugin
    registry    *PluginRegistry
}

type PluginRegistry struct {
    plugins map[string]PluginInfo
    mu      sync.RWMutex
}

type PluginInfo struct {
    Name        string
    Version     string
    Path        string
    Loaded      bool
    Provider    ServiceProvider
}

func NewPluginModuleLoader(app Application, pluginDir string) *PluginModuleLoader {
    return &PluginModuleLoader{
        BasicModuleLoader: *NewBasicModuleLoader(app),
        pluginDir:         pluginDir,
        plugins:           make(map[string]*plugin.Plugin),
        registry:          &PluginRegistry{
            plugins: make(map[string]PluginInfo),
        },
    }
}

func (loader *PluginModuleLoader) LoadPlugin(name string) error {
    pluginPath := filepath.Join(loader.pluginDir, name+".so")
    
    p, err := plugin.Open(pluginPath)
    if err != nil {
        return fmt.Errorf("failed to open plugin %s: %w", name, err)
    }
    
    // Look for NewProvider symbol
    symbol, err := p.Lookup("NewProvider")
    if err != nil {
        return fmt.Errorf("plugin %s missing NewProvider symbol: %w", name, err)
    }
    
    // Type assert to provider factory
    newProvider, ok := symbol.(func() ServiceProvider)
    if !ok {
        return fmt.Errorf("plugin %s NewProvider has wrong signature", name)
    }
    
    provider := newProvider()
    
    // Register plugin
    loader.registry.Register(name, PluginInfo{
        Name:     name,
        Path:     pluginPath,
        Loaded:   true,
        Provider: provider,
    })
    
    loader.plugins[name] = p
    
    return loader.LoadModule(provider)
}

func (loader *PluginModuleLoader) LoadPluginsFromDirectory() error {
    files, err := filepath.Glob(filepath.Join(loader.pluginDir, "*.so"))
    if err != nil {
        return fmt.Errorf("failed to scan plugin directory: %w", err)
    }
    
    for _, file := range files {
        name := strings.TrimSuffix(filepath.Base(file), ".so")
        if err := loader.LoadPlugin(name); err != nil {
            return fmt.Errorf("failed to load plugin %s: %w", name, err)
        }
    }
    
    return nil
}
```

## Advanced Features

### 1. Dependency Resolution

```go
func (loader *BasicModuleLoader) resolveDependencyOrder(modules []interface{}) ([]interface{}, error) {
    // Create dependency graph
    graph := make(map[string][]string)
    moduleMap := make(map[string]interface{})
    
    for _, module := range modules {
        name := loader.getModuleName(module)
        deps := loader.getModuleDependencies(module)
        
        graph[name] = deps
        moduleMap[name] = module
    }
    
    // Topological sort
    sorted, err := topologicalSort(graph)
    if err != nil {
        return nil, fmt.Errorf("circular dependency detected: %w", err)
    }
    
    // Return sorted modules
    result := make([]interface{}, 0, len(sorted))
    for _, name := range sorted {
        if module, exists := moduleMap[name]; exists {
            result = append(result, module)
        }
    }
    
    return result, nil
}

func topologicalSort(graph map[string][]string) ([]string, error) {
    visited := make(map[string]bool)
    temp := make(map[string]bool)
    result := make([]string, 0)
    
    var visit func(string) error
    visit = func(node string) error {
        if temp[node] {
            return fmt.Errorf("circular dependency involving %s", node)
        }
        if visited[node] {
            return nil
        }
        
        temp[node] = true
        for _, dep := range graph[node] {
            if err := visit(dep); err != nil {
                return err
            }
        }
        temp[node] = false
        visited[node] = true
        result = append([]string{node}, result...)
        
        return nil
    }
    
    for node := range graph {
        if !visited[node] {
            if err := visit(node); err != nil {
                return nil, err
            }
        }
    }
    
    return result, nil
}
```

### 2. Module Validation

```go
func (loader *BasicModuleLoader) validateModule(module interface{}) error {
    if module == nil {
        return fmt.Errorf("module cannot be nil")
    }
    
    switch m := module.(type) {
    case ServiceProvider:
        return loader.validateServiceProvider(m)
    case ServiceProviderDeferred:
        return loader.validateDeferredProvider(m)
    default:
        return fmt.Errorf("unsupported module type: %T", module)
    }
}

func (loader *BasicModuleLoader) validateServiceProvider(provider ServiceProvider) error {
    // Check if provider implements required methods
    if provider == nil {
        return fmt.Errorf("service provider cannot be nil")
    }
    
    // Validate dependencies
    deps := provider.Requires()
    for _, dep := range deps {
        if !loader.isDependencyAvailable(dep) {
            return fmt.Errorf("dependency %s not available", dep)
        }
    }
    
    // Validate provided services
    services := provider.Providers()
    for _, service := range services {
        if service == "" {
            return fmt.Errorf("empty service name in provider")
        }
    }
    
    return nil
}

func (loader *BasicModuleLoader) isDependencyAvailable(dep string) bool {
    // Check if dependency is registered in container
    return loader.app.Container().Bound(dep)
}
```

### 3. Hot Module Reloading

```go
type HotReloadLoader struct {
    BasicModuleLoader
    watcher    *fsnotify.Watcher
    reloadChan chan string
    hotModules map[string]interface{}
}

func NewHotReloadLoader(app Application) *HotReloadLoader {
    watcher, _ := fsnotify.NewWatcher()
    
    return &HotReloadLoader{
        BasicModuleLoader: *NewBasicModuleLoader(app),
        watcher:           watcher,
        reloadChan:        make(chan string, 10),
        hotModules:        make(map[string]interface{}),
    }
}

func (loader *HotReloadLoader) EnableHotReload(modulePath string) error {
    return loader.watcher.Add(modulePath)
}

func (loader *HotReloadLoader) StartWatching() {
    go func() {
        for {
            select {
            case event := <-loader.watcher.Events:
                if event.Op&fsnotify.Write == fsnotify.Write {
                    loader.reloadChan <- event.Name
                }
            case err := <-loader.watcher.Errors:
                log.Printf("Watcher error: %v", err)
            }
        }
    }()
    
    go func() {
        for modulePath := range loader.reloadChan {
            if err := loader.reloadModule(modulePath); err != nil {
                log.Printf("Failed to reload module %s: %v", modulePath, err)
            }
        }
    }()
}

func (loader *HotReloadLoader) reloadModule(modulePath string) error {
    // Unload old module
    if oldModule, exists := loader.hotModules[modulePath]; exists {
        if err := loader.unloadModule(oldModule); err != nil {
            return fmt.Errorf("failed to unload old module: %w", err)
        }
    }
    
    // Load new module
    newModule, err := loader.loadModuleFromPath(modulePath)
    if err != nil {
        return fmt.Errorf("failed to load new module: %w", err)
    }
    
    if err := loader.LoadModule(newModule); err != nil {
        return fmt.Errorf("failed to register new module: %w", err)
    }
    
    loader.hotModules[modulePath] = newModule
    return nil
}
```

## Error Handling & Recovery

### 1. Graceful Error Handling

```go
type ErrorRecoveryLoader struct {
    BasicModuleLoader
    errorHandler ErrorHandler
    recovery     RecoveryStrategy
}

type ErrorHandler interface {
    HandleLoadError(module interface{}, err error) error
    HandleBootError(provider ServiceProvider, err error) error
}

type RecoveryStrategy interface {
    Recover(failedModules []interface{}) error
    Rollback(loadedModules []interface{}) error
}

func (loader *ErrorRecoveryLoader) LoadModules(modules ...interface{}) error {
    var loadedModules []interface{}
    var errors []error
    
    for _, module := range modules {
        if err := loader.LoadModule(module); err != nil {
            errors = append(errors, err)
            
            // Handle error
            if handlerErr := loader.errorHandler.HandleLoadError(module, err); handlerErr != nil {
                // Recovery failed, rollback
                if rollbackErr := loader.recovery.Rollback(loadedModules); rollbackErr != nil {
                    return fmt.Errorf("load failed and rollback failed: %w, rollback: %w", err, rollbackErr)
                }
                return fmt.Errorf("module load failed: %w", err)
            }
            
            // Error handled successfully, continue
            continue
        }
        
        loadedModules = append(loadedModules, module)
    }
    
    if len(errors) > 0 {
        // Some modules failed but were handled
        log.Printf("Some modules failed to load but were handled: %v", errors)
    }
    
    return nil
}
```

### 2. Health Checks

```go
func (loader *BasicModuleLoader) HealthCheck() error {
    // Check if application is bootstrapped
    if !loader.bootstrapped {
        return fmt.Errorf("application not bootstrapped")
    }
    
    // Check core providers
    coreServices := []string{"config", "logger", "events"}
    for _, service := range coreServices {
        if !loader.app.Container().Bound(service) {
            return fmt.Errorf("core service %s not available", service)
        }
    }
    
    // Check module health
    for _, module := range loader.modules {
        if healthChecker, ok := module.(HealthChecker); ok {
            if err := healthChecker.HealthCheck(); err != nil {
                return fmt.Errorf("module %T health check failed: %w", module, err)
            }
        }
    }
    
    return nil
}

type HealthChecker interface {
    HealthCheck() error
}
```

## Testing Strategies

### 1. Mock Module Loader

```go
func TestModuleLoader(t *testing.T) {
    mockApp := new(mocks.Application)
    mockContainer := di.New()
    
    mockApp.On("Container").Return(mockContainer)
    mockApp.On("Register", mock.AnythingOfType("ServiceProvider")).Return()
    mockApp.On("RegisterWithDependencies").Return(nil)
    mockApp.On("Boot").Return(nil)
    
    loader := NewBasicModuleLoader(mockApp)
    
    err := loader.BootstrapApplication()
    assert.NoError(t, err)
    
    mockApp.AssertExpectations(t)
}
```

### 2. Integration Testing

```go
func TestModuleLoaderIntegration(t *testing.T) {
    config := &Config{
        Environment: "test",
        Database: DatabaseConfig{
            Driver: "sqlite",
            DSN:    ":memory:",
        },
    }
    
    app := NewApp(config)
    loader := NewConfigModuleLoader(app, &LoaderConfig{
        CoreProviders: []string{"config", "logger"},
        Modules: []ModuleConfig{
            {
                Name: "database",
                Type: "provider",
                Package: "github.com/app/providers/database",
            },
        },
    })
    
    err := loader.BootstrapApplication()
    require.NoError(t, err)
    
    err = loader.LoadModulesFromConfig()
    require.NoError(t, err)
    
    // Verify services are available
    db := app.MustMake("database")
    assert.NotNil(t, db)
}
```

## Performance Considerations

### 1. Lazy Loading

```go
type LazyModuleLoader struct {
    BasicModuleLoader
    lazyModules map[string]func() (interface{}, error)
}

func (loader *LazyModuleLoader) RegisterLazyModule(name string, factory func() (interface{}, error)) {
    loader.lazyModules[name] = factory
}

func (loader *LazyModuleLoader) LoadModuleOnDemand(name string) error {
    factory, exists := loader.lazyModules[name]
    if !exists {
        return fmt.Errorf("lazy module %s not registered", name)
    }
    
    module, err := factory()
    if err != nil {
        return fmt.Errorf("failed to create lazy module %s: %w", name, err)
    }
    
    return loader.LoadModule(module)
}
```

### 2. Concurrent Loading

```go
func (loader *BasicModuleLoader) LoadModulesConcurrent(modules ...interface{}) error {
    var wg sync.WaitGroup
    errChan := make(chan error, len(modules))
    
    for _, module := range modules {
        wg.Add(1)
        go func(m interface{}) {
            defer wg.Done()
            if err := loader.LoadModule(m); err != nil {
                errChan <- err
            }
        }(module)
    }
    
    wg.Wait()
    close(errChan)
    
    // Check for errors
    for err := range errChan {
        return err
    }
    
    return nil
}
```

## Best Practices

### 1. Module Organization

```go
// Organize modules by category
type ModuleCategory string

const (
    CoreModule         ModuleCategory = "core"
    InfrastructureModule ModuleCategory = "infrastructure"
    DomainModule       ModuleCategory = "domain"
    WebModule          ModuleCategory = "web"
    PluginModule       ModuleCategory = "plugin"
)

func (loader *BasicModuleLoader) LoadModulesByCategory(category ModuleCategory, modules ...interface{}) error {
    log.Printf("Loading %s modules...", category)
    
    if err := loader.LoadModules(modules...); err != nil {
        return fmt.Errorf("failed to load %s modules: %w", category, err)
    }
    
    log.Printf("Successfully loaded %d %s modules", len(modules), category)
    return nil
}
```

### 2. Environment-Specific Loading

```go
func (loader *ConfigModuleLoader) LoadEnvironmentModules() error {
    env := loader.config.Environment
    
    var modules []interface{}
    for _, moduleConfig := range loader.config.Modules {
        if loader.isEnvironmentMatch(moduleConfig.Environment) {
            module, err := loader.createModule(moduleConfig)
            if err != nil {
                return err
            }
            modules = append(modules, module)
        }
    }
    
    log.Printf("Loading %d modules for environment: %s", len(modules), env)
    return loader.LoadModules(modules...)
}

func (loader *ConfigModuleLoader) isEnvironmentMatch(envs []string) bool {
    if len(envs) == 0 {
        return true // Load in all environments if not specified
    }
    
    for _, env := range envs {
        if env == loader.config.Environment {
            return true
        }
    }
    
    return false
}
```

## Tài liệu liên quan

- [container.md](./container.md) - Chi tiết về DI Container
- [application.md](./application.md) - Application interface  
- [provider.md](./provider.md) - Service Provider pattern
- [deferred.md](./deferred.md) - Deferred service providers
- [examples/](../examples/) - Các ví dụ module loading thực tế

## Changelog

Xem [CHANGELOG.md](../CHANGELOG.md) để biết lịch sử thay đổi của ModuleLoader interface.
