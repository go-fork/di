# Migration Guide - v0.1.2 - Container Interface

Hướng dẫn cập nhật từ phiên bản v0.1.1 lên v0.1.2 của go.fork.vn/di.

## Tổng quan thay đổi

Phiên bản v0.1.2 giới thiệu thay đổi quan trọng nhưng không phá vỡ API hiện có:

1. **Container Interface**: `Container` giờ đây là interface thay vì struct cụ thể
2. **Mocks Regeneration**: Các mock objects được tạo lại để hỗ trợ Container interface mới

Phiên bản này không phá vỡ API hiện tại, nhưng cung cấp kiến trúc linh hoạt hơn và khả năng test tốt hơn.

## Bước 1: Cập nhật thư viện

```bash
go get -u go.fork.vn/di@v0.1.2
```

## Bước 2: Kiểm tra code sử dụng trực tiếp `container` struct

Nếu bạn đã truy cập trực tiếp vào struct `container` thay vì thông qua interface, hãy cập nhật code để sử dụng interface `Container`:

**Trước đây (không khuyến nghị):**
```go
import "go.fork.vn/di"

// Truy cập trực tiếp vào container struct - không còn hoạt động
c := di.New().(*di.container)
```

**Bây giờ (khuyến nghị):**
```go
import "go.fork.vn/di"

// Luôn sử dụng Container interface
c := di.New() // c là di.Container interface
```

## Bước 3: Cập nhật các mocks

Nếu bạn đã sử dụng mock objects từ gói `go.fork.vn/di/mocks`, các mock objects giờ đây đã được regenerate để phù hợp với Container interface mới:

```go
import (
    "go.fork.vn/di/mocks"
)

func TestWithMocks() {
    // Container mock đã được cập nhật để implement Container interface mới
    mockContainer := new(mocks.Container)
    
    // Thiết lập expectations như trước đây
    mockContainer.On("Bind", "service", mock.AnythingOfType("BindingFunc")).Return()
    mockContainer.On("Make", "service").Return(myService, nil)
    
    // ...
}
```

## Bước 4: Kiểm tra lại tests

Sau khi cập nhật:
- Chạy kiểm thử unit tests
- Kiểm tra các mocks đã cập nhật

```bash
go test ./...
```

## Lợi ích của Container Interface

Việc chuyển đổi `Container` từ struct sang interface mang lại nhiều lợi ích:

1. **Mock Testing**: Dễ dàng tạo mocks cho testing
2. **Flexibility**: Cho phép tạo nhiều implementation khác nhau (ví dụ: cache-aware container, logging container)
3. **Dependency Inversion**: Tuân theo Dependency Inversion Principle triệt để hơn
4. **Decoupling**: Giảm sự phụ thuộc vào chi tiết hiện thực cụ thể
5. **Extensibility**: Dễ dàng mở rộng hệ thống với các container tùy chỉnh

## Không có Breaking Changes

Tất cả code hiện có sử dụng API công khai của `Container` sẽ tiếp tục hoạt động mà không cần thay đổi. Chỉ có code truy cập trực tiếp vào struct nội bộ `container` (không khuyến nghị) mới cần được cập nhật.

---

Liên hệ team nếu bạn gặp vấn đề khi cập nhật lên v0.1.2.
