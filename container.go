// Package di cung cấp dependency injection container cho framework Fork.
//
// Container quản lý tất cả các dependencies và cung cấp các phương thức để
// đăng ký, resolve và bind các dependency.
package di

import (
	"fmt"
	"reflect"
	"sync"
)

// Container là "trái tim" của hệ thống Dependency Injection (DI) trong Fork framework.
//
// # Định nghĩa & Tiêu chuẩn
//
// Container là hiện thực chuẩn DI hiện đại, lấy cảm hứng từ các framework như Laravel, Spring, nhưng được tối ưu hóa cho Go:
//   - Tuân thủ nghiêm ngặt nguyên tắc SOLID, đặc biệt là Dependency Inversion Principle.
//   - Đảm bảo separation of concerns: tách biệt logic khởi tạo, quản lý vòng đời, và resolve dependency.
//   - Hỗ trợ Service-Repository pattern, Adapter pattern, Service Provider pattern.
//   - Cho phép mở rộng, kiểm soát, testability và maintainability tối đa cho ứng dụng Go.
//
// # Vai trò
//
//   - Là registry trung tâm cho toàn bộ dependency, service, adapter, repository, ...
//   - Quản lý binding, singleton, instance, alias, và tự động resolve dependency qua reflection.
//   - Đảm bảo mọi thành phần trong hệ thống có thể được inject, mock, hoặc thay thế dễ dàng.
//   - Là nền tảng cho mọi module, provider, middleware, ... trong Fork.
//
// # Định nghĩa cấu trúc
//
// Container quản lý các dependency thông qua các trường:
//   - bindings: map[string]BindingFunc — ánh xạ abstract type (tên logic) tới factory function khởi tạo instance.
//   - instances: map[string]interface{} — lưu trữ các singleton instance đã được khởi tạo.
//   - aliases: map[string]string — ánh xạ alias tới abstract type gốc, hỗ trợ truy cập đa tên.
//   - mu: sync.RWMutex — đảm bảo an toàn concurrent cho mọi thao tác đăng ký/resolve.
type Container struct {
	// bindings chứa các factory function tạo dependency theo abstract type.
	bindings map[string]BindingFunc

	// instances lưu trữ các singleton instance đã được khởi tạo.
	instances map[string]interface{}

	// aliases ánh xạ alias tới abstract type gốc.
	aliases map[string]string

	// mu bảo vệ mọi thao tác concurrent trên container.
	mu sync.RWMutex
}

// New khởi tạo một DI container rỗng, sẵn sàng cho việc đăng ký binding, instance, alias.
// Trả về: *Container instance mới.
func New() *Container {
	return &Container{
		bindings:  make(map[string]BindingFunc),
		instances: make(map[string]interface{}),
		aliases:   make(map[string]string),
	}
}

// Bind đăng ký một binding (factory function) cho abstract type.
//
//   - Mục đích: Cho phép đăng ký cách khởi tạo một dependency động, phục vụ cho việc resolve về sau.
//   - Logic: Lưu factory function vào map bindings, override nếu đã tồn tại.
//   - Tham số:
//   - abstract: string — tên logic của dependency (thường là interface hoặc service name).
//   - concrete: BindingFunc — factory function nhận container, trả về instance.
//   - Trả về: Không trả về.
//   - Lỗi: Nếu abstract rỗng hoặc nil, panic hoặc silent error (tùy implement).
func (c *Container) Bind(abstract string, concrete BindingFunc) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.bindings[abstract] = concrete
}

// BindIf đăng ký binding chỉ khi chưa tồn tại.
//
//   - Mục đích: Đảm bảo không override binding đã có, dùng cho module mở rộng.
//   - Logic: Kiểm tra tồn tại, chỉ bind nếu chưa có.
//   - Tham số: như Bind.
//   - Trả về: true nếu đăng ký thành công, false nếu đã tồn tại.
func (c *Container) BindIf(abstract string, concrete BindingFunc) bool {
	c.mu.RLock()
	_, exists := c.bindings[abstract]
	c.mu.RUnlock()

	if !exists {
		c.Bind(abstract, concrete)
		return true
	}

	return false
}

// singletonResolver là hàm nội bộ xử lý logic của singleton để dễ test
func (c *Container) singletonResolver(abstract string, concrete BindingFunc) interface{} {
	c.mu.RLock()
	instance, exists := c.instances[abstract]
	c.mu.RUnlock()

	if exists {
		return instance
	}

	instance = concrete(c)

	c.mu.Lock()
	c.instances[abstract] = instance
	c.mu.Unlock()

	return instance
}

// Singleton đăng ký binding singleton (chỉ tạo một instance duy nhất).
//
//   - Mục đích: Đảm bảo dependency chỉ được khởi tạo một lần duy nhất trong suốt vòng đời container.
//   - Logic: Factory function được wrap lại, lưu instance vào map instances khi lần đầu resolve.
//   - Tham số: như Bind.
//   - Trả về: Không trả về.
func (c *Container) Singleton(abstract string, concrete BindingFunc) {
	c.Bind(abstract, func(c *Container) interface{} {
		return c.singletonResolver(abstract, concrete)
	})
}

// Instance đăng ký một instance đã khởi tạo sẵn.
//
//   - Mục đích: Cho phép inject các giá trị đã tồn tại (config, logger, ...), không cần factory.
//   - Logic: Lưu trực tiếp vào map instances.
//   - Tham số:
//   - abstract: string — tên logic.
//   - instance: interface{} — giá trị đã khởi tạo.
//   - Trả về: Không trả về.
func (c *Container) Instance(abstract string, instance interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.instances[abstract] = instance
}

// Alias đăng ký một alias cho abstract type.
//
//   - Mục đích: Cho phép truy cập dependency qua nhiều tên khác nhau (ví dụ: "log" và "logger").
//   - Logic: Lưu alias vào map aliases, resolve alias sẽ trả về instance của abstract gốc.
//   - Tham số:
//   - abstract: string — tên gốc.
//   - alias: string — tên alias.
//   - Trả về: Không trả về.
func (c *Container) Alias(abstract, alias string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.aliases[alias] = abstract
}

// Make resolve một dependency từ container.
//
//   - Mục đích: Resolve instance theo abstract type đã đăng ký (binding, instance, alias).
//   - Logic: Ưu tiên alias, instance, binding. Nếu không có, trả về error.
//   - Tham số:
//   - abstract: string — tên logic.
//   - Trả về:
//   - interface{}: instance đã resolve.
//   - error: nếu không tìm thấy hoặc binding lỗi.
func (c *Container) Make(abstract string) (interface{}, error) {
	return c.make(abstract)
}

// MustMake resolve một dependency, panic nếu lỗi.
//
//   - Mục đích: Resolve instance, panic nếu không tìm thấy hoặc binding lỗi (dùng cho critical dependency).
//   - Tham số: như Make.
//   - Trả về: interface{} instance đã resolve.
//   - Lỗi: panic nếu không resolve được.
func (c *Container) MustMake(abstract string) interface{} {
	instance, err := c.make(abstract)
	if err != nil {
		panic(err)
	}
	return instance
}

// make là hiện thực nội bộ của Make
func (c *Container) make(abstract string) (interface{}, error) {
	c.mu.RLock()
	// Nếu có alias thì resolve alias trước
	if alias, exists := c.aliases[abstract]; exists {
		abstract = alias
	}

	// Nếu đã có instance thì trả về luôn
	if instance, exists := c.instances[abstract]; exists {
		c.mu.RUnlock()
		return instance, nil
	}

	// Resolve từ binding
	concrete, exists := c.bindings[abstract]
	c.mu.RUnlock()

	if !exists {
		return nil, fmt.Errorf("bind not found for: %s", abstract)
	}

	return concrete(c), nil
}

// Bound kiểm tra một abstract đã được đăng ký binding/instance/alias chưa.
//
//   - Mục đích: Hỗ trợ kiểm tra trạng thái container, phục vụ cho module động.
//   - Tham số: abstract: string.
//   - Trả về: true nếu đã đăng ký, false nếu chưa.
func (c *Container) Bound(abstract string) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()

	_, boundAsBinding := c.bindings[abstract]
	_, boundAsInstance := c.instances[abstract]
	_, boundAsAlias := c.aliases[abstract]

	return boundAsBinding || boundAsInstance || boundAsAlias
}

// Reset xóa toàn bộ binding, instance, alias khỏi container.
//
//   - Mục đích: Làm sạch container, thường dùng cho test hoặc reload.
//   - Trả về: Không trả về.
func (c *Container) Reset() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.bindings = make(map[string]BindingFunc)
	c.instances = make(map[string]interface{})
	c.aliases = make(map[string]string)
}

// Call gọi một hàm và tự động resolve các dependency qua reflection.
//
//   - Mục đích: Tự động inject các dependency vào callback function, hỗ trợ DI cho hàm tự do.
//   - Logic: Phân tích các tham số của callback, resolve từ container hoặc lấy từ additionalParams.
//   - Tham số:
//   - callback: interface{} — function cần gọi.
//   - additionalParams: ...interface{} — các tham số bổ sung (ưu tiên inject).
//   - Trả về:
//   - []interface{}: kết quả trả về của callback.
//   - error: nếu không resolve được tham số hoặc callback không hợp lệ.
//   - Lỗi: Trả về error nếu callback không phải function, hoặc không resolve được dependency.
func (c *Container) Call(callback interface{}, additionalParams ...interface{}) ([]interface{}, error) {
	callbackType := reflect.TypeOf(callback)
	if callbackType.Kind() != reflect.Func {
		return nil, fmt.Errorf("callback must be a function")
	}

	var args []reflect.Value
	for i := 0; i < callbackType.NumIn(); i++ {
		paramType := callbackType.In(i)

		// Kiểm tra xem có trong additionalParams không
		found := false
		for _, param := range additionalParams {
			paramValue := reflect.ValueOf(param)
			if paramValue.Type().AssignableTo(paramType) {
				args = append(args, paramValue)
				found = true
				break
			}
		}

		if !found {
			// Thử resolve từ container
			typeName := paramType.String()
			instance, err := c.make(typeName)
			if err != nil {
				return nil, fmt.Errorf("cannot resolve parameter %s: %v", typeName, err)
			}
			args = append(args, reflect.ValueOf(instance))
		}
	}

	returnValues := reflect.ValueOf(callback).Call(args)
	var result []interface{}
	for _, v := range returnValues {
		result = append(result, v.Interface())
	}

	return result, nil
}
