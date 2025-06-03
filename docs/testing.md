# Testing với Container Interface

## Tổng quan

Tài liệu này hướng dẫn cách test các components sử dụng `Container` interface trong Fork framework. Với sự giới thiệu của Container interface từ phiên bản v0.1.2, việc mock và test các components sử dụng DI container trở nên dễ dàng và linh hoạt hơn.

## Testing với Mocks

### Sử dụng Container Mock

```go
import (
    "testing"
    "go.fork.vn/di"
    "go.fork.vn/di/mocks"
    "github.com/stretchr/testify/mock"
)

func TestService(t *testing.T) {
    // Tạo mock cho Container interface
    mockContainer := new(mocks.Container)
    
    // Thiết lập expectations cho mock
    mockContainer.On("MustMake", "logger").Return(&MockLogger{})
    mockContainer.On("MustMake", "database").Return(&MockDatabase{})
    
    // Sử dụng mock container trong service cần test
    service := NewService(mockContainer)
    
    // Thực hiện test
    result := service.DoSomething()
    
    // Kiểm tra kết quả
    assert.Equal(t, expectedResult, result)
    
    // Verify expectations
    mockContainer.AssertExpectations(t)
}
```

### Tạo Service với Container Injection

```go
// Service sử dụng Container interface
type Service struct {
    container di.Container
}

// Khởi tạo service với container injection
func NewService(container di.Container) *Service {
    return &Service{
        container: container,
    }
}

// Method sử dụng container để resolve dependencies
func (s *Service) DoSomething() Result {
    logger := s.container.MustMake("logger").(Logger)
    db := s.container.MustMake("database").(Database)
    
    // Logic sử dụng logger và database...
}
```

## Testing với Custom Container Implementation

Bạn có thể tạo custom implementation của Container interface để sử dụng trong tests:

```go
// Custom test container
type TestContainer struct {
    instances map[string]interface{}
}

func NewTestContainer() *TestContainer {
    return &TestContainer{
        instances: make(map[string]interface{}),
    }
}

// Implementation của Container interface
func (c *TestContainer) Bind(abstract string, concrete di.BindingFunc) {
    c.instances[abstract] = concrete(c)
}

func (c *TestContainer) BindIf(abstract string, concrete di.BindingFunc) bool {
    if _, exists := c.instances[abstract]; exists {
        return false
    }
    c.instances[abstract] = concrete(c)
    return true
}

func (c *TestContainer) Singleton(abstract string, concrete di.BindingFunc) {
    c.instances[abstract] = concrete(c)
}

func (c *TestContainer) Instance(abstract string, instance interface{}) {
    c.instances[abstract] = instance
}

func (c *TestContainer) Alias(abstract, alias string) {
    c.instances[alias] = c.instances[abstract]
}

func (c *TestContainer) Make(abstract string) (interface{}, error) {
    instance, exists := c.instances[abstract]
    if !exists {
        return nil, fmt.Errorf("binding not found: %s", abstract)
    }
    return instance, nil
}

func (c *TestContainer) MustMake(abstract string) interface{} {
    instance, err := c.Make(abstract)
    if err != nil {
        panic(err)
    }
    return instance
}

func (c *TestContainer) Bound(abstract string) bool {
    _, exists := c.instances[abstract]
    return exists
}

func (c *TestContainer) Reset() {
    c.instances = make(map[string]interface{})
}

func (c *TestContainer) Call(callback interface{}, additionalParams ...interface{}) ([]interface{}, error) {
    // Implement Call cho test container nếu cần
    // ...
    return nil, nil
}

// Sử dụng trong tests
func TestWithCustomContainer(t *testing.T) {
    container := NewTestContainer()
    
    // Thiết lập test dependencies
    container.Instance("logger", &MockLogger{})
    container.Instance("database", &MockDatabase{})
    
    // Khởi tạo service với custom container
    service := NewService(container)
    
    // Thực hiện test
    // ...
}
```

## Tích hợp với Application Mock

Khi cần test component sử dụng Application interface:

```go
func TestWithApplicationAndContainer(t *testing.T) {
    // Tạo mock cho Container và Application
    mockContainer := new(mocks.Container)
    mockApp := new(mocks.Application)
    
    // Thiết lập expectations cho Application
    mockApp.On("Container").Return(mockContainer)
    
    // Thiết lập expectations cho Container
    mockContainer.On("MustMake", "config").Return(&MockConfig{})
    
    // Sử dụng trong provider hoặc component cần test
    provider := NewProvider()
    provider.Register(mockApp)
    
    // Verify expectations
    mockApp.AssertExpectations(t)
    mockContainer.AssertExpectations(t)
}
```

## Best Practices

1. **Khai báo rõ ràng dependencies**: Component nên khai báo rõ những dependency nào sẽ được resolve từ container
2. **Sử dụng constructor injection**: Inject container qua constructor thay vì truy cập global container
3. **Mock riêng biệt**: Tạo mock cho từng dependency thay vì mock toàn bộ container nếu có thể
4. **Tránh over-mocking**: Chỉ mock những thành phần cần thiết
5. **Verify expectations**: Luôn kiểm tra rằng tất cả expectations được gọi đúng cách

## Lưu ý quan trọng

- `Container` là interface từ v0.1.2, cho phép dễ dàng mock và thay thế trong tests
- Mock objects cho `Container` được tạo sẵn trong package `go.fork.vn/di/mocks`
- Với complex interactions, cân nhắc sử dụng real Container với test instances
- Luôn assertExpectations để đảm bảo mocks được sử dụng đúng cách

---

Xem thêm:
- [Container Documentation](container.md)
- [Mocks Documentation](../mocks/README.md)
- [Application Documentation](application.md)
