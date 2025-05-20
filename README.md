# di

Một container Dependency Injection (DI) hiện đại, nhẹ nhàng cho Go, lấy cảm hứng từ Laravel và các framework hiện đại. Cung cấp đăng ký dịch vụ, giải quyết phụ thuộc và quản lý vòng đời cho việc xây dựng ứng dụng Go có khả năng mở rộng.

## Tính năng
- Mẫu Service Provider
- Tự động giải quyết phụ thuộc
- Binding singleton và transient
- Tải dịch vụ trì hoãn (deferred)
- API đơn giản, không cần tạo mã

## Cài đặt

```
go get github.com/go-fork/di
```

## Sử dụng

### Cơ bản

```go
import "github.com/go-fork/di"

container := di.New()
container.Bind("service", func(c *di.Container) interface{} {
    return &MyService{}
})
service := container.MustMake("service").(*MyService)
```

### Singleton

```go
// Đăng ký dịch vụ singleton (chỉ khởi tạo một lần)
container.Singleton("database", func(c *di.Container) interface{} {
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
container.Singleton("logger", func(c *di.Container) interface{} {
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
container.Singleton("userRepo", func(c *di.Container) interface{} {
    return &UserRepository{}
})
container.Singleton("userService", func(c *di.Container) interface{} {
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

Xem [doc.go](./doc.go) để biết thêm chi tiết.

## Giấy phép
MIT
