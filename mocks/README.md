# Mocks

Thư mục này chứa các mock object được tạo tự động từ các interface trong package `di` bằng công cụ [mockery](https://github.com/vektra/mockery).

## Mục đích

Các mock object này giúp:
- Tạo mocks cho testing
- Thực hiện Unit Testing với các dependency giả lập
- Kiểm tra tương tác giữa các component mà không cần triển khai đầy đủ

## Các Mock đã tạo

1. **Application** - Mock cho interface `Application`
2. **ServiceProvider** - Mock cho interface `ServiceProvider`
3. **ServiceProviderDeferred** - Mock cho interface `ServiceProviderDeferred`
4. **ModuleLoaderContract** - Mock cho interface `ModuleLoaderContract`

## Cách sử dụng

```go
import (
	"testing"
	"github.com/go-fork/di/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestExample(t *testing.T) {
	// Tạo mock object
	mockApp := new(mocks.Application)
	
	// Thiết lập kỳ vọng cho mock
	mockApp.On("Container").Return(&di.Container{})
	mockApp.On("Make", "service").Return("mock-service", nil)
	
	// Sử dụng mock trong test
	container := mockApp.Container()
	service, err := mockApp.Make("service")
	
	// Kiểm tra kết quả
	assert.NotNil(t, container)
	assert.Equal(t, "mock-service", service)
	assert.Nil(t, err)
	
	// Kiểm tra rằng các phương thức được gọi đúng cách
	mockApp.AssertExpectations(t)
}
```

## Tái tạo Mock

Để tạo lại các mock, sử dụng công cụ mockery:

```bash
go install github.com/vektra/mockery/v2@latest
cd /path/to/di
mockery --name Application --output ./mocks --outpkg mocks
mockery --name ServiceProvider --output ./mocks --outpkg mocks
mockery --name ServiceProviderDeferred --output ./mocks --outpkg mocks
mockery --name ModuleLoaderContract --output ./mocks --outpkg mocks
```
